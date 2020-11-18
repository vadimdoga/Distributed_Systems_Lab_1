package tools

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func RabbitMQConnect() (*amqp.Connection, *amqp.Channel, amqp.Queue) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := declareQueue(ch)

	return conn, ch, q
}

func declareQueue(ch *amqp.Channel) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		"ORDER_CREATED", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	return q, err
}

func QueuePublish(ch *amqp.Channel, q amqp.Queue, jsonBody string) {
	body := []byte(jsonBody)
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		failOnError(err, "Failed to publish a message")
	}
}

func QueueReceive(ch *amqp.Channel, q amqp.Queue, recvChannel chan []byte) {
	for {
		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		failOnError(err, "Failed to register a consumer")
		for d := range msgs {
			recvChannel <- d.Body
		}
	}
}
