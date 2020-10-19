package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Ctx ...
var Ctx context.Context

// EstablishConnection ...
func EstablishConnection() *mongo.Database {
	mongodbURI := os.Getenv("MongoDbURI")
	// Database Config
	clientOptions := options.Client().ApplyURI(mongodbURI)
	client, err := mongo.NewClient(clientOptions)
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
		log.Println("Connected!")
	}
	// Connect to the database
	db := client.Database("products_service")

	return db
}
