package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vadimdoga/PAD_Products_Service/db"
	"github.com/vadimdoga/PAD_Products_Service/tools"
	"github.com/vadimdoga/PAD_Products_Service/utils"
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

	var product db.Products
	checkRes := db.ProductCollection.FindOne(tools.Ctx, bson.M{"_id": objID})

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
	response, err := db.ProductCollection.CountDocuments(tools.Ctx, bson.M{"$or": []bson.M{{"status": "building"}, {"status": "delivering"}}})
	if err != nil {
		log.Fatal(err)
	}

	var count utils.CountResponse

	count.Count = response

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	countResponse, _ := json.Marshal(count)
	w.Write(countResponse)
	return
}
