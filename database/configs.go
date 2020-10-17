package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Ctx ...
var Ctx context.Context

// EstablishConnection ...
func EstablishConnection() *mongo.Database {
	// Database Config
	clientOptions := options.Client().ApplyURI("mongodb+srv://admin:admin@cluster0.cnciz.mongodb.net/storage_service?retryWrites=true&w=majority")
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
	db := client.Database("storage_service")

	return db
}
