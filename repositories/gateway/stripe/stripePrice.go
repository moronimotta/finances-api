package stripeRepository

import (
	"finances-api/entities"
	"fmt"
	"strings"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/price"
)

// TODO: Figure it out how to make customers set an specific amount to pay
// even if the subscription is for 3 months, the customer can pay for everything in the first month.

type StripePrice interface {
	CreatePrice(productID string, unitAmount int64, currency string) (string, error)
	ChangePrice(oldPriceID, productID string, unitAmount int64, currency string) (string, error)
	GetPrice(priceID []string) ([]entities.TransactionItem, error)
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
	if strings.TrimSpace(currency) == "" {
		currency = string(stripe.CurrencyUSD)
	}

	priceParams := &stripe.PriceParams{
		Currency:   stripe.String(currency),
		Product:    stripe.String(productID),
		UnitAmount: stripe.Int64(unitAmount),
	}
	starterPrice, err := price.New(priceParams)
	if err != nil {
		return "", err
	}
	if starterPrice == nil {
		return "", fmt.Errorf("Stripe returned an empty price")
	}

	return starterPrice.ID, nil
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

func (r *stripePriceRepository) GetPrice(priceID []string) ([]entities.TransactionItem, error) {
	stripe.Key = r.stripeKey

	var output []entities.TransactionItem

	for _, id := range priceID {
		price, err := price.Get(id, nil)
		if err != nil {
			return nil, err
		}

		item := entities.TransactionItem{
			ProductExternalID: price.Product.ID,
			UnitAmount:        price.UnitAmount,
			Currency:          string(price.Currency),
			Quantity:          1,
		}
		output = append(output, item)
	}

	return output, nil
}
