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

// GetWalletTokenBalances retrieves all token balances (including native tokens) for a wallet address
// Supports multi-chain and multi-token with fiat values
func (c *Controller) GetWalletTokenBalances(ctx *gin.Context) {
	address := ctx.Param("address")
	if address == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
		return
	}

	// Get optional query parameters
	cursor := ctx.Query("cursor")
	limitStr := ctx.DefaultQuery("limit", "100")
	excludeSpamStr := ctx.DefaultQuery("exclude_spam", "true")
	chain := ctx.DefaultQuery("chain", c.chain) // Allow override of default chain

	// Parse limit parameter
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	// Parse exclude_spam parameter
	excludeSpam, err := strconv.ParseBool(excludeSpamStr)
	if err != nil {
		excludeSpam = true // Default to excluding spam
	}

	// Call Moralis API to get token balances
	balances, err := c.service.GetWalletTokenBalances(address, chain, cursor, limit, excludeSpam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map Moralis token balances to standard format
	var mappedBalances []models.WalletTokenBalance
	for _, token := range balances.Result {
		mappedBalance := MapTokenBalanceToStandard(token, chain)
		mappedBalances = append(mappedBalances, mappedBalance)
	}

	// Prepare response with pagination info
	response := models.WalletTokenBalancesResponse{
		Success:  true,
		Address:  address,
		Chain:    chain,
		Balances: mappedBalances,
		Cursor:   balances.Cursor,
		HasMore:  balances.HasMore,
	}

	ctx.JSON(http.StatusOK, response)
}
