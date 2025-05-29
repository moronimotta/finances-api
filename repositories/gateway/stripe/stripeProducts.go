package stripeRepository

import (
	"finances-api/entities"
	"finances-api/utils/meta"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/product"
)

type StripeProducts interface {
	CreateProduct(name, description string, localProduct entities.Products) (string, error)
	UpdateProduct(productID, name, description string, meta meta.Meta) error
	DeactivateProduct(productID string) error
}

type stripeProductsRepository struct {
	stripeKey string
}

func NewStripeProductsRepository(key string) StripeProducts {
	return &stripeProductsRepository{

		stripeKey: key,
	}
}

func (r *stripeProductsRepository) CreateProduct(name, description string, localProduct entities.Products) (string, error) {
	stripe.Key = r.stripeKey

	product_params := &stripe.ProductParams{
		Name:        stripe.String(name),
		Description: stripe.String(description),
		Metadata: map[string]string{
			"local_product_id": localProduct.ID,
		},
	}
	starter_product, _ := product.New(product_params)
	if starter_product == nil {
		return "", nil
	}
	return starter_product.ID, nil
}

func (r *stripeProductsRepository) UpdateProduct(productID, name, description string, meta meta.Meta) error {
	stripe.Key = r.stripeKey

	params := &stripe.ProductParams{}

	if name == "" && description == "" && meta == nil {
		return nil
	}

	if name != "" {
		params.Name = stripe.String(name)
	}
	if description != "" {
		params.Description = stripe.String(description)
	}

	if meta != nil {
		params.Metadata = meta
	}

	_, err := product.Update(productID, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *stripeProductsRepository) DeactivateProduct(productID string) error {
	stripe.Key = r.stripeKey

	params := &stripe.ProductParams{
		Active: stripe.Bool(false),
	}

	_, err := product.Update(productID, params)
	if err != nil {
		return err
	}
	return nil
}
