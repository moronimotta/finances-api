package stripeRepository

import (
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

type stripeCheckoutRepository struct {
	stripeKey string
}

func NewStripeCheckoutRepository(key string) StripeCheckout {
	return &stripeCheckoutRepository{
		stripeKey: key,
	}
}

func (r *stripeCheckoutRepository) CreateCheckoutSession(priceID, customerID, successURL, cancelURL string) (string, error) {
	stripe.Key = r.stripeKey

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Customer:   stripe.String(customerID),
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
	}

	result, err := session.New(params)
	if err != nil {
		return "", err
	}

	return result.URL, nil
}
