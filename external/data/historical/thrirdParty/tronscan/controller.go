package tronscan

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Controller struct {
	service *Service
}

func NewController() *Controller {
	return &Controller{
		service: NewService(),
	}
}

// GetAddressTransactions handles the API request for Tron address transactions
func (c *Controller) GetAddressTransactions(ctx *gin.Context) {
	address := ctx.Param("address")
	if address == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
		return
	}

	if !ValidateTronAddress(address) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid Tron address format"})
		return
	}

	// Parse query parameters
	limitStr := ctx.DefaultQuery("limit", "50")
	startStr := ctx.DefaultQuery("start", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200 // Max limit
	}

	start, err := strconv.Atoi(startStr)
	if err != nil || start < 0 {
		start = 0
	}

	// Fetch transactions from TronScan
	response, err := c.service.GetAddressTransactions(address, limit, start)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map transactions to standard format
	var standardTransactions []interface{}
	for _, tx := range response.Data {
		standardTx := MapToStandardTransaction(tx, address)
		standardTransactions = append(standardTransactions, standardTx)
	}

	ctx.JSON(http.StatusOK, standardTransactions)
}

// GetWalletTokenBalances handles the API request for Tron wallet token balances
func (c *Controller) GetWalletTokenBalances(ctx *gin.Context) {
	address := ctx.Param("address")
	if address == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
		return
	}

	if !ValidateTronAddress(address) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid Tron address format"})
		return
	}

	// Parse query parameters
	limitStr := ctx.DefaultQuery("limit", "20")
	startStr := ctx.DefaultQuery("start", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	if limit > 200 {
		limit = 200 // Max limit
	}

	start, err := strconv.Atoi(startStr)
	if err != nil || start < 0 {
		start = 0
	}

	// Fetch account info to get native TRX balance and TRC20 tokens
	accountInfo, err := c.service.GetAccountInfo(address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map token balances to standard format
	var balances []interface{}

	// Add native TRX balance first (assume TRX price, could be fetched from price API)
	trxPrice := 0.10 // Placeholder - should fetch real price
	nativeTRX := CreateNativeTRXBalance(accountInfo.Balance, trxPrice)
	balances = append(balances, nativeTRX)

	// Add TRC20 token balances
	for _, token := range accountInfo.TRC20TokenBalances {
		standardBalance := MapTokenBalanceToStandard(token, address)
		balances = append(balances, standardBalance)
	}

	// Add TRC10 token balances (if any)
	for _, token := range accountInfo.TokenBalances {
		standardBalance := MapTokenBalanceToStandard(token, address)
		balances = append(balances, standardBalance)
	}

	// Prepare response
	response := map[string]interface{}{
		"success":  true,
		"address":  address,
		"chain":    "tron",
		"balances": balances,
		"total":    len(balances),
	}

	ctx.JSON(http.StatusOK, response)
}
