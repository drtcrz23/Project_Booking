package model

type Hotel struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Price      string `json:"Price"`
	HotelierId int    `json:"hotelier_id"`
	Rooms      []Room `json:"rooms"` // Список комнат, которые принадлежат отелю
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
