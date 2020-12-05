package tools

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/vadimdoga/PAD_Products_Service/utils"
)

func ReceiveOrder() {
	orderRecvChannel := make(chan []byte)

	go QueueReceive("ORDER_CREATED", orderRecvChannel)

	for {
		bytesBody := <-orderRecvChannel
		var jsonBody utils.EventReceive
		err := json.Unmarshal([]byte(bytesBody), &jsonBody)
		FailOnJsonError(err, "Casting error", bytesBody)

		receiveMsg := fmt.Sprintf("Received transaction: %s", jsonBody.TransactionID)
		fmt.Println(receiveMsg)

		totalPrice, newJsonProducts := ProcessOrderEvent(jsonBody)

		var newJsonBody utils.EventPublish

		newJsonBody.Products = newJsonProducts
		newJsonBody.TotalPrice = totalPrice
		newJsonBody.TransactionID = jsonBody.TransactionID
		newJsonBody.UserID = jsonBody.UserID

		body, err := json.Marshal(newJsonBody)
		utils.FailOnError(err, "Value Error")

		QueuePublish("PRODUCTS_CHECKING", body)
	}
}

func PublishCompensateOrder(oldJson []byte, error_msg string) {
	newJsonBody := utils.EventCompensate{
		ErrorMsg:       error_msg,
		ProductDetails: string(oldJson),
	}

	body, err := json.Marshal(newJsonBody)
	utils.FailOnError(err, "Cast error")

	QueuePublish("COMPENSATION_ORDER_CREATED", body)
}

func ReceiveCompensateProducts() {
	cmpProductsRecvChannel := make(chan []byte)

	go QueueReceive("COMPENSATION_PRODUCTS_CHECKING", cmpProductsRecvChannel)

	for {
		bytesBody := <-cmpProductsRecvChannel

		var jsonBody utils.EventPublish

		err := json.Unmarshal([]byte(bytesBody), &jsonBody)
		utils.FailOnError(err, "Casting error")

		CompensateProducts(jsonBody)

	}
}

func FailOnJsonError(err error, msg string, recvBody []byte) {
	if err != nil {
		PublishCompensateOrder(recvBody, msg)
		log.Println(fmt.Sprintf("%s : %s", msg, err.Error()))
	}
}
