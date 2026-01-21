package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sachatarba/rsoi_hotels/internal/reservations/domain"
	"github.com/sachatarba/rsoi_hotels/internal/reservations/domain/entity"
	"log/slog"
)

type HotelRepository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewHotelRepository(pool *pgxpool.Pool, logger *slog.Logger) *HotelRepository {
	return &HotelRepository{
		pool:   pool,
		logger: logger,
	}
}

func (h *HotelRepository) GetHotelByUuid(ctx context.Context, uid uuid.UUID) (*entity.Hotel, error) {
	query := `SELECT * FROM hotels WHERE hotel_uid = $1`

	rows, _ := h.pool.Query(ctx, query, uid)
	hotel, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Hotel])

	if errors.Is(err, pgx.ErrNoRows) {
		h.logger.Error("No hotel found with hotel_uid", "hotel_uid", uid)
		return nil, domain.ErrHotelNotFound
	}
	if err != nil {
		h.logger.Error("Error getting hotel by hotel_uid:", "hotel_uid", uid, "err", err)
		return nil, err
	}

	return &hotel, nil
}

func (h *HotelRepository) GetHotels(ctx context.Context, page int, size int) ([]entity.Hotel, error) {
	query := `SELECT * FROM hotels LIMIT $1 OFFSET $2`

	offset := (page - 1) * size
	limit := size

	rows, err := h.pool.Query(ctx, query, limit, offset)
	if err != nil {
		h.logger.Error("Error getting hotels", "page", page, "size", size, "err", err)
		return nil, err
	}
	defer rows.Close()

	hotels, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Hotel])
	if err != nil {
		h.logger.Error("Error collecting hotels", "page", page, "size", size, "err", err)
		return nil, err
	}

	return hotels, nil
}

func (h *HotelRepository) GetHotelIdByUid(ctx context.Context, hotelUid uuid.UUID) (int, error) {
	query := `SELECT id FROM hotels WHERE hotel_uid = $1`
	var id int
	err := h.pool.QueryRow(ctx, query, hotelUid).Scan(&id)
	if err != nil {
		h.logger.Error("Error getting hotel id by uid", "uid", hotelUid, "err", err)
		return 0, err
	}
	return id, nil
}
