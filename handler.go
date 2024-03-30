package main

import (
	"net/http"
	"time"
)

const (
	charset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLength = 6
)

type URL struct {
	origin string
	short  string
	expiry time.Time
}

func (u *URL) handler(w http.ResponseWriter, r *http.Request) {
	// 	if u.origin == "" {
	// 		http.Error(w, "Empty input URL", http.StatusBadRequest)
	// 		return
	// 	}

	// 	res, err := validateUrl(u.origin)
	// 	if err != nil {
	// 		http.Error(w, "Invalid URL", http.StatusBadRequest)
	// 		return
	// 	}

	// 	// db, err := sql.Open(dbdriver, dbInfo())
	// 	// if err != nil {
	// 	// 	http.Error(w, "Database connection failure: "+err.Error(), http.StatusInternalServerError)
	// 	// 	return
	// 	// }

	// 	// defer db.Close()

	// 	// if err := db.Ping(); err != nil {
	// 	// 	http.Error(w, "Database connection interruption: "+err.Error(), http.StatusInternalServerError)
	// 	// 	return
	// 	// }

	// 	// var row URL
	// 	// err = db.QueryRow("SELECT origin, short, expiry FROM "+tablename+" WHERE origin = '"+u.origin+"';").Scan(&row.origin, &row.short, &row.expiry)
	// 	if err != nil {
	// 		if !errors.Is(err, sql.ErrNoRows) {
	// 			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
	// 			return
	// 		}

	// 		u.short = generateShortUrl(res)
	// 		u.expiry = time.Now()
	// 		k, err := db.Exec("INSERT INTO " + tableName + " (origin, short, expiry) VALUES ('" + u.origin + "', '" + u.short + "', CURRENT_TIMESTAMP);")
	// 		if err != nil {
	// 			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
	// 			return
	// 		}
	// 		m, err := k.RowsAffected()
	// 		if err != nil || m != int64(1) {
	// 			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
	// 			return
	// 		}
}

// 	if timeExpired(row.expiry) {
// 		u.short = generateShortUrl(res)
// 		u.expiry = time.Now()
// 		_, err := db.Exec("UPDATE " + tableName + " SET short = '" + u.short + "', " + " expiry = CURRENT_TIMESTAMP;")
// 		if err != nil {
// 			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	} else {
// 		u.expiry = row.expiry
// 		u.short = row.short
// 	}
// 	http.Redirect(w, r, u.origin, http.StatusSeeOther)
// }

// func generateShortUrl(u *url.URL) string {
// 	rand.New(rand.NewSource(time.Now().UnixNano()))
// 	shortKey := make([]byte, keyLength)
// 	for i := range shortKey {
// 		shortKey[i] = charset[rand.Intn(len(charset))]
// 	}
// 	return u.Scheme + "://" + u.Host + "/" + string(shortKey)
// }

// func validateUrl(inputUrl string) (*url.URL, error) {
// 	return url.ParseRequestURI(inputUrl)
// }

// func timeExpired(expiry time.Time) bool {
// 	if (time.Now()).Sub(expiry) < (time.Hour * 5) {
// 		return false
// 	}
// 	return true
// }
