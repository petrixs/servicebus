package ServiceBus

import (
	"github.com/petrixs/servicebus/messages"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQClient struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Exchange   string
	Queue      string
	Serializer *JSONSerializer
}

func NewRabbitMQClient(amqpURL, exchange, queue string) (*RabbitMQClient, error) {

	// todo: check that exchange is not empty

	conn, err := amqp.Dial(amqpURL)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	//todo: Вынеси создание эксченжа, канала, очереди и привязки очереди к эксченжу в отдельные методы

	if exchange != "" {

		err := ch.ExchangeDeclare(
			exchange,
			"direct",
			true,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			ch.Close()
			conn.Close()
			return nil, err
		}

	}

	if queue != "" {

		_, err := ch.QueueDeclare(
			queue,
			true,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			ch.Close()
			conn.Close()
			return nil, err
		}

	}

	err = ch.QueueBind(
		queue,
		"",
		exchange,
		false,
		nil,
	)

	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	serializer := &JSONSerializer{}

	return &RabbitMQClient{
		Connection: conn,
		Channel:    ch,
		Exchange:   exchange,
		Queue:      queue,
		Serializer: serializer,
	}, nil

}

func (client *RabbitMQClient) Send(message Message) error {

	//todo: првоерь что у нас уже есть открытый канал и открытое соединение

	body, err := client.Serializer.Marshal(message)

	if err != nil {
		return err
	}

	err = client.Channel.Publish(
		client.Exchange,
		message.GetRoutingKey(),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	return err
}

func (client *RabbitMQClient) Consume(handler func(data interface{}) error) {

	msgs, err := client.Channel.Consume(
		client.Queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("Failed register consumer: %s", err)
	}

	forever := make(chan bool)

	go func() {

		for d := range msgs {

			var msg interface{}

			switch d.RoutingKey {

			case "crypto.rate":
				msg = &messages.CryptoCurrencyRate{}
			case "crypto.top":
				msg = &messages.TopCryptoCurrencies{}
			default:
				log.Printf("Unknown message type: %s", d.RoutingKey)

			}

			if err := client.Serializer.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("Failed to decode message %s", err)
			}

			if err := handler(msg); err != nil {
				log.Printf("Failed to handle message %s", err)
			}

		}

	}()

	<-forever

}
