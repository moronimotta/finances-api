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
	params := &stripe.ChargeParams{}
	params.AddExpand("balance_transaction")
	ch, err := charge.Get(chargeID, params)
	if err != nil {
		return output, err
	}

	if ch.PaymentMethodDetails != nil {
		output.PaymentMethod = string(ch.PaymentMethodDetails.Type)
	}
	output.Currency = string(ch.Currency)
	output.AmountTotal = ch.Amount
	output.AmountPayed = ch.AmountCaptured
	output.AmountRefunded = ch.AmountRefunded
	if ch.BalanceTransaction != nil {
		output.StripeFee = ch.BalanceTransaction.Fee
		output.StripeFeeKnown = true
	}
	output.NetAmount = output.AmountTotal - output.AmountRefunded - output.StripeFee
	if ch.Refunded {
		output.Status = "refunded"
	}
	output.Meta = ch.Metadata
	output.ReceiptURL = ch.ReceiptURL
	return output, nil
}
