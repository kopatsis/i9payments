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

func UpdateSubscriptionPaymentMethod(auth *auth.Client, database *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req struct {
			PaymentMethodID string `json:"payment_method_id"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		_, userid, err := multipass.BothIDsFromCookie(c, auth, database)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		userpayment, err := getUserPayment(database, userid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		s, err := sub.Get(userpayment.SubscriptionID, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		params := &stripe.PaymentMethodAttachParams{
			Customer: stripe.String(s.Customer.ID),
		}
		_, err = paymentmethod.Attach(req.PaymentMethodID, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to attach payment method"})
			return
		}

		customerParams := &stripe.CustomerParams{
			InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
				DefaultPaymentMethod: stripe.String(req.PaymentMethodID),
			},
		}
		_, err = customer.Update(s.Customer.ID, customerParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
			return
		}

		_, err = sub.Update(s.ID, &stripe.SubscriptionParams{
			DefaultPaymentMethod: stripe.String(req.PaymentMethodID),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Payment method updated"})
	}
}
