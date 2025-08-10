package stripeRepository

import (
	"finances-api/entities"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/charge"
)

type StripeCharge interface {
	GetCharge(chargeID string) (entities.Transactions, error)
}

type stripeChargeRepository struct {
	stripeKey string
}

func NewStripeChargeRepository(key string) StripeCharge {
	return &stripeChargeRepository{
		stripeKey: key,
	}
}
func (r *stripeChargeRepository) GetCharge(chargeID string) (entities.Transactions, error) {
	var output entities.Transactions
	stripe.Key = r.stripeKey
	ch, err := charge.Get(chargeID, nil)
	if err != nil {
		return output, err
	}

	output.PaymentMethod = string(ch.PaymentMethodDetails.Type)
	output.Currency = string(ch.Currency)
	output.AmountRefunded = ch.AmountRefunded
	if ch.Refunded {
		output.Status = "refunded"
	}
	output.ReceiptURL = ch.ReceiptURL
	return output, nil
}
