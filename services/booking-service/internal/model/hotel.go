package model

type Hotel struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Price      string `json:"Price"`
	HotelierId int    `json:"hotelier_id"`
	Rooms      []Room `json:"rooms"` // Список комнат, которые принадлежат отелю
}

type Room struct {
	ID         int     `json:"id"`
	HotelId    int     `json:"hotel_id"`
	RoomNumber string  `json:"room_number"`
	Type       string  `json:"type_room"` // Тип комнаты
	Price      float64 `json:"price"`
	Status     string  `json:"status"` // Статус комнаты (available, booked, maintenance)
}
