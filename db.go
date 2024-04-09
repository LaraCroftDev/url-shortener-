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

// Every shortened url would be expired within 24 hours after it is generated
func insertUrl(origin, short string) error {
	db := dbConnect()
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO " + tableName + " (origin, short, expiry) VALUES ($1, $2, CURRENT_TIMESTAMP + INTERVAL '24 HOUR')")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(origin, short)
	return err
}

func getUrl(short string) (URL, error) {
	db := dbConnect()
	defer db.Close()
	var row URL
	err := db.QueryRow("SELECT origin, short, expiry FROM "+tableName+" WHERE short = '"+short+"';").Scan(&row.Origin, &row.Short, &row.Expiry)
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
