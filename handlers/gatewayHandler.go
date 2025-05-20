package handlers

import (
	"errors"
	usescases "finances-api/usecases/gateway"
)

type GatewayHttpHandler struct {
	usescases.GatewayUsecase
}

func NewGatewayHttpHandler(gatewayName string) (*GatewayHttpHandler, error) {
	var usecaseInput usescases.GatewayUsecase

	switch gatewayName {
	case "stripe":
		usecaseInput = *usescases.NewStripeUsecase()
	default:
		return nil, errors.New("unsupported payment gateway")
	}
	return &GatewayHttpHandler{
		usecaseInput,
	}, nil
}

func (h *GatewayHttpHandler) EventBus(payload []byte, signature string) error {
	switch h.Repository.(type) {
	case *usescases.StripeUsecase:
		return h.EventBus(payload, signature)
	default:
		return errors.New("unsupported payment gateway")
	}
}
