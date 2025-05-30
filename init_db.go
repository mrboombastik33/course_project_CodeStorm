package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func initDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		return nil, err
	}

	// Create rooms table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS rooms (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			esp_id INTEGER UNIQUE NOT NULL,
			status TEXT DEFAULT 'closed'
		)
	`)
	if err != nil {
		return nil, err
	}

	// Create users table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			role TEXT NOT NULL,
			card_uid TEXT UNIQUE
		)
	`)
	if err != nil {
		return nil, err
	}

	// Create bookings table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS bookings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			room_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			status TEXT DEFAULT 'active',
			FOREIGN KEY (room_id) REFERENCES rooms(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		return nil, err
	}

	// Create access_logs table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS access_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			room_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			access_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			action TEXT NOT NULL,
			FOREIGN KEY (room_id) REFERENCES rooms(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}