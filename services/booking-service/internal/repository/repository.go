package repository

import (
	"BookingService/internal/model"
	"BookingService/internal/parser_data"
	"database/sql"
	"errors"
	"fmt"
)

func CreateDBConnection(dbName string) (*sql.DB, error) {
	connStr := dbName
	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS booking (
			id BIGSERIAL PRIMARY KEY
			hotel_id BIGINT
		    room_id BIGINT
			user_id BIGINT /*REFERENCES users(id)*/
			start_date TEXT NOT NULL
			end_date TEXT NOT NULL
		    Price BIGINT NOT NULL
			status TEXT NOT NULL
			payment_status TEXT NOT NULL
		)
	`)
	return err
}

func AddBooking(booking *model.Booking, room model.Room, db *sql.DB) error {
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
	booking.Price = int(price)
	_, err := db.Exec(`INSERT INTO booking (hotel_id, roomd_id, user_id, start_date, end_date, price, status, payment_status)
						VALUES (?, ?, ?, ?, ?, ?, ?)`,
		booking.HotelId, booking.RoomId,
		booking.UserId, booking.StartDate,
		booking.EndDate, price,
		booking.Status, booking.PaymentStatus)
	if err != nil {
		return fmt.Errorf("ошибка добавления бронирования: %w", err)
	}
	return err
}

func DeleteBooking(booking model.DeleteBooking, db *sql.DB) error {
	_, err := db.Exec(`DELETE FROM booking WHERE id = ?`, booking.ID)
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
		SET hotel_id = ?, room_id = ?, user_id = ?, start_date = ?, end_date = ?, price = ?, status = ?, payment_status = ?
		WHERE id = ?
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
	rows, err := db.Query("SELECT id, hotel_id, room_id, user_id, start_date, end_date, price, status, payment_status FROM booking WHERE user_id = ?", userId)
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
