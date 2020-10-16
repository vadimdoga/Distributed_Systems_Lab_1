package database

import "time"

// Products ...
type Products struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	StorageAddress string    `json:"storage_address"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Status         string    `json:"status"`
}
