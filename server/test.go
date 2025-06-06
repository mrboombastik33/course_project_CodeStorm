package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// func formHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	group := r.FormValue("group")
// 	name := r.FormValue("name")
// 	subject := r.FormValue("subject")
// 	day := r.FormValue("day")

// 	fmt.Println("🔹 Отримано бронювання:")
// 	fmt.Println("Група:", group)
// 	fmt.Println("Ім'я:", name)
// 	fmt.Println("Предмет:", subject)
// 	fmt.Println("День:", day)

// 	// w.Write([]byte("Бронювання прийнято!"))
// }

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

// func main() {
// 	// Обробка POST-запиту з форми
// 	http.HandleFunc("/book", formHandler)

// 	// Статичні файли: html, css, js
// 	fs := http.FileServer(http.Dir("../frontend"))
// 	http.Handle("/", fs)

// 	fmt.Println("Сервер працює на http://localhost:8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
//cd server && go run test.go
//Захоститься сайт і в консольці виведе дані з форми
// Перевірка: http://localhost:8080
