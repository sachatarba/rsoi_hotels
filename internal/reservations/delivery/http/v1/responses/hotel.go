package responses

type HotelResponse struct {
	Id       int    `json:"id"`
	HotelUid int    `json:"hotel_uid"`
	Name     string `json:"name"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Stars    int    `json:"stars"`
	Price    int    `json:"price"`
}
