package main

import (
	"github.com/petrixs/servicebus"
	"log"
)

func main() {
	//Create RebbitMQ
	rabbit, err := ServiceBus.NewRabbitMQClient(
		"amqp://guest:guest@localhost:5672/",
		"test1",
		"q2")

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	text := &ServiceBus.TestMessage{
		Key:  "key",
		Body: "body",
	}

	err = rabbit.Send(text)

	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}
	log.Println("ok")

}
