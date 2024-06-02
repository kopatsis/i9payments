package pay

import (
	"context"
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
		var request struct {
			PaymentMethodID string `json:"paymentMethodId"`
			PriceID         string `json:"priceId"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		uid, userid, err := multipass.BothIDsFromCookie(c, auth, database)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		userRecord, err := auth.GetUser(context.Background(), uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		email := userRecord.Email

		customerParams := &stripe.CustomerParams{
			Email: stripe.String(email),
		}
		customerParams.Metadata = map[string]string{
			"userId": userid,
		}

		stripeCustomer, err := customer.New(customerParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		attachParams := &stripe.PaymentMethodAttachParams{
			Customer: stripe.String(stripeCustomer.ID),
		}
		_, err = paymentmethod.Attach(request.PaymentMethodID, attachParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		customerUpdateParams := &stripe.CustomerParams{
			InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
				DefaultPaymentMethod: stripe.String(request.PaymentMethodID),
			},
		}
		_, err = customer.Update(stripeCustomer.ID, customerUpdateParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		subscriptionParams := &stripe.SubscriptionParams{
			Customer: stripe.String(stripeCustomer.ID),
			Items: []*stripe.SubscriptionItemsParams{
				{
					Price: stripe.String(request.PriceID),
				},
			},
		}
		subscriptionParams.AddMetadata("userId", userid)

		newSub, err := sub.New(subscriptionParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		length := "monthly"
		if request.PriceID == "price_1PJfbpIstWH7VBmu1nToVdC9" {
			length = "yearly"
		}

		if err := setUserPayingPartial(database, newSub.ID, uid, length, userid); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"subscriptionId": newSub.ID,
			"status":         newSub.Status,
		})
	}
}
