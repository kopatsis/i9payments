package multipass

import (
	"i9pay/platform/login"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Multipass(authClient *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		refresh := c.Query("multipass")
		code := c.Query("code")

		if status := checkSpecialCode(code, database); !status || refresh == "" {
			c.Redirect(http.StatusFound, "/pay")
			go func() {
				deleteSpecialCode(code, database)
			}()
			return
		}

		refresh, token, err := login.GetNewIDToken(refresh)
		if err != nil {
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
