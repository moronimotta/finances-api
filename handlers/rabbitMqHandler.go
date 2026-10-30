package handlers

import (
	"errors"
	"finances-api/entities"
	usecases "finances-api/usecases"
	"finances-api/utils/meta"
	"fmt"
	"log"
	"log/slog"
	"math"
	"os"
	"strconv"
	"strings"

	messageWorker "github.com/moronimotta/message-worker-module"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
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
	case "service.created":
		dataMap, ok := event.Data.(map[string]interface{})
		if !ok {
			return errors.New("event data is not a map[string]interface{}")
		}
		return h.createServicePaymentProduct(dataMap)
	case "user.created":
		dataMap, ok := event.Data.(map[string]interface{})
		if !ok {
			slog.Error("Event data type assertion failed for event user.created")
			return errors.New("event data is not a map[string]interface{}")
		}
		name := dataMap["name"].(string)
		email := dataMap["email"].(string)
		userID := dataMap["user_id"].(string)

		externalID, err := h.usecases.Gateway.CreateCustomer(name, email, userID)
		if err != nil {
			slog.Error("Error creating customer for event user.created", err)
		}
		if externalID == "" {
			slog.Error("External ID is empty for event user.created")
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
			slog.Error("Event data type assertion failed for event user.updated")
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
			slog.Error("Error updating customer for event user.updated", err)
			log.Printf("Error updating customer: %v", err)
		}
	case "course.deleted":
		dataMap, ok := event.Data.(map[string]interface{})
		if !ok {
			slog.Error("Event data type assertion failed for event course.deleted")
			return errors.New("event data is not a map[string]interface{}")
		}

		id := dataMap["product_id"].(string)
		log.Printf("Received course.deleted event for product_id: %s", id)

		product, err := h.usecases.Db.GetProductByID(id)
		if err != nil {
			slog.Error("Failed to get product for deletion", "error", err, "id", id)
			return errors.New("error getting product")
		}

		if err := h.usecases.Gateway.DeactivateProduct(product.GatewayProductExternalID); err != nil {
			slog.Error("Failed to deactivate product in gateway", "error", err)
			return errors.New("error deactivating product in gateway")
		}
		if err := h.usecases.Db.DeleteProduct(product.ID); err != nil {
			slog.Error("Failed to delete product from database", "error", err, "id", product.ID)
			return errors.New("error deleting product")
		}
		log.Printf("Product with ID %s deleted successfully", id)

	default:
		return errors.New("event not found")
	}
	return nil
}

func (h *RabbitMqHandler) createServicePaymentProduct(data map[string]interface{}) error {
	serviceID, err := eventString(data, "service_id")
	if err != nil {
		return err
	}
	name, err := eventString(data, "name")
	if err != nil {
		return err
	}
	description, _ := data["description"].(string)
	currency, err := eventString(data, "currency")
	if err != nil {
		return err
	}
	currency = strings.ToLower(currency)
	amount, err := eventAmount(data["price"])
	if err != nil {
		return err
	}

	localProduct, lookupErr := h.usecases.Db.GetProductByBarberServiceID(serviceID)
	if lookupErr != nil && !errors.Is(lookupErr, gorm.ErrRecordNotFound) {
		return lookupErr
	}
	if localProduct == nil {
		localProduct = &entities.Products{
			BarberServiceID: serviceID,
			Name:            name,
			Description:     description,
			Price:           amount,
			Currency:        currency,
			GatewayName:     "stripe",
		}
		productID, err := h.usecases.Gateway.CreateProduct(name, description, *localProduct)
		if err != nil {
			return err
		}
		if strings.TrimSpace(productID) == "" {
			return errors.New("Stripe did not return a product ID for service")
		}
		priceID, err := h.usecases.Gateway.CreatePrice(productID, amount, currency)
		if err != nil {
			return err
		}
		if strings.TrimSpace(priceID) == "" {
			return errors.New("Stripe did not return a price ID for service")
		}
		localProduct.GatewayProductExternalID = productID
		localProduct.GatewayPriceExternalID = priceID
		if err := h.usecases.Db.CreateProduct(localProduct); err != nil {
			return fmt.Errorf("persist synchronized service %s: %w", serviceID, err)
		}
	} else if localProduct.Price != amount || !strings.EqualFold(localProduct.Currency, currency) || localProduct.GatewayPriceExternalID == "" {
		priceID := localProduct.GatewayPriceExternalID
		if priceID == "" {
			priceID, err = h.usecases.Gateway.CreatePrice(localProduct.GatewayProductExternalID, amount, currency)
		} else {
			priceID, err = h.usecases.Gateway.ChangePrice(priceID, localProduct.GatewayProductExternalID, amount, currency)
		}
		if err != nil {
			return err
		}
		localProduct.Name = name
		localProduct.Description = description
		localProduct.Price = amount
		localProduct.Currency = currency
		localProduct.GatewayPriceExternalID = priceID
		localProduct.GatewayName = "stripe"
		if err := h.usecases.Db.UpdateProduct(localProduct); err != nil {
			return fmt.Errorf("update synchronized service %s: %w", serviceID, err)
		}
	}

	metadata := meta.Meta{
		"barber_service_id":        serviceID,
		"gateway_current_price_id": localProduct.GatewayPriceExternalID,
	}
	if err := h.usecases.Gateway.UpdateProduct(localProduct.GatewayProductExternalID, name, description, metadata); err != nil {
		return fmt.Errorf("update Stripe metadata for service %s: %w", serviceID, err)
	}

	return h.PublishMessage("barber-app", "service.payment_details_created", map[string]string{
		"service_id":        serviceID,
		"external_id":       localProduct.GatewayProductExternalID,
		"external_price_id": localProduct.GatewayPriceExternalID,
	})
}

func eventString(data map[string]interface{}, key string) (string, error) {
	value, ok := data[key].(string)
	if !ok || strings.TrimSpace(value) == "" {
		return "", errors.New("event field " + key + " must be a non-empty string")
	}
	return strings.TrimSpace(value), nil
}

func eventAmount(value interface{}) (int64, error) {
	var amount int64
	switch typed := value.(type) {
	case int:
		amount = int64(typed)
	case int64:
		amount = typed
	case float64:
		if math.Trunc(typed) != typed || typed > math.MaxInt64 || typed < math.MinInt64 {
			return 0, errors.New("event field price must be an integer amount")
		}
		amount = int64(typed)
	case string:
		parsed, err := strconv.ParseInt(typed, 10, 64)
		if err != nil {
			return 0, errors.New("event field price must be an integer amount")
		}
		amount = parsed
	default:
		return 0, errors.New("event field price must be an integer amount")
	}
	if amount <= 0 {
		return 0, errors.New("event field price must be greater than zero")
	}
	return amount, nil
}
