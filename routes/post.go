package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	dtb "github.com/vadimdoga/Distributed_Systems_Lab_1/database"
	"github.com/vadimdoga/Distributed_Systems_Lab_1/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddProducts ...
func AddProducts(w http.ResponseWriter, r *http.Request) {

	checkLimit := utils.CheckPostLimit()
	if checkLimit == false {
		utils.JSONError(w, "Service Unavailable. Limit can't be exceeded!", 503)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)

	var products dtb.Products
	timestamp := time.Now()

	json.Unmarshal(reqBody, &products)

	prd := &products
	prd.CreatedAt = timestamp
	prd.Status = "building"

	res, err := dtb.ProductCollection.InsertOne(dtb.Ctx, products)
	if err != nil {
		utils.JSONError(w, err, 500)
		return
	}
	productID := res.InsertedID.(primitive.ObjectID).String()

	utils.JSONResponse(w, "success", productID, prd.Status, 201)
	return
}
