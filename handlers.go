package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	charset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLength = 6
)

var responseHTML = `
	<h2>URL Shortener</h2>
    <form method="post" action="/api">
        <input type="text" " name="url" placeholder="Enter a URL">
        <input type="submit" value="Shorten">
    </form>
`

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(responseHTML)
	fmt.Fprintf(w, responseHTML)
}

func handleShortener(w http.ResponseWriter, r *http.Request) {
	var url URL
	url.Origin = r.PostFormValue("url")
	if url.Origin == "" {
		http.Error(w, "Empty input URL", http.StatusBadRequest)
		return
	}

	if _, err := validateURL(url.Origin); err != nil {
		http.Error(w, "Invalid input URL", http.StatusBadRequest)
		return
	}

	short := generateShortURL()
	if err := insertUrl(url.Origin, short); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fullPath := fmt.Sprintf("%v://%v%v", r.URL.Scheme, r.Host, short)
	if strings.Contains(r.Host, "localhost") {
		fullPath = "http" + fullPath
	}
	res := response{
		ShortenedURL: fullPath,
	}
	json.NewEncoder(w).Encode(res)
}

func handlerRedirect(w http.ResponseWriter, r *http.Request) {
	// Purge db from expired shortened URLs before to begin
	if err := deleteExpiredUrls(); err != nil {
		fmt.Println("Error deleting expired URLs: ", err)
	}

	read, err := getUrl(r.URL.String())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// If the shortened URL is expired instead of sending back an error, user'll be redirected to homepage
	if errors.Is(err, sql.ErrNoRows) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, read.Origin, http.StatusSeeOther)
}

func generateShortURL() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return "/api/short/" + string(shortKey)
}

// validateURL makes sure the input url is a valid url
func validateURL(inputURL string) (*url.URL, error) {
	return url.ParseRequestURI(inputURL)
}
