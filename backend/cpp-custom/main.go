package main

import (
	"cpp-custom/logger"
	"cpp-custom/middleware"
	"cpp-custom/router"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error.Println("Error loading .env file")
	}
	port := os.Getenv("GO_APP_PORT")
	middleware.Port = port
	r := router.Router()
	logger.Info.Println("Starting server on the port " + port + " ...") // TODO: added runserver configurations
	log.Fatal(http.ListenAndServe(":" + port, r))
}

