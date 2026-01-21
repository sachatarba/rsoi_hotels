package entity

import "github.com/google/uuid"

type Hotel struct {
	Id        int       `db:"id" json:"id"`
	HotelUuid uuid.UUID `db:"hotel_uid" json:"hotelUid"`
	Name      string    `db:"name" json:"name"`
	Country   string    `db:"country" json:"country"`
	City      string    `db:"city" json:"city"`
	Address   string    `db:"address" json:"address"`
	Stars     int       `db:"stars" json:"stars"`
	Price     int       `db:"price" json:"price"`
}
