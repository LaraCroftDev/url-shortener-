package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	_ "github.com/lib/pq"
)

const (
	dbdriver  = "postgres"
	host      = "localhost"
	port      = 5432
	user      = "postgres"
	password  = "irfan"
	dbname    = "urlshortener"
	tablename = "urls"
	charset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLength = 6
)

type Url struct {
	OrigUrl  string
	ShortUrl string
	Expiry   time.Time
}

func (u *Url) handler(w http.ResponseWriter, r *http.Request) {
	if u.OrigUrl == "" {
		log.Fatalln(errors.New("Empty input url."))
		return
	}

	res, err := validateUrl(u.OrigUrl)
	if err != nil {
		log.Fatalln(err)
		return
	}

	db, err := sql.Open(dbdriver, dbInfo())
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
		return
	}

	row := Url{}
	err = db.QueryRow("SELECT * FROM "+tablename+" WHERE origurl LIKE '"+u.OrigUrl+"';").Scan(&row.OrigUrl, &row.ShortUrl, &row.Expiry)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatalln(err)
			return
		}

		u.ShortUrl = generateShortUrl(res)
		u.Expiry = time.Now()
		k, err := db.Exec("INSERT INTO " + tablename + " (origurl, shorturl, expiry) VALUES ('" + u.OrigUrl + "', '" + u.ShortUrl + "', CURRENT_TIMESTAMP);")
		if err != nil {
			log.Fatalln(err)
			return
		}
		m, err := k.RowsAffected()
		if err != nil || m != int64(1) {
			log.Fatalln(err)
			return
		}
	}

	if timeExpired(row.Expiry) {
		u.ShortUrl = generateShortUrl(res)
		u.Expiry = time.Now()
		_, err := db.Exec("UPDATE " + tablename + " SET shorturl = '" + u.ShortUrl + "', " + " expiry = CURRENT_TIMESTAMP;")
		if err != nil {
			log.Fatalln(err)
			return
		}
	} else {
		u.Expiry = row.Expiry
		u.ShortUrl = row.ShortUrl
	}
	http.Redirect(w, r, u.ShortUrl, http.StatusSeeOther)
}

func generateShortUrl(u *url.URL) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return u.Scheme + "://" + u.Host + "/" + string(shortKey)
}

func validateUrl(inputUrl string) (*url.URL, error) {
	return url.ParseRequestURI(inputUrl)
}

func dbInfo() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}

func timeExpired(expiry time.Time) bool {
	if (time.Now()).Sub(expiry) < (time.Hour * 5) {
		return false
	}
	return true
}
