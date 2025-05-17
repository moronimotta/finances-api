package usecases

import (
	"finances-api/repositories"
)

type FinancialUsecase struct {
	Repository repositories.FinancialRepository
}

func NewFinancialUsecase(repository repositories.FinancialRepository) *FinancialUsecase {
	return &FinancialUsecase{
		Repository: repository,
	}
}

func (f *FinancialUsecase) CreateProduct(name, description string) (string, error) {
	productID, err := f.Repository.CreateProduct(name, description)
	if err != nil {
		return "", err
	}
	return productID, nil
}

func (f *FinancialUsecase) CreatePrice(productID string, unitAmount int64, currency string) (string, error) {
	priceID, err := f.Repository.CreatePrice(productID, unitAmount, currency)
	if err != nil {
		return "", err
	}
	return priceID, nil
}
