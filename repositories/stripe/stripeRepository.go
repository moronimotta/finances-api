package stripeRepository

type StripeRepository interface {
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
