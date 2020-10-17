package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

	// connect to gateway
	// gatewayConnection(serviceAddress)
}

func handleRequests() string {
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	addr := ip + ":" + port
	router := mux.NewRouter()

	storedProductsRouter := router.PathPrefix("/products").Subrouter()

	storedProductsRouter.HandleFunc("/{id}", routes.GetProducts).Methods("GET")
	storedProductsRouter.HandleFunc("", routes.AddProducts).Methods("POST")
	storedProductsRouter.HandleFunc("/{id}", routes.UpdateProducts).Methods("PUT")
	storedProductsRouter.HandleFunc("/deliver/{id}", routes.DeliverProducts).Methods("PUT")

	storedProductsRouter.HandleFunc("", routes.GetCountStatus).Methods("GET")

	log.Println("Starting server on", addr)
	log.Fatal(http.ListenAndServe(addr, router))

	return addr
}

func gatewayConnection(serviceAddress string) string {
	gatewayAddress := os.Getenv("GATEWAY_ADDR")

	requestBody, err := json.Marshal(map[string]string{
		"address": serviceAddress,
	})

	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(gatewayAddress+"/", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}
