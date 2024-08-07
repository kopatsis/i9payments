package main

import (
	"log"
	"net/http"
	"os"

	"i9pay/db"
	"i9pay/platform"
	"i9pay/platform/login"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/account"
)

func main() {
	if err := godotenv.Load(); err != nil {
		if os.Getenv("APP_ENV") != "production" {
			log.Fatalf("Failed to load the env vars: %v", err)
		}
	}

	auth := login.InitFirebase()

	client, database, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Error while connecting to mongoDB: %s.\nExiting.", err)
	}
	defer db.DisConnectDB(client)

	stripe.Key = os.Getenv("STRIPE_SECRET")

	acct, err := account.Get()
	if err != nil {
		log.Fatalf("Stripe API key test failed: %v", err)
	}
	log.Printf("Stripe API key test succeeded: Account ID = %s, Email = %s", acct.ID, acct.Email)

	apiKey := os.Getenv("SENDGRID_KEY")
	if apiKey == "" {
		log.Fatal("SENDGRID_API_KEY environment variable is not set")
	}

	sendclient := sendgrid.NewSendClient(apiKey)

	rtr := platform.New(auth, database, sendclient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(":"+port, rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
