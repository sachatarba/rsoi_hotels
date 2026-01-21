package models

type Loyalty struct {
	Id               int
	Username         string
	ReservationCount int
	Status           string
	Discount         int
}
