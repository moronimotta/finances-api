package handlers

import (
	"errors"
	usescases "finances-api/usecases"
)

type HttpHandler struct {
	usescases.FinancialUsecase
}

func NewHttpHandler(gatewayName string) (*HttpHandler, error) {
	var usecaseInput usescases.FinancialUsecase

	switch gatewayName {
	case "stripe":
		usecaseInput = *usescases.NewStripeUsecase()
	default:
		return nil, errors.New("unsupported payment gateway")
	}
	return &HttpHandler{
		usecaseInput,
	}, nil
}

func (h *HttpHandler) EventBus(payload []byte, signature string) error {
	switch h.Repository.(type) {
	case *usescases.StripeUsecase:
		return h.EventBus(payload, signature)
	default:
		return errors.New("unsupported payment gateway")
	}
}
