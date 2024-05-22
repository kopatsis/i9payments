package login

import (
	"context"
	"net/http"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func Cookie(IDToken string, authClient *auth.Client, c *gin.Context) error {
	expiresIn := time.Hour * 24 * 14
	sessionCookie, err := authClient.SessionCookie(context.Background(), IDToken, expiresIn)
	if err != nil {
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

	return nil
}
