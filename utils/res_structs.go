package utils

type msgResponse struct {
	Message string `json:"message"`
	ID      string `json:"product_id"`
	Status  string `json:"status"`
}

// CountResponse ...
type CountResponse struct {
	Count int64 `json:"counted_processes"`
}
