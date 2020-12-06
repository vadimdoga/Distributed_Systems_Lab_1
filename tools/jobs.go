package tools

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/vadimdoga/PAD_Products_Service/db"
	"github.com/vadimdoga/PAD_Products_Service/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TimeoutTasks ...
func TimeoutTasks() {
	timeoutEnv, err := strconv.ParseInt(os.Getenv("TIMEOUT"), 10, 32)
	if err != nil {
		log.Fatal(err)
	}

	for {
		var products []db.Products

		timeout := time.Duration(timeoutEnv)
		time.Sleep(timeout * time.Second)

		response, err := db.ProductCollection.Find(Ctx, bson.M{"status": "building"})
		utils.FailOnError(err, "No Products with status: building")

		response.All(Ctx, &products)

		for _, pr := range products {
			currentTime := time.Now()
			diff := currentTime.Sub(pr.CreatedAt)
			if diff.Seconds() >= 1 {
				_, err := db.ProductCollection.DeleteOne(Ctx, bson.M{"_id": pr.ID})

				utils.FailOnError(err, "Delete unsuccesfull")
				if err == nil {
					log.Printf("Task %s timeout reached!. Succesfull delete!", pr.ID.String())
				}
			}
		}
	}
}

func priorityCountDocuments() (int64, int64) {
	highPriority, highErr := db.ProductCollection.CountDocuments(Ctx, bson.M{
		"$and": []bson.M{
			{"priority": "high"},
		}, "$or": []bson.M{
			{"status": "building"},
			{"status": "delivering"},
		},
	})
	utils.FailOnError(highErr, "High priority error")

	lowPriority, lowErr := db.ProductCollection.CountDocuments(Ctx, bson.M{
		"$and": []bson.M{
			{"priority": "low"},
		}, "$or": []bson.M{
			{"status": "building"},
			{"status": "delivering"},
		},
	})
	utils.FailOnError(lowErr, "Low priority error")

	return highPriority, lowPriority
}

// CheckPostLimit ...
func CheckPostLimit() (bool, bool) {
	var high bool = false
	var low bool = false

	highLimitEnv := os.Getenv("HIGH_LIMIT")
	lowLimitEnv := os.Getenv("LOW_LIMIT")

	highLimit, _ := strconv.ParseInt(highLimitEnv, 10, 64)
	lowLimit, _ := strconv.ParseInt(lowLimitEnv, 10, 64)

	highCountedLimit, lowCountedLimit := priorityCountDocuments()

	if highCountedLimit < highLimit {
		high = true
	}

	if lowCountedLimit < lowLimit {
		low = true
	}

	return high, low
}

func ProcessOrderEvent(jsonBody utils.EventReceive) (float64, []utils.ProductsPublishList) {
	var totalPrice float64 = 0
	var newJsonProducts []utils.ProductsPublishList

	for _, productJson := range jsonBody.Products {
		var productDB db.Products
		var newProductItem utils.ProductsPublishList

		err := db.ProductCollection.FindOne(Ctx, bson.M{"title": productJson.ProductTitle}).Decode(&productDB)

		FailOnJsonError(err, "No product with such title", jsonBody)
		if err != nil {
			return 0, nil
		}

		if productDB.Quantity < int(productJson.Amount) {
			errMsg := fmt.Errorf("Not enough products in stock! In stock are %d product(s)", int(productDB.Quantity))
			FailOnJsonError(errMsg, "Not enough quantity", jsonBody)
			return 0, nil
		}

		newQuantity := productDB.Quantity - int(productJson.Amount)

		_, err = db.ProductCollection.UpdateOne(
			Ctx,
			productDB,
			bson.M{
				"$set": bson.M{
					"quantity": newQuantity,
				},
			},
		)

		FailOnJsonError(err, "Update error", jsonBody)
		if err != nil {
			return 0, nil
		}

		newProductItem.Amount = productJson.Amount
		newProductItem.ProductTitle = productJson.ProductTitle
		newProductItem.ProductID = productDB.ID.Hex()

		newJsonProducts = append(newJsonProducts, newProductItem)

		totalPrice += productDB.Price
	}

	return totalPrice, newJsonProducts
}

func CompensateProducts(jsonBody utils.EventPublish) {
	for _, productJson := range jsonBody.Products {
		var productDB db.Products

		objID, err := primitive.ObjectIDFromHex(productJson.ProductID)
		utils.FailOnError(err, "Cast to objID")

		err = db.ProductCollection.FindOne(Ctx, bson.M{"_id": objID}).Decode(&productDB)
		utils.FailOnError(err, "Product Not Found")

		newQuantity := productDB.Quantity + int(productJson.Amount)

		_, err = db.ProductCollection.UpdateOne(
			Ctx,
			productDB,
			bson.M{
				"$set": bson.M{
					"quantity": newQuantity,
				},
			},
		)

		if err != nil {
			log.Println(fmt.Sprintf("Successful compensate on transaction: %s", jsonBody.TransactionID))
		}
	}
}
