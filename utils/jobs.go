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
)

// TimeoutTasks ...
func TimeoutTasks() {
	time.Sleep(1 * time.Minute)

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
