package services

import (
	"context"
	"errors"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/domain"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/domain/entity"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/domain/repository"
	"log/slog"
)

type LoyaltyService struct {
	repo   repository.ILoyaltyRepository
	logger *slog.Logger
}

func NewLoyaltyService(loyaltyRepository repository.ILoyaltyRepository,
	logger *slog.Logger) *LoyaltyService {
	return &LoyaltyService{
		repo:   loyaltyRepository,
		logger: logger,
	}
}

func (s *LoyaltyService) GetLoyaltyByUsername(ctx context.Context,
	username string) (*entity.Loyalty, error) {

	s.logger.Info("Getting loyalty for user", "username", username)

	loyalty, err := s.repo.GetLoyaltyByUsername(ctx, username)
	if errors.Is(err, domain.ErrLoyaltyNotFound) {
		loyalty, err := s.CreateLoyalty(ctx, username, 0)
		if err != nil {
			s.logger.Error("Error creating loyalty in loyalty service", "username", username, "err", err)
			return nil, err
		}

		return loyalty, nil
	}
	if err != nil {
		s.logger.Error("Error getting loyalty by username", "username", username, "err", err)
		return nil, err
	}

	return loyalty, nil
}

func (s *LoyaltyService) CreateLoyalty(ctx context.Context,
	username string, reservationCount int) (*entity.Loyalty, error) {
	loyalty, err := s.repo.GetLoyaltyByUsername(ctx, username)
	if err != nil && !errors.Is(err, domain.ErrLoyaltyNotFound) {
		s.logger.Error("Error getting loyalty by username", "username", username, "err", err)
		return nil, err
	}

	if loyalty != nil {
		s.logger.Error("Loyalty already exists", "username", username, "id", loyalty.Id)
		return nil, domain.ErrLoyaltyAlreadyExists
	}

	loyalty = &entity.Loyalty{
		Username: username,
	}

	loyalty.AddReservations(reservationCount)

	res, err := s.repo.CreateLoyalty(ctx, loyalty.Username, loyalty.ReservationCount(),
		loyalty.Status(), loyalty.Discount())
	if err != nil {
		s.logger.Error("Error creating loyalty", "username", username, "err", err)
		return nil, err
	}

	return res, nil
}

func (s *LoyaltyService) AddReservationCount(ctx context.Context,
	username string, reservationCount int) error {
	loyalty, err := s.repo.GetLoyaltyByUsername(ctx, username)
	if err != nil {
		s.logger.Error("Error getting loyalty by username", "username", username, "err", err)
		return err
	}

	loyalty.AddReservations(reservationCount)

	err = s.repo.UpdateLoyaltyById(ctx, loyalty.Id,
		loyalty.Username,
		loyalty.ReservationCount(),
		loyalty.Status(),
		loyalty.Discount(),
	)
	if err != nil {
		s.logger.Error("Error updating loyalty by username", "username", username, "err", err)
		return err
	}

	return nil
}
