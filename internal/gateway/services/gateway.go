package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sachatarba/rsoi_hotels/internal/gateway/domain/entity"
	"github.com/sachatarba/rsoi_hotels/internal/gateway/domain/repository"
	"github.com/sachatarba/rsoi_hotels/pkg/circuitbreaker"
	"github.com/sachatarba/rsoi_hotels/pkg/queue"
	"log/slog"
	"math"
	"time"
)

type GatewayService struct {
	reservationRepo repository.ReservationRepo
	paymentRepo     repository.PaymentRepo
	loyaltyRepo     repository.LoyaltyRepo
	logger          *slog.Logger
	queue           *queue.Queue
}

func NewGatewayService(
	reservationRepo repository.ReservationRepo,
	paymentRepo repository.PaymentRepo,
	loyaltyRepo repository.LoyaltyRepo,
	logger *slog.Logger,
	queue *queue.Queue,
) *GatewayService {
	return &GatewayService{
		reservationRepo: reservationRepo,
		paymentRepo:     paymentRepo,
		loyaltyRepo:     loyaltyRepo,
		logger:          logger,
		queue:           queue,
	}
}

func (s *GatewayService) GetHotels(ctx context.Context, page, size int) (*entity.PaginationResponse, error) {
	hotels, err := s.reservationRepo.GetHotels(ctx, page, size)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(hotels))
	for i, v := range hotels {
		items[i] = v
	}

	return &entity.PaginationResponse{
		Page:          page,
		PageSize:      size,
		TotalElements: len(hotels),
		Items:         items,
	}, nil
}

func (s *GatewayService) GetUserInfo(ctx context.Context, username string) (*entity.UserInfoResponse, error) {
	loyalty, err := s.loyaltyRepo.GetLoyalty(ctx, username)
	if err != nil {
		s.logger.Error("Failed to get loyalty for GetUserInfo", "error", err)
		loyalty = &entity.LoyaltyInfoResponse{
			Status: "UNAVAILABLE",
		}
	}

	rawReservations, err := s.reservationRepo.GetUserReservations(ctx, username)
	if err != nil {
		return nil, err
	}

	reservations := make([]entity.ReservationResponse, 0, len(rawReservations))
	for _, raw := range rawReservations {
		enriched, _ := s.enrichReservation(ctx, raw)
		reservations = append(reservations, enriched)
	}

	return &entity.UserInfoResponse{
		Reservations: reservations,
		Loyalty:      *loyalty,
	}, nil
}

func (s *GatewayService) GetReservation(ctx context.Context, username string, reservationUid uuid.UUID) (*entity.ReservationResponse, error) {
	raw, err := s.reservationRepo.GetReservation(ctx, username, reservationUid)
	if err != nil {
		return nil, err
	}

	enriched, err := s.enrichReservation(ctx, *raw)
	if err != nil {
		if errors.Is(err, circuitbreaker.ErrCircuitOpen) {
			return &enriched, nil
		}
		return nil, err
	}
	return &enriched, nil
}

func (s *GatewayService) BookHotel(ctx context.Context, username string, req entity.CreateReservationRequest) (*entity.CreateReservationResponse, error) {
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}
	endDate, err := time.Parse(layout, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	hotel, err := s.reservationRepo.GetHotel(ctx, req.HotelUid)
	if err != nil {
		return nil, fmt.Errorf("hotel is critical, not found: %w", err)
	}

	nights := int(math.Ceil(endDate.Sub(startDate).Hours() / 24))
	if nights <= 0 {
		return nil, errors.New("invalid date range")
	}
	totalPrice := nights * hotel.Price

	loyalty, err := s.loyaltyRepo.GetLoyalty(ctx, username)
	discount := 0
	if err == nil && loyalty != nil {
		discount = loyalty.Discount
	} else {
		s.logger.Warn("Could not get loyalty info, proceeding without discount", "username", username, "error", err)
	}

	finalPrice := int(float64(totalPrice) * (1 - float64(discount)/100.0))

	var paymentUid uuid.UUID
	compensations := make([]func(context.Context), 0)

	paymentUid, err = s.paymentRepo.CreatePayment(ctx, finalPrice)
	if err != nil {
		s.rollback(compensations)
		return nil, fmt.Errorf("payment is critical, failed: %w", err)
	}
	compensations = append(compensations, func(ctx context.Context) {
		_ = s.paymentRepo.CancelPayment(ctx, paymentUid)
	})

	resRaw, err := s.reservationRepo.CreateReservation(ctx, username, req.HotelUid, startDate, endDate, paymentUid)
	if err != nil {
		s.rollback(compensations)
		return nil, err
	}

	err = s.loyaltyRepo.UpdateLoyaltyCount(context.Background(), username, 1)
	if err != nil {
		s.logger.Error("Loyalty service failed, rolling back payment", "username", username, "error", err)
		s.rollback(compensations)
		return nil, fmt.Errorf("loyalty service failed: %w", err)
	}

	return &entity.CreateReservationResponse{
		ReservationUid: resRaw.ReservationUid,
		HotelUid:       req.HotelUid,
		StartDate:      entity.Date(startDate),
		EndDate:        entity.Date(endDate),
		Discount:       discount,
		Status:         "PAID",
		Payment: entity.PaymentInfo{
			Status: "PAID",
			Price:  finalPrice,
		},
	}, nil
}

func (s *GatewayService) CancelReservation(ctx context.Context, username string, reservationUid uuid.UUID) error {
	raw, err := s.reservationRepo.GetReservation(ctx, username, reservationUid)
	if err != nil {
		return err
	}

	err = s.reservationRepo.CancelReservation(ctx, username, reservationUid)
	if err != nil {
		return err
	}

	if raw.PaymentUid != uuid.Nil {
		err = s.paymentRepo.CancelPayment(ctx, raw.PaymentUid)
		if err != nil {
			s.logger.Error("Failed to cancel payment, rolling back reservation status", "uid", raw.PaymentUid, "error", err)
			return fmt.Errorf("payment service is critical, failed to cancel: %w", err)
		}
	}

	loyaltyUpdateTask := func() error {
		return s.loyaltyRepo.UpdateLoyaltyCount(context.Background(), username, -1)
	}
	err = loyaltyUpdateTask()
	if err != nil {
		s.logger.Error("Failed to decrease loyalty, adding to background queue", "error", err)
		s.queue.Add(loyaltyUpdateTask)
	}

	return nil
}

func (s *GatewayService) GetLoyalty(ctx context.Context, username string) (*entity.LoyaltyInfoResponse, error) {
	return s.loyaltyRepo.GetLoyalty(ctx, username)
}

func (s *GatewayService) enrichReservation(ctx context.Context, raw entity.ReservationServiceResponse) (entity.ReservationResponse, error) {
	res := entity.ReservationResponse{
		ReservationUid: raw.ReservationUid,
		Status:         raw.Status,
		StartDate:      entity.Date(raw.StartDate),
		EndDate:        entity.Date(raw.EndDate),
		PaymentUid:     raw.PaymentUid,
		HotelUid:       raw.HotelUid,
	}

	if raw.HotelUid != uuid.Nil {
		hotel, err := s.reservationRepo.GetHotel(ctx, raw.HotelUid)
		if err == nil {
			res.Hotel = entity.HotelInfo{
				HotelUid:    hotel.HotelUid,
				Name:        hotel.Name,
				FullAddress: fmt.Sprintf("%s, %s, %s", hotel.Country, hotel.City, hotel.Address),
				Stars:       hotel.Stars,
			}
		}
	}

	if raw.PaymentUid != uuid.Nil {
		payment, err := s.paymentRepo.GetPayment(ctx, raw.PaymentUid)
		if err != nil {
			s.logger.Warn("Failed to get payment info on enrichment", "payment_uid", raw.PaymentUid, "error", err)
		} else {
			res.Payment = *payment
		}
	}

	return res, nil
}

func (s *GatewayService) GetUserReservations(ctx context.Context, username string) ([]entity.ReservationResponse, error) {
	rawReservations, err := s.reservationRepo.GetUserReservations(ctx, username)
	if err != nil {
		return nil, err
	}

	reservations := make([]entity.ReservationResponse, 0, len(rawReservations))
	for _, raw := range rawReservations {
		enriched, _ := s.enrichReservation(ctx, raw)
		reservations = append(reservations, enriched)
	}

	return reservations, nil
}

func (s *GatewayService) rollback(compensations []func(context.Context)) {
	ctx := context.Background()
	for i := len(compensations) - 1; i >= 0; i-- {
		compensations[i](ctx)
	}
}
