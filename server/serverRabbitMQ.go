package server

import (
	"finances-api/db"
	"finances-api/handlers"
	"log"
	"os"

	messageWorker "github.com/moronimotta/message-worker-module"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQServer struct {
	db            db.Database
	connectionUrl string
	queueName     string
	topicName     string
}

func NewRabbitMQServer(db db.Database) *RabbitMQServer {

	return &RabbitMQServer{
		db:            db,
		connectionUrl: os.Getenv("RABBITMQ_URL"),
		queueName:     os.Getenv("RABBITMQ_QUEUE_NAME"),
		topicName:     os.Getenv("RABBITMQ_TOPIC_NAME"),
	}
}
func (s *RabbitMQServer) Start() {
	// Setup repositories and handler
	rabbitMqHandler := handlers.NewRabbitMqHandler(s.db, "stripe")

	// CONSUMER
	var consumerInput messageWorker.Consumer
	consumerInput.ConnectionURL = os.Getenv("RABBITMQ_URL")
	consumerInput.QueueName = s.queueName
	consumerInput.TopicName = s.topicName

	messageWorker.Listen(consumerInput,
		func(msg amqp.Delivery) {
			var event messageWorker.Event
			err := event.Unmarshal(msg.Body)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				return
			}
			err = rabbitMqHandler.EventBus(event)
			if err != nil {
				log.Printf("EventBus error: %v", err)
			}
		},
	)

}
