package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/drtcrz23/Project_Booking/services/hotel-service/internal/model"

	_ "github.com/lib/pq"
)

func CreateDBConnection(dbName string) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=localhost port=5433 user=postgres password=123 dbname=hotel-db sslmode=disable")
	db, err := sql.Open("postgres", psqlInfo)
	//connStr := dbName
	// db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS hoteliers (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS hotel (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		price INTEGER NOT NULL,
		hotelier_id INTEGER NOT NULL,
		CONSTRAINT fk_hotelier 
			FOREIGN KEY (hotelier_id) 
			REFERENCES hoteliers(id)
	);

	CREATE TABLE IF NOT EXISTS rooms (
		id SERIAL PRIMARY KEY,
		hotel_id INTEGER NOT NULL,
		room_number TEXT NOT NULL,
		type_room TEXT NOT NULL,
		price FLOAT NOT NULL,
		status TEXT NOT NULL,
		CONSTRAINT fk_hotel 
			FOREIGN KEY (hotel_id) 
			REFERENCES hotel(id)
	);
    `)
	return err
}

func InsertInTable(db *sql.DB) error {
	_, err := db.Exec(`
        INSERT INTO hoteliers
		(name)
		VALUES('test_name');
    `)
	return err
}

func AddHotel(hotelier model.Hotelier, hotel model.HotelData, db *sql.DB) error {
	_, err := db.Exec(`INSERT INTO hotel (name, price, hotelier_id) VALUES ($1, $2, $3)`, hotel.Name, hotel.Price, hotelier.HotelierId)
	if err != nil {
		return err
	}
	return err
}

func SetHotel(hotel model.UpdateData, db *sql.DB) error {
	_, err := db.Exec(`
    UPDATE hotel 
    SET name = $1, price = $2
    WHERE id = $3
`, hotel.Name, hotel.Price, hotel.ID)
	if err != nil {
		return err
	}
	return err
}

func GetHotelsByHotelier(hotelierID int, db *sql.DB) ([]model.Hotel, error) {
	rows, err := db.Query(`
        SELECT id, name, price, hotelier_id FROM hotel WHERE hotelier_id = $1
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
		rooms, err := GetRoomByHotel(hotel.ID, db)
		if err != nil {
			return nil, fmt.Errorf("ошибка при получении комнат для отеля с ID %d: %w", hotel.ID, err)
		}

		hotel.Rooms = rooms
		hotels = append(hotels, hotel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении результатов запроса: %w", err)
	}

	return hotels, nil
}

func GetHotelById(id int, db *sql.DB) (*model.Hotel, error) {
	var hotel model.Hotel
	query := `SELECT id, name, price, hotelier_id FROM hotel WHERE id = $1`
	row := db.QueryRow(query, id)
	err := row.Scan(&hotel.ID, &hotel.Name, &hotel.Price, &hotel.HotelierId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("отель с таким ID не найден")
		}
		return nil, fmt.Errorf("ошибка при получении отеля: %w", err)
	}

	rooms, err := GetRoomByHotel(id, db)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении комнат для отеля с ID %d: %w", id, err)
	}

	hotel.Rooms = rooms

	return &hotel, nil
}

func GetRoomByHotel(hotelID int, db *sql.DB) ([]model.Room, error) {
	roomQuery := `
		SELECT id, hotel_id, room_number, type_room, price, status
		FROM rooms
		WHERE hotel_id = $1
	`
	rows, err := db.Query(roomQuery, hotelID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе списка комнат: %w", err)
	}
	defer rows.Close()

	var rooms []model.Room
	for rows.Next() {
		var room model.Room
		if err := rows.Scan(&room.ID, &room.HotelId, &room.RoomNumber, &room.Type, &room.Price, &room.Status); err != nil {
			return nil, fmt.Errorf("ошибка при чтении данных о комнате: %w", err)
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении результатов запроса: %w", err)
	}

	return rooms, nil
}

func AddRoom(room model.Room, db *sql.DB) error {
	query := `INSERT INTO rooms (hotel_id, room_number, type_room, price, status) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, room.HotelId, room.RoomNumber, room.Type, room.Price, room.Status)
	if err != nil {
		return err
	}
	return nil
}

func SetRoom(room model.Room, db *sql.DB) error {
	query := `UPDATE rooms
			  SET room_number = $1, type_room = $2, price = $3, status = $4
			  WHERE id = $5 AND hotel_id = $6`
	result, err := db.Exec(query, room.RoomNumber, room.Type, room.Price, room.Status, room.ID, room.HotelId)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении комнаты: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при проверке затронутых строк: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("комната с ID %d не найдена или данные не изменены", room.ID)
	}
	// RowsAffected проверяет, была ли удалена строка. Если 0, значит, комната с указанным ID не найдена.
	return nil
}

func DeleteRoom(roomID int, hotelID int, db *sql.DB) error {
	query := `DELETE FROM rooms
			  WHERE id = $1 AND hotel_id = $2`
	result, err := db.Exec(query, roomID, hotelID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении комнаты: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при проверке затронутых строк: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("комната с ID %d не найдена", roomID)
	}
	// RowsAffected проверяет, была ли удалена строка. Если 0, значит, комната с указанным ID не найдена.
	return nil
}
