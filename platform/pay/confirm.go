package pay

import (
	"encoding/json"
	"net/http"
	"os"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/webhook"
	"go.mongodb.org/mongo-driver/mongo"
)

func Webhook(auth *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		const MaxBodyBytes = int64(65536)
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

		payload, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading request body"})
			return
		}

		endpointSecret := os.Getenv("ENDPOINT_SECR")

		event, err := webhook.ConstructEvent(payload, c.GetHeader("Stripe-Signature"), endpointSecret)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error verifying webhook signature"})
			return
		}

		if event.Type == "invoice.payment_succeeded" {
			var invoice stripe.Invoice
			if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing webhook JSON"})
				return
			}

			subscriptionID := invoice.Subscription
			customerID := invoice.Customer.ID

			customerReal, err := customer.Get(customerID, nil)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error customer doesn't exist"})
				return
			}

			paymentIntent, _ := paymentintent.Get(invoice.PaymentIntent.ID, nil)
			userId := paymentIntent.Metadata["userId"]

			if userId == "" {
				userId = customerReal.Metadata["userId"]
			}

			if userId == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error user id doesn't exist in either metafield"})
				return
			}

			if err := setUserPaying(database, subscriptionID.ID, userId); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error updating the user"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "success"})

		}
	}
}
