package tee

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"strings"
)

type TEEService struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func NewTEEService() *TEEService {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	return &TEEService{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}
}

// GetPublicKey returns hex-encoded public key that users will use
func (t *TEEService) GetPublicKey() string {
	pubKeyBytes := crypto.CompressPubkey(t.publicKey)
	return hex.EncodeToString(pubKeyBytes)
}

// ProcessEncryptedData decrypts and processes the data
func (t *TEEService) ProcessEncryptedData(encryptedData []byte) (int64, error) {

	decryptedData, err := ecies.ImportECDSA(t.privateKey).Decrypt(encryptedData, nil, nil)
	if err != nil {
		return 0, err
	}

	switch strings.ToLower(string(decryptedData)) {
	case "extremely high risk":
		return 4, nil
	case "high risk":
		return 3, nil
	case "slightly high risk":
		return 2, nil
	case "low risk":
		return 1, nil
	default:
		return 0, nil
	}
}
