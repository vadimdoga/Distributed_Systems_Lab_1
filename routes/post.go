package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/vadimdoga/PAD_Products_Service/db"
	"github.com/vadimdoga/PAD_Products_Service/tools"
	"github.com/vadimdoga/PAD_Products_Service/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddProducts ...
func AddProducts(w http.ResponseWriter, r *http.Request) {

	highCheckLimit, lowCheckLimit := tools.CheckPostLimit()

	reqBody, _ := ioutil.ReadAll(r.Body)

	var products db.Products
	timestamp := time.Now()

	json.Unmarshal(reqBody, &products)

	if products.Priority == "high" {
		if highCheckLimit == false {
			utils.JSONError(w, "Service Unavailable. High priority limit can't be exceeded!", 503)
			return
		}
	} else if products.Priority == "low" {
		if lowCheckLimit == false {
			utils.JSONError(w, "Service Unavailable. Low priority limit can't be exceeded!", 503)
			return
		}
	}

	products.CreatedAt = timestamp
	products.Status = "building"

	res, err := db.ProductCollection.InsertOne(tools.Ctx, products)
	if err != nil {
		utils.JSONError(w, err, 500)
		return
	}
	productID := res.InsertedID.(primitive.ObjectID)

	utils.JSONResponse(w, "success", productID.Hex(), products.Status, 201)
	return
}
