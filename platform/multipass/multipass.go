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
			c.Redirect(http.StatusFound, "/pay")
			go func() {
				deleteSpecialCode(code, database)
			}()
			return
		}

		if _, err := authClient.VerifyIDToken(context.Background(), token); err != nil {
			c.Redirect(http.StatusFound, "/pay")
			go func() {
				deleteSpecialCode(code, database)
			}()
			return
		}

		if err := login.Cookie(token, refresh, authClient, c); err != nil {
			c.Redirect(http.StatusFound, "/pay")
		}

		go func() {
			deleteSpecialCode(code, database)
		}()

		c.Redirect(http.StatusFound, "/pay")

	}
}
