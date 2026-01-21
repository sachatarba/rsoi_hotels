package requests

type ReservationRequest struct {
	Username         string `json:"username"`
	ReservationCount int    `json:"reservation_count"`
}
