package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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
	for {
		var products []dtb.Products

		time.Sleep(1 * time.Second)

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
func CheckPostLimit() bool {
	limitEnv := os.Getenv("LIMIT")

	limit, _ := strconv.ParseInt(limitEnv, 10, 64)

	countRes := CountDocuments()

	if countRes < limit {
		return true
	}

	return false
}

// GatewayConnection ...
func GatewayConnection(serviceAddress string) string {
	gatewayAddress := os.Getenv("GATEWAY_ADDR")

	requestBody, err := json.Marshal(map[string]string{
		"address": serviceAddress,
	})

	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(gatewayAddress+"/", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}
