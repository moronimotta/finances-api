package usecases

import (
	"finances-api/entities"
	"finances-api/repositories"
	stripeRepository "finances-api/repositories/gateway/stripe"
	"finances-api/utils/meta"
)

type GatewayUsecase struct {
	Repository repositories.GatewayRepository
}

// func NewGatewayUsecase(repository repositories.GatewayRepository) *GatewayUsecase {
// 	return &GatewayUsecase{
// 		Repository: repository,
// 	}
// }

func NewGatewayUsecase(gatewayName string) *GatewayUsecase {
	var repository repositories.GatewayRepository

	switch gatewayName {
	case "stripe":
		repository = stripeRepository.NewStripeRepository()
	default:
		repository = nil
	}

	return &GatewayUsecase{
		Repository: repository,
	}
}

func (f *GatewayUsecase) CreateProduct(name, description string, localProduct entities.Products) (string, error) {
	productID, err := f.Repository.CreateProduct(name, description, localProduct)
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

func (f *GatewayUsecase) ChangePrice(oldPriceID, productID string, unitAmount int64, currency string) (string, error) {
	newPriceID, err := f.Repository.ChangePrice(oldPriceID, productID, unitAmount, currency)
	if err != nil {
		return "", err
	}
	return newPriceID, nil
}

func (f *GatewayUsecase) DeactivateProduct(productID string) error {
	err := f.Repository.DeactivateProduct(productID)
	if err != nil {
		return err
	}
	return nil
}

func (f *GatewayUsecase) UpdateProduct(productID, name, description string, meta meta.Meta) error {
	err := f.Repository.UpdateProduct(productID, name, description, meta)
	if err != nil {
		return err
	}
	return nil
}

func (f *GatewayUsecase) CreateCheckoutSession(priceID []string, customerID string, meta meta.Meta) (string, error) {
	sessionID, err := f.Repository.CreateCheckoutSession(priceID, customerID, meta)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func (f *GatewayUsecase) CreateCustomer(name, email, localUserID string) (string, error) {
	customerID, err := f.Repository.CreateCustomer(name, email, localUserID)
	if err != nil {
		return "", err
	}
	return customerID, nil
}

func (f *GatewayUsecase) UpdateCustomer(customerID, name, email string) error {
	err := f.Repository.UpdateCustomer(customerID, name, email)
	if err != nil {
		return err
	}
	return nil
}

// get charge
func (f *GatewayUsecase) GetCharge(chargeID string) (entities.Transactions, error) {
	charge, err := f.Repository.GetCharge(chargeID)
	if err != nil {
		return entities.Transactions{}, err
	}
	return charge, nil
}

func (f *GatewayUsecase) GetPrice(priceID []string) ([]entities.TransactionItem, error) {
	items, err := f.Repository.GetPrice(priceID)
	if err != nil {
		return nil, err
	}
	return items, nil
}
