package stripeRepository

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"finances-api/entities"
	"finances-api/utils/meta"
	"os"
	"strings"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

type stripeCheckoutRepository struct {
	stripeKey string
}

type StripeCheckout interface {
	CreateCheckoutSession(priceID []string, customerID string, meta meta.Meta) (string, error)
	CreateAppointmentCheckoutSession(priceIDs []string, appointmentID string, metadata meta.Meta) (entities.CheckoutSession, error)
}

func (r *stripeCheckoutRepository) CreateAppointmentCheckoutSession(priceIDs []string, appointmentID string, metadata meta.Meta) (entities.CheckoutSession, error) {
	stripe.Key = r.stripeKey
	lineItems := make([]*stripe.CheckoutSessionLineItemParams, 0, len(priceIDs))
	for _, priceID := range priceIDs {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(priceID),
			Quantity: stripe.Int64(1),
		})
	}

	successURL := os.Getenv("STRIPE_SUCCESS_URL")
	if successURL == "" {
		successURL = "http://localhost:3000/payment/success?session_id={CHECKOUT_SESSION_ID}"
	}
	cancelURL := os.Getenv("STRIPE_CANCEL_URL")
	if cancelURL == "" {
		cancelURL = "http://localhost:3000/booking?checkout=cancelled"
	}

	params := &stripe.CheckoutSessionParams{
		Mode:              stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems:         lineItems,
		SuccessURL:        stripe.String(successURL),
		CancelURL:         stripe.String(cancelURL),
		Metadata:          metadata,
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{Metadata: metadata},
	}
	keyMaterial := appointmentID + ":" + strings.Join(priceIDs, ",") + ":" + metadata["checkout_attempt_id"]
	keyDigest := sha256.Sum256([]byte(keyMaterial))
	params.SetIdempotencyKey("barber-appointment-" + hex.EncodeToString(keyDigest[:]))

	result, err := session.New(params)
	if err != nil {
		return entities.CheckoutSession{}, err
	}
	if result == nil || result.ID == "" || result.URL == "" {
		return entities.CheckoutSession{}, errors.New("Stripe returned an incomplete Checkout Session")
	}
	return entities.CheckoutSession{ID: result.ID, URL: result.URL}, nil
}

func NewStripeCheckoutRepository(key string) StripeCheckout {
	return &stripeCheckoutRepository{
		stripeKey: key,
	}
}

func (r *stripeCheckoutRepository) CreateCheckoutSession(priceID []string, customerID string, meta meta.Meta) (string, error) {
	stripe.Key = r.stripeKey
	lineItems := []*stripe.CheckoutSessionLineItemParams{}
	for _, id := range priceID {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(id),
			Quantity: stripe.Int64(1),
		})
	}

	paymentIntentData := &stripe.CheckoutSessionPaymentIntentDataParams{
		Metadata: meta,
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
		Customer:           stripe.String(customerID),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		UIMode:             stripe.String("embedded"),
		PaymentIntentData:  paymentIntentData,
	}

	if automaticTaxEnabled() {
		params.AutomaticTax = &stripe.CheckoutSessionAutomaticTaxParams{Enabled: stripe.Bool(true)}
	}

	if os.Getenv("STRIPE_SUCCESS_URL") != "" {
		params.ReturnURL = stripe.String(os.Getenv("STRIPE_SUCCESS_URL"))
	} else {
		params.ReturnURL = stripe.String("http://localhost:3000/payment/success")
	}

	result, err := session.New(params)
	if err != nil {
		return "", err
	}

	return result.ClientSecret, nil
}
