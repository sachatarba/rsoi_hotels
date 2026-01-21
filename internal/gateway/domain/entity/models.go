package entity

import (
	"github.com/google/uuid"
	"time"
)

type Date time.Time

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(d).Format("2006-01-02") + `"`), nil
}

type CreateReservationRequest struct {
	HotelUid  uuid.UUID `json:"hotelUid" binding:"required"`
	StartDate string    `json:"startDate" binding:"required"`
	EndDate   string    `json:"endDate" binding:"required"`
}

type PaginationResponse struct {
	Page          int           `json:"page"`
	PageSize      int           `json:"pageSize"`
	TotalElements int           `json:"totalElements"`
	TotalPages    int           `json:"totalPages"`
	Items         []interface{} `json:"items"`
}

type HotelResponse struct {
	HotelUid uuid.UUID `json:"hotelUid"`
	Name     string    `json:"name"`
	Country  string    `json:"country"`
	City     string    `json:"city"`
	Address  string    `json:"address"`
	Stars    int       `json:"stars"`
	Price    int       `json:"price"`
}

type HotelInfo struct {
	HotelUid    uuid.UUID `json:"hotelUid"`
	Name        string    `json:"name"`
	FullAddress string    `json:"fullAddress"`
	Stars       int       `json:"stars"`
}

type LoyaltyInfoResponse struct {
	Status           string `json:"status"`
	Discount         int    `json:"discount"`
	ReservationCount int    `json:"reservationCount"`
}

type PaymentInfo struct {
	Status string `json:"status"`
	Price  int    `json:"price"`
}

type ReservationResponse struct {
	ReservationUid uuid.UUID   `json:"reservationUid"`
	Hotel          HotelInfo   `json:"hotel"`
	StartDate      Date        `json:"startDate"`
	EndDate        Date        `json:"endDate"`
	Status         string      `json:"status"`
	Payment        PaymentInfo `json:"payment"`

	PaymentUid uuid.UUID `json:"-"`
	HotelUid   uuid.UUID `json:"-"`
}

type CreateReservationResponse struct {
	ReservationUid uuid.UUID   `json:"reservationUid"`
	HotelUid       uuid.UUID   `json:"hotelUid"`
	StartDate      Date        `json:"startDate"`
	EndDate        Date        `json:"endDate"`
	Discount       int         `json:"discount"`
	Status         string      `json:"status"`
	Payment        PaymentInfo `json:"payment"`
}

type UserInfoResponse struct {
	Reservations []ReservationResponse `json:"reservations"`
	Loyalty      LoyaltyInfoResponse   `json:"loyalty"`
}

type ReservationServiceResponse struct {
	ReservationUid uuid.UUID `json:"reservationUid"`
	HotelUid       uuid.UUID `json:"hotelUid"`
	PaymentUid     uuid.UUID `json:"paymentUid"`
	Status         string    `json:"status"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	HotelId        int       `json:"hotelId"`
}
