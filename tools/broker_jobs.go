package tools

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/vadimdoga/PAD_Products_Service/utils"
	"github.com/valyala/fastjson"
)

func ReceiveOrder() {
	orderRecvChannel := make(chan []byte)

	go QueueReceive("ORDER_CREATED", orderRecvChannel)

	for {
		bytesBody := <-orderRecvChannel
		var jsonBody utils.EventReceive
		err := json.Unmarshal(bytesBody, &jsonBody)
		FailOnCastError(err, "Casting error", bytesBody)
		if err != nil {
			continue
		}

		receiveMsg := fmt.Sprintf("Received transaction: %s", jsonBody.TransactionID)
		fmt.Println(receiveMsg)

		totalPrice, newJsonProducts := ProcessOrderEvent(jsonBody)
		if totalPrice == 0 && newJsonProducts == nil {
			continue
		}

		var newJsonBody utils.EventPublish

		newJsonBody.Products = newJsonProducts
		newJsonBody.TotalPrice = totalPrice
		newJsonBody.TransactionID = jsonBody.TransactionID
		newJsonBody.UserID = jsonBody.UserID

		body, err := json.Marshal(newJsonBody)
		FailOnCastError(err, "Value error", bytesBody)
		if err != nil {
			continue
		}

		QueuePublish("PRODUCTS_CHECKING", body)
	}
}

func PublishCompensateOrder(oldJson utils.EventReceive, error_msg string) {
	newJsonBody := utils.EventCompensate{
		TransactionID: oldJson.TransactionID,
		UserID:        oldJson.UserID,
		Products:      oldJson.Products,
		ErrorMsg:      error_msg,
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

func FailOnJsonError(err error, msg string, recvBody utils.EventReceive) {
	if err != nil {
		PublishCompensateOrder(recvBody, err.Error())
		log.Println(fmt.Sprintf("%s : %s", msg, err.Error()))
	}
}

func FailOnCastError(err error, msg string, recvBody []byte) {
	if err != nil {
		var p fastjson.Parser
		v, _ := p.Parse(string(recvBody))

		transactionID := v.GetStringBytes("transaction_id")

		newBody := utils.EventReceive{
			TransactionID: string(transactionID),
		}

		PublishCompensateOrder(newBody, err.Error())

		log.Printf("%s : %s", msg, err.Error())
	}
}
