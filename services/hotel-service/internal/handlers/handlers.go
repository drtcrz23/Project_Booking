package handlers

import (
	"HotelService/internal/model"
	"HotelService/internal/repository"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
)

type Handler struct {
	DB *sql.DB
}

type QueryStatus struct {
	Status string `json:"status"`
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) AddHotel(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}
	var hotel model.HotelData
	err = json.Unmarshal(body, &hotel)
	if err != nil {
		http.Error(w, "Ошибка при парсинге JSON", http.StatusBadRequest)
		return
	}

	err = repository.AddHotel(hotel, h.DB)
	if err != nil {
		http.Error(w, "Ошибка при добавление отеля", http.StatusBadRequest)
		return
	}

	version := QueryStatus{
		Status: "Done",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version)
}

func (h *Handler) SetHotel(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}
	var hotel model.UpdateData
	err_json := json.Unmarshal(body, &hotel)
	if err_json != nil {
		http.Error(w, "Ошибка при парсинге JSON", http.StatusBadRequest)
		return
	}

	err = repository.SetHotel(hotel, h.DB)
	if err != nil {
		http.Error(w, "Ошибка при добавление отеля", http.StatusBadRequest)
		return
	}

	version := QueryStatus{
		Status: "Done",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version)
}

func (h *Handler) GetAllHotels(w http.ResponseWriter, r *http.Request) ([]model.Hotel, error) {
	var hotels []model.Hotel
	rows, err := h.DB.Query("SELECT * FROM hotel")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var hotel model.Hotel
		if err := rows.Scan(&hotel.ID, &hotel.Name, &hotel.Price); err != nil {
			return nil, err
		}
		hotels = append(hotels, hotel)
	}

	return hotels, nil
}
