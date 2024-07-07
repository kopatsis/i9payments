package login

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

func InitFirebase() *auth.Client {
	firebaseConfigBase64 := os.Getenv("FIREBASE_CONFIG_BASE64")
	if firebaseConfigBase64 == "" {
		log.Fatal("FIREBASE_CONFIG_BASE64 environment variable is not set.")
	}

	configJSON, err := base64.StdEncoding.DecodeString(firebaseConfigBase64)
	if err != nil {
		log.Fatalf("Error decoding FIREBASE_CONFIG_BASE64: %v", err)
	}
	opt := option.WithCredentialsJSON(configJSON)
	fmt.Println(opt)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("Failed to initialize Firebase app")
	}

	auth, err := app.Auth(context.Background())
	if err != nil {
		panic("Failed to initialize Firebase auth")
	}

	return auth
}

func ExtractUIDFromSession(c *gin.Context, auth *auth.Client) (string, int64, error) {
	cookie, err := c.Cookie("session")
	if err != nil {
		fmt.Println("no cookie")
		return "", 0, fmt.Errorf("no session cookie found: %w", err)
	}

	token, err := auth.VerifySessionCookie(context.Background(), cookie)
	if err != nil {
		fmt.Println("invalid cookie")
		return "", 0, fmt.Errorf("invalid session cookie: %w", err)
	}

	return token.UID, token.IssuedAt, nil
}
