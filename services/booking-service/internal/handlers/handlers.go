package handlers

import (
	"BookingService/internal/kafkaProduceTools"
	"BookingService/internal/model"
	"BookingService/internal/repository"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	pb "github.com/drtcrz23/Project_Booking/services/grpc"
	"io"
	"net/http"
	"strconv"
)

type Handler struct {
	DB          *sql.DB
	Producer    *kafkaProduceTools.Producer
	HotelClient pb.HotelServiceClient
}

type QueryStatus struct {
	Status string `json:"status"`
}

func NewHandler(db *sql.DB, producer *kafkaProduceTools.KafkaProducer, hotelClient pb.HotelServiceClient) *Handler {
	return &Handler{DB: db, Producer: producer, HotelClient: hotelClient}
}

func (handler *Handler) AddBooking(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}
	// TODO
	// Буду добавлять метод получения отеля по grpc, он будет юзаться еще и в update
	var booking model.Booking
	err = json.Unmarshal(body, &booking)
	if err != nil {
		http.Error(w, "Ошибка при парсинге JSON", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	hotel, err := handler.HotelClient.GetHotelById(ctx, &pb.GetHotelRequest{HotelId: int32(booking.HotelId)})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve hotel: %v", err), http.StatusInternalServerError)
		return
	}
	var room model.Room
	for _, cur := range hotel.Rooms {
		if cur.Id == int32(booking.RoomId) {
			room = ConvertToModelRoom(cur)
			break
		}
	}
	fmt.Printf("Retrieved hotel: %v\n", hotel.Name)
	//hotel, err := handler.getHotelByGRPC(booking.HotelId)
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("Failed to retrieve hotel: %v", err), http.StatusInternalServerError)
	//	return
	//}

	err = repository.AddBooking(booking, room, handler.DB)
	if err != nil {
		http.Error(w, "Ошибка при добавление бронирования", http.StatusBadRequest)
		return
	}

	// ДОБАВИТЬ ИНИЦИАЛИЗАЦИЮ ЮЗЕРА

	user := model.User{Id:999,
	Name:"fakename",
	Surname :"fakesurname",
	Phone:"fakenumber",
	Email:"fakeemail@gmail.com",
	Balance: 999}

	message_text := CreateTextMessageForBookingEvent(&booking, &room, hotel.Name, &user)

	err = handler.Producer.SendMessage(user.Email, message_text)
	if err != nil {
		http.Error(w, "Ошибка при отправке события в Kafka", http.StatusInternalServerError)
		return
	}

	version := QueryStatus{
		Status: "Done",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version)
}

func ConvertToModelRoom(grpcRoom *pb.Room) model.Room {
	return model.Room{
		ID:         int(grpcRoom.Id),
		HotelId:    int(grpcRoom.HotelId),
		RoomNumber: string(grpcRoom.RoomNumber),
		Type:       string(grpcRoom.Type),
		Price:      float64(grpcRoom.Price),
		Status:     string(grpcRoom.Status),
	}
}

func CreateTextMessageForBookingEvent(booking *model.Booking, room *model.Room, hotel_name string, user *model.User) string {
	return string("Dear, " + user.Name + " we notify you, that you have booked a room " + room.Type +
	" with number " + string(room.RoomNumber) + " in hotel " + hotel_name + " for the dates from " +
	booking.StartDate + " to " + booking.EndDate + ".\n" + "It's cost is " + string(booking.Price))
}

func (handler *Handler) SetBooking(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}

	var booking model.UpdateBooking
	err_json := json.Unmarshal(body, &booking)
	if err_json != nil {
		http.Error(w, "Ошибка при парсинге JSON", http.StatusBadRequest)
		return
	}

	var room model.Room

	err = repository.UpdateBooking(booking, room, handler.DB)
	if err != nil {
		http.Error(w, "Ошибка при обновлении бронирования", http.StatusBadRequest)
		return
	}

	version := QueryStatus{
		Status: "Done",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version)
}

func (handler *Handler) GetBookings(w http.ResponseWriter, r *http.Request) {
	bookings, err := repository.GetAllBookings(handler.DB)
	if err != nil {
		http.Error(w, "Ошибка при получении данных о бронированиях", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(bookings)
	if err != nil {
		http.Error(w, "Ошибка при сериализации данных о бронированиях", http.StatusInternalServerError)
		return
	}
}

func (handler *Handler) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}

	var booking model.DeleteBooking
	err_json := json.Unmarshal(body, &booking)
	if err_json != nil {
		http.Error(w, "Ошибка при парсинге JSON", http.StatusBadRequest)
		return
	}

	err = repository.DeleteBooking(booking, handler.DB)
	if err != nil {
		http.Error(w, "Ошибка при удалении бронирования", http.StatusBadRequest)
		return
	}

	version := QueryStatus{
		Status: "Done",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version)
}

func (handler *Handler) GetBookingByUser(w http.ResponseWriter, r *http.Request) {
	userIdParam := r.URL.Query().Get("userId")
	if userIdParam == "" {
		http.Error(w, "Не указан userId", http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		http.Error(w, "Некорректный userId", http.StatusBadRequest)
		return
	}

	bookings, err := repository.GetBookingByUser(handler.DB, userId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения бронирований: %v", err), http.StatusInternalServerError)
		return
	}

	if len(bookings) == 0 {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(map[string]string{"message": "Бронирования не найдены"})
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bookings); err != nil {
		http.Error(w, "Ошибка при отправке данных", http.StatusInternalServerError)
		return
	}
}
