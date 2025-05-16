package usecases

import (
	"errors"
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

func (f *FinancialUsecase) EventBus(event string) error {
	switch event {
	case "customer.created":
		// Handle customer creation
	default:
		return errors.New("event not found")
	}
	return nil
}
