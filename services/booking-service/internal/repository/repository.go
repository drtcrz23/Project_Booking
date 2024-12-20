package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/drtcrz23/Project_Booking/services/booking-service/internal/model"
	"github.com/drtcrz23/Project_Booking/services/booking-service/internal/parser_data"

	_ "github.com/lib/pq"
)

func CreateDBConnection(dbName string) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=localhost port=5434 user=postgres password=123 dbname=booking-db sslmode=disable")
	db, err := sql.Open("postgres", psqlInfo)
	// connStr := dbName
	// db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS booking (
		id SERIAL PRIMARY KEY,
		hotel_id BIGINT NOT NULL,
		room_id BIGINT NOT NULL,
		user_id BIGINT NOT NULL,
		start_date TEXT NOT NULL,
		end_date TEXT NOT NULL,
		price BIGINT NOT NULL,
		status TEXT NOT NULL,
		payment_status TEXT NOT NULL
	);

	`)
	return err
}

func AddBooking(booking *model.Booking, room model.Room, db *sql.DB) (int, error) {
	if room.Status != "available" {
		return -1, errors.New("room is not available")
	}

	startDate := booking.StartDate
	endDate := booking.EndDate

	days, err_data := parser_data.ParseAndCalculateDays(startDate, endDate)
	if err_data != nil {
		return -1, err_data
	}

	price := room.Price * days

	var bookingId int
	err := db.QueryRow(`INSERT INTO booking (hotel_id, room_id, user_id, start_date, end_date, price, status, payment_status)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		booking.HotelId, booking.RoomId,
		booking.UserId, booking.StartDate,
		booking.EndDate, price,
		booking.Status, booking.PaymentStatus).Scan(&bookingId)
	if err != nil {
		return -1, fmt.Errorf("ошибка добавления бронирования: %w", err)
	}
	booking.ID = bookingId
	return bookingId, err
}

func DeleteBooking(booking model.DeleteBooking, db *sql.DB) error {
	_, err := db.Exec(`DELETE FROM booking WHERE id = $1`, booking.ID)
	if err != nil {
		return fmt.Errorf("ошибка удаления бронирования: %w", err)
	}
	return nil
}

func UpdateBooking(booking model.UpdateBooking, room model.Room, db *sql.DB) error {
	if room.Status != "available" {
		return errors.New("room is not available")
	}

	startDate := booking.StartDate
	endDate := booking.EndDate

	days, err_data := parser_data.ParseAndCalculateDays(startDate, endDate)
	if err_data != nil {
		return err_data
	}

	price := room.Price * days

	_, err := db.Exec(`
		UPDATE booking
		SET hotel_id = $1, room_id = $2, user_id = $3, start_date = $4, end_date = $5, price = $6, status = $7, payment_status = $8
		WHERE id = $9
	`,
		booking.HotelId, booking.RoomId,
		booking.UserId, booking.StartDate,
		booking.EndDate, price,
		booking.Status, booking.PaymentStatus,
		booking.ID,
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления бронирования: %w", err)
	}
	return nil
}

func GetAllBookings(db *sql.DB) ([]model.Booking, error) {
	rows, err := db.Query("SELECT hotel_id, room_id, user_id, start_date, end_date, price, status, payment_status FROM booking")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []model.Booking

	for rows.Next() {
		var booking model.Booking

		err := rows.Scan(
			&booking.HotelId,
			&booking.RoomId,
			&booking.UserId,
			&booking.StartDate,
			&booking.EndDate,
			&booking.Price,
			&booking.Status,
			&booking.PaymentStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки GetAllBookings: %w", err)
		}
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка обработки строк GetAllBookings: %w", err)
	}

	return bookings, nil
}

func GetBookingByUser(db *sql.DB, userId int) ([]model.Booking, error) {
	rows, err := db.Query("SELECT id, hotel_id, room_id, user_id, start_date, end_date, price, status, payment_status FROM booking WHERE user_id = $1", userId)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса к базе данных: %w", err)
	}
	defer rows.Close()

	var bookings []model.Booking

	for rows.Next() {
		var booking model.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.HotelId,
			&booking.RoomId,
			&booking.UserId,
			&booking.StartDate,
			&booking.EndDate,
			&booking.Price,
			&booking.Status,
			&booking.PaymentStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении данных из результата запроса: %w", err)
		}
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при завершении чтения строк: %w", err)
	}

	return bookings, nil
}
func GetBookingById(db *sql.DB, id int) (model.Booking, error) {
	// SQL-запрос для получения данных о бронировании
	query := `
	  SELECT id, hotel_id, room_id, user_id, start_date, end_date, price, status, payment_status 
	  FROM booking 
	  WHERE id = $1`

	// Выполняем запрос
	row := db.QueryRow(query, id)

	// Создаем переменную для хранения результата
	var booking model.Booking

	// Сканируем результат запроса в структуру
	err := row.Scan(
		&booking.ID,
		&booking.HotelId,
		&booking.RoomId,
		&booking.UserId,
		&booking.StartDate,
		&booking.EndDate,
		&booking.Price,
		&booking.Status,
		&booking.PaymentStatus,
	)
	if err != nil {
		return booking, err
	}

	return booking, nil
}
