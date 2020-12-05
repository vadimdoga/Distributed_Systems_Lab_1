package main

import (
	"log"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/vadimdoga/PAD_Products_Service/routes"
	"github.com/vadimdoga/PAD_Products_Service/tools"
)

func main() {
	// start config
	tools.APIConfig()

	// start timeout check
	go tools.TimeoutTasks()

	// connect to rabbitmq
	tools.RabbitMQConnect()

	runtime.GOMAXPROCS(2)

	//start events
	go tools.ReceiveOrder()
	go tools.ReceiveCompensateProducts()

	// handle routes
	handleRequests(tools.SERVICE_ADDRESS)
}

func handleRequests(serviceAddress string) {
	router := mux.NewRouter()
	storedProductsRouter := router.PathPrefix("/api/products").Subrouter()

	storedProductsRouter.HandleFunc("/{id}", routes.GetProducts).Methods("GET")
	storedProductsRouter.HandleFunc("", routes.AddProducts).Methods("POST")
	storedProductsRouter.HandleFunc("/{id}", routes.UpdateProducts).Methods("PUT")
	storedProductsRouter.HandleFunc("/deliver/{id}", routes.DeliverProducts).Methods("PUT")

	storedProductsRouter.HandleFunc("", routes.GetCountStatus).Methods("GET")

	log.Println("Starting server on", serviceAddress)
	log.Fatal(http.ListenAndServe(serviceAddress, router))
	return
}
