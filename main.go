package main

import (
	"fmt"
	"net/http"
)

func main() {
	r := Router()
	fmt.Println("URL Shortener is running on :8181")
	http.ListenAndServe(":8181", r)
}
