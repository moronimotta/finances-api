package usecases

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"

	"finances-api/entities"
	"finances-api/utils/meta"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidAppointmentCheckout    = errors.New("invalid appointment checkout request")
	ErrAppointmentCheckoutExists     = errors.New("appointment already has a pending or completed payment")
	ErrAppointmentServiceUnavailable = errors.New("appointment service is unavailable for payment")
)

// CreateAppointmentCheckoutSession uses the Barber service catalog synchronized
// by service.created events. The browser supplies identifiers only; prices and
// Stripe Price IDs always come from Payment API's persisted catalog.
func (s *StripeUsecase) CreateAppointmentCheckoutSession(appointmentID string, serviceIDs []string) (entities.CheckoutSession, int64, string, error) {
	parsedAppointmentID, err := uuid.Parse(strings.TrimSpace(appointmentID))
	if err != nil {
		return entities.CheckoutSession{}, 0, "", fmt.Errorf("%w: appointment_id must be a UUID", ErrInvalidAppointmentCheckout)
	}
	appointmentID = parsedAppointmentID.String()
	if len(serviceIDs) == 0 || len(serviceIDs) > 25 {
		return entities.CheckoutSession{}, 0, "", fmt.Errorf("%w: provide between 1 and 25 service IDs", ErrInvalidAppointmentCheckout)
	}

	canonicalIDs := make([]string, 0, len(serviceIDs))
	seen := make(map[string]struct{}, len(serviceIDs))
	for _, serviceID := range serviceIDs {
		parsedID, err := uuid.Parse(strings.TrimSpace(serviceID))
		if err != nil {
			return entities.CheckoutSession{}, 0, "", fmt.Errorf("%w: every service_id must be a UUID", ErrInvalidAppointmentCheckout)
		}
		serviceID = parsedID.String()
		if _, exists := seen[serviceID]; exists {
			return entities.CheckoutSession{}, 0, "", fmt.Errorf("%w: duplicate service IDs are not allowed", ErrInvalidAppointmentCheckout)
		}
		seen[serviceID] = struct{}{}
		canonicalIDs = append(canonicalIDs, serviceID)
	}
	sort.Strings(canonicalIDs)
	serviceIDsValue := strings.Join(canonicalIDs, ",")

	previous, previousErr := s.DB.GetTransactionByAppointmentID(appointmentID)
	if previousErr != nil && !errors.Is(previousErr, gorm.ErrRecordNotFound) {
		return entities.CheckoutSession{}, 0, "", fmt.Errorf("load appointment payment: %w", previousErr)
	}
	if previousErr == nil {
		switch previous.PaymentStatus {
		case "paid", "refunded":
			return entities.CheckoutSession{}, 0, "", ErrAppointmentCheckoutExists
		case "pending":
			if previous.Meta.Get("service_ids") != serviceIDsValue {
				return entities.CheckoutSession{}, 0, "", ErrAppointmentCheckoutExists
			}
			if previous.CheckoutURL != "" {
				if err := s.publishAppointmentPaymentPending(appointmentID, previous.StripeSessionID, previous.AmountTotal, previous.Currency); err != nil {
					return entities.CheckoutSession{}, 0, "", fmt.Errorf("publish pending payment status: %w", err)
				}
				return entities.CheckoutSession{ID: previous.StripeSessionID, URL: previous.CheckoutURL}, previous.AmountTotal, previous.Currency, nil
			}
		}
	}

	linePriceIDs := make([]string, 0, len(canonicalIDs))
	var amountTotal int64
	currency := ""
	for _, serviceID := range canonicalIDs {
		product, err := s.DB.GetProductByBarberServiceID(serviceID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return entities.CheckoutSession{}, 0, "", fmt.Errorf("%w: %s", ErrAppointmentServiceUnavailable, serviceID)
			}
			return entities.CheckoutSession{}, 0, "", fmt.Errorf("load trusted Barber service %s: %w", serviceID, err)
		}
		if product.GatewayPriceExternalID == "" || product.Price <= 0 || product.DeletedAt != "" {
			return entities.CheckoutSession{}, 0, "", fmt.Errorf("%w: %s", ErrAppointmentServiceUnavailable, serviceID)
		}
		serviceCurrency := strings.ToLower(strings.TrimSpace(product.Currency))
		if serviceCurrency == "" {
			return entities.CheckoutSession{}, 0, "", fmt.Errorf("%w: %s has no currency", ErrAppointmentServiceUnavailable, serviceID)
		}
		if currency == "" {
			currency = serviceCurrency
		} else if currency != serviceCurrency {
			return entities.CheckoutSession{}, 0, "", fmt.Errorf("%w: appointment services must use a single currency", ErrInvalidAppointmentCheckout)
		}
		if amountTotal > math.MaxInt64-product.Price {
			return entities.CheckoutSession{}, 0, "", errors.New("appointment total exceeds supported amount")
		}
		amountTotal += product.Price
		linePriceIDs = append(linePriceIDs, product.GatewayPriceExternalID)
	}

	metadata := meta.New()
	metadata.Add(barberPaymentContextKey, barberPaymentContext)
	metadata.Add(appointmentIDMetadata, appointmentID)
	attemptID := uuid.NewString()
	metadata.Add("checkout_attempt_id", attemptID)
	session, err := s.GatewayUsecase.CreateAppointmentCheckoutSession(linePriceIDs, appointmentID, metadata)
	if err != nil {
		return entities.CheckoutSession{}, 0, "", fmt.Errorf("create Stripe Checkout Session: %w", err)
	}

	paymentMeta := meta.New()
	paymentMeta.Add("service_ids", serviceIDsValue)
	payment := entities.Transactions{
		Type:              "barber_appointment",
		PaymentStatus:     "pending",
		AppointmentID:     appointmentID,
		StripeSessionID:   session.ID,
		CheckoutAttemptID: attemptID,
		CheckoutURL:       session.URL,
		AmountTotal:       amountTotal,
		Currency:          currency,
		NetAmount:         amountTotal,
		Meta:              paymentMeta,
	}
	if err := s.DB.CreateTransaction(&payment); err != nil {
		return entities.CheckoutSession{}, 0, "", fmt.Errorf("save pending appointment payment: %w", err)
	}
	if err := s.publishAppointmentPaymentPending(appointmentID, session.ID, amountTotal, currency); err != nil {
		return entities.CheckoutSession{}, 0, "", fmt.Errorf("publish pending payment status: %w", err)
	}
	return session, amountTotal, currency, nil
}

func (s *StripeUsecase) publishAppointmentPaymentPending(appointmentID, sessionID string, amountTotal int64, currency string) error {
	if s.Pub == nil {
		return nil
	}
	return s.Pub.Publish(barberEventsTopic, "appointment_payment_status_update", map[string]interface{}{
		"appointment_id": appointmentID,
		"payment_status": "pending",
		"session_id":     sessionID,
		"amount_total":   amountTotal,
		"currency":       currency,
	})
}
