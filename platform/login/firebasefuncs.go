package login

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

func InitFirebase() *auth.Client {
	opt := option.WithCredentialsFile("i9auth-firebase-adminsdk-dgzg6-f59f9349ed.json")
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

func ExtractUIDFromSession(c *gin.Context, auth *auth.Client) (string, error) {
	cookie, err := c.Cookie("session")
	if err != nil {
		fmt.Println("no cookie")
		return "", fmt.Errorf("no session cookie found: %w", err)
	}

	token, err := auth.VerifySessionCookie(context.Background(), cookie)
	if err != nil {
		fmt.Println("invalid cookie")
		return "", fmt.Errorf("invalid session cookie: %w", err)
	}

	return token.UID, nil
}
