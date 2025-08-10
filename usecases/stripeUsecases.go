package usecases

import (
	"encoding/json"
	"finances-api/entities"
	"finances-api/repositories"
	"fmt"
	"log"
	"os"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

type StripeUsecase struct {
	GatewayUsecase
	DB repositories.FinancialRepository
}

func NewStripeUsecase(usecases *PaymentAPIUsecases) *StripeUsecase {
	return &StripeUsecase{
		GatewayUsecase: *NewGatewayUsecase("stripe"),
		DB:             usecases.Db,
	}
}

func (s *StripeUsecase) EventBus(payload []byte, signature string) error {

	// Verify the webhook signature
	event, err := webhook.ConstructEvent(payload, signature, os.Getenv("STRIPE_SECRET_ENDPOINT"))
	if err != nil {
		log.Printf("Error verifying webhook signature: %v\n", err)
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			return err
		}

		var input entities.Transactions
		input.ExternalID = paymentIntent.ID
		input.Status = string(paymentIntent.Status)
		input.AmountTotal = paymentIntent.Amount
		input.AmountPayed = paymentIntent.AmountReceived

		if err := s.DB.CreateTransaction(&input); err != nil {
			return err
		}

	case "customer.subscription.created":
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			return err
		}
	case "customer.created":
		fmt.Printf("Customer created: %s\n", event.Data.Object)
		return nil
	case "charge.refunded":
		var charge stripe.Charge
		err := json.Unmarshal(event.Data.Raw, &charge)
		if err != nil {
			return err
		}
	default:
		log.Printf("Unhandled ev1ent type: %s\n", event.Type)
	}
	return nil
}

// func (s *StripeUsecase) CreateTransaction(input entities.Transactions) error {
// 	stripeChargeInfo, err := s.Repository.GetCharge(input.ExternalID)
// 	if err != nil {
// 		return err
// 	}

// 	input.PaymentMethod = stripeChargeInfo.PaymentMethod
// 	input.Currency = stripeChargeInfo.Currency
// 	input.AmountRefunded = stripeChargeInfo.AmountRefunded
// 	input.Status = stripeChargeInfo.Status
// 	input.ReceiptURL = stripeChargeInfo.ReceiptURL

// 	err = s.Repository.CreateTransaction(input)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
