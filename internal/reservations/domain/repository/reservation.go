package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/sachatarba/rsoi_hotels/internal/reservations/domain/entity"
	"time"
)

type IReservationRepository interface {
	CreateReservation(ctx context.Context,
		reservationUid uuid.UUID,
		username string,
		paymentUid uuid.UUID,
		hotelId int,
		status entity.ReservationStatus,
		startDate time.Time,
		endDate time.Time) (*entity.Reservation, error)

	GetReservations(ctx context.Context, username string) ([]entity.Reservation, error)
	UpdateReservation(ctx context.Context, uid uuid.UUID, status entity.ReservationStatus) error
	GetReservationByUid(ctx context.Context, uid uuid.UUID) (*entity.Reservation, error)
}
