package login

import (
	"fmt"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func VerifyToken(authClient *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			IDToken      string `json:"idToken"`
			RefreshToken string `json:"refreshToken"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			fmt.Println("wrong req")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}

		if err := Cookie(request.IDToken, request.RefreshToken, authClient, c); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to create a session cookie"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Session cookie set successfully"})
	}
}
