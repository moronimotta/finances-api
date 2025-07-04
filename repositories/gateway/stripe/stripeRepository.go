package stripeRepository

import (
	"finances-api/repositories"
	"os"
)

type StripeRepository struct {
	StripeProducts
	StripePrice
	StripeCheckout
	StripeCustomer
}

func NewStripeRepository() repositories.GatewayRepository {
	key := os.Getenv("STRIPE_KEY")
	return StripeRepository{
		StripeProducts: NewStripeProductsRepository(key),
		StripePrice:    NewStripePriceRepository(key),
		StripeCheckout: NewStripeCheckoutRepository(key),
		StripeCustomer: NewStripeCustomerRepository(key),
	}
}
