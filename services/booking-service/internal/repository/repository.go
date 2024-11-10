package repository

import (
	"BookingService/internal/model"
	"database/sql"
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
		booking.Price, // нужно прайс нормально подгружать в handlers по grpc
		booking.Status, booking.PaymentStatus)
	if err != nil {
		return err
	}
	return err
}
