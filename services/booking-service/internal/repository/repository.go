package repository

import (
	"BookingService/internal/model"
	"database/sql"
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
			hotel_id BIGINT REFERENCES hotel(id)
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

func AddBooking(booking model.Booking, db *sql.DB) error {
	_, err := db.Exec(`INSERT INTO booking (hotel_id, user_id, start_date, end_date, price, status, payment_status)
						VALUES (?, ?, ?, ?, ?, ?, ?)`,
		booking.HotelId, booking.UserId,
		booking.StartDate, booking.EndDate,
		booking.Price,
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

func UpdateBooking(booking model.UpdateBooking, db *sql.DB) error {
	_, err := db.Exec(`
		UPDATE booking
		SET hotel_id = ?, user_id = ?, start_date = ?, end_date = ?, price = ?, status = ?, payment_status = ?
		WHERE id = ?
	`,
		booking.HotelId, booking.UserId,
		booking.StartDate, booking.EndDate,
		booking.Price,
		booking.Status, booking.PaymentStatus,
		booking.ID,
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления бронирования: %w", err)
	}
	return nil
}
