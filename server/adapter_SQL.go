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


// func book(db *sql.DB, name string, ESP_ID int) (keyUID string, err error) {
// 	row := db.QueryRow("SELECT uid FROM access_inf WHERE name = ?", name)
// 	var keyFromDB string
// 	err = row.Scan(&keyFromDB)

// 	if err != nil {
// 		return "", err
// 	}

// 	row = db.QueryRow("UPDATE uid FROM key_inf WHERE ESP_ID = ?", ESP_ID)

// 	return keyFromDB, nil
// }


func book(db *sql.DB, name string, ESP_ID int) (keyUID string, err error) {
	// Отримуємо ключ за ім'ям
	row := db.QueryRow("SELECT uid FROM access_inf WHERE name = ?", name)
	var keyFromDB string
	err = row.Scan(&keyFromDB)
	if err != nil {
		return "", err
	}

	// Оновлюємо key_inf, встановлюючи uid = keyFromDB для ESP_ID
	_, err = db.Exec("UPDATE key_inf SET uid = ? WHERE ESP_ID = ?", keyFromDB, ESP_ID)
	if err != nil {
		return "", err
	}

	return keyFromDB, nil
}
