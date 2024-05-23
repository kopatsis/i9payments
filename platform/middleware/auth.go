package middleware

import (
	"context"
	"i9pay/platform/login"
	"net/http"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authClient *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.URL.Path == "/login" {
			c.Next()
			return
		}

		cookie, err := c.Cookie("session")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		token, err := authClient.VerifySessionCookie(context.Background(), cookie)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		const maxAge = 14 * 24 * 60 * 60
		issuedAt := time.Unix(token.IssuedAt, 0)
		expirationTime := issuedAt.Add(time.Duration(maxAge) * time.Second)
		if time.Until(expirationTime) < (10 * 24 * time.Hour) {
			refreshToken, err := c.Cookie("refreshToken")
			if err != nil {
				c.Next()
				return
			}

			newID, newRefresh, err := login.GetNewIDToken(refreshToken)
			if err != nil {
				c.Next()
				return
			}

			login.Cookie(newID, newRefresh, authClient, c)
		}

		c.Next()
	}
}
