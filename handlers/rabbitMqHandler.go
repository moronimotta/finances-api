package handlers

import (
	"os"

	messageWorker "github.com/moronimotta/message-worker-module"
)

type RabbitMqHandler struct {
	// Repository repositories.FinancialRepository
	// Usecase     usecases.DbUsecase
}

// For now, we are not using the repository in the handler
// but in the future, if a consumer is needed, set the handler
// as DbHttpHandler type.
// Never use gateway handler in the repository.
// Gateway requests should be made by user by http requests. Never as messages btw services.
func NewRabbitMqHandler() *RabbitMqHandler {
	return &RabbitMqHandler{
		// Repository: repository,
	}
}

func (h *RabbitMqHandler) SendMessage(eventName, topicName string, eventData interface{}) {
	event := messageWorker.Event{
		Event: eventName,
		Data:  eventData,
	}

	pbInput := messageWorker.Publisher{
		ConnectionURL: os.Getenv("RABBITMQ_URL"),
		TopicName:     topicName,
	}

	messageWorker.SendMessage(pbInput, event)
}
