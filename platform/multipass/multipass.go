package multipass

import (
	"context"
	"i9pay/platform/login"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Multipass(authClient *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("multipass")
		refresh := c.Query("multipass")
		code := c.Query("code")

		if status := checkSpecialCode(code, database); !status || token == "" {
			c.JSON(400, gin.H{
				"error": "invalid code",
			})
		}

		if _, err := authClient.VerifyIDToken(context.Background(), token); err != nil {
			c.JSON(400, gin.H{
				"error": "invalid token",
			})
		}

		if err := login.Cookie(token, refresh, authClient, c); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to create a session cookie"})
			return
		}

		go func() {
			deleteSpecialCode(code, database)
		}()

		c.Redirect(http.StatusFound, "/pay")

	}
}
