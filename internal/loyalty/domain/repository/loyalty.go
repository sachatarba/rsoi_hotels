package repository

import (
	"context"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/domain/entity"
)

type ILoyaltyRepository interface {
	GetLoyaltyByUsername(ctx context.Context, username string) (*entity.Loyalty, error)
	UpdateLoyaltyById(ctx context.Context, id int, username string, reservationCount int,
		status entity.LoyaltyStatus, discount int) error

	CreateLoyalty(ctx context.Context, username string, reservationCount int,
		status entity.LoyaltyStatus, discount int) (*entity.Loyalty, error)
}
