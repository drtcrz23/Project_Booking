package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/drtcrz23/Project_Booking/services/booking-service/internal/kafka_producer"
	"github.com/drtcrz23/Project_Booking/services/booking-service/internal/model"
	"github.com/drtcrz23/Project_Booking/services/booking-service/internal/parser_data"
	"github.com/drtcrz23/Project_Booking/services/booking-service/internal/repository"
	pb "github.com/drtcrz23/Project_Booking/services/hotel-service/pkg/api"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	DB          *sql.DB
	Producer    *kafka_producer.KafkaProducer
	HotelClient pb.HotelServiceClient
}

type QueryStatus struct {
	Status string `json:"status"`
}

func NewHandler(db *sql.DB, producer *kafka_producer.KafkaProducer, hotelClient pb.HotelServiceClient) *Handler {
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

	fmt.Printf("Retrieved hotel: %v\n", hotel.Name)
	var room model.Room
	for _, cur := range hotel.Rooms {
		if cur.Id == int32(booking.RoomId) {
			room = ConvertToModelRoom(cur)
			break
		}
	}

	startDate := booking.StartDate
	endDate := booking.EndDate

	days, err_data := parser_data.ParseAndCalculateDays(startDate, endDate)
	if err_data != nil {
		return
	}

	price := room.Price * days
	booking.Price = int(price)

	booking_id, err := repository.AddBooking(&booking, room, handler.DB)
	if err != nil {
		errorMessage := fmt.Sprintf("Ошибка при добавление бронирования: %v", err)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	paymentRequest := map[string]interface{}{
		"booking_id": booking_id,
		"price":      booking.Price,
	}

	err = handler.sendPaymentRequest(paymentRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Payment service error: %v", err), http.StatusInternalServerError)
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

func (handler *Handler) sendPaymentRequest(request map[string]interface{}) error {
	url := "http://localhost:8083/pay"

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-OK response: %s", resp.Status)
	}

	log.Println("Payment request sent successfully")
	return nil
}

func (handler *Handler) SendMessageAfterSuccessfullyPay(w http.ResponseWriter, r *http.Request) {
	var paymentResp PaymentResponse
	err := json.NewDecoder(r.Body).Decode(&paymentResp)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if paymentResp.Status != "ok" {
		http.Error(w, paymentResp.Message, http.StatusBadRequest)
	}

	//booking, err := repository.GetBookingById(handler.DB, paymentResp.BookingId)
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("Failed to get booking: %v", err), http.StatusInternalServerError)
	//	return
	//}
}

type PaymentResponse struct {
	BookingId int     `json:"booking_id"`
	Price     float64 `json:"price"`
	Status    string  `json:"status"`
	Message   string  `json:"message"`
}
