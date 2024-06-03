package login

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func Cookie(IDToken, refreshToken string, authClient *auth.Client, c *gin.Context) error {
	expiresIn := time.Hour * 24 * 14

	fmt.Println(IDToken)
	token, err := authClient.VerifyIDToken(context.Background(), IDToken)
	if err != nil {
		log.Printf("Error verifying ID token: %v", err)
	} else {
		log.Printf("Token verified successfully: %v", token)
	}

	sessionCookie, err := authClient.SessionCookie(context.Background(), IDToken, expiresIn)
	if err != nil {
		fmt.Println("here? issue: ", err.Error())
		return err
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session",
		Value:    sessionCookie,
		MaxAge:   int(expiresIn.Seconds()),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})

	return nil
}
