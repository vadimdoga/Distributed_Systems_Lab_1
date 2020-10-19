package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	dtb "github.com/vadimdoga/Distributed_Systems_Lab_1/database"
	"go.mongodb.org/mongo-driver/bson"
)

// TimeoutTasks ...
func TimeoutTasks() {
	timeoutEnv := os.Getenv("TIMEOUT")
	for {
		var products []dtb.Products

		timeout, err := time.ParseDuration(timeoutEnv)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(timeout)

		response, err := dtb.ProductCollection.Find(dtb.Ctx, bson.M{"status": "delivering"})
		if err != nil {
			log.Fatal(err)
		}

		response.All(dtb.Ctx, &products)

		for _, pr := range products {
			currentTime := time.Now()
			diff := currentTime.Sub(pr.CreatedAt)
			if diff.Seconds() >= 1 {
				_, err := dtb.ProductCollection.DeleteOne(dtb.Ctx, bson.M{"_id": pr.ID})
				if err != nil {
					log.Fatal(err)
				} else {
					log.Printf("Task %s timeout reached!. Succesfull delete!", pr.ID.String())
				}
			}
		}
	}
}

// CheckPostLimit ...
func CheckPostLimit() (bool, bool) {
	var high bool = false
	var low bool = false

	highLimitEnv := os.Getenv("HIGH_LIMIT")
	lowLimitEnv := os.Getenv("LOW_LIMIT")

	highLimit, _ := strconv.ParseInt(highLimitEnv, 10, 64)
	lowLimit, _ := strconv.ParseInt(lowLimitEnv, 10, 64)

	highCountedLimit, lowCountedLimit := PriorityCountDocuments()

	if highCountedLimit < highLimit {
		high = true
	}

	if lowCountedLimit < lowLimit {
		low = true
	}

	return high, low
}

// GatewayConnection ...
func GatewayConnection(serviceAddress string) string {
	gatewayAddress := os.Getenv("GATEWAY_ADDR")

	requestBody, err := json.Marshal(map[string]string{
		"address": serviceAddress,
		"serviceType": "products",
	})

	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(gatewayAddress+"/register", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	return resp.Status
}
