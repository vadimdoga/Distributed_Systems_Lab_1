package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	dtb "github.com/vadimdoga/Distributed_Systems_Lab_1/database"
	"github.com/vadimdoga/Distributed_Systems_Lab_1/routes"
	"github.com/vadimdoga/Distributed_Systems_Lab_1/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	var db *mongo.Database
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	serviceAddress := ip + ":" + port

	// start the database
	db = dtb.EstablishConnection()

	// start models
	dtb.ProductsCollection(db)

	// connect to gateway
	resp := utils.GatewayConnection(serviceAddress)
	if len(resp) != 0 {
		log.Println("Connected to gateway")
	}

	// handle routes
	handleRequests(serviceAddress)
}

func handleRequests(serviceAddress string) {
	router := mux.NewRouter()

	storedProductsRouter := router.PathPrefix("/products").Subrouter()

	storedProductsRouter.HandleFunc("/{id}", routes.GetProducts).Methods("GET")
	storedProductsRouter.HandleFunc("", routes.AddProducts).Methods("POST")
	storedProductsRouter.HandleFunc("/{id}", routes.UpdateProducts).Methods("PUT")
	storedProductsRouter.HandleFunc("/deliver/{id}", routes.DeliverProducts).Methods("PUT")

	storedProductsRouter.HandleFunc("", routes.GetCountStatus).Methods("GET")

	log.Println("Starting server on", serviceAddress)
	log.Fatal(http.ListenAndServe(serviceAddress, router))
	return
}
