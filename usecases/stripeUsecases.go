package usecases

import (
	"encoding/json"
	"finances-api/entities"
	"finances-api/messaging"
	"finances-api/repositories"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

type StripeUsecase struct {
	GatewayUsecase
	DB  repositories.FinancialRepository
	Pub messaging.Publisher
}

func NewStripeUsecase(usecases *PaymentAPIUsecases) *StripeUsecase {
	return &StripeUsecase{
		GatewayUsecase: *NewGatewayUsecase("stripe"),
		DB:             usecases.Db,
		Pub:            usecases.Pub,
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
		input.ExternalID = paymentIntent.LatestCharge.ID
		input.Status = string(paymentIntent.Status)
		input.AmountTotal = paymentIntent.Amount
		input.AmountPayed = paymentIntent.AmountReceived
		input.UserExternalID = paymentIntent.Customer.ID
		// Send via meta the local user id, product_id and product external id

		if err := s.CreateTransaction(input); err != nil {
			return err
		}

		// TODO: send a message to rabbitmq to akademia api to create user progress or append with the new course!!!

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

func (s *StripeUsecase) CreateTransaction(input entities.Transactions) error {
	stripeChargeInfo, err := s.Repository.GetCharge(input.ExternalID)
	if err != nil {
		return err
	}

	input.PaymentMethod = stripeChargeInfo.PaymentMethod
	input.Currency = stripeChargeInfo.Currency
	input.AmountRefunded = stripeChargeInfo.AmountRefunded
	if stripeChargeInfo.Status != "" {
		input.Status = stripeChargeInfo.Status
	}
	input.ReceiptURL = stripeChargeInfo.ReceiptURL

	err = s.DB.CreateTransaction(&input)
	// create transaction items
	if err != nil {
		return err
	}

	var priceIDs []string
	for key, value := range stripeChargeInfo.Meta {
		if strings.HasPrefix(key, "external_id_") {
			priceIDs = append(priceIDs, value)
		}
	}

	items, err := s.GetPrice(priceIDs)
	if err != nil {
		return err
	}

	for i := range items {
		items[i].TransactionID = input.ID
	}

	s.DB.CreateTransactionItems(items)

	// get the local_product_id_...
	var localProductIDs []string
	for key, value := range stripeChargeInfo.Meta {
		if strings.HasPrefix(key, "local_product_id_") {
			localProductIDs = append(localProductIDs, value)
		}
	}

	// get user id
	var userID string
	for key, value := range stripeChargeInfo.Meta {
		if strings.HasPrefix(key, "user_id") {
			userID = value
		}
	}

	if s.Pub != nil {
		_ = s.Pub.Publish(
			"akademia-api",
			"user.new_course",
			map[string]interface{}{
				"user_id":           userID,
				"local_product_ids": localProductIDs,
			},
		)
	}

	return nil
}
