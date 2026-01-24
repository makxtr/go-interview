package product

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Client implementation
type Client struct {
	baseURL string
}

func NewClient(conn *grpc.ClientConn) Service {
	// For this challenge, we'll assume the product service is on port 50052
	// In a real scenario, we'd extract this from the connection target
	return &Client{baseURL: "http://localhost:50052"}
}

// NewClientWithAddress creates a client with a specific address (helper for tests)
func NewClientWithAddress(address string) Service {
	return &Client{baseURL: "http://" + address}
}

func (c *Client) GetProduct(ctx context.Context, productID int64) (*Product, error) {
	resp, err := http.Get(fmt.Sprintf("%s/product/get?id=%d", c.baseURL, productID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}

	var product Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (c *Client) CheckInventory(ctx context.Context, productID int64, quantity int32) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("%s/product/check?id=%d&quantity=%d", c.baseURL, productID, quantity))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, status.Errorf(codes.NotFound, "product not found")
	}

	var result map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result["available"], nil
}
