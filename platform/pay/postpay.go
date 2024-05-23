package pay

import (
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/sub"
)

func PostPayment(auth *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.PostForm("email")
		userId := c.PostForm("userId")
		token := c.PostForm("token")
		subscription := c.PostForm("subscription")

		priceID := "price_1PJfbQIstWH7VBmuNNsoLTN2"
		if subscription == "yearly" {
			priceID = "price_1PJfbpIstWH7VBmu1nToVdC9"
		}

		customerParams := &stripe.CustomerParams{
			Email: stripe.String(email),
		}
		customerParams.SetSource(token)
		customerParams.Metadata = map[string]string{
			"userId": userId,
		}

		stripeCustomer, err := customer.New(customerParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		subscriptionParams := &stripe.SubscriptionParams{
			Customer: stripe.String(stripeCustomer.ID),
			Items: []*stripe.SubscriptionItemsParams{
				{
					Price: stripe.String(priceID),
				},
			},
		}

		newsub, err := sub.New(subscriptionParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"subscriptionId": newsub.ID,
			"status":         newsub.Status,
		})
	}
}
