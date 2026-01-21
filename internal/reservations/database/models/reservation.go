package models

import (
	"github.com/google/uuid"
	"time"
)

type Reservation struct {
	Id              int
	ReservationUuid uuid.UUID
	Username        string
	PaymentUuid     uuid.UUID
	HotelId         int
	Status          string
	StartDate       time.Time
	EndDate         time.Time
}
