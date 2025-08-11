package handlers

import (
	"os"

	messageWorker "github.com/moronimotta/message-worker-module"
)

type RabbitPublisher struct{}

func NewRabbitPublisher() *RabbitPublisher {
	return &RabbitPublisher{}
}

func (p *RabbitPublisher) Publish(topicName, eventName string, data map[string]interface{}) error {
	input := messageWorker.Publisher{}
	input.ConnectionURL = os.Getenv("RABBITMQ_URL")
	input.TopicName = topicName

	msg := messageWorker.Event{
		Event: eventName,
		Data:  data,
	}
	messageWorker.SendMessage(input, msg)
	return nil
}
