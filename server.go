package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Room struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	ESPID  int    `json:"esp_id"`
	Status string `json:"status"`
}

type Booking struct {
	ID        int       `json:"id"`
	RoomID    int       `json:"room_id"`
	UserID    int       `json:"user_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
}

var db *sql.DB

func init() {
	var err error
	db, err = initDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
}

func checkAccess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cardUID := r.URL.Query().Get("card_uid")
	espID, err := strconv.Atoi(r.URL.Query().Get("esp_id"))
	if err != nil {
		io.WriteString(w, "ERR")
		log.Printf("%s: invalid ESP_ID: %v\n", ctx.Value("ServerAddr"), err)
		return
	}

	// Check if user exists and has access to the room
	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE card_uid = ?", cardUID).Scan(&userID)
	if err != nil {
		io.WriteString(w, "DENIED")
		return
	}

	// Check if there's an active booking for this room and user
	var bookingID int
	err = db.QueryRow(`
		SELECT b.id FROM bookings b
		JOIN rooms r ON b.room_id = r.id
		WHERE r.esp_id = ? AND b.user_id = ? AND b.status = 'active'
		AND datetime('now') BETWEEN b.start_time AND b.end_time`,
		espID, userID).Scan(&bookingID)

	if err != nil {
		io.WriteString(w, "DENIED")
		return
	}

	// Log the access
	_, err = db.Exec(`
		INSERT INTO access_logs (room_id, user_id, action)
		SELECT r.id, ?, 'access' FROM rooms r WHERE r.esp_id = ?`,
		userID, espID)
	if err != nil {
		log.Printf("Failed to log access: %v\n", err)
	}

	io.WriteString(w, "GRANTED")
}

func getRooms(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, esp_id, status FROM rooms")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.ID, &room.Name, &room.ESPID, &room.Status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rooms = append(rooms, room)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

func createBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var booking Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if room is available
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM bookings
			WHERE room_id = ? AND status = 'active'
			AND (
				(start_time BETWEEN ? AND ?)
				OR (end_time BETWEEN ? AND ?)
				OR (? BETWEEN start_time AND end_time)
			)
		)`,
		booking.RoomID, booking.StartTime, booking.EndTime,
		booking.StartTime, booking.EndTime, booking.StartTime).Scan(&exists)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "Room is already booked for this time", http.StatusConflict)
		return
	}

	// Create booking
	result, err := db.Exec(`
		INSERT INTO bookings (room_id, user_id, start_time, end_time, status)
		VALUES (?, ?, ?, ?, 'active')`,
		booking.RoomID, booking.UserID, booking.StartTime, booking.EndTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	booking.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/access", checkAccess)
	mux.HandleFunc("/rooms", getRooms)
	mux.HandleFunc("/bookings", createBooking)

	ctx := context.Background()
	server := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, "ServerAddr", l.Addr().String())
			return ctx
		},
	}

	log.Printf("Server starting on port 3333...")
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("server closed\n")
	} else if err != nil {
		log.Printf("error listening for server: %s\n", err)
	}
}