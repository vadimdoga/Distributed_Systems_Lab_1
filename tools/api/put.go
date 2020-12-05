package api

import (
	"fmt"
	"time"

	"github.com/vadimdoga/PAD_Products_Service/db"
	"github.com/vadimdoga/PAD_Products_Service/tools"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CheckStatus ...
func CheckStatus(productID primitive.ObjectID, status string) bool {
	var product db.Products
	db.ProductCollection.FindOne(tools.Ctx, bson.M{"_id": productID}).Decode(&product)
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

	_, err := db.ProductCollection.UpdateOne(
		tools.Ctx,
		filter,
		update,
	)

	if err != nil {
		fmt.Println(err)
	}
}
