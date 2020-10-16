package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// AddStoredProducts ...
func AddStoredProducts(w http.ResponseWriter, r *http.Request) {
	//*return id of the process
	reqBody, _ := ioutil.ReadAll(r.Body)

	var products Products
	timestamp := time.Now()

	json.Unmarshal(reqBody, &products)

	prd := &products
	prd.CreatedAt = timestamp

	_, err := productsCollection.InsertOne(ctx, products)
	if err != nil {
		JSONError(w, err, 500)
	}

	JSONResponse(w, "success", 201)
	return
}
