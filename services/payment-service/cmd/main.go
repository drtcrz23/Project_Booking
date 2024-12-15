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
	UserID int     `json:"user_id"`
	Price  float64 `json:"price"`
}

type PaymentResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
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
		Status:  status,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func processPayment(req PaymentRequest) (string, string) {
	success := true //rand.Intn(2) == 0

	if success {
		return "ok", fmt.Sprintf("Payment of %.2f for user %d successful", req.Price, req.UserID)
	}

	return "failed", "Payment failed due to a random error"
}
