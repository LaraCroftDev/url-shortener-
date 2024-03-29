package main

import (
	"fmt"
	"net/http"
)

const ( // move these to test later
	localhost = "http://localhost:8181/test"
)

func main() {
	l := &Url{
		OrigUrl: localhost,
	}
	http.HandleFunc("/test", l.handler)
	fmt.Println("URL Shortener is running on :8181")
	http.ListenAndServe(":8181", nil)

}
