package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sachatarba/rsoi_hotels/internal/payments/domain/entity"
	"log/slog"
)

type PaymentRepository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewPaymentRepository(pool *pgxpool.Pool, logger *slog.Logger) *PaymentRepository {
	return &PaymentRepository{
		pool:   pool,
		logger: logger,
	}
}

func (r *PaymentRepository) CreatePayment(ctx context.Context, paymentUid uuid.UUID, price int, status entity.PaymentStatus) error {
	query := `INSERT INTO payment (payment_uid, price, status) VALUES ($1, $2, $3)`
	_, err := r.pool.Exec(ctx, query, paymentUid, price, status)
	return err
}

func (r *PaymentRepository) GetPaymentByUid(ctx context.Context, uid uuid.UUID) (*entity.Payment, error) {
	query := `SELECT id, payment_uid, status, price FROM payment WHERE payment_uid = $1`

	var p entity.Payment
	var statusStr string

	err := r.pool.QueryRow(ctx, query, uid).Scan(&p.Id, &p.PaymentUuid, &statusStr, &p.Price)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}
	p.Status = entity.PaymentStatus(statusStr)
	return &p, nil
}

func (r *PaymentRepository) UpdatePaymentStatus(ctx context.Context, uid uuid.UUID, status entity.PaymentStatus) error {
	query := `UPDATE payment SET status = $1 WHERE payment_uid = $2`
	tag, err := r.pool.Exec(ctx, query, status, uid)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}
