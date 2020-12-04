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
	iniProfile, err := config.GetSection("DATABASE")
	if err != nil {
		log.Fatalln("failed to get profile DATABASE", err)
	}

	mongodbHost, _ := iniProfile.GetKey("MONGODB_HOST")
	mongodbDB, _ := iniProfile.GetKey("MONGODB_DB")
	mongodbPort, _ := iniProfile.GetKey("MONGODB_PORT")
	mongodbUsername, _ := iniProfile.GetKey("MONGODB_USERNAME")
	mongodbPassword, _ := iniProfile.GetKey("MONGODB_PASSWORD")

	os.Setenv("MONGODB_HOST", mongodbHost.String())
	os.Setenv("MONGODB_DB", mongodbDB.String())
	os.Setenv("MONGODB_PORT", mongodbPort.String())
	os.Setenv("MONGODB_USERNAME", mongodbUsername.String())
	os.Setenv("MONGODB_PASSWORD", mongodbPassword.String())
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
