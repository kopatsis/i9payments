package pay

import (
	"i9pay/platform/multipass"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/sub"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostPayment(auth *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.PostForm("email")

		paymentMethodID := c.PostForm("paymentMethod")
		subscription := c.PostForm("subscription")

		uid, userid, err := multipass.BothIDsFromCookie(c, auth, database)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		priceID := "price_1PJfbQIstWH7VBmuNNsoLTN2"
		if subscription == "yearly" {
			priceID = "price_1PJfbpIstWH7VBmu1nToVdC9"
		}

		customerParams := &stripe.CustomerParams{
			Email:         stripe.String(email),
			PaymentMethod: stripe.String(paymentMethodID),
		}
		customerParams.Metadata = map[string]string{
			"userId": userid,
		}

		stripeCustomer, err := customer.New(customerParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		params := &stripe.PaymentMethodAttachParams{
			Customer: stripe.String(stripeCustomer.ID),
		}
		_, err = paymentmethod.Attach(paymentMethodID, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		customerParamsUpdate := &stripe.CustomerParams{
			InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
				DefaultPaymentMethod: stripe.String(paymentMethodID),
			},
		}
		_, err = customer.Update(stripeCustomer.ID, customerParamsUpdate)
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
		subscriptionParams.AddMetadata("userId", userid)

		newsub, err := sub.New(subscriptionParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := setUserPayingPartial(database, newsub.ID, uid, subscription, userid); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"subscriptionId": newsub.ID,
			"status":         newsub.Status,
		})
	}
}
