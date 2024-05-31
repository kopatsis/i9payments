package pay

import (
	"fmt"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/invoice"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/sub"
)

func getPaymentMethodDetails(subscriptionID string) (string, string, string, error) {
	params := &stripe.InvoiceListParams{
		Subscription: stripe.String(subscriptionID),
	}
	params.Filters.AddFilter("limit", "", "1")

	i := invoice.List(params)
	if !i.Next() {
		return "", "", "", fmt.Errorf("no invoices found for subscription ID: %s", subscriptionID)
	}

	inv := i.Invoice()
	paymentIntent, err := paymentintent.Get(inv.PaymentIntent.ID, nil)
	if err != nil {
		return "", "", "", fmt.Errorf("error retrieving payment intent: %v", err)
	}

	paymentMethod := paymentIntent.PaymentMethod
	if paymentMethod != nil {
		switch paymentMethod.Type {
		case stripe.PaymentMethodTypeCard:
			card := paymentMethod.Card
			return "Card", string(card.Brand), card.Last4, nil
		default:
			return string(paymentMethod.Type), "", "", nil
		}
	}
	return "", "", "", fmt.Errorf("no payment method found for payment intent: %s", inv.PaymentIntent.ID)
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
		ProrationBehavior: stripe.String("none"),
	}

	_, err = sub.Update(subscriptionID, updateParams)
	if err != nil {
		return err
	}

	return nil
}
