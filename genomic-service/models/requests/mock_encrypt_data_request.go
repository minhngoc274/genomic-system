package requests

// MockEncryptDataRequest ...
type MockEncryptDataRequest struct {
	PublicKey string `json:"public_key" binding:"required"`
	Data      string `json:"data" binding:"required"`
}
