package user

import "context"

// Service interface
type Service interface {
	GetUser(ctx context.Context, userID int64) (*User, error)
	ValidateUser(ctx context.Context, userID int64) (bool, error)
}
