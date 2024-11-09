package model

type Hotel struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price string `json:"Price"`
}

type HotelData struct {
	Name  string `json:"name"`
	Price int    `json:"Price"`
}

type UpdateData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"Price"`
}
