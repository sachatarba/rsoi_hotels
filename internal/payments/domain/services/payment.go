package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sachatarba/rsoi_hotels/internal/payments/domain/entity"
)

type IPaymentService interface {
	Create(ctx context.Context, price int) (uuid.UUID, error)
	Cancel(ctx context.Context, paymentUid uuid.UUID) error
	GetDetails(ctx context.Context, paymentUid uuid.UUID) (*entity.Payment, error)
}
