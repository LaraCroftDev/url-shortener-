package main

import (
	"fmt"
	"net/http"
)

const ( // move these to test later
	localhost = "http://localhost:8181/hello"
)

func main() {
	l := URL{
		origin: localhost,
	}
	http.HandleFunc("/", l.handleShortener)
	http.HandleFunc("/short/", l.handlerRedirect)

	fmt.Println("URL Shortener is running on :8181")
	http.ListenAndServe(":8181", nil)

}
