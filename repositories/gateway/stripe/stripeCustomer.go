package stripeRepository

import (
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
)

type StripeCustomer interface {
	CreateCustomer(name, email string) (string, error)
}

type stripeCustomerRepository struct {
	stripeKey string
}

func NewStripeCustomerRepository(key string) StripeCustomer {
	return &stripeCustomerRepository{
		stripeKey: key,
	}
}

func (r *stripeCustomerRepository) CreateCustomer(name, email string) (string, error) {
	params := &stripe.CustomerParams{
		Name:  stripe.String(name),
		Email: stripe.String(email),
	}
	result, err := customer.New(params)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}
