package main

import (
	"log"
	"net/http"
	"os"

	"i9pay/platform"
	"i9pay/platform/login"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		if os.Getenv("APP_ENV") != "production" {
			log.Fatalf("Failed to load the env vars: %v", err)
		}
	}

	auth := login.InitFirebase()

	rtr := platform.New(auth)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4040"
	}

	if err := http.ListenAndServe(":"+port, rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
