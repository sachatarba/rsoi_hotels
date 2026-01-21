package errors

import "errors"

var ErrLoyaltyServiceUnavailable = errors.New("Loyalty Service unavailable")
var ErrPaymentServiceUnavailable = errors.New("Payment Service unavailable")
var ErrReservationServiceUnavailable = errors.New("Reservation Service unavailable")
