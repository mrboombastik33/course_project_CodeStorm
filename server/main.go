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
	mux.HandleFunc("/key", getKey)

	ctx := context.Background()
	server := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, "ServerAddr", l.Addr().String())
			return ctx
		},
	}
	// err := server.ListenAndServeTLS("/home/pavlo/Стільниця/cert/cert.pem", "/home/pavlo/Стільниця/cert/key.pem")
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
}