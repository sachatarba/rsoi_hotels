package models

import "github.com/google/uuid"

type Hotel struct {
	Id        int
	HotelUuid uuid.UUID
	Name      string
	Country   string
	City      string
	Address   string
	Stars     int
	Price     int
}
