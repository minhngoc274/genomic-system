package blockchain

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/minhngoc274/genomic-system/genomic-service/adapters/blockchain/contracts"
	"github.com/minhngoc274/genomic-system/genomic-service/config"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BlockchainService struct {
	client     *ethclient.Client
	controller *contracts.Controller
	nft        *contracts.GeneNFT
	token      *contracts.PCSPToken
	props      *config.ChainProperties
}

func NewBlockchainService(props *config.ChainProperties) (*BlockchainService, error) {
	client, err := ethclient.Dial(props.RPC)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to network: %v", err)
	}

	controller, err := contracts.NewController(common.HexToAddress(props.ControllerAddress), client)
	if err != nil {
		return nil, fmt.Errorf("failed to load controller: %v", err)
	}

	nft, err := contracts.NewGeneNFT(common.HexToAddress(props.NFTAddress), client)
	if err != nil {
		return nil, fmt.Errorf("failed to load NFT contract: %v", err)
	}

	token, err := contracts.NewPCSPToken(common.HexToAddress(props.TokenAddress), client)
	if err != nil {
		return nil, fmt.Errorf("failed to load token contract: %v", err)
	}

	return &BlockchainService{
		props:      props,
		client:     client,
		controller: controller,
		nft:        nft,
		token:      token,
	}, nil
}

func (s *BlockchainService) GetTransactOpts() (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(s.props.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}

	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(s.props.ChainID))
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %v", err)
	}

	return opts, nil
}

// UploadData starts the upload data on blockchain
func (s *BlockchainService) UploadData(docID string) error {
	opts, err := s.GetTransactOpts()
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %v", err)
	}

	tx, err := s.controller.UploadData(opts, docID)
	if err != nil {
		return fmt.Errorf("failed to initiate upload: %v", err)
	}

	_, err = bind.WaitMined(context.Background(), s.client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction: %v", err)
	}

	return fmt.Errorf("failed to get session ID from event")
}

// Confirm confirm session data and reward for user
func (s *BlockchainService) Confirm(sessionID, riskScore *big.Int, docID, contentHash, proof string) error {
	opts, err := s.GetTransactOpts()
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %v", err)
	}

	tx, err := s.controller.Confirm(
		opts,
		docID,
		contentHash,
		proof,
		sessionID,
		riskScore,
	)
	if err != nil {
		return fmt.Errorf("failed to confirm upload: %v", err)
	}

	_, err = bind.WaitMined(context.Background(), s.client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction: %v", err)
	}

	return nil
}
