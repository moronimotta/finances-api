package usecases

import (
	"finances-api/repositories"
)

type GatewayUsecase struct {
	Repository repositories.GatewayRepository
}

func NewGatewayUsecase(repository repositories.GatewayRepository) *GatewayUsecase {
	return &GatewayUsecase{
		Repository: repository,
	}
}

func (f *GatewayUsecase) CreateProduct(name, description string) (string, error) {
	productID, err := f.Repository.CreateProduct(name, description)
	if err != nil {
		return "", err
	}
	return productID, nil
}

func (f *GatewayUsecase) CreatePrice(productID string, unitAmount int64, currency string) (string, error) {
	priceID, err := f.Repository.CreatePrice(productID, unitAmount, currency)
	if err != nil {
		return "", err
	}
	return priceID, nil
}

func (f *GatewayUsecase) CreateCheckoutSession(priceID, customerID, successURL, cancelURL string) (string, error) {
	sessionID, err := f.Repository.CreateCheckoutSession(priceID, customerID, successURL, cancelURL)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}
