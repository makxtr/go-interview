package user

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
	users map[int64]*User
}

// NewServer creates a new Server
func NewServer() *Server {
	users := map[int64]*User{
		1: {ID: 1, Username: "alice", Email: "alice@example.com", Active: true},
		2: {ID: 2, Username: "bob", Email: "bob@example.com", Active: true},
		3: {ID: 3, Username: "charlie", Email: "charlie@example.com", Active: false},
	}
	return &Server{users: users}
}

// GetUser retrieves a user by ID
func (s *Server) GetUser(ctx context.Context, userID int64) (*User, error) {
	user, exists := s.users[userID]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return user, nil
}

// ValidateUser checks if a user exists and is active
func (s *Server) ValidateUser(ctx context.Context, userID int64) (bool, error) {
	user, exists := s.users[userID]
	if !exists {
		return false, status.Errorf(codes.NotFound, "user not found")
	}
	return user.Active, nil
}

// StartService starts the user service on the given port
func StartService(port string) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(common.LoggingInterceptor))
	userServer := NewServer()

	// Register HTTP handlers for gRPC methods
	mux := http.NewServeMux()
	mux.HandleFunc("/user/get", func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("id")
		userID, _ := strconv.ParseInt(userIDStr, 10, 64)

		user, err := userServer.GetUser(r.Context(), userID)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	mux.HandleFunc("/user/validate", func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("id")
		userID, _ := strconv.ParseInt(userIDStr, 10, 64)

		valid, err := userServer.ValidateUser(r.Context(), userID)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"valid": valid})
	})

	go func() {
		log.Printf("User service HTTP server listening on %s", port)
		if err := http.Serve(lis, mux); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	return s, nil
}
