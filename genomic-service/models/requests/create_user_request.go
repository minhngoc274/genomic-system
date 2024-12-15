package requests

// CreateUserRequest ...
type CreateUserRequest struct {
	Address string `json:"address" binding:"required,eth_addr"`
}
