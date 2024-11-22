package handlers

import (
	"HotelService/internal/model"
	"HotelService/internal/repository"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
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

func (h *Handler) GetHotelById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	// Преобразуем ID в число
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	hotel, err := repository.GetHotelById(id, h.DB)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, fmt.Sprintf("Hotel with ID %d not found", id), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve hotel", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(hotel); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
