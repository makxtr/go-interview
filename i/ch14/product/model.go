package product

// Product represents a product in the catalog
type Product struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Inventory int32   `json:"inventory"`
}

// DTOs

type GetProductRequest struct {
	ProductId int64 `json:"product_id"`
}

type GetProductResponse struct {
	Product *Product `json:"product"`
}

type CheckInventoryRequest struct {
	ProductId int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

type CheckInventoryResponse struct {
	Available bool `json:"available"`
}
