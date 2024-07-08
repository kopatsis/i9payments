package home

import (
	"context"
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
			c.Redirect(http.StatusFound, "/login")
			return
		}

		userRecord, err := auth.GetUser(context.Background(), uid)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		email := userRecord.Email
		user, err := multipass.UserFromUID(uid, database)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		issuedTime := time.Unix(iat, 0)
		resetTime := user.ResetDate.Time()

		if issuedTime.Before(resetTime) {
			login.CookieLogout(c)
			return
		}

		c.HTML(200, "admin.tmpl", gin.H{
			"Mobile": frommobile,
			"Verify": userRecord.EmailVerified,
			"Email":  email,
			"Paying": user.Paying,
			"Name":   user.Name,
		})
	}
}
