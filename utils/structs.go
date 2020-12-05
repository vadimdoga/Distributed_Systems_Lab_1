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

type productsList struct {
	ProductTitle string `json:"product_title"`
	Amount       int64  `json:"amount"`
}

type ProductsPublishList struct {
	ProductID    string `json:"product_id"`
	ProductTitle string `json:"product_title"`
	Amount       int64  `json:"amount"`
}

type EventReceive struct {
	TransactionID string         `json:"transaction_id"`
	UserID        int64          `json:"user_id"`
	Products      []productsList `json:"products"`
}

type EventPublish struct {
	TransactionID string                `json:"transaction_id"`
	UserID        int64                 `json:"user_id"`
	Products      []ProductsPublishList `json:"products"`
	TotalPrice    float64               `json:"total_price"`
}

type EventCompensate struct {
	ProductDetails string `json:"product_details"`
	ErrorMsg       string `json:"error_msg"`
}
