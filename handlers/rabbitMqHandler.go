package handlers

import (
	"errors"
	usecases "finances-api/usecases"
	"log"
	"os"

	messageWorker "github.com/moronimotta/message-worker-module"
	"github.com/redis/go-redis/v9"
)

type RabbitMqHandler struct {
	usecases    *usecases.PaymentAPIUsecases
	redisClient *redis.Client
}

func NewRabbitMqHandler(usecases *usecases.PaymentAPIUsecases, redisClient *redis.Client) *RabbitMqHandler {

	return &RabbitMqHandler{
		usecases:    usecases,
		redisClient: redisClient,
	}
}

func (h *RabbitMqHandler) PublishMessage(topicName, eventName string, data map[string]string) error {

	input := messageWorker.Publisher{}
	input.ConnectionURL = os.Getenv("RABBITMQ_URL")
	input.TopicName = topicName

	messageInput := messageWorker.Event{
		Event: eventName,
		Data:  data,
	}

	messageWorker.SendMessage(input, messageInput)
	return nil
}
func (h *RabbitMqHandler) EventBus(event messageWorker.Event) error {

	switch event.Event {
	case "user.created":
		dataMap, ok := event.Data.(map[string]interface{})
		if !ok {
			return errors.New("event data is not a map[string]interface{}")
		}
		name := dataMap["name"].(string)
		email := dataMap["email"].(string)
		userID := dataMap["user_id"].(string)

		externalID, err := h.usecases.Gateway.CreateCustomer(name, email, userID)
		if err != nil {
			log.Printf("Error creating customer: %v", err)
		}
		if externalID == "" {
			return errors.New("external ID is empty")
		}

		// send message back
		h.PublishMessage(
			"user-api",
			"user.updated",
			map[string]string{
				"id":          userID,
				"external_id": externalID,
			},
		)
	case "user.updated":
		dataMap, ok := event.Data.(map[string]interface{})
		if !ok {
			return errors.New("event data is not a map[string]interface{}")
		}
		var name, email string
		if v, ok := dataMap["name"].(string); ok {
			name = v
		}
		if v, ok := dataMap["email"].(string); ok {
			email = v
		}
		externalID := dataMap["external_id"].(string)
		if err := h.usecases.Gateway.UpdateCustomer(externalID, name, email); err != nil {
			log.Printf("Error updating customer: %v", err)
		}
	default:
		return errors.New("event not found")
	}
	return nil
}
