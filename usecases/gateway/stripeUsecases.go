package usecases

import (
	"encoding/json"
	stripeRepository "finances-api/repositories/gateway/stripe"
	"log"
	"os"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

type StripeUsecase struct {
	GatewayUsecase
}

func NewStripeUsecase() *GatewayUsecase {
	stripeRepo := stripeRepository.NewStripeRepository()
	return &GatewayUsecase{
		Repository: NewGatewayUsecase(stripeRepo),
	}
}

func (s *StripeUsecase) EventBus(payload []byte, signature string) error {

	// Verify the webhook signature
	event, err := webhook.ConstructEvent(payload, signature, os.Getenv("STRIPE_WEBHOOK_SECRET"))
	if err != nil {
		log.Printf("Error verifying webhook signature: %v\n", err)
	}

	// Handle the event based on its type
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			return err
		}
	case "customer.subscription.created":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			return err
		}
	case "customer.created":

	default:
		log.Printf("Unhandled event type: %s\n", event.Type)
	}
	return nil
}
