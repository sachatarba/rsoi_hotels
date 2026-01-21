package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sachatarba/rsoi_hotels/internal/gateway/domain/entity"
	domainErrors "github.com/sachatarba/rsoi_hotels/internal/gateway/domain/errors"
	"github.com/sachatarba/rsoi_hotels/pkg/circuitbreaker"
	"net/http"
	"net/url"
	"time"
)

type LoyaltyRepository struct {
	client  *http.Client
	baseURL string
	cb      *circuitbreaker.CircuitBreaker
}

func NewLoyaltyRepository(baseURL string) *LoyaltyRepository {
	return &LoyaltyRepository{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: baseURL,
		cb:      circuitbreaker.New(3, 1, 30*time.Second),
	}
}

func (r *LoyaltyRepository) GetLoyalty(ctx context.Context, username string) (*entity.LoyaltyInfoResponse, error) {
	op := func() (interface{}, error) {
		safeUsername := url.QueryEscape(username)
		fullUrl := fmt.Sprintf("%s/api/v1/loyalties?username=%s", r.baseURL, safeUsername)
		req, err := http.NewRequestWithContext(ctx, "GET", fullUrl, nil)
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
		var loyalty entity.LoyaltyInfoResponse
		if err := json.NewDecoder(resp.Body).Decode(&loyalty); err != nil {
			return nil, err
		}
		return &loyalty, nil
	}
	res, err := r.cb.Execute(op)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", domainErrors.ErrLoyaltyServiceUnavailable, err)
	}
	return res.(*entity.LoyaltyInfoResponse), nil
}

func (r *LoyaltyRepository) UpdateLoyaltyCount(ctx context.Context, username string, countChange int) error {
	body := map[string]interface{}{
		"username":          username,
		"reservation_count": countChange,
	}
	jsonBody, _ := json.Marshal(body)
	fullUrl := fmt.Sprintf("%s/api/v1/loyalties/reservations", r.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", fullUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("%w: %w", domainErrors.ErrLoyaltyServiceUnavailable, err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", domainErrors.ErrLoyaltyServiceUnavailable, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: status %d", domainErrors.ErrLoyaltyServiceUnavailable, resp.StatusCode)
	}
	return nil
}
