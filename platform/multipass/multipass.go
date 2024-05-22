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

		mongoID, email, err := emailAndIDfromToken(token, database)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "invalid user",
			})
		}

		if err := login.Cookie(token, authClient, c); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to create a session cookie"})
			return
		}

		go func() {
			deleteSpecialCode(code, database)
		}()

		c.HTML(http.StatusOK, "pay.tmpl", gin.H{
			"UserID":    mongoID,
			"UserEmail": email,
		})

	}
}
