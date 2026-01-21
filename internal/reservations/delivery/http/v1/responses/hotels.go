package responses

type HotelsPageResponse struct {
	Hotels     []HotelResponse
	Page       int `json:"page"`
	Size       int `json:"size"`
	PagesCount int `json:"total_pages"`
}
