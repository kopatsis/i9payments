package home

import (
	"context"
	"i9pay/platform/login"
	"i9pay/platform/multipass"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AdminPanel(frommobile, verifmessage bool, auth *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := login.ExtractUIDFromSession(c, auth)
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

		c.HTML(200, "admin.tmpl", gin.H{
			"NotMobile": !frommobile,
			"Verify":    verifmessage,
			"Email":     email,
			"Paying":    user.Paying,
			"Name":      user.Name,
		})
	}
}
