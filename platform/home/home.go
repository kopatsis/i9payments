package home

import (
	"context"
	"html"
	"i9pay/platform/login"
	"i9pay/platform/multipass"
	"net/http"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AdminPanel(frommobile bool, auth *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, iat, err := login.ExtractUIDFromSession(c, auth)
		if err != nil {
			if frommobile {
				c.Redirect(http.StatusFound, "/login?returnTo=mobile")
			} else {
				c.Redirect(http.StatusFound, "/login")
			}
			return
		}

		userRecord, err := auth.GetUser(context.Background(), uid)
		if err != nil {
			if frommobile {
				c.Redirect(http.StatusFound, "/login?returnTo=mobile")
			} else {
				c.Redirect(http.StatusFound, "/login")
			}
			return
		}

		email := userRecord.Email
		user, err := multipass.UserFromUID(uid, database)
		if err != nil {
			if frommobile {
				c.Redirect(http.StatusFound, "/login?returnTo=mobile")
			} else {
				c.Redirect(http.StatusFound, "/login")
			}
			return
		}

		issuedTime := time.Unix(iat, 0)
		resetTime := user.ResetDate.Time()

		if issuedTime.Before(resetTime) {
			login.CookieLogout(c)
			return
		}

		properLogin := true

		cookie, err := c.Cookie("properLogin")
		if err != nil {
			properLogin = false
		} else {
			if cookie != "TRUE" {
				properLogin = false
			}
		}

		c.HTML(200, "admin.tmpl", gin.H{
			"ClientOn": properLogin,
			"Mobile":   frommobile,
			"Verify":   userRecord.EmailVerified,
			"Email":    html.EscapeString(email),
			"Paying":   user.Paying,
			"Name":     html.EscapeString(user.Name),
		})
	}
}
