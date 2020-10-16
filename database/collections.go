package database

import "go.mongodb.org/mongo-driver/mongo"

// ProductCollection ...
var ProductCollection *mongo.Collection

// ProductsCollection ...
func ProductsCollection(c *mongo.Database) {
	ProductCollection = c.Collection("Products")
	return
}
