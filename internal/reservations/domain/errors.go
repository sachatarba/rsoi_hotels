package domain

import "errors"

var ErrHotelNotFound = errors.New("hotel not found by hotel_uid")
var ErrReservationNotFound = errors.New("reservation not found by reservation_uid")
