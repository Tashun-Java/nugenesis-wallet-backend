package staticControllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tashunc/nugenesis-wallet-backend/static/staticModels"
	"github.com/tashunc/nugenesis-wallet-backend/static/staticServices"
)

// BlockchainController handles HTTP requests for blockchain operations
type BlockchainController struct {
	blockchainService *staticServices.BlockchainService
}

// NewBlockchainController creates a new BlockchainController instance
func NewBlockchainController(blockchainService *staticServices.BlockchainService) *BlockchainController {
	return &BlockchainController{
		blockchainService: blockchainService,
	}
}

// GetAllBlockchains handles GET requests to get all blockchains with their IDs
func (c *BlockchainController) GetAllBlockchains(ctx *gin.Context) {
	blockchains, err := c.blockchainService.GetAllBlockchains()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, staticModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	response := staticModels.AllBlockchainsResponse{
		Blockchains: blockchains,
		Count:       len(blockchains),
	}
	ctx.JSON(http.StatusOK, response)
}

// GetBlockchainByID handles GET requests to get a blockchain by its ID
func (c *BlockchainController) GetBlockchainByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, staticModels.ErrorResponse{
			Error: "ID parameter is required",
		})
		return
	}

	blockchain, err := c.blockchainService.GetBlockchainByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, staticModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if blockchain == nil {
		ctx.JSON(http.StatusNotFound, staticModels.ErrorResponse{
			Error: "Blockchain not found",
		})
		return
	}

	response := staticModels.BlockchainByIDResponse{
		Blockchain: *blockchain,
	}
	ctx.JSON(http.StatusOK, response)
}

// GetBlockchainID handles GET requests to get the ID for a blockchain by name
func (c *BlockchainController) GetBlockchainID(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, staticModels.ErrorResponse{
			Error: "Name parameter is required",
		})
		return
	}

	id, exists := c.blockchainService.GetBlockchainID(name)
	if !exists {
		ctx.JSON(http.StatusNotFound, staticModels.ErrorResponse{
			Error: "Blockchain not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"blockchain_name": name,
		"id":              id,
	})
}
