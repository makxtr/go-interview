package main

import (
	"context"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestUserService(t *testing.T) {
	// Start the user service
	server, err := StartUserService(":50051")
	if err != nil {
		t.Fatalf("Failed to start user service: %v", err)
	}
	defer server.Stop()

	// Wait for the server to start
	time.Sleep(200 * time.Millisecond)

	// Connect to the service
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to user service: %v", err)
	}
	defer conn.Close()

	// Create a client from the connection
	userClient := NewUserServiceClient(conn)

	t.Run("GetUser", func(t *testing.T) {
		// Test getting an existing user
		user, err := userClient.GetUser(context.Background(), 1)
		if err != nil {
			t.Errorf("GetUser failed: %v", err)
		} else {
			if user.ID != 1 || user.Username != "alice" {
				t.Errorf("Expected user with ID 1 and username 'alice', got ID %d and username '%s'", user.ID, user.Username)
			}
		}

		// Test getting a non-existent user
		_, err = userClient.GetUser(context.Background(), 999)
		if err == nil {
			t.Errorf("Expected error for non-existent user, got nil")
		} else if status.Code(err) != codes.NotFound {
			t.Errorf("Expected NotFound error, got %v", err)
		}
	})

	t.Run("ValidateUser", func(t *testing.T) {
		// Test validating an active user
		valid, err := userClient.ValidateUser(context.Background(), 1)
		if err != nil {
			t.Errorf("ValidateUser failed: %v", err)
		} else if !valid {
			t.Errorf("Expected user 1 to be valid")
		}

		// Test validating an inactive user
		valid, err = userClient.ValidateUser(context.Background(), 3)
		if err != nil {
			t.Errorf("ValidateUser failed: %v", err)
		} else if valid {
			t.Errorf("Expected user 3 to be invalid")
		}

		// Test validating a non-existent user
		_, err = userClient.ValidateUser(context.Background(), 999)
		if err == nil {
			t.Errorf("Expected error for non-existent user, got nil")
		} else if status.Code(err) != codes.NotFound {
			t.Errorf("Expected NotFound error, got %v", err)
		}
	})
}

func TestProductService(t *testing.T) {
	// Start the product service
	server, err := StartProductService(":50052")
	if err != nil {
		t.Fatalf("Failed to start product service: %v", err)
	}
	defer server.Stop()

	// Wait for the server to start
	time.Sleep(200 * time.Millisecond)

	// Connect to the service
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to product service: %v", err)
	}
	defer conn.Close()

	// Create a client from the connection
	productClient := NewProductServiceClient(conn)

	t.Run("GetProduct", func(t *testing.T) {
		// Test getting an existing product
		product, err := productClient.GetProduct(context.Background(), 1)
		if err != nil {
			t.Errorf("GetProduct failed: %v", err)
		} else {
			if product.ID != 1 || product.Name != "Laptop" {
				t.Errorf("Expected product with ID 1 and name 'Laptop', got ID %d and name '%s'", product.ID, product.Name)
			}
		}

		// Test getting a non-existent product
		_, err = productClient.GetProduct(context.Background(), 999)
		if err == nil {
			t.Errorf("Expected error for non-existent product, got nil")
		} else if status.Code(err) != codes.NotFound {
			t.Errorf("Expected NotFound error, got %v", err)
		}
	})

	t.Run("CheckInventory", func(t *testing.T) {
		// Test checking inventory for a product with sufficient inventory
		available, err := productClient.CheckInventory(context.Background(), 1, 5)
		if err != nil {
			t.Errorf("CheckInventory failed: %v", err)
		} else if !available {
			t.Errorf("Expected product 1 to be available in quantity 5")
		}

		// Test checking inventory for a product with insufficient inventory
		available, err = productClient.CheckInventory(context.Background(), 1, 15)
		if err != nil {
			t.Errorf("CheckInventory failed: %v", err)
		} else if available {
			t.Errorf("Expected product 1 to be unavailable in quantity 15")
		}

		// Test checking inventory for a product with zero inventory
		available, err = productClient.CheckInventory(context.Background(), 3, 1)
		if err != nil {
			t.Errorf("CheckInventory failed: %v", err)
		} else if available {
			t.Errorf("Expected product 3 to be unavailable")
		}

		// Test checking inventory for a non-existent product
		_, err = productClient.CheckInventory(context.Background(), 999, 1)
		if err == nil {
			t.Errorf("Expected error for non-existent product, got nil")
		} else if status.Code(err) != codes.NotFound {
			t.Errorf("Expected NotFound error, got %v", err)
		}
	})
}

func TestOrderService(t *testing.T) {
	// Start both services
	var wg sync.WaitGroup
	wg.Add(2)

	var userServer, productServer *grpc.Server

	go func() {
		defer wg.Done()
		server, err := StartUserService(":50053")
		if err != nil {
			t.Errorf("Failed to start user service: %v", err)
			return
		}
		userServer = server
		// Keep service running for the duration of the test
		time.Sleep(3 * time.Second)
	}()

	go func() {
		defer wg.Done()
		server, err := StartProductService(":50054")
		if err != nil {
			t.Errorf("Failed to start product service: %v", err)
			return
		}
		productServer = server
		// Keep service running for the duration of the test
		time.Sleep(3 * time.Second)
	}()

	// Wait for services to start
	time.Sleep(300 * time.Millisecond)

	// Connect to both services
	orderService, err := ConnectToServices("localhost:50053", "localhost:50054")
	if err != nil {
		t.Fatalf("Failed to connect to services: %v", err)
	}

	t.Run("CreateOrder_Success", func(t *testing.T) {
		// Test creating an order with valid user and product
		order, err := orderService.CreateOrder(context.Background(), 1, 1, 2)
		if err != nil {
			t.Errorf("CreateOrder failed: %v", err)
		} else {
			if order.UserID != 1 || order.ProductID != 1 || order.Quantity != 2 {
				t.Errorf("Expected order with UserID 1, ProductID 1, and Quantity 2, got UserID %d, ProductID %d, and Quantity %d",
					order.UserID, order.ProductID, order.Quantity)
			}
			if order.Total != 1999.98 { // 2 * 999.99
				t.Errorf("Expected total 1999.98, got %f", order.Total)
			}
		}
	})

	t.Run("CreateOrder_InvalidUser", func(t *testing.T) {
		// Test creating an order with an inactive user
		_, err := orderService.CreateOrder(context.Background(), 3, 1, 2)
		if err == nil {
			t.Errorf("Expected error for inactive user, got nil")
		}
	})

	t.Run("CreateOrder_NonExistentUser", func(t *testing.T) {
		// Test creating an order with a non-existent user
		_, err := orderService.CreateOrder(context.Background(), 999, 1, 2)
		if err == nil {
			t.Errorf("Expected error for non-existent user, got nil")
		}
	})

	t.Run("CreateOrder_InsufficientInventory", func(t *testing.T) {
		// Test creating an order with insufficient inventory
		_, err := orderService.CreateOrder(context.Background(), 1, 1, 15)
		if err == nil {
			t.Errorf("Expected error for insufficient inventory, got nil")
		}
	})

	t.Run("CreateOrder_NonExistentProduct", func(t *testing.T) {
		// Test creating an order with a non-existent product
		_, err := orderService.CreateOrder(context.Background(), 1, 999, 2)
		if err == nil {
			t.Errorf("Expected error for non-existent product, got nil")
		}
	})

	// Clean up servers
	if userServer != nil {
		userServer.Stop()
	}
	if productServer != nil {
		productServer.Stop()
	}

	// Wait for all services to complete
	wg.Wait()
}

func TestInterceptors(t *testing.T) {
	t.Run("LoggingInterceptor", func(t *testing.T) {
		// Test that logging interceptor can be called without error
		ctx := context.Background()
		req := &GetUserRequest{UserId: 1}
		info := &grpc.UnaryServerInfo{
			FullMethod: "/user.UserService/GetUser",
		}
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return &GetUserResponse{User: &User{ID: 1, Username: "test"}}, nil
		}

		resp, err := LoggingInterceptor(ctx, req, info, handler)
		if err != nil {
			t.Errorf("LoggingInterceptor failed: %v", err)
		}
		if resp == nil {
			t.Errorf("LoggingInterceptor returned nil response")
		}
	})

	t.Run("AuthInterceptor", func(t *testing.T) {
		// Test that auth interceptor can be called without error
		ctx := context.Background()
		method := "/user.UserService/GetUser"
		req := &GetUserRequest{UserId: 1}
		reply := &GetUserResponse{}
		invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			return nil
		}

		err := AuthInterceptor(ctx, method, req, reply, nil, invoker)
		if err != nil {
			t.Errorf("AuthInterceptor failed: %v", err)
		}
	})
}
