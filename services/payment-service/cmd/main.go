package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type PaymentRequest struct {
	BookingID int     `json:"booking_id"`
	Price     float64 `json:"price"`
}

type PaymentResponse struct {
	BookingId int     `json:"booking_id"`
	Price     float64 `json:"price"`
	Status    string  `json:"status"`
	Message   string  `json:"message"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/pay", handlePayment)

	port := "8083"
	fmt.Printf("Payment Service running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handlePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var paymentReq PaymentRequest
	err := json.NewDecoder(r.Body).Decode(&paymentReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	status, message := processPayment(paymentReq)

	response := PaymentResponse{
		BookingId: paymentReq.BookingID,
		Price:     paymentReq.Price,
		Status:    status,
		Message:   message,
	}

	err = forwardPaymentResponse(response)
	if err != nil {
		log.Printf("Failed to forward payment response: %v", err)
		http.Error(w, "Failed to forward payment response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func processPayment(req PaymentRequest) (string, string) {
	return "ok", fmt.Sprintf("Payment of %.2f for user %d successful", req.Price, req.BookingID)
}

func forwardPaymentResponse(response PaymentResponse) error {
	url := "http://localhost:8081/booking/call"
	responseData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(responseData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("received non-OK response: %s, body: %s", resp.Status, string(body))
	}

	return nil
}
