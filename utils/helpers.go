package utils

import (
	"fmt"
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

	os.Setenv("MONGODB_HOST", mongodbHost.String())
	os.Setenv("MONGODB_DB", mongodbDB.String())
	os.Setenv("MONGODB_PORT", mongodbPort.String())

}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Println(fmt.Sprintf("%s : %s", msg, err.Error()))
	}
}

func SuccessOrError(err error, successMsg string, errorMsg string) {
	if err != nil {
		log.Fatal(fmt.Sprintf("%s : %s", errorMsg, err.Error()))
	} else {
		log.Println(successMsg)
	}
}
