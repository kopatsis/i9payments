package login

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func VerifyToken(authClient *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			IDToken string `json:"idToken"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			fmt.Println("wrong req")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}

		expiresIn := time.Hour * 24 * 14
		sessionCookie, err := authClient.SessionCookie(context.Background(), request.IDToken, expiresIn)
		if err != nil {
			fmt.Println("cookie err")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to create a session cookie"})
			return
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "session",
			Value:    sessionCookie,
			MaxAge:   int(expiresIn.Seconds()),
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
		})

		c.JSON(http.StatusOK, gin.H{"message": "Session cookie set successfully"})
	}
}
