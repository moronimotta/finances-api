package handlers

import (
	"errors"
	usescases "finances-api/usecases"
)

type WebhookHandler struct {
	usescases.PaymentAPIUsecases
	gatewayName string
}

func NewWebhookHandler(gatewayName string, usecases *usescases.PaymentAPIUsecases) (*WebhookHandler, error) {
	return &WebhookHandler{
		PaymentAPIUsecases: *usecases,
		gatewayName:        gatewayName,
	}, nil
}

func (h *WebhookHandler) EventBus(payload []byte, signature string) error {
	switch h.gatewayName {
	case "stripe":
		su := usescases.NewStripeUsecase(&h.PaymentAPIUsecases)
		return su.EventBus(payload, signature)
	default:
		return errors.New("unsupported payment gateway")
	}
}
