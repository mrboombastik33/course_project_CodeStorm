package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
)

var credentialsSQL = struct {
	user     string
	password string
	DBName   string
}{
	user:     "ESP-usr",
	password: "passwd123",
	DBName:   "rfid_project",
}

// http://192.168.1.104:3333/table?key=

func getKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	hasESP_ID := r.URL.Query().Has("ESP_ID")
	ESP_ID, err := strconv.Atoi(r.URL.Query().Get("ESP_ID"))
	if err != nil {
		io.WriteString(w, "ERR")
		log.Printf("%s: got BAD request: ESP(%t)=%d\n", ctx.Value("ServerAddr"), hasESP_ID, ESP_ID)
		return
	}
	fmt.Printf("%s: got request: ESP_ID(%t)=%d\n", ctx.Value("ServerAddr"), hasESP_ID, ESP_ID)

	db, SQLerr := startMySQL(credentialsSQL.user, credentialsSQL.password, credentialsSQL.DBName)
	if SQLerr != nil {
		log.Panic("MySQL error: ", SQLerr)
	}

	defer db.Close()

	keyUID, err := getKeyUID(db, ESP_ID)
	if err != nil {
		log.Print("Extracting UID from db error: ", err)
		io.WriteString(w, "ERR")
		return
	}
	io.WriteString(w, keyUID)
	fmt.Printf("Query OK, UID: %s\n", keyUID)
}

func main() {
	mux := http.NewServeMux()

	// Add CORS middleware
	mux.HandleFunc("/api/book", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		bookingHandler(w, r)
	})

	mux.HandleFunc("/key", getKey)
	// 🔽 Додаємо хостинг фронтенду
	fs := http.FileServer(http.Dir("../frontend"))
	mux.Handle("/", fs)
	ctx := context.Background()
	server := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, "ServerAddr", l.Addr().String())
			return ctx
		},
	}

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
}
