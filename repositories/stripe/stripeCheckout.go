package stripeRepository

import (
	"finances-api/db"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

type stripeCheckoutRepository struct {
	db        db.Database
	stripeKey string
}

func NewStripeCheckoutRepository(db db.Database, key string) StripeCheckout {
	return &stripeCheckoutRepository{
		db:        db,
		stripeKey: key,
	}
}

func (r *stripeCheckoutRepository) CreateCheckoutSession(priceID, successURL, cancelURL string) (string, error) {
	stripe.Key = r.stripeKey

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
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
