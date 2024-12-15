package controllers

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	baseEx "github.com/golibs-starter/golib/exception"
	"github.com/golibs-starter/golib/web/response"
	"github.com/minhngoc274/genomic-system/genomic-service/adapters/blockchain"
	"github.com/minhngoc274/genomic-system/genomic-service/adapters/repositories"
	"github.com/minhngoc274/genomic-system/genomic-service/adapters/tee"
	"github.com/minhngoc274/genomic-system/genomic-service/models"
	"github.com/minhngoc274/genomic-system/genomic-service/models/requests"
	"github.com/minhngoc274/genomic-system/genomic-service/utils"
	"net/http"
	"strings"
)

// UserController represents ...
type UserController struct {
	userRepository        *repositories.UserRepository
	geneticDataRepository *repositories.GeneticDataRepository
	blockchainService     *blockchain.BlockchainService
	teeService            *tee.TEEService
}

// NewUserController return a new UserController instance
func NewUserController(
	userRepository *repositories.UserRepository,
	geneticDataRepository *repositories.GeneticDataRepository,
	blockchainService *blockchain.BlockchainService,
	teeService *tee.TEEService,
) *UserController {
	return &UserController{
		userRepository:        userRepository,
		geneticDataRepository: geneticDataRepository,
		blockchainService:     blockchainService,
		teeService:            teeService,
	}
}

// Register register user then return tee public key so that user can use this to encrypt data
func (c *UserController) Register(ctx *gin.Context) {
	var req requests.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			response.WriteError(ctx.Writer, baseEx.New(http.StatusBadRequest, validationErrors[0].Error()))
			return
		}
		response.WriteError(ctx.Writer, baseEx.New(http.StatusBadRequest, fmt.Sprintf("unmarshall request error: %v", err)))
		return
	}

	if err := c.userRepository.Save(models.User{Address: req.Address}); err != nil {
		if strings.Contains(err.Error(), "user already exists") {
			response.WriteError(ctx.Writer, baseEx.New(http.StatusBadRequest, err.Error()))
		} else {
			response.WriteError(ctx.Writer, baseEx.New(http.StatusInternalServerError, err.Error()))
		}
		return
	}

	response.Write(ctx.Writer, response.Ok(map[string]interface{}{
		"public_key": c.teeService.GetPublicKey(),
	}))
}

// UploadData upload encrypt data
func (c *UserController) UploadData(ctx *gin.Context) {
	var req requests.UploadDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			response.WriteError(ctx.Writer, baseEx.New(http.StatusBadRequest, validationErrors[0].Error()))
			return
		}
		response.WriteError(ctx.Writer, baseEx.New(http.StatusBadRequest, fmt.Sprintf("unmarshall request error: %v", err)))
		return
	}

	if !utils.VerifySignature(req.Address, req.Signature, req.EncryptedData) {
		response.WriteError(ctx.Writer, baseEx.New(http.StatusBadRequest, errors.New("verify signature failed").Error()))
		return
	}

	dataHash := crypto.Keccak256Hash(req.EncryptedData).Bytes()
	fileID := hex.EncodeToString(dataHash[:16])
	if err := c.geneticDataRepository.Create(models.GeneticData{
		FileID:        fileID,
		UserAddress:   req.Address,
		DataHash:      dataHash,
		EncryptedData: req.EncryptedData,
	}); err != nil {
		if strings.Contains(err.Error(), "file ID already exists") {
			response.WriteError(ctx.Writer, baseEx.New(http.StatusBadRequest, err.Error()))
		} else {
			response.WriteError(ctx.Writer, baseEx.New(http.StatusInternalServerError, err.Error()))
		}
		return
	}

	if err := c.blockchainService.UploadData(fileID); err != nil {
		response.WriteError(ctx.Writer, baseEx.New(http.StatusInternalServerError, err.Error()))
	}

	response.Write(ctx.Writer, response.Ok("Upload data successfully"))
}
