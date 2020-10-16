package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

// HandleRequests ...
func HandleRequests() {
	addr := ":5000"
	router := mux.NewRouter()

	storedProductsRouter := router.PathPrefix("/products").Subrouter()

	storedProductsRouter.HandleFunc("/{id}", GetStoredProducts).Methods("GET")
	storedProductsRouter.HandleFunc("", AddStoredProducts).Methods("POST")
	storedProductsRouter.HandleFunc("/{id}", UpdateStoredProducts).Methods("PUT")
	storedProductsRouter.HandleFunc("/finalize/{id}", FinalizeStoredProducts).Methods("GET")

	log.Println("Starting server on", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
