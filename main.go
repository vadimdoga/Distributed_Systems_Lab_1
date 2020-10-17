package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	dtb "github.com/vadimdoga/Distributed_Systems_Lab_1/database"
	"github.com/vadimdoga/Distributed_Systems_Lab_1/routes"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	var db *mongo.Database

	// start the database
	db = dtb.EstablishConnection()

	// start models
	dtb.ProductsCollection(db)

	// handle routes
	handleRequests()
}

func handleRequests() {
	addr := ":5000"
	router := mux.NewRouter()

	storedProductsRouter := router.PathPrefix("/products").Subrouter()

	storedProductsRouter.HandleFunc("/{id}", routes.GetProducts).Methods("GET")
	storedProductsRouter.HandleFunc("", routes.AddProducts).Methods("POST")
	storedProductsRouter.HandleFunc("/{id}", routes.UpdateProducts).Methods("PUT")
	storedProductsRouter.HandleFunc("/deliver/{id}", routes.DeliverProducts).Methods("PUT")

	log.Println("Starting server on", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
