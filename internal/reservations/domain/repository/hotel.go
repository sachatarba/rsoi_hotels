package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/sachatarba/rsoi_hotels/internal/reservations/domain/entity"
)

type IHotelRepository interface {
	GetHotelByUuid(ctx context.Context, id uuid.UUID) (*entity.Hotel, error)
	GetHotels(ctx context.Context, page int, size int) ([]entity.Hotel, error)
	GetHotelIdByUid(ctx context.Context, hotelUid uuid.UUID) (int, error)
}
