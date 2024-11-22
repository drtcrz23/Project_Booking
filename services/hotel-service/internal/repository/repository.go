package repository

import (
	"HotelService/internal/model"
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
        CREATE TABLE IF NOT EXISTS hotel (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            price SERIAL NOT NULL
        )
    `)
	return err
}

func AddHotel(hotel model.HotelData, db *sql.DB) error {
	_, err := db.Exec(`INSERT INTO hotel (name, price) VALUES (?, ?)`, hotel.Name, hotel.Price)
	if err != nil {
		return err
	}
	return err
}

func SetHotel(hotel model.UpdateData, db *sql.DB) error {
	_, err := db.Exec(`
    UPDATE hotel 
    SET name = ?, price = ?
    WHERE id = ?
`, hotel.Name, hotel.Price, hotel.ID)
	if err != nil {
		return err
	}
	return err
}

func GetHotelById(id int, db *sql.DB) (model.Hotel, error) {
	var hotel model.Hotel
	query := `SELECT id, name, price FROM hotel WHERE id = ?`
	row := db.QueryRow(query, id)
	err := row.Scan(&hotel.ID, &hotel.Name, &hotel.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return hotel, fmt.Errorf("hotel with ID %d not found", id)
		}
		return hotel, err
	}
	return hotel, nil
}
