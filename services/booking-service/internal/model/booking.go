package model

type Booking struct {
	ID            int    `json:"id"`
	HotelId       int    `json:"hotel_id"`
	UserId        int    `json:"user_id"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	Price         int    `json:"price"`
	Status        string `json:"status"`
	PaymentStatus string `json:"payment_status"`
}

type UpdateBooking struct {
	ID            int    `json:"id"`
	HotelId       int    `json:"hotel_id"`
	UserId        int    `json:"user_id"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	Price         int    `json:"price"`
	Status        string `json:"status"`
	PaymentStatus string `json:"payment_status"`
}

type DeleteBooking struct {
	ID int `json:"id"`
}
