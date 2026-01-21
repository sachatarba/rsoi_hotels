package services

import (
	"context"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/domain/entity"
)

type ILoyaltyService interface {
	GetLoyaltyByUsername(ctx context.Context, username string) (*entity.Loyalty, error)
	CreateLoyalty(ctx context.Context, username string, reservationCount int) (*entity.Loyalty, error)
	AddReservationCount(ctx context.Context, username string, reservationCount int) error
}
