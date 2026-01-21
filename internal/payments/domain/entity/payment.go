package entity

import "github.com/google/uuid"

type Payment struct {
	Id          int
	PaymentUuid uuid.UUID
	Status      PaymentStatus
	Price       int
}
