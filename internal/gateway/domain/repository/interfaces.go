package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/sachatarba/rsoi_hotels/internal/gateway/domain/entity"
	"time"
)

type ReservationRepo interface {
	GetHotels(ctx context.Context, page, size int) ([]entity.HotelResponse, error)
	GetHotel(ctx context.Context, hotelUid uuid.UUID) (*entity.HotelResponse, error)

	GetUserReservations(ctx context.Context, username string) ([]entity.ReservationServiceResponse, error)
	GetReservation(ctx context.Context, username string, reservationUid uuid.UUID) (*entity.ReservationServiceResponse, error)
	CreateReservation(ctx context.Context, username string, hotelUid uuid.UUID, startDate, endDate time.Time, paymentUid uuid.UUID) (*entity.ReservationServiceResponse, error)
	CancelReservation(ctx context.Context, username string, reservationUid uuid.UUID) error
}

type PaymentRepo interface {
	CreatePayment(ctx context.Context, price int) (uuid.UUID, error)
	CancelPayment(ctx context.Context, paymentUid uuid.UUID) error
	GetPayment(ctx context.Context, paymentUid uuid.UUID) (*entity.PaymentInfo, error)
}

type LoyaltyRepo interface {
	GetLoyalty(ctx context.Context, username string) (*entity.LoyaltyInfoResponse, error)
	UpdateLoyaltyCount(ctx context.Context, username string, countChange int) error
}
