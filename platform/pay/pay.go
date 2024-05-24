package pay

import (
	"context"
	"i9pay/platform/login"
	"i9pay/platform/multipass"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Subscription(auth *auth.Client, database *mongo.Database) gin.HandlerFunc {
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

		if user.Paying {

			external := ""
			cardBrand := ""
			lastFour := ""
			if user.Provider == "Apple" || user.Provider == "Android" {
				external = user.Provider
			} else {
				_, cardBrand, lastFour, err = getPaymentMethodDetails(user.Provider)
				if err != nil {
					c.HTML(200, "error.tmpl", nil)
					return
				}
			}

			c.HTML(200, "alreadypaying.tmpl", gin.H{
				"Email":    email,
				"UserID":   user.ID.Hex(),
				"External": external,
				"Brand":    cardBrand,
				"Four":     lastFour,
			})
			return
		}

		c.HTML(200, "pay.tmpl", gin.H{
			"UserEmail": email,
			"UserID":    user.ID.Hex(),
		})

	}
}
