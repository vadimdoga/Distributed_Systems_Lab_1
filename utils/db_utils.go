package utils

import (
	"fmt"
	"log"
	"time"

	dtb "github.com/vadimdoga/Distributed_Systems_Lab_1/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CheckStatus ...
func CheckStatus(productID primitive.ObjectID, status string) bool {
	var product dtb.Products
	dtb.ProductCollection.FindOne(dtb.Ctx, bson.M{"_id": productID}).Decode(&product)
	if product.Status == status {
		return true
	}

	return false
}

// UpdateStatusDelivered ...
func UpdateStatusDelivered(objID primitive.ObjectID) {
	time.Sleep(10 * time.Second)

	filter := bson.M{"_id": bson.M{"$eq": objID}}

	update := bson.M{
		"$set": bson.M{
			"status": "delivered",
		},
	}

	_, err := dtb.ProductCollection.UpdateOne(
		dtb.Ctx,
		filter,
		update,
	)

	if err != nil {
		fmt.Println(err)
	}
}

// PriorityCountDocuments ...
func PriorityCountDocuments() (int64, int64) {
	highPriority, highErr := dtb.ProductCollection.CountDocuments(dtb.Ctx, bson.M{
		"$and": []bson.M{
			{"priority": "high"},
		}, "$or": []bson.M{
			{"status": "building"},
			{"status": "delivering"},
		},
	})
	if highErr != nil {
		log.Fatal(highErr)
	}

	lowPriority, lowErr := dtb.ProductCollection.CountDocuments(dtb.Ctx, bson.M{
		"$and": []bson.M{
			{"priority": "low"},
		}, "$or": []bson.M{
			{"status": "building"},
			{"status": "delivering"},
		},
	})
	if lowErr != nil {
		log.Fatal(lowErr)
	}

	return highPriority, lowPriority
}
