package messaging

// Publisher is an abstraction over a message bus (e.g., RabbitMQ)
// to avoid import cycles between usecases and handlers.
type Publisher interface {
	Publish(topicName, eventName string, data map[string]interface{}) error
}
