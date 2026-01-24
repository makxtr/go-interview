package user

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
	// Extract address from connection for HTTP calls
	// In a real gRPC implementation, this would use the connection directly
	return &Client{baseURL: "http://localhost:50051"}
}

// NewClientWithAddress creates a client with a specific address (helper for tests)
func NewClientWithAddress(address string) Service {
	return &Client{baseURL: "http://" + address}
}

func (c *Client) GetUser(ctx context.Context, userID int64) (*User, error) {
	resp, err := http.Get(fmt.Sprintf("%s/user/get?id=%d", c.baseURL, userID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Client) ValidateUser(ctx context.Context, userID int64) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("%s/user/validate?id=%d", c.baseURL, userID))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, status.Errorf(codes.NotFound, "user not found")
	}

	var result map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result["valid"], nil
}
