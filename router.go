package main

import "github.com/gorilla/mux"

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", homepageHandler)
	router.HandleFunc("/api", handleShortener).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/short/{short}", handlerRedirect).Methods("GET", "OPTIONS")
	return router
}
