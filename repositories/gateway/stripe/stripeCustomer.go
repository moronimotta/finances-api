package stripeRepository

import (
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
)

type StripeCustomer interface {
	CreateCustomer(name, email, localUserID string) (string, error)
	UpdateCustomer(customerID, name, email string) error
}

type stripeCustomerRepository struct {
	stripeKey string
}

func NewStripeCustomerRepository(key string) StripeCustomer {
	return &stripeCustomerRepository{
		stripeKey: key,
	}
}

func (r *stripeCustomerRepository) CreateCustomer(name, email, localUserID string) (string, error) {

	stripe.Key = r.stripeKey

	params := &stripe.CustomerParams{
		Name:  stripe.String(name),
		Email: stripe.String(email),
		Metadata: map[string]string{
			"local_user_id": localUserID,
		},
	}
	result, err := customer.New(params)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

func (r *stripeCustomerRepository) UpdateCustomer(customerID, name, email string) error {
	stripe.Key = r.stripeKey

	params := &stripe.CustomerParams{}
	if name != "" {
		params.Name = stripe.String(name)
	}
	if email != "" {
		params.Email = stripe.String(email)
	}

	_, err := customer.Update(customerID, params)
	if err != nil {
		return err
	}
	return nil
}
