package stripeRepository

import "finances-api/repositories"

type StripeRepository struct {
	StripeProducts
	StripePrice
}

type StripeProducts interface {
	CreateProduct(name, description string) (string, error)
}

type StripePrice interface {
	CreatePrice(productID string, unitAmount int64, currency string) (string, error)
}

type StripeCheckout interface {
	CreateCheckoutSession(priceID, successURL, cancelURL string) (string, error)
}

func NewStripeRepository(key string) repositories.FinancialRepository {
	return StripeRepository{
		StripeProducts: NewStripeProductsRepository(key),
		StripePrice:    NewStripePriceRepository(key),
	}
}
