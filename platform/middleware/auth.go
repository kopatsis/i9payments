package middleware

import (
	"context"
	"i9pay/platform/login"
	"net/http"
	"strings"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authClient *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		if strings.HasPrefix(c.Request.URL.Path, "/static") {
			c.Next()
			return
		}

		if c.Request.URL.Path == "/login" || c.Request.URL.Path == "/verifyToken" || c.Request.URL.Path == "/new" || c.Request.URL.Path == "/code" || c.Request.URL.Path == "/multipass" || c.Request.URL.Path == "/confirmationwh" || c.Request.URL.Path == "/failedwh" || c.Request.URL.Path == "/resetdate" {
			c.Next()
			return
		}

		cookie, err := c.Cookie("session")
		if err != nil {
			if c.Request.Method == http.MethodGet {
				returnTo := c.Request.URL.Path
				if len(returnTo) > 0 && returnTo[0] == '/' {
					returnTo = returnTo[1:]
				}
				c.Redirect(http.StatusFound, "/login?returnTo="+returnTo)
			} else {
				c.Redirect(http.StatusFound, "/login")
			}
			c.Abort()
			return
		}

		token, err := authClient.VerifySessionCookie(context.Background(), cookie)
		if err != nil {
			if c.Request.Method == http.MethodGet {
				returnTo := c.Request.URL.Path
				if len(returnTo) > 0 && returnTo[0] == '/' {
					returnTo = returnTo[1:]
				}
				c.Redirect(http.StatusFound, "/login?returnTo="+returnTo)
			} else {
				c.Redirect(http.StatusFound, "/login")
			}
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
