package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	dtb "github.com/vadimdoga/Distributed_Systems_Lab_1/database"
	"github.com/vadimdoga/Distributed_Systems_Lab_1/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateProducts ...
func UpdateProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	reqBody, _ := ioutil.ReadAll(r.Body)

	var products dtb.Products
	timestamp := time.Now()

	json.Unmarshal(reqBody, &products)
	storageAddress := &products.StorageAddress

	objID, idErr := primitive.ObjectIDFromHex(productID)
	if idErr != nil {
		utils.JSONError(w, idErr, 404)
		return
	}

	checkStatusField := utils.CheckStatus(objID, "building")

	if checkStatusField == false {
		utils.JSONError(w, "This obj is not in building stage", 400)
		return
	}

	filter := bson.M{"_id": bson.M{"$eq": objID}}

	update := bson.M{
		"$set": bson.M{
			"updatedAt":      timestamp,
			"storageAddress": storageAddress,
		},
	}

	_, err := dtb.ProductCollection.UpdateOne(
		dtb.Ctx,
		filter,
		update,
	)

	if err != nil {
		utils.JSONError(w, err, 500)
		return
	}

	utils.JSONResponse(w, "successfully updated!", productID, "building", 200)
	return
}

// DeliverProducts ...
func DeliverProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	objID, idErr := primitive.ObjectIDFromHex(productID)
	if idErr != nil {
		utils.JSONError(w, idErr, 404)
		return
	}

	checkStatusField := utils.CheckStatus(objID, "building")
	if checkStatusField == false {
		utils.JSONError(w, "This obj is not in building stage", 400)
		return
	}

	filter := bson.M{"_id": bson.M{"$eq": objID}}

	update := bson.M{
		"$set": bson.M{
			"status": "delivering",
		},
	}

	_, err := dtb.ProductCollection.UpdateOne(
		dtb.Ctx,
		filter,
		update,
	)

	if err != nil {
		utils.JSONError(w, err, 500)
		return
	}

	go utils.UpdateStatusDelivered(objID)

	utils.JSONResponse(w, "successfully finalized!", productID, "delivering", 200)
	return
}
