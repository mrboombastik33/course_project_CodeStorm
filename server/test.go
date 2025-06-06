package main

import (
	"fmt"
	"net/http"
)

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	group := r.FormValue("group")
	name := r.FormValue("name")
	subject := r.FormValue("subject")
	day := r.FormValue("day")

	fmt.Println("🔹 Отримано бронювання:")
	fmt.Println("Група:", group)
	fmt.Println("Ім'я:", name)
	fmt.Println("Предмет:", subject)
	fmt.Println("День:", day)

	// w.Write([]byte("Бронювання прийнято!"))
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