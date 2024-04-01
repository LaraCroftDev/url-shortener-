package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	dbDriver  = "postgres"
	host      = "localhost"
	port      = 5432
	dbUser    = "postgres"
	password  = "irfan"
	dbName    = "urlshortener"
	tableName = "urls"
)

func dbConnect() *sql.DB {
	db, err := sql.Open(dbDriver, dbInfo())
	if err != nil {
		fmt.Println("Error DB connection: ", err)
		return nil
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Error DB connection interruption: ", err)
		return nil
	}
	fmt.Println("Successful DB connection")
	return db
}

func insertUrl(origin, short string) error {
	db := dbConnect()
	defer db.Close()
	if _, err := db.Exec("INSERT INTO " + tableName + " (origin, short, expiry) VALUES ('" + origin + "', '" + short + "', CURRENT_TIMESTAMP + INTERVAL '24 HOUR');"); err != nil {
		return err
	}
	return nil
}

func getUrl(short string) (URL, error) {
	var row URL
	db := dbConnect()
	err := db.QueryRow("SELECT origin, short, expiry FROM "+tableName+" WHERE short = '"+short+"';").Scan(&row.origin, &row.short, &row.expiry)
	db.Close()
	return row, err
}

func deleteExpiredUrls() error {
	db := dbConnect()
	defer db.Close()
	result, err := db.Exec("DELETE FROM " + tableName + " WHERE expiry < NOW();")
	if err != nil {
		return err
	}
	count, _ := result.RowsAffected()
	fmt.Println("Number of expired URLs deleted: ", count)
	return nil
}

func dbInfo() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, dbUser, password, dbName)
}
