package utils

import (
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
