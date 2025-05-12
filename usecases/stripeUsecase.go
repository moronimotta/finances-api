package usecases

import stripeRepository "finances-api/repositories/stripe"

type StripeUsecase struct {
	Repository stripeRepository.StripeRepository
}

func NewStripeUsecase(repository stripeRepository.StripeRepository) *StripeUsecase {
	return &StripeUsecase{
		Repository: repository,
	}
}

// CreateProduct
// UpdateProduct
// CreateCheckoutSession
// CreateCustomer
// UpdateCustomer
// DeleteCustomer
