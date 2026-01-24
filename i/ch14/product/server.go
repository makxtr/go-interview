package product

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"ch14/common"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the Service
type Server struct {
	products map[int64]*Product
}

// NewServer creates a new Server
func NewServer() *Server {
	products := map[int64]*Product{
		1: {ID: 1, Name: "Laptop", Price: 999.99, Inventory: 10},
		2: {ID: 2, Name: "Phone", Price: 499.99, Inventory: 20},
		3: {ID: 3, Name: "Headphones", Price: 99.99, Inventory: 0},
	}
	return &Server{products: products}
}

// GetProduct retrieves a product by ID
func (s *Server) GetProduct(ctx context.Context, productID int64) (*Product, error) {
	product, exists := s.products[productID]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}

	return product, nil
}

// CheckInventory checks if a product is available in the requested quantity
func (s *Server) CheckInventory(ctx context.Context, productID int64, quantity int32) (bool, error) {
	product, exists := s.products[productID]
	if !exists {
		return false, status.Errorf(codes.NotFound, "product not found")
	}
	return product.Inventory >= quantity, nil
}

// StartService starts the product service on the given port
func StartService(port string) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(common.LoggingInterceptor))
	productServer := NewServer()

	// Hint: create listener, gRPC server with interceptor, register service, serve
	mux := http.NewServeMux()
	mux.HandleFunc("/product/get", func(w http.ResponseWriter, r *http.Request) {
		productIDStr := r.URL.Query().Get("id")
		productID, _ := strconv.ParseInt(productIDStr, 10, 64)

		product, err := productServer.GetProduct(r.Context(), productID)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	})

	mux.HandleFunc("/product/check", func(w http.ResponseWriter, r *http.Request) {
		productIDStr := r.URL.Query().Get("id")
		productID, _ := strconv.ParseInt(productIDStr, 10, 64)
		quantityStr := r.URL.Query().Get("quantity")
		quantity, _ := strconv.ParseInt(quantityStr, 10, 32)

		available, err := productServer.CheckInventory(r.Context(), productID, int32(quantity))
		if err != nil {
			if status.Code(err) == codes.NotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"available": available})
	})

	go func() {
		log.Printf("Product service HTTP server listening on %s", port)
		if err := http.Serve(lis, mux); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	return s, nil
}
