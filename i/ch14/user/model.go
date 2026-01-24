package user

// User represents a user in the system
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Active   bool   `json:"active"`
}

// DTOs

type GetUserRequest struct {
	UserId int64 `json:"user_id"`
}

type GetUserResponse struct {
	User *User `json:"user"`
}

type ValidateUserRequest struct {
	UserId int64 `json:"user_id"`
}

type ValidateUserResponse struct {
	Valid bool `json:"valid"`
}
