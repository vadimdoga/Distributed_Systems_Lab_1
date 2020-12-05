package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/vadimdoga/PAD_Products_Service/db"
	"github.com/vadimdoga/PAD_Products_Service/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Ctx context.Context
var ConnDB *mongo.Database
var SERVICE_ADDRESS string

func development() {
	utils.ReadConfigFile("./config.ini")
}

func production() {

}

func gatewayConnection(serviceAddress string) {
	var resp *http.Response
	var connected bool = false
	gatewayAddress := os.Getenv("GATEWAY_ADDR")

	requestBody, _ := json.Marshal(map[string]string{
		"address":     serviceAddress,
		"serviceType": "products",
	})

	for i := 1; i <= 3; i++ {
		resp, _ = http.Post(gatewayAddress+"/register", "application/json", bytes.NewBuffer(requestBody))
		timeDelay := math.Pow(2, float64(i))
		if resp.StatusCode != 200 {
			time.Sleep(time.Duration(timeDelay) * time.Second)
			connected = false
			log.Println("Gateway connection failed. Retrying...")
		} else {
			connected = true
			break
		}
	}

	if connected {
		log.Println("Succesfull connection to gateway!")
	} else {
		log.Println("Couldn't connect to gateway!")
	}

	defer resp.Body.Close()
}

// EstablishConnection ...
func establishConnection() *mongo.Database {
	mongodbHost := os.Getenv("MONGODB_HOST")
	mongodbPort := os.Getenv("MONGODB_PORT")
	mongodbDB := os.Getenv("MONGODB_DB")

	// Database Config
	uri := fmt.Sprintf("mongodb://%s:%s", mongodbHost, mongodbPort)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	//Set up a context required by mongo.Connect
	Ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(Ctx)
	//Cancel context to avoid memory leak
	defer cancel()

	// Ping our db connection
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connection to DB established!")
	}
	// Connect to the database
	db := client.Database(mongodbDB)

	return db
}

func APIConfig() {
	appEnv := os.Getenv("APP_ENVIRONMENT")

	if appEnv == "development" {
		development()
	} else if appEnv == "production" {
		production()
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	SERVICE_ADDRESS = host + ":" + port
	// start the database
	ConnDB = establishConnection()

	// start models
	db.ProductsCollection(ConnDB)

	// connect to gateway
	// gatewayConnection(SERVICE_ADDRESS)
}
