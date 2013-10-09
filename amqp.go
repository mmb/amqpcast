package amqpcast

import (
	"log"

	"github.com/streadway/amqp"
)

func handleDeliveries(deliveries <-chan amqp.Delivery, cstr *Caster) {
	for d := range deliveries {
		log.Printf("received AMQP message (%dB): %q", len(d.Body), d.Body)
		cstr.Outbound <- string(d.Body[:])
	}
}

func InitAmqp(amqpUrl *string, amqpExchange *string, amqpKey *string, cstr *Caster) {
	amqpConn, err := amqp.Dial(*amqpUrl)
	if err != nil {
		log.Fatal(err)
	}

	amqpChannel, err := amqpConn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	queue, err := amqpChannel.QueueDeclare(
		"",    // name of the queue
		false, // durable
		true,  // autodelete
		true,  // exclusive
		false, // nowait
		nil,   // args
	)
	if err != nil {
		log.Fatal(err)
	}

	amqpChannel.QueueBind(
		queue.Name,
		*amqpKey,
		*amqpExchange,
		false, // nowait
		nil,   // args
	)

	deliveries, err := amqpChannel.Consume(
		queue.Name,
		"",    // consumer
		false, // autoack
		false, // exclusive
		false, // nolocal
		false, // nowait
		nil,   // args
	)

	go handleDeliveries(deliveries, cstr)
}
