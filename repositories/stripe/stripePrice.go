package stripeRepository

import (
	"finances-api/db"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/price"
)

type stripePriceRepository struct {
	db        db.Database
	stripeKey string
}

func NewStripePriceRepository(db db.Database, key string) StripePrice {
	return &stripePriceRepository{
		db:        db,
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
