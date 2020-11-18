package tools

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"github.com/vadimdoga/Distributed_Systems_Lab_1/utils"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func RabbitMQConnect() (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return conn, ch
}

func declareQueue(ch *amqp.Channel, eventName string) amqp.Queue {
	q, err := ch.QueueDeclare(
		eventName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	return q
}

func queuePublish(ch *amqp.Channel, q amqp.Queue, jsonBody string) {
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

func queueReceive(ch *amqp.Channel, q amqp.Queue, recvChannel chan []byte) {
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

func WaitForMQ(ch *amqp.Channel) {
	recvChannel := make(chan []byte)
	orderCreatedQueue := declareQueue(ch, "ORDER_CREATED")
	productsCheckingQueue := declareQueue(ch, "PRODUCTS_CHECKING")
	go queueReceive(ch, orderCreatedQueue, recvChannel)

	for {
		bytesBody := <-recvChannel
		jsonBody := utils.DecodeJSON(bytesBody)
		fmt.Print(jsonBody["magic"])
		//todo: use jsonBody to compute smth

		queuePublish(ch, productsCheckingQueue, "I am from products service")
	}

}
