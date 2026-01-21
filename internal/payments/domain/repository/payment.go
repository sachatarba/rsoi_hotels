package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/sachatarba/rsoi_hotels/internal/payments/domain/entity"
)

type IPaymentRepository interface {
	CreatePayment(ctx context.Context, paymentUid uuid.UUID, price int, status entity.PaymentStatus) error
	GetPaymentByUid(ctx context.Context, uid uuid.UUID) (*entity.Payment, error)
	UpdatePaymentStatus(ctx context.Context, uid uuid.UUID, status entity.PaymentStatus) error
}
