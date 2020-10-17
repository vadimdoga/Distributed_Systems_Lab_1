package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	dtb "github.com/vadimdoga/Distributed_Systems_Lab_1/database"
	"github.com/vadimdoga/Distributed_Systems_Lab_1/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetProducts ...
func GetProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	objID, idErr := primitive.ObjectIDFromHex(productID)
	if idErr != nil {
		utils.JSONError(w, idErr, 404)
		return
	}

	var product dtb.Products
	checkRes := dtb.ProductCollection.FindOne(dtb.Ctx, bson.M{"_id": objID})

	if checkRes.Err() != nil {
		utils.JSONError(w, checkRes.Err().Error(), 404)
		return
	}

	checkRes.Decode(&product)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	res, _ := json.Marshal(product)
	w.Write(res)
}

// GetCountStatus ...
func GetCountStatus(w http.ResponseWriter, r *http.Request) {
	response := utils.CountDocuments()

	var count utils.CountResponse

	count.Count = response

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	countResponse, _ := json.Marshal(count)
	w.Write(countResponse)
	return
}
