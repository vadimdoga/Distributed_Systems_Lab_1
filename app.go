package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/vadimdoga/Distributed_Systems_Lab_1/routes"
	"github.com/vadimdoga/Distributed_Systems_Lab_1/tools"
)

func main() {
	_, ch, q := tools.RabbitMQConnect()

	recvChannel := make(chan []byte)

	go tools.QueueReceive(ch, q, recvChannel)

	tools.QueuePublish(ch, q, `{
			"id" : 11,
			"name" : "Irshad",
			"department" : "IT",
			"designation" : "Product Manager"
	}`)

	// start config
	tools.APIConfig()

	// start timeout check
	go tools.TimeoutTasks()

	// handle routes
	handleRequests(tools.SERVICE_ADDRESS)
}

func handleRequests(serviceAddress string) {
	router := mux.NewRouter()
	BASE_PATH := os.Getenv("BASE_PATH")
	storedProductsRouter := router.PathPrefix(BASE_PATH).Subrouter()

	storedProductsRouter.HandleFunc("/{id}", routes.GetProducts).Methods("GET")
	storedProductsRouter.HandleFunc("", routes.AddProducts).Methods("POST")
	storedProductsRouter.HandleFunc("/{id}", routes.UpdateProducts).Methods("PUT")
	storedProductsRouter.HandleFunc("/deliver/{id}", routes.DeliverProducts).Methods("PUT")

	storedProductsRouter.HandleFunc("", routes.GetCountStatus).Methods("GET")

	log.Println("Starting server on", serviceAddress)
	log.Fatal(http.ListenAndServe(serviceAddress, router))
	return
}
