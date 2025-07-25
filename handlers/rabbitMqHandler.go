package handlers

import (
	"errors"
	"finances-api/db"
	usecases "finances-api/usecases/db"
	usescases "finances-api/usecases/gateway"
	"log"
	"os"

	messageWorker "github.com/moronimotta/message-worker-module"
	"github.com/redis/go-redis/v9"
)

type RabbitMqHandler struct {
	DbUsecase      usecases.DbUsecase
	GatewayUsecase usescases.GatewayUsecase
	redisClient    *redis.Client
}

func NewRabbitMqHandler(db db.Database, gatewayName string, redisClient *redis.Client) *RabbitMqHandler {
	var usecaseDb usecases.DbUsecase
	var usecaseGateway usescases.GatewayUsecase

	switch db.GetDB().Dialector.Name() {
	case "postgres":
		usecaseDb = *usecases.NewPgUsecase(db)
	default:
		usecaseDb = usecases.DbUsecase{}
	}

	switch gatewayName {
	case "stripe":
		usecaseGateway = *usescases.NewStripeUsecase()
	default:
		usecaseGateway = usescases.GatewayUsecase{}
	}

	return &RabbitMqHandler{
		DbUsecase:      usecaseDb,
		GatewayUsecase: usecaseGateway,
		redisClient:    redisClient,
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

		externalID, err := h.GatewayUsecase.Repository.CreateCustomer(name, email, userID)
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
		if err := h.GatewayUsecase.Repository.UpdateCustomer(externalID, name, email); err != nil {
			log.Printf("Error updating customer: %v", err)
		}
	default:
		return errors.New("event not found")
	}
	return nil
}
