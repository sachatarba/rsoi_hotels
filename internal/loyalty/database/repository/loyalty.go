package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/database/models"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/domain"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/domain/entity"
	"log/slog"
)

type LoyaltyRepository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewLoyaltyRepository(pool *pgxpool.Pool, logger *slog.Logger) *LoyaltyRepository {
	return &LoyaltyRepository{
		pool:   pool,
		logger: logger,
	}
}

func (l *LoyaltyRepository) GetLoyaltyByUsername(ctx context.Context, username string) (*entity.Loyalty, error) {
	query := `SELECT id, username, reservation_count, status, discount FROM loyalty WHERE username = $1`

	var result models.Loyalty
	err := l.pool.
		QueryRow(ctx, query, username).
		Scan(&result.Id, &result.Username, &result.ReservationCount, &result.Status, &result.Discount)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLoyaltyNotFound
		}
		l.logger.Error("error getting loyalty by username", "username", username, "err", err)
		return nil, err
	}

	loyalty, err := entity.NewLoyalty(result.Id,
		result.Username,
		result.ReservationCount,
		entity.NewLoyaltyStatus(result.Status),
		result.Discount,
	)

	return loyalty, err
}

func (l *LoyaltyRepository) UpdateLoyaltyById(ctx context.Context, id int, username string, reservationCount int,
	status entity.LoyaltyStatus, discount int) error {
	query := `UPDATE loyalty SET username = $1, reservation_count = $2, status = $3, discount = $4 WHERE id = $5`

	tag, err := l.pool.Exec(ctx, query, username, reservationCount, status, discount, id)
	if err != nil {
		l.logger.Error("error updating loyalty by id", "id", id, "err", err)
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrLoyaltyNotFound
	}

	return nil
}

func (l *LoyaltyRepository) CreateLoyalty(ctx context.Context, username string, reservationCount int,
	status entity.LoyaltyStatus, discount int) (*entity.Loyalty, error) {
	query := `INSERT INTO loyalty (username, reservation_count, status, discount) VALUES ($1, $2, $3, $4) RETURNING id`

	var id int
	err := l.pool.
		QueryRow(ctx, query, username, reservationCount, status, discount).
		Scan(&id)
	if err != nil {
		l.logger.Error("error creating loyalty by username", "username", username, "err", err)
		return nil, err
	}

	loyalty, err := entity.NewLoyalty(id, username, reservationCount, status, discount)
	if err != nil {
		return nil, err
	}

	return loyalty, nil
}
