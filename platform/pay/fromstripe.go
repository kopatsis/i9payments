package pay

import (
	"fmt"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/invoice"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/sub"
)

func getPaymentMethodDetails(subscriptionID string) (string, string, string, error) {
	params := &stripe.InvoiceListParams{
		Subscription: stripe.String(subscriptionID),
	}
	params.Filters.AddFilter("limit", "", "10") // Retrieve more than one invoice

	i := invoice.List(params)

	for i.Next() {
		inv := i.Invoice()
		if inv.PaymentIntent != nil {
			paymentIntent, err := paymentintent.Get(inv.PaymentIntent.ID, nil)
			if err != nil {
				return "", "", "", fmt.Errorf("error retrieving payment intent: %v", err)
			}

			if paymentIntent.PaymentMethod != nil {
				paymentMethodID := paymentIntent.PaymentMethod.ID
				if paymentMethodID == "" {
					return "", "", "", fmt.Errorf("no payment method associated with the payment intent: %s", inv.PaymentIntent.ID)
				}

				paymentMethod, err := paymentmethod.Get(paymentMethodID, nil)
				if err != nil {
					return "", "", "", fmt.Errorf("error retrieving payment method: %v", err)
				}

				switch paymentMethod.Type {
				case stripe.PaymentMethodTypeCard:
					card := paymentMethod.Card
					return "Card", string(card.Brand), card.Last4, nil
				default:
					return string(paymentMethod.Type), "", "", nil
				}
			}
		}
	}

	return "", "", "", fmt.Errorf("no valid payment intent found for subscription ID: %s", subscriptionID)
}

func UpdateSubscriptionPlan(subscriptionID, newPriceID string) error {
	currentSub, err := sub.Get(subscriptionID, nil)
	if err != nil {
		return err
	}

	updateParams := &stripe.SubscriptionParams{
		Items: []*stripe.SubscriptionItemsParams{
			{
				ID:    stripe.String(currentSub.Items.Data[0].ID),
				Price: stripe.String(newPriceID),
			},
		},
	}

	_, err = sub.Update(subscriptionID, updateParams)
	if err != nil {
		return err
	}

	return nil
}
