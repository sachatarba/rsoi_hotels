package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sachatarba/rsoi_hotels/internal/gateway/domain/entity"
	"github.com/sachatarba/rsoi_hotels/pkg/circuitbreaker"
	"net/http"
	"time"
)

type ReservationRepository struct {
	client  *http.Client
	baseURL string
	cb      *circuitbreaker.CircuitBreaker
}

func NewReservationRepository(baseURL string) *ReservationRepository {
	return &ReservationRepository{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: baseURL,
		cb:      circuitbreaker.New(3, 1, 30*time.Second),
	}
}

func (r *ReservationRepository) GetHotels(ctx context.Context, page, size int) ([]entity.HotelResponse, error) {
	op := func() (interface{}, error) {
		url := fmt.Sprintf("%s/api/v1/hotels?page=%d&size=%d", r.baseURL, page, size)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := r.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("reservation service error (get hotels): status %d", resp.StatusCode)
		}

		var hotels []entity.HotelResponse
		if err := json.NewDecoder(resp.Body).Decode(&hotels); err != nil {
			return nil, err
		}
		return hotels, nil
	}

	res, err := r.cb.Execute(op)
	if err != nil {
		return nil, err
	}

	return res.([]entity.HotelResponse), nil
}

func (r *ReservationRepository) GetHotel(ctx context.Context, hotelUid uuid.UUID) (*entity.HotelResponse, error) {
	op := func() (interface{}, error) {
		url := fmt.Sprintf("%s/api/v1/hotels/%s", r.baseURL, hotelUid)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := r.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("reservation service error (get hotel): status %d", resp.StatusCode)
		}

		var hotel entity.HotelResponse
		if err := json.NewDecoder(resp.Body).Decode(&hotel); err != nil {
			return nil, err
		}
		return &hotel, nil
	}

	res, err := r.cb.Execute(op)
	if err != nil {
		return nil, err
	}

	return res.(*entity.HotelResponse), nil
}

func (r *ReservationRepository) GetUserReservations(ctx context.Context, username string) ([]entity.ReservationServiceResponse, error) {
	op := func() (interface{}, error) {
		fullUrl := fmt.Sprintf("%s/api/v1/reservations", r.baseURL)
		req, err := http.NewRequestWithContext(ctx, "GET", fullUrl, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-User-Name", username)

		resp, err := r.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("reservation service error (get user reservations): status %d", resp.StatusCode)
		}

		var rawReservations []entity.ReservationServiceResponse
		if err := json.NewDecoder(resp.Body).Decode(&rawReservations); err != nil {
			return nil, err
		}
		return rawReservations, nil
	}

	res, err := r.cb.Execute(op)
	if err != nil {
		return nil, err
	}

	return res.([]entity.ReservationServiceResponse), nil
}

func (r *ReservationRepository) GetReservation(ctx context.Context, username string, reservationUid uuid.UUID) (*entity.ReservationServiceResponse, error) {
	op := func() (interface{}, error) {
		url := fmt.Sprintf("%s/api/v1/reservations/%s", r.baseURL, reservationUid)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-User-Name", username)

		resp, err := r.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("reservation service error (get reservation): status %d", resp.StatusCode)
		}

		var resSvc entity.ReservationServiceResponse
		if err := json.NewDecoder(resp.Body).Decode(&resSvc); err != nil {
			return nil, err
		}
		return &resSvc, nil
	}

	res, err := r.cb.Execute(op)
	if err != nil {
		return nil, err
	}

	return res.(*entity.ReservationServiceResponse), nil
}

func (r *ReservationRepository) CreateReservation(ctx context.Context, username string, hotelUid uuid.UUID, startDate, endDate time.Time,
	paymentUid uuid.UUID) (*entity.ReservationServiceResponse, error) {

	payload := map[string]interface{}{
		"username":    username,
		"payment_uid": paymentUid,
		"hotel_uid":   hotelUid,
		"start_date":  startDate.Format("2006-01-02T15:04:05Z07:00"),
		"end_date":    endDate.Format("2006-01-02T15:04:05Z07:00"),
	}
	jsonBody, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/api/v1/reservations", r.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-Name", username)

	resp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("reservation service create error: status %d", resp.StatusCode)
	}

	var res entity.ReservationServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *ReservationRepository) CancelReservation(ctx context.Context, username string, reservationUid uuid.UUID) error {
	url := fmt.Sprintf("%s/api/v1/reservations/%s", r.baseURL, reservationUid)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-User-Name", username)

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("reservation service cancel error: status %d", resp.StatusCode)
	}

	return nil
}
