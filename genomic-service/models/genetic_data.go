package models

type GeneticData struct {
	FileID        string `json:"file_id"`
	UserAddress   string `json:"user_address"`
	DataHash      []byte `json:"data_hash"`
	EncryptedData []byte `json:"encrypted_data"`
	IsConfirmed   bool   `json:"is_confirmed"`
}
