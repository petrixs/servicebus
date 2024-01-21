package ServiceBus

import (
	"github.com/streadway/amqp"
)

type RabbitMQClient struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Exchange   string
	Queue      string
	Serializer *JSONSerializer
}
