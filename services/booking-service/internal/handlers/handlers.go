package handlers

import (
	"BookingService/internal/kafka_producer"
	"BookingService/internal/model"
	"BookingService/internal/repository"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
)

type Handler struct {
	DB       *sql.DB
	Producer *kafka_producer.KafkaProducer
}

type QueryStatus struct {
	Status string `json:"status"`
}

func NewHandler(db *sql.DB, producer *kafka_producer.KafkaProducer) *Handler {
	return &Handler{DB: db, Producer: producer}
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

	//hotel, err := handler.getHotelByGRPC(booking.HotelId)
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("Failed to retrieve hotel: %v", err), http.StatusInternalServerError)
	//	return
	//}

	err = repository.AddBooking(booking, handler.DB)
	if err != nil {
		http.Error(w, "Ошибка при добавление бронирования", http.StatusBadRequest)
		return
	}

	event := map[string]interface{}{
		"hotel_id":       booking.HotelId,
		"user_id":        booking.UserId,
		"start_date":     booking.StartDate,
		"end_date":       booking.EndDate,
		"price":          booking.Price,
		"status":         booking.Status,
		"payment_status": booking.PaymentStatus,
	}

	message, err := json.Marshal(event)
	if err != nil {
		http.Error(w, "Ошибка при подготовке события Kafka", http.StatusInternalServerError)
		return
	}

	err = handler.Producer.Publish(r.Context(), "booking_event", message)
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

	err = repository.UpdateBooking(booking, handler.DB)
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
