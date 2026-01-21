package requests

import (
	"github.com/google/uuid"
	"time"
)

type CreateReservationRequest struct {
	Username   string    `json:"username"`
	PaymentUid uuid.UUID `json:"payment_uid"`
	HotelId    int       `json:"hotel_id"`
	HotelUid   uuid.UUID `json:"hotel_uid"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
}
