package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sachatarba/rsoi_hotels/internal/gateway/domain/entity"
)

type IGatewayService interface {
	GetUserReservations(ctx context.Context, username string) ([]entity.ReservationResponse, error)
	GetHotels(ctx context.Context, page, size int) (*entity.PaginationResponse, error)
	GetUserInfo(ctx context.Context, username string) (*entity.UserInfoResponse, error)
	GetReservation(ctx context.Context, username string, reservationUid uuid.UUID) (*entity.ReservationResponse, error)
	BookHotel(ctx context.Context, username string, req entity.CreateReservationRequest) (*entity.CreateReservationResponse, error)
	CancelReservation(ctx context.Context, username string, reservationUid uuid.UUID) error
	GetLoyalty(ctx context.Context, username string) (*entity.LoyaltyInfoResponse, error)
}
