package pay

import (
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
)

func PostPayment(auth *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.PostForm("email")
		userId := c.PostForm("userId")
		token := c.PostForm("token")
		subscription := c.PostForm("subscription")

		priceID := "price_1JHkW2LJHkW2LJHkW2LJHkW2" // Replace with your actual monthly price ID
		if subscription == "yearly" {
			priceID = "price_1JHkW2LJHkW2LJHkW2LJHkW3" // Replace with your actual yearly price ID
		}

		// Create a new customer in Stripe
		customerParams := &stripe.CustomerParams{
			Email: stripe.String(email),
		}
		customerParams.SetSource(token) // set the token for the customer's payment method
		customerParams.Metadata = map[string]string{
			"userId": userId, // include user ID in metadata
		}
		stripeCustomer, err := customer.New(customerParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Create a subscription for the customer
		subscriptionParams := &stripe.SubscriptionParams{
			Customer: stripe.String(stripeCustomer.ID),
			Items: []*stripe.SubscriptionItemsParams{
				{
					Price: stripe.String(priceID),
				},
			},
		}
		subscription, err := sub.New(subscriptionParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Respond with the subscription details
		c.JSON(http.StatusOK, gin.H{
			"subscriptionId": subscription.ID,
			"status":         subscription.Status,
		})
	}
}
