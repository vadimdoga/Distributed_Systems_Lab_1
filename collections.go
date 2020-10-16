package main

import "go.mongodb.org/mongo-driver/mongo"

var productsCollection *mongo.Collection

// ProductsCollection ...
func ProductsCollection(c *mongo.Database) {
	productsCollection = c.Collection("Products")
	return
}
