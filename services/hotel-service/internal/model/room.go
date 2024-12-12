package model

type Room struct {
	ID         int     `json:"id"`
	HotelId    int     `json:"hotel_id"`
	RoomNumber string  `json:"room_number"`
	Type       string  `json:"type_room"` // Тип комнаты
	Price      float64 `json:"price"`
	Status     string  `json:"status"` // Статус комнаты (available, booked, maintenance)
}
