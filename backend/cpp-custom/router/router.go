package router

import (
	"cpp-custom/middleware"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	// common
	router.HandleFunc("/api/ping", middleware.Ping).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/check", middleware.CheckForErrors).Methods("POST", "OPTIONS")
	return router
}