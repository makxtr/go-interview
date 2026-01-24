package order

// Order represents an order in the system
type Order struct {
	ID        int64   `json:"id"`
	UserID    int64   `json:"user_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int32   `json:"quantity"`
	Total     float64 `json:"total"`
}
