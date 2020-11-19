package utils

import (
	"encoding/json"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

func ReadConfigFile(filename string) {
	config, err := ini.Load(filename)
	if err != nil {
		log.Fatalln("failed to load shared credentials file", err)
	}
	iniProfile, err := config.GetSection("ENVIRONMENT")
	if err != nil {
		log.Fatalln("failed to get profile ENVIRONMENT", err)
	}

	highLimit, _ := iniProfile.GetKey("HIGH_LIMIT")
	lowLimit, _ := iniProfile.GetKey("LOW_LIMIT")
	timeout, _ := iniProfile.GetKey("TIMEOUT")
	gatewayAddr, _ := iniProfile.GetKey("GATEWAY_ADDR")
	ip, _ := iniProfile.GetKey("IP")
	port, _ := iniProfile.GetKey("PORT")
	mongodbURI, _ := iniProfile.GetKey("MongoDbURI")
	basePath, _ := iniProfile.GetKey("BASE_PATH")

	os.Setenv("HIGH_LIMIT", highLimit.String())
	os.Setenv("LOW_LIMIT", lowLimit.String())
	os.Setenv("TIMEOUT", timeout.String())
	os.Setenv("GATEWAY_ADDR", gatewayAddr.String())
	os.Setenv("IP", ip.String())
	os.Setenv("PORT", port.String())
	os.Setenv("MongoDbURI", mongodbURI.String())
	os.Setenv("BASE_PATH", basePath.String())
}

func DecodeReceiver(body []byte) EventReceive {
	var decodedJSON EventReceive

	err := json.Unmarshal([]byte(body), &decodedJSON)
	if err != nil {
		log.Fatal(err.Error())
	}

	return decodedJSON
}

// func EncodePublisher(body ) EventPublish {
// 	var encodedJSON EventReceive

// 	err := json.Marshal(EventPublish)
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	return decodedJSON
// }
