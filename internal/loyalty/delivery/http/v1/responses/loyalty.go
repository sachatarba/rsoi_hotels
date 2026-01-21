package responses

type LoyaltyResponse struct {
	Id               int    `json:"id"`
	Username         string `json:"username"`
	ReservationCount int    `json:"reservationCount"`
	Status           string `json:"status"`
	Discount         int    `json:"discount"`
}
