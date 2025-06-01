package main

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

func startMySQL(user, password, DBname string) (DBinst *sql.DB, err error) {
	cfg := mysql.NewConfig()
	cfg.User = user
	cfg.Passwd = password
	cfg.DBName = DBname

	DBinst, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		return nil, fmt.Errorf("unable to open connection: %v", err)
	}
	err = DBinst.Ping()
	if err != nil {
		return nil, fmt.Errorf("unable to ping: %v", err)
	}

	return
}

func getKeyUID(db *sql.DB, ESP_ID int) (keyUID string, err error) {
	row := db.QueryRow("SELECT uid FROM key_inf WHERE esp_id = ?", ESP_ID)
	var keyFromDB string
	err = row.Scan(&keyFromDB)
	if err != nil {
		return "", err
	}
	return keyFromDB, nil
}