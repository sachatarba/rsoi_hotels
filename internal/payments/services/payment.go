package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sachatarba/rsoi_hotels/internal/payments/domain/entity"
	"github.com/sachatarba/rsoi_hotels/internal/payments/domain/repository"
	"log/slog"
)

type PaymentService struct {
	repo   repository.IPaymentRepository
	logger *slog.Logger
}

func NewPaymentService(repo repository.IPaymentRepository, logger *slog.Logger) *PaymentService {
	return &PaymentService{
		repo:   repo,
		logger: logger,
	}
}

func (s *PaymentService) Create(ctx context.Context, price int) (uuid.UUID, error) {
	uid := uuid.New()
	err := s.repo.CreatePayment(ctx, uid, price, entity.Paid)
	if err != nil {
		s.logger.Error("Failed to create payment", "error", err)
		return uuid.Nil, err
	}
	return uid, nil
}

func (s *PaymentService) Cancel(ctx context.Context, paymentUid uuid.UUID) error {
	err := s.repo.UpdatePaymentStatus(ctx, paymentUid, entity.Canceled)
	if err != nil {
		s.logger.Error("Failed to cancel payment", "uid", paymentUid, "error", err)
		return err
	}
	return nil
}

func (s *PaymentService) GetDetails(ctx context.Context, paymentUid uuid.UUID) (*entity.Payment, error) {
	return s.repo.GetPaymentByUid(ctx, paymentUid)
}
