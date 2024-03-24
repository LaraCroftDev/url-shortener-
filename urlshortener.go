package main

import (
	"fmt"
	"net/http"
)

func redirectTo(NewUrl string) {
	http.Handle("/", http.RedirectHandler(NewUrl, http.StatusMovedPermanently))
	err := http.ListenAndServe("0.0.0.0:9000", nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
