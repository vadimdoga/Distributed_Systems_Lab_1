package tools

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/vadimdoga/Distributed_Systems_Lab_1/db"
	"go.mongodb.org/mongo-driver/bson"
)

// TimeoutTasks ...
func TimeoutTasks() {
	timeoutEnv, err := strconv.ParseInt(os.Getenv("TIMEOUT"), 10, 32)
	if err != nil {
		log.Fatal(err)
	}

	for {
		var products []db.Products

		timeout := time.Duration(timeoutEnv)
		time.Sleep(timeout * time.Second)

		response, err := db.ProductCollection.Find(Ctx, bson.M{"status": "building"})
		if err != nil {
			log.Fatal(err)
		}

		response.All(Ctx, &products)

		for _, pr := range products {
			currentTime := time.Now()
			diff := currentTime.Sub(pr.CreatedAt)
			if diff.Seconds() >= 1 {
				_, err := db.ProductCollection.DeleteOne(Ctx, bson.M{"_id": pr.ID})
				if err != nil {
					log.Fatal(err)
				} else {
					log.Printf("Task %s timeout reached!. Succesfull delete!", pr.ID.String())
				}
			}
		}
	}
}

func priorityCountDocuments() (int64, int64) {
	highPriority, highErr := db.ProductCollection.CountDocuments(Ctx, bson.M{
		"$and": []bson.M{
			{"priority": "high"},
		}, "$or": []bson.M{
			{"status": "building"},
			{"status": "delivering"},
		},
	})
	if highErr != nil {
		log.Fatal(highErr)
	}

	lowPriority, lowErr := db.ProductCollection.CountDocuments(Ctx, bson.M{
		"$and": []bson.M{
			{"priority": "low"},
		}, "$or": []bson.M{
			{"status": "building"},
			{"status": "delivering"},
		},
	})
	if lowErr != nil {
		log.Fatal(lowErr)
	}

	return highPriority, lowPriority
}

// CheckPostLimit ...
func CheckPostLimit() (bool, bool) {
	var high bool = false
	var low bool = false

	highLimitEnv := os.Getenv("HIGH_LIMIT")
	lowLimitEnv := os.Getenv("LOW_LIMIT")

	highLimit, _ := strconv.ParseInt(highLimitEnv, 10, 64)
	lowLimit, _ := strconv.ParseInt(lowLimitEnv, 10, 64)

	highCountedLimit, lowCountedLimit := priorityCountDocuments()

	if highCountedLimit < highLimit {
		high = true
	}

	if lowCountedLimit < lowLimit {
		low = true
	}

	return high, low
}
