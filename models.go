package main

import "time"

type URL struct {
	Origin string
	Short  string
	Expiry time.Time
}

type response struct {
	ShortenedURL string `json:"short,omitempty"`
}
