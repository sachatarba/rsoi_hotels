package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sachatarba/rsoi_hotels/internal/reservations/domain/entity"
	"time"
)

type IBookingService interface {
	GetHotels(ctx context.Context, page int, size int) ([]entity.Hotel, error)
	GetHotelById(ctx context.Context, uid uuid.UUID) (*entity.Hotel, error)
	BookHotel(ctx context.Context,
		hotelUid uuid.UUID,
		username string,
		paymentUid uuid.UUID,
		hotelId int,
		startDate time.Time,
		endDate time.Time) (*entity.Reservation, error)
	GetReservations(ctx context.Context, username string) ([]entity.Reservation, error)
	GetReservationByUid(ctx context.Context, reservationUid uuid.UUID) (*entity.Reservation, error)
	CancelReservation(ctx context.Context, reservationUid uuid.UUID) error
}
