package tools

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/vadimdoga/Distributed_Systems_Lab_1/utils"
)

func ReceiveOrder() {
	orderRecvChannel := make(chan []byte)

	go QueueReceive("ORDER_CREATED", orderRecvChannel)

	for {
		bytesBody := <-orderRecvChannel
		jsonBody := utils.DecodeReceiver(bytesBody)
		receiveMsg := fmt.Sprintf("Received transaction: %s", jsonBody.TransactionID)
		fmt.Println(receiveMsg)

		totalPrice, newJsonProducts := ProcessOrderEvent(jsonBody)

		var newJsonBody utils.EventPublish

		newJsonBody.Products = newJsonProducts
		newJsonBody.TotalPrice = totalPrice
		newJsonBody.TransactionID = jsonBody.TransactionID
		newJsonBody.UserID = jsonBody.UserID

		body, err := json.Marshal(newJsonBody)
		utils.FailOnError(err, "Value Error", body)

		QueuePublish("PRODUCTS_CHECKING", body)
	}
}

func PublishCompensateOrder(jsonBody utils.EventCompensate) {
	body, err := json.Marshal(jsonBody)
	if err != nil {
		log.Fatalf("%s", err)
	}

	QueuePublish("COMPENSATION_ORDER_CREATED", body)
}

func ReceiveCompensateProducts() {
	cmpProductsRecvChannel := make(chan []byte)

	go QueueReceive("COMPENSATION_PRODUCTS_CHECKING", cmpProductsRecvChannel)
}

func FailOnJsonError(err error, msg string, body []byte) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)

		PublishCompensateOrder(body)
	}
}
