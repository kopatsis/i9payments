package pay

import (
	"context"
	"fmt"
	"i9pay/platform/login"
	"i9pay/platform/multipass"
	"net/http"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/setupintent"
	"github.com/stripe/stripe-go/v72/sub"
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
			c.HTML(200, "error.tmpl", gin.H{"Error": err.Error()})
			return
		}

		if userpayment == nil {

			params := &stripe.SetupIntentParams{
				PaymentMethodTypes: stripe.StringSlice([]string{
					"card",
				}),
			}
			si, err := setupintent.New(params)
			if err != nil {
				c.HTML(200, "error.tmpl", gin.H{"Error": err.Error()})
				return
			}

			c.HTML(200, "pay.tmpl", gin.H{
				"ClientSecret": si.ClientSecret,
				"Email":        email,
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
				c.HTML(200, "error.tmpl", gin.H{"Error": err.Error()})
				return
			}

			if userpayment.Ending {
				c.HTML(200, "ending.html", gin.H{
					"Date": userpayment.EndDate.Time().Format("01/02/2006"),
				})
			}

			paymentType, cardBrand, lastFour, err := getPaymentMethodDetails(userpayment.SubscriptionID)
			if err != nil {
				c.HTML(200, "error.tmpl", gin.H{"Error": err.Error()})
				return
			}

			params := &stripe.SetupIntentParams{
				PaymentMethodTypes: stripe.StringSlice([]string{
					"card",
				}),
			}
			si, err := setupintent.New(params)
			if err != nil {
				c.HTML(200, "error.tmpl", gin.H{"Error": err.Error()})
				return
			}

			fmt.Println(paymentType, cardBrand, lastFour)

			if paymentType != "Card" {
				c.HTML(200, "alreadypaying.tmpl", gin.H{
					"ClientSecret": si.ClientSecret,
					"Date":         time.Unix(s.CurrentPeriodEnd, 0).Format("01/02/2006"),
					"Email":        email,
					"External":     paymentType,
					"Length":       userpayment.SubLength,
				})
				return
			}

			c.HTML(200, "alreadypaying.tmpl", gin.H{
				"ClientSecret": si.ClientSecret,
				"Date":         time.Unix(s.CurrentPeriodEnd, 0).Format("01/02/2006"),
				"Email":        email,
				"Brand":        cardBrand,
				"Four":         lastFour,
				"Length":       userpayment.SubLength,
			})
			return
		}

		c.HTML(200, "error.tmpl", gin.H{"Error": "User payment exists and processing false, but user obj paying false (shouldn't be possible)"})

	}
}
