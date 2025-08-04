package stripeRepository

import (
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/price"
)

type StripePrice interface {
	CreatePrice(productID string, unitAmount int64, currency string) (string, error)
	ChangePrice(oldPriceID, productID string, unitAmount int64, currency string) (string, error)
}

type stripePriceRepository struct {
	stripeKey string
}

func NewStripePriceRepository(key string) StripePrice {
	return &stripePriceRepository{
		stripeKey: key,
	}
}

func (r *stripePriceRepository) CreatePrice(productID string, unitAmount int64, currency string) (string, error) {
	stripe.Key = r.stripeKey

	price_params := &stripe.PriceParams{
		Currency:   stripe.String(string(stripe.CurrencyUSD)),
		Product:    stripe.String(productID),
		UnitAmount: stripe.Int64(unitAmount),
	}
	starter_price, _ := price.New(price_params)
	if starter_price == nil {
		return "", nil
	}

	return starter_price.ID, nil
}

func (r *stripePriceRepository) UpdatePrice(priceID string, unitAmount int64, currency string) error {
	stripe.Key = r.stripeKey

	params := &stripe.PriceParams{
		UnitAmount: stripe.Int64(unitAmount),
		Currency:   stripe.String(currency),
	}
	_, err := price.Update(priceID, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *stripePriceRepository) ChangePrice(oldPriceID, productID string, unitAmount int64, currency string) (string, error) {
	stripe.Key = r.stripeKey
	params := &stripe.PriceParams{
		Active: stripe.Bool(false),
	}
	_, err := price.Update(oldPriceID, params)
	if err != nil {
		return "", err
	}

	new_price, err := r.CreatePrice(productID, unitAmount, currency)
	if err != nil {
		return "", err
	}

	return new_price, nil
}
