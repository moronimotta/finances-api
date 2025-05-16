package handlers

import (
	"errors"
	stripeRepository "finances-api/repositories/stripe"
	usescases "finances-api/usecases"
)

type HttpHandler struct {
	Repo usescases.FinancialUsecase
}

func NewHttpHandler(gatewayName, gatewayAccessKey string) (*HttpHandler, error) {
	var usecaseInput usescases.FinancialUsecase

	switch gatewayName {
	case "stripe":
		stripeRepo := stripeRepository.NewStripeRepository(gatewayAccessKey)
		usecaseInput = *usescases.NewFinancialUsecase(stripeRepo)
	default:
		return nil, errors.New("unsupported payment gateway")
	}
	return &HttpHandler{
		Repo: usecaseInput,
	}, nil
}

// func NewUserRabbitMQHandler(usecaseInput usescases.UserEventUsecase) *UserRabbitMQHandler {
// 	return &UserRabbitMQHandler{
// 		Repo: usecaseInput,
// 	}
// }
