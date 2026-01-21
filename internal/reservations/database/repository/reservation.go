package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sachatarba/rsoi_hotels/internal/reservations/domain"
	"github.com/sachatarba/rsoi_hotels/internal/reservations/domain/entity"
	"log/slog"
	"time"
)

type ReservationRepository struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewReservationRepository(pool *pgxpool.Pool, logger *slog.Logger) *ReservationRepository {
	return &ReservationRepository{
		pool:   pool,
		logger: logger,
	}
}

func (r *ReservationRepository) CreateReservation(ctx context.Context,
	reservationUid uuid.UUID,
	username string,
	paymentUid uuid.UUID,
	hotelId int,
	status entity.ReservationStatus,
	startDate time.Time, endDate time.Time) (*entity.Reservation, error) {

	query := "INSERT INTO reservation (reservation_uid, username, payment_uid, hotel_id, status, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	var id int
	err := r.pool.
		QueryRow(ctx, query, reservationUid, username, paymentUid, hotelId, status, startDate, endDate).
		Scan(&id)
	if err != nil {
		r.logger.Error("Error inserting reservation", "err", err)
		return nil, err
	}

	return &entity.Reservation{
		Id:              id,
		ReservationUuid: reservationUid,
		Username:        username,
		PaymentUuid:     paymentUid,
		HotelId:         hotelId,
		Status:          status,
		StartDate:       startDate,
		EndDate:         endDate,
	}, nil
}

func (r *ReservationRepository) GetReservations(ctx context.Context, username string) ([]entity.Reservation, error) {
	query := `
		SELECT r.id, r.reservation_uid, r.username, r.payment_uid, r.hotel_id, r.status, r.start_date, r.end_date, h.hotel_uid
		FROM reservation r
		JOIN hotels h ON r.hotel_id = h.id
		WHERE r.username = $1
	`

	rows, err := r.pool.Query(ctx, query, username)
	if err != nil {
		r.logger.Error("Error getting reservations", "err", err)
		return nil, err
	}

	reservations, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Reservation])

	return reservations, nil
}

func (r *ReservationRepository) UpdateReservation(ctx context.Context, uid uuid.UUID, status entity.ReservationStatus) error {
	query := "UPDATE reservation SET status = $1 WHERE reservation_uid = $2"

	tag, err := r.pool.Exec(ctx, query, status, uid)
	if err != nil {
		r.logger.Error("Error updating reservation by reservation_uid", "reservation_uid", uid, "err", err)
		return fmt.Errorf("error updating reservation by reservation_uid: %s: %w", uid, err)
	}

	if tag.RowsAffected() == 0 {
		r.logger.Error("Error reservation to update not found in database", "reservation_uid", uid, "err", err)
		return domain.ErrReservationNotFound
	}

	return nil
}

func (r *ReservationRepository) GetReservationByUid(ctx context.Context, uid uuid.UUID) (*entity.Reservation, error) {
	query := `
		SELECT r.id, r.reservation_uid, r.username, r.payment_uid, r.hotel_id, r.status, r.start_date, r.end_date, h.hotel_uid
		FROM reservation r
		JOIN hotels h ON r.hotel_id = h.id
		WHERE r.reservation_uid = $1
	`

	var res entity.Reservation
	err := r.pool.QueryRow(ctx, query, uid).Scan(
		&res.Id, &res.ReservationUuid, &res.Username, &res.PaymentUuid, &res.HotelId, &res.Status, &res.StartDate, &res.EndDate, &res.HotelUid,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrReservationNotFound
	}
	if err != nil {
		r.logger.Error("Error getting reservation by uid", "uid", uid, "err", err)
		return nil, err
	}

	res.Status = entity.NewReservationStatus(res.Status.String())
	return &res, nil
}
