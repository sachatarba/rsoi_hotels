package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sachatarba/rsoi_hotels/internal/gateway/domain/entity"
	domainErrors "github.com/sachatarba/rsoi_hotels/internal/gateway/domain/errors"
	"github.com/sachatarba/rsoi_hotels/pkg/circuitbreaker"
	"net/http"
	"time"
)

type PaymentRepository struct {
	client  *http.Client
	baseURL string
	cb      *circuitbreaker.CircuitBreaker
}

func NewPaymentRepository(baseURL string) *PaymentRepository {
	return &PaymentRepository{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: baseURL,
		cb:      circuitbreaker.New(3, 1, 30*time.Second),
	}
}

func (r *PaymentRepository) CreatePayment(ctx context.Context, price int) (uuid.UUID, error) {
	body := map[string]int{"price": price}
	jsonBody, _ := json.Marshal(body)
	url := fmt.Sprintf("%s/api/v1/payments", r.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %w", domainErrors.ErrPaymentServiceUnavailable, err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.client.Do(req)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %w", domainErrors.ErrPaymentServiceUnavailable, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return uuid.Nil, fmt.Errorf("%w: status %d", domainErrors.ErrPaymentServiceUnavailable, resp.StatusCode)
	}
	var response struct {
		PaymentUid uuid.UUID `json:"payment_uid"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return uuid.Nil, fmt.Errorf("%w: %w", domainErrors.ErrPaymentServiceUnavailable, err)
	}
	return response.PaymentUid, nil
}

func (r *PaymentRepository) CancelPayment(ctx context.Context, paymentUid uuid.UUID) error {
	url := fmt.Sprintf("%s/api/v1/payments/%s", r.baseURL, paymentUid)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("%w: %w", domainErrors.ErrPaymentServiceUnavailable, err)
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", domainErrors.ErrPaymentServiceUnavailable, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: status %d", domainErrors.ErrPaymentServiceUnavailable, resp.StatusCode)
	}
	return nil
}

func (r *PaymentRepository) GetPayment(ctx context.Context, paymentUid uuid.UUID) (*entity.PaymentInfo, error) {
	op := func() (interface{}, error) {
		url := fmt.Sprintf("%s/api/v1/payments/%s", r.baseURL, paymentUid)
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
			return nil, fmt.Errorf("status %d", resp.StatusCode)
		}
		var tempResp struct {
			Status string `json:"Status"`
			Price  int    `json:"Price"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&tempResp); err != nil {
			return nil, err
		}
		return &entity.PaymentInfo{Status: tempResp.Status, Price: tempResp.Price}, nil
	}
	res, err := r.cb.Execute(op)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domainErrors.ErrPaymentServiceUnavailable, err)
	}
	return res.(*entity.PaymentInfo), nil
}
