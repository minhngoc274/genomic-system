package config

import (
	"github.com/golibs-starter/golib/config"
)

// ChainProperties represents ...
type ChainProperties struct {
	RPC               string
	ChainID           int64
	ControllerAddress string
	TokenAddress      string
	NFTAddress        string
	PrivateKey        string
}

// NewChainProperties return a new ChainProperties instance
func NewChainProperties(loader config.Loader) (*ChainProperties, error) {
	props := ChainProperties{}
	err := loader.Bind(&props)
	return &props, err
}

// Prefix return config prefix
func (t *ChainProperties) Prefix() string {
	return "app.chains"
}
