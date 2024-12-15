package model

type User struct {
	Id      int    `json:"user_id"`
	Name    string `json:"user_name"`
	Surname string `json:"surname"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Balance int32  `json:"balance"`
}
