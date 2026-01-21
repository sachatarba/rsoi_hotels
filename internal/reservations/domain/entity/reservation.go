package entity

import (
	"github.com/google/uuid"
	"time"
)

type Reservation struct {
	Id              int               `db:"id" json:"-"`
	ReservationUuid uuid.UUID         `db:"reservation_uid" json:"reservationUid"`
	Username        string            `db:"username" json:"username"`
	PaymentUuid     uuid.UUID         `db:"payment_uid" json:"paymentUid"`
	HotelId         int               `db:"hotel_id" json:"hotelId"`
	Status          ReservationStatus `db:"status" json:"status"`
	StartDate       time.Time         `db:"start_date" json:"startDate"`
	EndDate         time.Time         `db:"end_date" json:"endDate"`
	HotelUid        uuid.UUID         `db:"hotel_uid" json:"hotelUid"`
}
