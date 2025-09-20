package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)


func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	espIDstr := r.FormValue("ESP_ID")
	ESP_ID, err := strconv.Atoi(espIDstr)
	if err != nil {
		http.Error(w, "Invalid ESP_ID", http.StatusBadRequest)
		return
	}

	fmt.Printf("🔹 Отримано бронювання: name=%s, ESP_ID=%d\n", name, ESP_ID)

	// Підключення до БД
	db, SQLerr := startMySQL(credentialsSQL.user, credentialsSQL.password, credentialsSQL.DBName)
	if SQLerr != nil {
		log.Printf("MySQL error: %v", SQLerr)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Виклик функції book
	keyUID, err := book(db, name, ESP_ID)
	if err != nil {
		log.Printf("Book error: %v", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	// Надсилаємо UID у відповідь
	w.Write([]byte(fmt.Sprintf("UID заброньовано: %s", keyUID)))
}
