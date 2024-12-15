package requests

// UploadDataRequest ...
type UploadDataRequest struct {
	Address       string `json:"address" binding:"required,eth_addr"`
	EncryptedData []byte `json:"encrypted_data" binding:"required"`
	Signature     string `json:"signature" binding:"required"`
}
