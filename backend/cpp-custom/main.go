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
		panic("error loading .env file")
	}
	if os.Getenv("ENABLE_LOGGING") == "true" {
		loggers := make(map[string]string)
		loggers["memory_l"] = "memory"
		loggers["procedures_tree_l"] = "procedures_tree"
		loggers["tree_l"] = "tree"
		err = logger.Init(loggers)
		if err != nil {
			panic("error logger initializing")
		}
	}
	port := os.Getenv("GO_APP_PORT")
	middleware.Port = port
	r := router.Router()
	logger.Info.Println("Starting server on the port " + port + " ...") // TODO: added runserver configurations
	log.Fatal(http.ListenAndServe(":"+port, r))
}
