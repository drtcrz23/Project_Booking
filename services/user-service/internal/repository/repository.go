package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/drtcrz23/Project_Booking/services/user-service/internal/model"
	_ "github.com/lib/pq"
)

func CReateDBConnection(dbName string) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=localhost port=5435 user=postgres password=123 dbname=user-db sslmode=disable")
	db, err := sql.Open("postgres", psqlInfo)
	// connStr := dbName
	// db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CReateTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		username VARCHAR(50) NOT NULL,
		surname VARCHAR(100) NOT NULL,
		phone VARCHAR(20) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		balance NUMERIC(10, 2) NOT NULL DEFAULT 0
	);

    `)
	return err
}

func AddUser(user model.User, db *sql.DB) (model.User, error) {
	query := `INSERT INTO users (user_name, surname, phone, email, balance) VALUES ($1, $2, $3, $4, $5)`
	result, err := db.Exec(query, user.Name, user.Surname, user.Phone, user.Email, user.Balance)
	if err != nil {
		return model.User{}, fmt.Errorf("ошибка при добавлении пользователя: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return model.User{}, fmt.Errorf("не удалось получить ID пользователя: %w", err)
	}

	user.Id = int(userID)
	return user, nil
}

func SetUser(user model.User, db *sql.DB) error {
	query := `UPDATE users
			  SET user_name = $1, surname = $2, phone = $3, email = $4, balance = $5
			  WHERE user_id = $6`
	result, err := db.Exec(query, user.Name, user.Surname, user.Phone, user.Email, user.Balance, user.Id)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении пользователя: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при проверке затронутых строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("пользователь с ID %d не найден или данные не изменены", user.Id)
	}

	return nil
}

func GetUserById(id int, db *sql.DB) (*model.User, error) {
	var user model.User
	query := `SELECT user_id, user_name, surname, phone, email, balance FROM users WHERE user_id = $1`
	row := db.QueryRow(query, id)

	err := row.Scan(&user.Id, &user.Name, &user.Surname, &user.Phone, &user.Email, &user.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("пользователь с таким ID не найден")
		}
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}

	return &user, nil
}

func GetAllUsers(db *sql.DB) ([]model.User, error) {
	query := `SELECT user_id, name, surname, phone, email, balance FROM users`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе всех пользователей: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Name, &user.Surname, &user.Phone, &user.Email, &user.Balance); err != nil {
			return nil, fmt.Errorf("ошибка при чтении данных о пользователе: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении результатов запроса: %w", err)
	}

	return users, nil
}

func DeleteUser(id int, db *sql.DB) error {
	query := `DELETE FROM users WHERE user_id = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении пользователя: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при проверке затронутых строк: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("пользователь с ID %d не найден", id)
	}

	return nil
}
