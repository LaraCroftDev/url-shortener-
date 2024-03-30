package main

import (
	"database/sql"
	"errors"
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
		fmt.Println("Error connecting to DB: ", err)
		return nil
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Database connection interruption: ", err)
		return nil
	}
	fmt.Println("Successfully connected!")
	return db
}

func insertUrl(origin, short string) error {
	db := dbConnect()
	defer db.Close()

	if _, err := db.Exec("INSERT INTO " + tableName + " (origin, short, expiry) VALUES ('" + origin + "', '" + short + "', CURRENT_TIMESTAMP);"); err != nil {
		return err
	}
	return nil
}

func getUrl(origin string) (URL, error) {
	db := dbConnect()
	defer db.Close()

	var row URL
	err := db.QueryRow("SELECT origin, short, expiry FROM "+tableName+" WHERE origin = '"+origin+"';").Scan(&row.origin, &row.short, &row.expiry)
	return row, err
}

func updateUrl(url URL) error {
	db := dbConnect()
	defer db.Close()

	res, err := db.Exec("UPDATE " + tableName + " SET short = '" + url.short + "', " + " expiry = CURRENT_TIMESTAMP WHERE origin = '" + url.origin + "';")
	if err != nil {
		return err
	}
	effect, _ := res.RowsAffected()
	if effect != int64(1) {
		return errors.New("Error: recent updated affected more than one row")
	}
	return nil
}

func deleteUrl(origin string) error {
	db := dbConnect()
	defer db.Close()

	res, err := db.Exec("DELETE FROM " + tableName + " WHERE origin = '" + origin + "';")
	if err != nil {
		return err
	}
	effect, _ := res.RowsAffected()
	if effect != int64(1) {
		return errors.New("Error: recent delete affected more than one row")
	}
	return nil
}

func dbInfo() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, dbUser, password, dbName)
}
