package main

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
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

func (u *URL) handlerRedirect(w http.ResponseWriter, r *http.Request) {
	// Cleanup db before to begin
	err := deleteExpiredUrls()
	if err != nil {
		fmt.Println("Error deleting expired URLs: ", err)
	}
	read, err := getUrl(u.short)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// The shortened URL is expired
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Invalid URL ", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, read.origin, http.StatusSeeOther)
}

func (u *URL) handleShortener(w http.ResponseWriter, r *http.Request) {
	if u.origin == "" {
		http.Error(w, "Empty input URL", http.StatusBadRequest)
		return
	}
	res, err := validateURL(u.origin)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	u.short = generateShortURL(res)
	if err := insertUrl(u.origin, u.short); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func generateShortURL(u *url.URL) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return u.Scheme + "://" + u.Host + "/short/" + string(shortKey)
}

func validateURL(inputURL string) (*url.URL, error) {
	return url.ParseRequestURI(inputURL)
}
