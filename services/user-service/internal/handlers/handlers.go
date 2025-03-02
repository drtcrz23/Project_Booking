package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/drtcrz23/Project_Booking/services/user-service/internal/model"
	"github.com/drtcrz23/Project_Booking/services/user-service/internal/repository"
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

func (h *Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}

	var user model.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Ошибка при парсинге JSON", http.StatusBadRequest)
		return
	}

	user, err = repository.AddUser(user, h.DB)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при добавлении пользователя: %v", err), http.StatusInternalServerError)
		return
	}

	version := QueryStatus{
		Status: "Done",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version)
}

func (h *Handler) SetUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}

	var user model.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Ошибка при парсинге JSON", http.StatusBadRequest)
		return
	}

	err = repository.SetUser(user, h.DB)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при обновлении пользователя: %v", err), http.StatusInternalServerError)
		return
	}

	version := QueryStatus{
		Status: "Done",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(version)
}

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}
	fmt.Println(id)
	user, err := repository.GetUserById(id, h.DB)
	if err != nil {
		if strings.Contains(err.Error(), "не найден") {
			http.Error(w, fmt.Sprintf("User with ID %d not found", id), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) ([]model.User, error) {
	var users []model.User
	rows, err := h.DB.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Surname, &user.Phone, &user.Email, &user.Balance); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	err = repository.DeleteUser(id, h.DB)
	if err != nil {
		if strings.Contains(err.Error(), "не найден") {
			http.Error(w, fmt.Sprintf("User with ID %d not found", id), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to delete user: %v", err), http.StatusInternalServerError)
		}
		return
	}

	version := QueryStatus{
		Status: "User deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(version); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
