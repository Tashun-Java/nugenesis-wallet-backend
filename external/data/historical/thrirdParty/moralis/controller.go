package moralis

import (
	"github.com/gin-gonic/gin"
	"github.com/tashunc/nugenesis-wallet-backend/external/models"
	"net/http"
	"strconv"
)

type Controller struct {
	service *Service
	chain   string
}

func NewController(chain string) *Controller {
	return &Controller{
		service: NewService(),
		chain:   chain,
	}
}

func (c *Controller) GetWalletHistory(ctx *gin.Context) {
	address := ctx.Param("address")
	if address == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
		return
	}

	cursor := ctx.Query("cursor")
	limitStr := ctx.DefaultQuery("limit", "100")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	history, err := c.service.GetWalletHistory(address, c.chain, cursor, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map Moralis transactions to standard transaction format
	var mappedTransactions []models.Transaction
	for _, tx := range history.Result {
		mappedTx := MapHistoryToTransaction(tx, address)
		mappedTransactions = append(mappedTransactions, mappedTx)
	}

	ctx.JSON(http.StatusOK, mappedTransactions)
}
