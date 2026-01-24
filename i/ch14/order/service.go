package order

import (
	"context"
	"fmt"

	"ch14/common"
	"ch14/product"
	"ch14/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// Service handles order creation
type Service struct {
	userClient    user.Service
	productClient product.Service
	orders        map[int64]*Order
	nextOrderID   int64
}

// NewService creates a new Service
func NewService(userClient user.Service, productClient product.Service) *Service {
	return &Service{
		userClient:    userClient,
		productClient: productClient,
		orders:        make(map[int64]*Order),
		nextOrderID:   1,
	}
}

// CreateOrder creates a new order
func (s *Service) CreateOrder(ctx context.Context, userID, productID int64, quantity int32) (*Order, error) {
	// 1. Validate user
	valid, err := s.userClient.ValidateUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, status.Errorf(codes.InvalidArgument, "user is not active")
	}

	// 2. Get product and check inventory
	prod, err := s.productClient.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	available, err := s.productClient.CheckInventory(ctx, productID, quantity)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, status.Errorf(codes.FailedPrecondition, "insufficient inventory")
	}

	// 3. Create order
	total := prod.Price * float64(quantity)
	order := &Order{
		ID:        s.nextOrderID,
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
		Total:     total,
	}
	s.orders[s.nextOrderID] = order
	s.nextOrderID++

	return order, nil
}

// GetOrder retrieves an order by ID
func (s *Service) GetOrder(orderID int64) (*Order, error) {
	order, exists := s.orders[orderID]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "order not found")
	}
	return order, nil
}

// ConnectToServices and return an OrderService
func ConnectToServices(userServiceAddr, productServiceAddr string) (*Service, error) {
	// Create gRPC connections with interceptors
	userConn, err := grpc.Dial(userServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(common.AuthInterceptor),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %v", err)
	}

	productConn, err := grpc.Dial(productServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(common.AuthInterceptor),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %v", err)
	}

	// Create clients
	userClient := user.NewClientWithAddress(userServiceAddr)
	productClient := product.NewClientWithAddress(productServiceAddr)

	// Note: We are leaking connections here (not closing them), but OrderService doesn't have a Close method.
	userConn.Close()
	productConn.Close()

	return NewService(userClient, productClient), nil
}
