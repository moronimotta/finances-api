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

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/webhook"
)

const (
	barberPaymentContextKey = "payment_context"
	barberPaymentContext    = "barber"
	appointmentIDMetadata   = "appointment_id"
	barberEventsTopic       = "barber-app"
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
	event, err := webhook.ConstructEventWithOptions(
		payload,
		signature,
		os.Getenv("STRIPE_SECRET_ENDPOINT"),
		webhook.ConstructEventOptions{
			// Local Stripe CLI events can use an older API version than this SDK.
			// Keep signature and timestamp verification enabled; opt into only the
			// API-version compatibility check through local configuration.
			IgnoreAPIVersionMismatch: os.Getenv("STRIPE_IGNORE_API_VERSION_MISMATCH") == "true",
		},
	)
	if err != nil {
		log.Printf("Error verifying webhook signature: %v\n", err)
		return err
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

		if err := s.CreateTransaction(input); err != nil {
			return err
		}

		if appointmentID, ok := barberAppointmentID(paymentIntent.Metadata); ok && s.Pub != nil {
			return s.Pub.Publish(barberEventsTopic, "appointment_status_update", map[string]interface{}{
				"appointment_id":    appointmentID,
				"status":            "scheduled",
				"payment_event":     "payment_intent.succeeded",
				"payment_intent_id": paymentIntent.ID,
				"charge_id":         paymentIntent.LatestCharge.ID,
			})
		}

		// TODO: Send a CRM communication

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

		metadata, err := paymentIntentMetadataForCharge(&charge)
		if err != nil {
			return err
		}
		if appointmentID, ok := barberAppointmentID(metadata); ok && s.Pub != nil {
			return s.Pub.Publish(barberEventsTopic, "appointment_status_update", map[string]interface{}{
				"appointment_id":    appointmentID,
				"status":            "canceled",
				"payment_event":     "charge.refunded",
				"payment_intent_id": paymentIntentID(&charge),
				"charge_id":         charge.ID,
			})
		}
	default:
		log.Printf("Unhandled ev1ent type: %s\n", event.Type)
	}
	return nil
}

// barberAppointmentID accepts only the explicit Barber payment context and a
// syntactically valid appointment UUID. Metadata without both values is not a
// Barber payment and must not produce Barber events.
func barberAppointmentID(metadata map[string]string) (string, bool) {
	if metadata[barberPaymentContextKey] != barberPaymentContext {
		return "", false
	}
	appointmentID := strings.TrimSpace(metadata[appointmentIDMetadata])
	parsedID, err := uuid.Parse(appointmentID)
	if err != nil {
		return "", false
	}
	return parsedID.String(), true
}

func paymentIntentMetadataForCharge(charge *stripe.Charge) (map[string]string, error) {
	if charge.PaymentIntent != nil && charge.PaymentIntent.Metadata != nil {
		return charge.PaymentIntent.Metadata, nil
	}
	if charge.PaymentIntent == nil || charge.PaymentIntent.ID == "" {
		return charge.Metadata, nil
	}

	stripe.Key = os.Getenv("STRIPE_KEY")
	intent, err := paymentintent.Get(charge.PaymentIntent.ID, nil)
	if err != nil {
		return nil, fmt.Errorf("retrieve payment intent %s for refund metadata: %w", charge.PaymentIntent.ID, err)
	}
	return intent.Metadata, nil
}

func paymentIntentID(charge *stripe.Charge) string {
	if charge.PaymentIntent == nil {
		return ""
	}
	return charge.PaymentIntent.ID
}

func (s *StripeUsecase) CreateTransaction(input entities.Transactions) error {
	// Check if transaction already exists (idempotency check)
	existingTx, err := s.DB.GetTransactionByExternalID(input.ExternalID)
	if err == nil && existingTx != nil {
		log.Printf("Transaction with external_id %s already exists, skipping duplicate creation", input.ExternalID)
		return nil // Idempotent - transaction already processed
	}

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

	var localProductIDs []string
	for key, value := range stripeChargeInfo.Meta {
		if strings.HasPrefix(key, "local_product_id_") {
			localProductIDs = append(localProductIDs, value)
		}
	}

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
