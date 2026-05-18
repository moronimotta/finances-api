package stripeRepository

import (
	"fmt"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/invoice"
	"github.com/stripe/stripe-go/v82/invoiceitem"
	"github.com/stripe/stripe-go/v82/price"
)

type StripeSubscription interface {
	CreateSubscription(customerID string, priceID string, rentDueTimestamp int64, autopay bool, contractMonths int64, propertyID, managerAccountID string) (string, error)
	CreateOneTimePaymentDue(customerID string, priceID string, dueTimestamp int64, propertyID, managerAccountID string) (string, error)
}

type StripeSubscriptionRepository struct {
	stripeKey string
}

func NewStripeSubscriptionRepository(key string) StripeSubscription {
	return &StripeSubscriptionRepository{
		stripeKey: key,
	}
}

func (r *StripeSubscriptionRepository) CreateSubscription(customerID string, priceID string, rentDueTimestamp int64, autopay bool, contractMonths int64, propertyID, managerAccountID string) (string, error) {
	return "", fmt.Errorf("stripe subscriptions are currently disabled")
}

func (r *StripeSubscriptionRepository) CreateOneTimePaymentDue(customerID string, priceID string, dueTimestamp int64, propertyID, managerAccountID string) (string, error) {
	stripe.Key = r.stripeKey

	priceDetails, err := price.Get(priceID, nil)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve price %s: %w", priceID, err)
	}

	invoiceItemParams := &stripe.InvoiceItemParams{
		Customer:    stripe.String(customerID),
		Amount:      stripe.Int64(priceDetails.UnitAmount),
		Currency:    stripe.String(string(priceDetails.Currency)),
		Description: stripe.String("Rent payment due"),
	}
	invoiceItemParams.AddMetadata("property_id", propertyID)
	invoiceItemParams.AddMetadata("property_manager_account", managerAccountID)
	invoiceItemParams.AddMetadata("billing_type", "one_time_due")
	invoiceItemParams.AddMetadata("original_price_id", priceID)

	_, err = invoiceitem.New(invoiceItemParams)
	if err != nil {
		return "", fmt.Errorf("failed to create invoice item: %w", err)
	}

	invoiceParams := &stripe.InvoiceParams{
		Customer:         stripe.String(customerID),
		CollectionMethod: stripe.String("send_invoice"),
		AutoAdvance:      stripe.Bool(true),
		AutomaticTax:     &stripe.InvoiceAutomaticTaxParams{Enabled: stripe.Bool(true)},
	}
	if dueTimestamp > 0 {
		invoiceParams.DueDate = stripe.Int64(dueTimestamp)
	} else {
		invoiceParams.DaysUntilDue = stripe.Int64(1)
	}

	invoiceParams.AddMetadata("property_id", propertyID)
	invoiceParams.AddMetadata("property_manager_account", managerAccountID)
	invoiceParams.AddMetadata("billing_type", "one_time_due")

	inv, err := invoice.New(invoiceParams)
	if err != nil {
		return "", fmt.Errorf("failed to create invoice: %w", err)
	}

	finalizedInvoice, err := invoice.FinalizeInvoice(inv.ID, nil)
	if err != nil {
		return "", fmt.Errorf("failed to finalize invoice: %w", err)
	}

	return finalizedInvoice.ID, nil
}
