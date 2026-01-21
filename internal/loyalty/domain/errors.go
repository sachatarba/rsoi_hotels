package domain

import "errors"

var ErrLoyaltyAlreadyExists = errors.New("loyalty with same username already exists")

var ErrLoyaltyStatusNotRespondReservations = errors.New("loyalty status do not respond reservations count")

var ErrDiscountNotRespondLoyaltyStatus = errors.New("discount status do not respond loyalty status")

var ErrLoyaltyNotFound = errors.New("loyalty not found by username")
