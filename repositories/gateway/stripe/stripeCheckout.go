package stripeRepository

import (
	"os"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

type stripeCheckoutRepository struct {
	stripeKey string
}

type StripeCheckout interface {
	CreateCheckoutSession(priceID, cutomerID string) (string, error)
}

func NewStripeCheckoutRepository(key string) StripeCheckout {
	return &stripeCheckoutRepository{
		stripeKey: key,
	}
}

func (r *stripeCheckoutRepository) CreateCheckoutSession(priceID, customerID string) (string, error) {
	stripe.Key = r.stripeKey

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Customer: stripe.String(customerID),
		Mode:     stripe.String(string(stripe.CheckoutSessionModePayment)),
		UIMode:   stripe.String("embedded"),
	}

	if os.Getenv("STRIPE_SUCCESS_URL") != "" {
		params.ReturnURL = stripe.String(os.Getenv("STRIPE_SUCCESS_URL"))
	} else {
		params.ReturnURL = stripe.String("https://example.com/success")
	}

	result, err := session.New(params)
	if err != nil {
		return "", err
	}

	return result.ClientSecret, nil
}
