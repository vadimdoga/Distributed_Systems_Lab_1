package tools

import (
	"os"

	"github.com/streadway/amqp"
	"github.com/vadimdoga/PAD_Products_Service/utils"
)

var MQChannel *amqp.Channel

func RabbitMQConnect() {
	port := os.Getenv("MQ_PORT")
	address := os.Getenv("MQ_ADDRESS")
	conn, err := amqp.Dial("amqp://guest:guest@" + address + ":" + port)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")

	MQChannel, err = conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
}

func declareQueue(eventName string) amqp.Queue {
	q, err := MQChannel.QueueDeclare(
		eventName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	return q
}

func QueuePublish(queueName string, body []byte) {
	declaredQueue := declareQueue(queueName)

	err := MQChannel.Publish(
		"",                 // exchange
		declaredQueue.Name, // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	utils.FailOnError(err, "Failed to publish a message")
}

func QueueReceive(queueName string, recvChannel chan []byte) {
	declaredQueue := declareQueue(queueName)
	for {
		msgs, err := MQChannel.Consume(
			declaredQueue.Name, // queue
			"",                 // consumer
			true,               // auto-ack
			false,              // exclusive
			false,              // no-local
			false,              // no-wait
			nil,                // args
		)
		utils.FailOnError(err, "Failed to register a consumer")

		for d := range msgs {
			d.Ack(false)
			recvChannel <- d.Body
		}
	}
}
