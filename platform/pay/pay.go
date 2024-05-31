package pay

import (
	"context"
	"i9pay/platform/login"
	"i9pay/platform/multipass"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/sub"
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

		userpayment, err := getUserPayment(database, user.ID.Hex())
		if err != nil {
			c.HTML(200, "error.tmpl", nil)
			return
		}

		if userpayment == nil {
			c.HTML(200, "pay.tmpl", gin.H{
				"Email": email,
			})
			return
		}

		if userpayment.Processing {
			c.HTML(200, "processing.tmpl", nil)
			return
		}

		if user.Paying {

			if user.Provider == "Apple" || user.Provider == "Android" {
				c.HTML(200, "external.tmpl", nil)
				return
			}

			s, err := sub.Get(userpayment.SubscriptionID, nil)
			if err != nil {
				c.HTML(200, "error.tmpl", nil)
				return
			}

			paymentType, cardBrand, lastFour, err := getPaymentMethodDetails(userpayment.SubscriptionID)
			if err != nil {
				c.HTML(200, "error.tmpl", nil)
				return
			}

			if paymentType != "Card" {
				c.HTML(200, "alreadypaying.tmpl", gin.H{
					"Email":        email,
					"External":     paymentType,
					"Customer":     s.Customer.ID,
					"Length":       userpayment.SubLength,
					"Subscription": userpayment.SubscriptionID,
				})
				return
			}

			c.HTML(200, "alreadypaying.tmpl", gin.H{
				"Email":        email,
				"Brand":        cardBrand,
				"Four":         lastFour,
				"Customer":     s.Customer.ID,
				"Length":       userpayment.SubLength,
				"Subscription": userpayment.SubscriptionID,
			})
			return
		}

		c.HTML(200, "error.tmpl", nil)

	}
}
