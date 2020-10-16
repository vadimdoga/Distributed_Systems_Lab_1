package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSONError ...
func JSONError(w http.ResponseWriter, err interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err)
}

// JSONResponse ...
func JSONResponse(w http.ResponseWriter, msg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	res := msgResponse{Message: msg}
	response, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Write(response)
}
