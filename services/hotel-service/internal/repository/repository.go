package repository

import (
	"HotelService/internal/model"
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
        CREATE TABLE IF NOT EXISTS hotel (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            price SERIAL NOT NULL,
  			hotelier_id INTEGER NOT NULL,
    		FOREIGN KEY (hotelier_id) REFERENCES hoteliers(id)
        );
        CREATE TABLE IF NOT EXISTS hoteliers (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    name TEXT NOT NULL
		);
    `)
	return err
}

func AddHotel(hotelier model.Hotelier, hotel model.HotelData, db *sql.DB) error {
	_, err := db.Exec(`INSERT INTO hotel (name, price, hotelier_id) VALUES (?, ?, ?)`, hotel.Name, hotel.Price, hotelier.HotelierId)
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

func GetHotelsByHotelier(hotelierID int, db *sql.DB) ([]model.Hotel, error) {
	rows, err := db.Query(`
        SELECT id, name, price, hotelier_id
        FROM hotel
        WHERE hotelier_id = ?
    `, hotelierID)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к базе данных: %w", err)
	}
	defer rows.Close()

	var hotels []model.Hotel
	for rows.Next() {
		var hotel model.Hotel
		err := rows.Scan(&hotel.ID, &hotel.Name, &hotel.Price, &hotel.HotelierId)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования результатов запроса: %w", err)
		}
		hotels = append(hotels, hotel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении результатов запроса: %w", err)
	}

	return hotels, nil
}
