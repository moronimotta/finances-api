package usecases

import (
	"encoding/json"
	"errors"
	"finances-api/entities"
	"finances-api/messaging"
	"finances-api/repositories"
	"finances-api/utils/meta"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/webhook"
	"gorm.io/gorm"
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

		if appointmentID, ok := barberAppointmentID(paymentIntent.Metadata); ok {
			if err := s.recordAppointmentPayment(&paymentIntent, appointmentID); err != nil {
				return err
			}
			return s.publishAppointmentPaymentStatus(appointmentID, "paid", paymentIntent.ID, "", paymentIntent.LatestCharge.ID)
		}

		var input entities.Transactions
		if paymentIntent.LatestCharge == nil || paymentIntent.LatestCharge.ID == "" {
			return errors.New("successful payment has no Stripe charge")
		}
		input.ExternalID = paymentIntent.LatestCharge.ID
		input.Status = string(paymentIntent.Status)
		input.PaymentStatus = "paid"
		input.AmountTotal = paymentIntent.Amount
		input.AmountPayed = paymentIntent.AmountReceived
		if paymentIntent.Customer != nil {
			input.UserExternalID = paymentIntent.Customer.ID
		}

		if err := s.CreateTransaction(input); err != nil {
			return err
		}

		// TODO: Send a CRM communication
	case "payment_intent.payment_failed":
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			return err
		}
		if appointmentID, ok := barberAppointmentID(paymentIntent.Metadata); ok {
			return s.setAppointmentPaymentStatus(appointmentID, "unpaid", paymentIntent.ID, "", paymentIntent.Metadata["checkout_attempt_id"])
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

		metadata, err := paymentIntentMetadataForCharge(&charge)
		if err != nil {
			return err
		}
		if appointmentID, ok := barberAppointmentID(metadata); ok {
			if err := s.recordAppointmentRefund(&charge, appointmentID); err != nil {
				return err
			}
			return s.publishAppointmentPaymentStatus(appointmentID, "refunded", paymentIntentID(&charge), "", charge.ID)
		}
		return s.recordRefund(&charge)
	case "charge.updated":
		var charge stripe.Charge
		if err := json.Unmarshal(event.Data.Raw, &charge); err != nil {
			return err
		}
		metadata, err := paymentIntentMetadataForCharge(&charge)
		if err != nil {
			return err
		}
		if appointmentID, ok := barberAppointmentID(metadata); ok {
			return s.recordAppointmentChargeUpdate(&charge, appointmentID)
		}
		return s.recordChargeUpdate(&charge)
	case "checkout.session.expired":
		var checkoutSession stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &checkoutSession); err != nil {
			return err
		}
		if appointmentID, ok := barberAppointmentID(checkoutSession.Metadata); ok {
			return s.setAppointmentPaymentStatus(appointmentID, "unpaid", "", checkoutSession.ID, checkoutSession.Metadata["checkout_attempt_id"])
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

func (s *StripeUsecase) recordAppointmentPayment(paymentIntent *stripe.PaymentIntent, appointmentID string) error {
	if paymentIntent.LatestCharge == nil || paymentIntent.LatestCharge.ID == "" {
		return errors.New("successful appointment payment has no Stripe charge")
	}
	chargeInfo, err := s.Repository.GetCharge(paymentIntent.LatestCharge.ID)
	if err != nil {
		return fmt.Errorf("load Stripe charge details: %w", err)
	}

	input := entities.Transactions{
		Type:                  "barber_appointment",
		PaymentStatus:         "paid",
		AppointmentID:         appointmentID,
		CheckoutAttemptID:     paymentIntent.Metadata["checkout_attempt_id"],
		StripePaymentIntentID: paymentIntent.ID,
		ExternalID:            paymentIntent.LatestCharge.ID,
		AmountTotal:           chargeInfo.AmountTotal,
		AmountPayed:           chargeInfo.AmountPayed,
		AmountRefunded:        chargeInfo.AmountRefunded,
		StripeFee:             chargeInfo.StripeFee,
		StripeFeeKnown:        chargeInfo.StripeFeeKnown,
		NetAmount:             chargeInfo.NetAmount,
		PaymentMethod:         chargeInfo.PaymentMethod,
		ReceiptURL:            chargeInfo.ReceiptURL,
		Currency:              chargeInfo.Currency,
	}
	if input.AmountTotal == 0 {
		input.AmountTotal = paymentIntent.Amount
	}
	if input.AmountPayed == 0 {
		input.AmountPayed = paymentIntent.AmountReceived
	}
	if input.AmountRefunded > 0 {
		input.PaymentStatus = "refunded"
	}
	input.NetAmount = input.AmountTotal - input.AmountRefunded - input.StripeFee
	if paymentIntent.Customer != nil {
		input.UserExternalID = paymentIntent.Customer.ID
	}
	input.Meta = meta.New()
	input.Meta.Add(barberPaymentContextKey, barberPaymentContext)
	input.Meta.Add(appointmentIDMetadata, appointmentID)

	existing, lookupErr := s.appointmentTransactionForWebhook(input.CheckoutAttemptID, input.StripePaymentIntentID, appointmentID)
	if lookupErr != nil && !errors.Is(lookupErr, gorm.ErrRecordNotFound) {
		return lookupErr
	}
	if lookupErr == nil {
		input.ID = existing.ID
		input.StripeSessionID = existing.StripeSessionID
		input.CheckoutURL = existing.CheckoutURL
		input.CreatedAt = existing.CreatedAt
		input.Meta.Merge(existing.Meta)
		return s.DB.UpdateTransaction(&input)
	}
	return s.DB.CreateTransaction(&input)
}

func (s *StripeUsecase) recordAppointmentRefund(charge *stripe.Charge, appointmentID string) error {
	chargeInfo, err := s.Repository.GetCharge(charge.ID)
	if err != nil {
		return fmt.Errorf("load refunded Stripe charge details: %w", err)
	}
	input := entities.Transactions{
		Type:                  "barber_appointment",
		PaymentStatus:         "refunded",
		AppointmentID:         appointmentID,
		StripePaymentIntentID: paymentIntentID(charge),
		ExternalID:            charge.ID,
		AmountTotal:           chargeInfo.AmountTotal,
		AmountPayed:           chargeInfo.AmountPayed,
		AmountRefunded:        chargeInfo.AmountRefunded,
		StripeFee:             chargeInfo.StripeFee,
		StripeFeeKnown:        chargeInfo.StripeFeeKnown,
		NetAmount:             chargeInfo.AmountTotal - chargeInfo.AmountRefunded - chargeInfo.StripeFee,
		PaymentMethod:         chargeInfo.PaymentMethod,
		ReceiptURL:            chargeInfo.ReceiptURL,
		Currency:              chargeInfo.Currency,
	}
	if charge.PaymentIntent != nil {
		input.CheckoutAttemptID = charge.PaymentIntent.Metadata["checkout_attempt_id"]
	}
	input.Meta = meta.New()
	input.Meta.Add(barberPaymentContextKey, barberPaymentContext)
	input.Meta.Add(appointmentIDMetadata, appointmentID)

	existing, lookupErr := s.appointmentTransactionForWebhook(input.CheckoutAttemptID, input.StripePaymentIntentID, appointmentID)
	if lookupErr != nil && !errors.Is(lookupErr, gorm.ErrRecordNotFound) {
		return lookupErr
	}
	if lookupErr == nil {
		input.ID = existing.ID
		input.StripeSessionID = existing.StripeSessionID
		input.CheckoutURL = existing.CheckoutURL
		input.CreatedAt = existing.CreatedAt
		input.Meta.Merge(existing.Meta)
		return s.DB.UpdateTransaction(&input)
	}
	return s.DB.CreateTransaction(&input)
}

func (s *StripeUsecase) appointmentTransactionForWebhook(attemptID, paymentIntentID, appointmentID string) (*entities.Transactions, error) {
	if attemptID != "" {
		transaction, err := s.DB.GetTransactionByCheckoutAttemptID(attemptID)
		if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
			return transaction, err
		}
	}
	if paymentIntentID != "" {
		transaction, err := s.DB.GetTransactionByPaymentIntentID(paymentIntentID)
		if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
			return transaction, err
		}
	}
	transaction, err := s.DB.GetTransactionByAppointmentID(appointmentID)
	if err == nil && transaction.PaymentStatus != "pending" {
		return nil, gorm.ErrRecordNotFound
	}
	return transaction, err
}

func (s *StripeUsecase) recordRefund(charge *stripe.Charge) error {
	transaction, err := s.DB.GetTransactionByExternalID(charge.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	chargeInfo, err := s.Repository.GetCharge(charge.ID)
	if err != nil {
		return err
	}
	transaction.PaymentStatus = "refunded"
	transaction.Status = "refunded"
	transaction.AmountRefunded = chargeInfo.AmountRefunded
	transaction.StripeFee = chargeInfo.StripeFee
	transaction.StripeFeeKnown = chargeInfo.StripeFeeKnown
	transaction.NetAmount = transaction.AmountTotal - chargeInfo.AmountRefunded - chargeInfo.StripeFee
	return s.DB.UpdateTransaction(transaction)
}

func (s *StripeUsecase) recordAppointmentChargeUpdate(charge *stripe.Charge, appointmentID string) error {
	transaction, err := s.appointmentTransactionForWebhook("", paymentIntentID(charge), appointmentID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	chargeInfo, err := s.Repository.GetCharge(charge.ID)
	if err != nil {
		return err
	}
	transaction.AmountRefunded = chargeInfo.AmountRefunded
	transaction.StripeFee = chargeInfo.StripeFee
	transaction.StripeFeeKnown = chargeInfo.StripeFeeKnown
	transaction.NetAmount = transaction.AmountTotal - transaction.AmountRefunded - transaction.StripeFee
	transaction.ReceiptURL = chargeInfo.ReceiptURL
	transaction.PaymentMethod = chargeInfo.PaymentMethod
	if transaction.AmountRefunded > 0 {
		transaction.PaymentStatus = "refunded"
	}
	return s.DB.UpdateTransaction(transaction)
}

func (s *StripeUsecase) recordChargeUpdate(charge *stripe.Charge) error {
	transaction, err := s.DB.GetTransactionByExternalID(charge.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	chargeInfo, err := s.Repository.GetCharge(charge.ID)
	if err != nil {
		return err
	}
	transaction.AmountRefunded = chargeInfo.AmountRefunded
	transaction.StripeFee = chargeInfo.StripeFee
	transaction.StripeFeeKnown = chargeInfo.StripeFeeKnown
	transaction.NetAmount = transaction.AmountTotal - transaction.AmountRefunded - transaction.StripeFee
	transaction.ReceiptURL = chargeInfo.ReceiptURL
	transaction.PaymentMethod = chargeInfo.PaymentMethod
	return s.DB.UpdateTransaction(transaction)
}

func (s *StripeUsecase) setAppointmentPaymentStatus(appointmentID, status, paymentIntentID, sessionID, attemptID string) error {
	transaction, err := s.appointmentTransactionForWebhook(attemptID, paymentIntentID, appointmentID)
	if err != nil {
		return fmt.Errorf("load appointment payment for status update: %w", err)
	}
	if sessionID != "" && transaction.StripeSessionID != sessionID {
		return nil // Ignore late expiration events for an older Checkout Session.
	}
	if attemptID != "" && transaction.CheckoutAttemptID != attemptID {
		return nil
	}
	if paymentIntentID != "" && transaction.StripePaymentIntentID != "" && transaction.StripePaymentIntentID != paymentIntentID {
		return nil
	}
	if status == "unpaid" && (transaction.PaymentStatus == "paid" || transaction.PaymentStatus == "refunded") {
		return nil
	}
	transaction.PaymentStatus = status
	if paymentIntentID != "" {
		transaction.StripePaymentIntentID = paymentIntentID
	}
	if err := s.DB.UpdateTransaction(transaction); err != nil {
		return err
	}
	return s.publishAppointmentPaymentStatus(appointmentID, status, paymentIntentID, sessionID, "")
}

func (s *StripeUsecase) publishAppointmentPaymentStatus(appointmentID, status, paymentIntentID, sessionID, chargeID string) error {
	if s.Pub == nil {
		return nil
	}
	return s.Pub.Publish(barberEventsTopic, "appointment_payment_status_update", map[string]interface{}{
		"appointment_id":    appointmentID,
		"payment_status":    status,
		"payment_intent_id": paymentIntentID,
		"session_id":        sessionID,
		"charge_id":         chargeID,
	})
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
	input.StripeFee = stripeChargeInfo.StripeFee
	input.StripeFeeKnown = stripeChargeInfo.StripeFeeKnown
	input.NetAmount = stripeChargeInfo.NetAmount
	if input.NetAmount == 0 && input.AmountTotal > 0 && input.AmountRefunded == 0 {
		input.NetAmount = input.AmountTotal - input.StripeFee
	}
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
