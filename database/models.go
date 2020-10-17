package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Products ...
type Products struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title          string             `json:"title" bson:"title"`
	StorageAddress string             `json:"storage_address" bson:"storageAddress"`
	Quantity       int                `json:"quantity" bson:"quantity"`
	CreatedAt      time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updatedAt"`
	Status         string             `json:"status" bson:"status"`
	Priority       string             `json:"priority" bson:"priority"`
}
