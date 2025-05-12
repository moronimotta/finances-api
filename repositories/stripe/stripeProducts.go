package stripeRepository

import (
	"finances-api/db"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/product"
)

type stripeProductsRepository struct {
	db        db.Database
	stripeKey string
}

func NewStripeProductsRepository(db db.Database, key string) StripeProducts {
	return &stripeProductsRepository{
		db:        db,
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
