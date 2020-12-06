package tools

import (
	"os"

	"github.com/streadway/amqp"
	"github.com/vadimdoga/PAD_Products_Service/utils"
)

var MQConn *amqp.Connection

func RabbitMQConnect() {
	port := os.Getenv("MQ_PORT")
	address := os.Getenv("MQ_ADDRESS")
	var err error
	MQConn, err = amqp.Dial("amqp://guest:guest@" + address + ":" + port)
	utils.SuccessOrError(err, "Successful Connection to RMQ", "Failed to connect to RabbitMQ")
}

func declareQueue(eventName string) (amqp.Queue, *amqp.Channel) {
	mqChannel, err := MQConn.Channel()
	utils.FailOnError(err, "Failed to open a channel")

	q, err := mqChannel.QueueDeclare(
		eventName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	return q, mqChannel
}

func QueuePublish(queueName string, body []byte) {
	declaredQueue, mqChannel := declareQueue(queueName)

	err := mqChannel.Publish(
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
	declaredQueue, mqChannel := declareQueue(queueName)
	for {
		msgs, err := mqChannel.Consume(
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
			recvChannel <- d.Body
		}
	}
}
