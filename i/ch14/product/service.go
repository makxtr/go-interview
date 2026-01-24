package product

import "context"

// Service interface
type Service interface {
	GetProduct(ctx context.Context, productID int64) (*Product, error)
	CheckInventory(ctx context.Context, productID int64, quantity int32) (bool, error)
}
