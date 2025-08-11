package server

import (
	"finances-api/db"
	"finances-api/handlers"
	"finances-api/usecases"
	"log"
	"log/slog"
	"os"

	messageWorker "github.com/moronimotta/message-worker-module"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type RabbitMQServer struct {
	db            db.Database
	usecases      *usecases.PaymentAPIUsecases
	connectionUrl string
	queueName     string
	topicName     string
	redisClient   *redis.Client
}

func NewRabbitMQServer(db db.Database, redisClient *redis.Client) *RabbitMQServer {

	usecases := usecases.NewPaymentAPIUsecases("stripe", db)
	// Inject a RabbitMQ publisher for outbound events
	var publisher handlers.RabbitPublisher
	usecases.Pub = &publisher

	return &RabbitMQServer{
		db:            db,
		usecases:      usecases,
		connectionUrl: os.Getenv("RABBITMQ_URL"),
		queueName:     os.Getenv("RABBITMQ_QUEUE_NAME"),
		topicName:     os.Getenv("RABBITMQ_TOPIC_NAME"),
		redisClient:   redisClient,
	}
}
func (s *RabbitMQServer) Start() {
	// Setup repositories and handler
	rabbitMqHandler := handlers.NewRabbitMqHandler(s.usecases, s.redisClient)

	// CONSUMER
	var consumerInput messageWorker.Consumer
	consumerInput.ConnectionURL = os.Getenv("RABBITMQ_URL")
	consumerInput.QueueName = s.queueName
	consumerInput.TopicName = s.topicName

	log.Println("Starting RabbitMQ consumer...")

	messageWorker.Listen(consumerInput,
		func(msg amqp.Delivery) {
			var event messageWorker.Event
			err := event.Unmarshal(msg.Body)
			if err != nil {
				slog.Error("Failed to unmarshal message", err)
				return
			}
			err = rabbitMqHandler.EventBus(event)
			if err != nil {
				slog.Error("Error processing event", err)
			}
		},
	)

}
