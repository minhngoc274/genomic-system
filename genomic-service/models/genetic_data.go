package models

type GeneticData struct {
	FileID        string
	UserAddress   string
	DataHash      []byte
	EncryptedData []byte
	IsConfirmed   bool
}
