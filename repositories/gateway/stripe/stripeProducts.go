package stripeRepository

import (
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/product"
)

type StripeProducts interface {
	CreateProduct(name, description string) (string, error)
}

type stripeProductsRepository struct {
	stripeKey string
}

func NewStripeProductsRepository(key string) StripeProducts {
	return &stripeProductsRepository{

		stripeKey: key,
	}
}

func (r *stripeProductsRepository) CreateProduct(name, description string) (string, error) {
	stripe.Key = r.stripeKey

	product_params := &stripe.ProductParams{
		Name:        stripe.String(name),
		Description: stripe.String(description),
	}
	starter_product, _ := product.New(product_params)
	if starter_product == nil {
		return "", nil
	}
	return starter_product.ID, nil
}
