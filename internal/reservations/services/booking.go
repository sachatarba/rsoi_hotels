package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sachatarba/rsoi_hotels/internal/reservations/domain/entity"
	"github.com/sachatarba/rsoi_hotels/internal/reservations/domain/repository"
	"log/slog"
	"time"
)

type BookingService struct {
	reservationRepo repository.IReservationRepository
	hotelsRepo      repository.IHotelRepository
	logger          *slog.Logger
}

func NewBookingService(reservationRepository repository.IReservationRepository,
	IHotelRepository repository.IHotelRepository, logger *slog.Logger) *BookingService {
	return &BookingService{
		reservationRepo: reservationRepository,
		hotelsRepo:      IHotelRepository,
		logger:          logger,
	}
}

func (b *BookingService) GetHotels(ctx context.Context, page int, size int) ([]entity.Hotel, error) {
	hotels, err := b.hotelsRepo.GetHotels(ctx, page, size)
	if err != nil {
		b.logger.Error("Error getting hotels in BookingService: ", "error", err)
		return nil, err
	}

	return hotels, nil
}

func (b *BookingService) GetHotelById(ctx context.Context, uid uuid.UUID) (*entity.Hotel, error) {
	hotel, err := b.hotelsRepo.GetHotelByUuid(ctx, uid)
	if err != nil {
		b.logger.Error("Error getting hotel by uid in BookingService: ", "uid", uid, "error", err)
		return nil, err
	}

	return hotel, nil
}

func (b *BookingService) BookHotel(ctx context.Context, hotelUid uuid.UUID, username string, paymentUid uuid.UUID, hotelId int, startDate time.Time, endDate time.Time) (*entity.Reservation, error) {
	realHotelId := hotelId
	if realHotelId == 0 {
		var err error
		realHotelId, err = b.hotelsRepo.GetHotelIdByUid(ctx, hotelUid)
		if err != nil {
			b.logger.Error("Failed to resolve hotel id by uid", "uid", hotelUid, "error", err)
			return nil, fmt.Errorf("hotel not found by uid: %w", err)
		}
	} else {
		_, err := b.hotelsRepo.GetHotelByUuid(ctx, hotelUid)
		if err != nil {
			return nil, fmt.Errorf("hotel not found: %w", err)
		}
	}

	reservationUid := uuid.New()

	res, err := b.reservationRepo.CreateReservation(ctx, reservationUid, username, paymentUid, realHotelId, entity.Paid, startDate, endDate)

	if err != nil {
		b.logger.Error("Error creating reservation in BookingService for username for hotel", "hotelId", realHotelId, "username", username, "error", err)
		return nil, err
	}

	return res, nil
}

func (b *BookingService) GetReservations(ctx context.Context, username string) ([]entity.Reservation, error) {
	reservations, err := b.reservationRepo.GetReservations(ctx, username)
	if err != nil {
		b.logger.Error("Error getting reservations in BookingService: ", "username", username, "error", err)
		return nil, err
	}

	return reservations, nil
}

func (b *BookingService) GetReservationByUid(ctx context.Context, reservationUid uuid.UUID) (*entity.Reservation, error) {
	reservation, err := b.reservationRepo.GetReservationByUid(ctx, reservationUid)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (b *BookingService) CancelReservation(ctx context.Context, reservationUid uuid.UUID) error {
	err := b.reservationRepo.UpdateReservation(ctx, reservationUid, entity.Canceled)
	if err != nil {
		return err
	}

	return nil
}
