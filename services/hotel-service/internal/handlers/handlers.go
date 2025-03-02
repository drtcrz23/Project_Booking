package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/drtcrz23/Project_Booking/services/hotel-service/internal/model"
	"github.com/drtcrz23/Project_Booking/services/hotel-service/internal/repository"
	pb "github.com/drtcrz23/Project_Booking/services/hotel-service/pkg/api"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	DB *sql.DB
	pb.UnimplementedHotelServiceServer
}

type QueryStatus struct {
	Status string `json:"status"`
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) GetHotelById(ctx context.Context, req *pb.GetHotelRequest) (*pb.Hotel, error) {
	// Получаем отель из базы данных
	hotel, err := repository.GetHotelById(int(req.HotelId), h.DB)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении отеля: %v", err)
	}

	var grpcRooms []*pb.Room
	for _, room := range hotel.Rooms {
		grpcRoom := &pb.Room{
			Id:         int32(room.ID),
			HotelId:    int32(room.HotelId),
			RoomNumber: room.RoomNumber,
			Type:       room.Type,
			Price:      float32(room.Price), // Конвертируем из float64 в float32
			Status:     room.Status,
		}
		grpcRooms = append(grpcRooms, grpcRoom)
	}

	// Возвращаем полученные данные в формате gRPC
	return &pb.Hotel{
		Id:         int32(hotel.ID),
		Name:       hotel.Name,
		Price:      hotel.Price,
		HotelierId: int32(hotel.HotelierId),
		Rooms:      grpcRooms,
	}, nil
}

func (h *Handler) AddHotel(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}
	var request struct {
		HotelierId int             `json:"hotelier_id"`
		Hotel      model.HotelData `json:"hotel"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, "Ошибка при парсинге JSON", http.StatusBadRequest)
		return
	}
	hotelier := model.Hotelier{HotelierId: request.HotelierId}
	err = repository.AddHotel(hotelier, request.Hotel, h.DB)
	if err != nil {
		http.Error(w, "Ошибка при добавлении отеля", http.StatusBadRequest)
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
		if err := rows.Scan(&hotel.ID, &hotel.Name, &hotel.Price, &hotel.HotelierId); err != nil {
			return nil, err
		}
		rooms, err := repository.GetRoomByHotel(hotel.ID, h.DB)
		if err != nil {
			return nil, fmt.Errorf("ошибка при получении комнат для отеля с ID %d: %w", hotel.ID, err)
		}

		hotel.Rooms = rooms
		hotels = append(hotels, hotel)
	}

	return hotels, nil
}

func (h *Handler) GetHotelsByHotelier(w http.ResponseWriter, r *http.Request) {
	hotelierIDStr := r.URL.Query().Get("hotelier_id")

	hotelierID, err := strconv.Atoi(hotelierIDStr)
	if err != nil {
		http.Error(w, "Неверный ID отельера", http.StatusBadRequest)
		return
	}
	hotels, err := repository.GetHotelsByHotelier(hotelierID, h.DB)
	if err != nil {
		http.Error(w, fmt.Errorf("Ошибка при получении отеле: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hotels)
}

func (h *Handler) GetHotelByIdUsers(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) GetRoomsByHotel(w http.ResponseWriter, r *http.Request) {
	// Получаем hotel_id из параметров запроса
	hotelIDStr := r.URL.Query().Get("hotel_id")
	if hotelIDStr == "" {
		http.Error(w, "Необходимо указать hotel_id", http.StatusBadRequest)
		return
	}

	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		http.Error(w, "Неверный формат hotel_id", http.StatusBadRequest)
		return
	}

	rooms, err := repository.GetRoomByHotel(hotelID, h.DB)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при получении списка комнат: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(rooms) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, `{"message": "Комнаты для указанного отеля не найдены"}`)
		return
	}

	if err := json.NewEncoder(w).Encode(rooms); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) AddRoom(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var request struct {
		Room model.Room `json:"room"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, "Ошибка при парсинге JSON", http.StatusBadRequest)
		return
	}

	if request.Room.RoomNumber == "" || request.Room.Type == "" || request.Room.Price <= 0 {
		http.Error(w, "Обязательные поля отсутствуют или некорректны", http.StatusBadRequest)
		return
	}

	err = repository.AddRoom(request.Room, h.DB)
	if err != nil {
		errorMessage := fmt.Sprintf("Ошибка при добавлении комнаты: %v", err)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	version := QueryStatus{
		Status: "Done",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version)
}

func (h *Handler) SetRoom(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка при чтении тела запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var room model.Room
	err = json.Unmarshal(body, &room)
	if err != nil {
		http.Error(w, "Ошибка при парсинге JSON", http.StatusBadRequest)
		return
	}

	if room.RoomNumber == "" || room.Type == "" || room.Price <= 0 {
		http.Error(w, "Обязательные поля отсутствуют или некорректны", http.StatusBadRequest)
		return
	}

	err = repository.SetRoom(room, h.DB)
	if err != nil {
		http.Error(w, "Ошибка при обновлении данных комнаты", http.StatusInternalServerError)
		return
	}

	version := QueryStatus{
		Status: "Done",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version)
}

func (h *Handler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	roomIDStr := r.URL.Query().Get("room_id")
	hotelIDStr := r.URL.Query().Get("hotel_id")

	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		http.Error(w, "Неверный ID комнаты", http.StatusBadRequest)
		return
	}
	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		http.Error(w, "Неверный ID отеля", http.StatusBadRequest)
		return
	}

	err = repository.DeleteRoom(roomID, hotelID, h.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := QueryStatus{
		Status: "Done",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
