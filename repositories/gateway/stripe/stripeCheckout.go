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
	CreateCheckoutSession(priceID []string, cutomerID string) (string, error)
}

func NewStripeCheckoutRepository(key string) StripeCheckout {
	return &stripeCheckoutRepository{
		stripeKey: key,
	}
}

func (r *stripeCheckoutRepository) CreateCheckoutSession(priceID []string, customerID string) (string, error) {
	stripe.Key = r.stripeKey
	lineItems := []*stripe.CheckoutSessionLineItemParams{}
	for _, id := range priceID {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(id),
			Quantity: stripe.Int64(1),
		})
	}
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
		Customer:           stripe.String(customerID),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		UIMode:             stripe.String("embedded"),
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
