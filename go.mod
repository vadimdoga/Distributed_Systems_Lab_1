module github.com/vadimdoga/PAD_Products_Service

go 1.15

require (
	github.com/gorilla/mux v1.8.0
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/streadway/amqp v1.0.0
	go.mongodb.org/mongo-driver v1.4.2
	gopkg.in/ini.v1 v1.62.0 // direct
)

replace github.com/vadimdoga/Distributed_Systems_Lab_1 => /PAD_Products_Service
