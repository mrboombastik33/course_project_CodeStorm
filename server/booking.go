package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type Booking struct {
	GroupName   string `json:"group"`
	StudentName string `json:"name"`
	Audience    string `json:"subject"`
	BookingTime string `json:"time"`
	DayOfWeek   string `json:"day"`
}

// CORRECTED FUNCTION
func createBooking(db *sql.DB, booking Booking) error {

	query := `
		INSERT INTO bookings (group_name, user_name, audience, booking_time, day_of_week)
		VALUES (?, ?, ?, STR_TO_DATE(?, '%H:%i'), ?)
	`

	// The strconv.Atoi conversion is no longer needed.
	// We pass booking.Audience (the original string) as the argument.
	_, err := db.Exec(query,
		booking.GroupName,
		booking.StudentName,
		booking.Audience,
		booking.BookingTime,
		booking.DayOfWeek,
	)

	if err != nil {
		// This will now provide a more accurate error if something goes wrong
		return fmt.Errorf("failed to insert booking: %v", err)
	}

	return nil
}

func bookingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received request: %s %s\n", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Log request headers
	fmt.Printf("Request headers: %v\n", r.Header)

	var booking Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		fmt.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received booking data: %+v\n", booking)

	if booking.GroupName == "" || booking.StudentName == "" ||
		booking.Audience == "" || booking.BookingTime == "" ||
		booking.DayOfWeek == "" {
		fmt.Printf("Missing required fields in booking: %+v\n", booking)
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// З'єднання з базою даних
	db, err := startMySQL(credentialsSQL.user, credentialsSQL.password, credentialsSQL.DBName)
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Створити бронювання
	if err := createBooking(db, booking); err != nil {
		fmt.Printf("Error creating booking: %v\n", err)
		http.Error(w, "Failed to create booking", http.StatusInternalServerError)
		return
	}

	// Повернути статус
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status":  "success",
		"message": "Booking created successfully",
	}
	fmt.Printf("Sending response: %+v\n", response)
	json.NewEncoder(w).Encode(response)
}
