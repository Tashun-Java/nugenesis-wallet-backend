package tronscan

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/tashunc/nugenesis-wallet-backend/external/data/fmv/thrirdParty/coingecko"
	"github.com/tashunc/nugenesis-wallet-backend/external/models"
)

// TokenIDServiceInterface defines the interface for token ID lookups
type TokenIDServiceInterface interface {
	GetTokenID(chain, address string) string
	GetTokenIDForNative(chain, symbol string) string
}

type Controller struct {
	service              *Service
	tokenIDServiceGetter func() TokenIDServiceInterface
}

func NewController() *Controller {
	return &Controller{
		service: NewService(),
	}
}

// SetTokenIDServiceGetter sets the function to get the token ID service
func (c *Controller) SetTokenIDServiceGetter(getter func() TokenIDServiceInterface) {
	c.tokenIDServiceGetter = getter
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

	// Get token ID service
	var tokenIDService TokenIDServiceInterface
	if c.tokenIDServiceGetter != nil {
		tokenIDService = c.tokenIDServiceGetter()
	}

	// Map token balances to standard format
	var balances []models.WalletTokenBalance

	// Add native TRX balance first (price will be enriched by CoinGecko)
	nativeTRX := CreateNativeTRXBalance(accountInfo.Balance, 0, tokenIDService)
	balances = append(balances, nativeTRX)

	// Add TRC20 token balances
	for _, token := range accountInfo.TRC20TokenBalances {
		standardBalance := MapTokenBalanceToStandard(token, address, tokenIDService)
		balances = append(balances, standardBalance)
	}

	// Add TRC10 token balances (if any)
	for _, token := range accountInfo.TokenBalances {
		standardBalance := MapTokenBalanceToStandard(token, address, tokenIDService)
		balances = append(balances, standardBalance)
	}

	// Filter balances by token_id if environment variable is set
	balances = filterBalancesByTokenID(balances)

	// Enrich balances with CoinGecko prices if UsdPrice/UsdValue are missing
	balances = enrichBalancesWithCoinGeckoPrices(balances)

	// Prepare response
	response := models.WalletTokenBalancesResponse{
		Success:  true,
		Address:  address,
		Chain:    "tron",
		Balances: balances,
		Cursor:   "",    // TronScan API doesn't use cursor pagination
		HasMore:  false, // We get all balances in one call
	}

	ctx.JSON(http.StatusOK, response)
}

// Helper functions

// shouldFilterTokensWithoutID checks if the environment variable is set to filter tokens without IDs
func shouldFilterTokensWithoutID() bool {
	filter := strings.ToLower(os.Getenv("FILTER_TOKENS_WITHOUT_ID"))
	return filter == "true" || filter == "1"
}

// filterBalancesByTokenID filters out tokens that don't have a token_id
func filterBalancesByTokenID(balances []models.WalletTokenBalance) []models.WalletTokenBalance {
	if !shouldFilterTokensWithoutID() {
		return balances
	}

	filtered := make([]models.WalletTokenBalance, 0)
	for _, balance := range balances {
		if balance.TokenID != "" {
			filtered = append(filtered, balance)
		}
	}
	return filtered
}

// enrichBalancesWithCoinGeckoPrices fetches missing prices from CoinGecko
func enrichBalancesWithCoinGeckoPrices(balances []models.WalletTokenBalance) []models.WalletTokenBalance {
	// Use the coingecko service to enrich balances
	service := coingecko.NewService()
	return service.EnrichBalancesWithPrices(balances)
}
