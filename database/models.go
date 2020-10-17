package database

import "time"

// Products ...
type Products struct {
	Title          string    `json:"title"`
	StorageAddress string    `json:"storage_address"`
	Quantity       int       `json:"quantity"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Status         string    `json:"status"`
}
