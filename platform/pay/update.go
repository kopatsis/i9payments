package pay

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/sub"
)

func UpdateSubscriptionPaymentMethod() gin.HandlerFunc {
	return func(c *gin.Context) {

		var req struct {
			SubscriptionID  string `json:"subscription_id"`
			PaymentMethodID string `json:"payment_method_id"`
			CustomerID      string `json:"customer_id"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		params := &stripe.PaymentMethodAttachParams{
			Customer: stripe.String(req.CustomerID),
		}
		_, err := paymentmethod.Attach(req.PaymentMethodID, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to attach payment method"})
			return
		}

		customerParams := &stripe.CustomerParams{
			InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
				DefaultPaymentMethod: stripe.String(req.PaymentMethodID),
			},
		}
		_, err = customer.Update(req.CustomerID, customerParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
			return
		}

		_, err = sub.Update(req.SubscriptionID, &stripe.SubscriptionParams{
			DefaultPaymentMethod: stripe.String(req.PaymentMethodID),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Payment method updated"})
	}
}
