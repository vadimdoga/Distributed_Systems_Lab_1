package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	dtb "github.com/vadimdoga/Distributed_Systems_Lab_1/database"
	"github.com/vadimdoga/Distributed_Systems_Lab_1/utils"
)

// AddStoredProducts ...
func AddStoredProducts(w http.ResponseWriter, r *http.Request) {
	//*return id of the process
	reqBody, _ := ioutil.ReadAll(r.Body)

	var products dtb.Products
	timestamp := time.Now()

	json.Unmarshal(reqBody, &products)

	prd := &products
	prd.CreatedAt = timestamp

	_, err := dtb.ProductCollection.InsertOne(dtb.Ctx, products)
	if err != nil {
		utils.JSONError(w, err, 500)
	}

	utils.JSONResponse(w, "success", 201)
	return
}
