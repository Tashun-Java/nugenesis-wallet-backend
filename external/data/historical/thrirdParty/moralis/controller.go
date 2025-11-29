package moralis

import (
	"github.com/gin-gonic/gin"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/moralis/moralis_models"
	"github.com/tashunc/nugenesis-wallet-backend/external/models"
	"net/http"
	"os"
	"strconv"
)

// TokenIDServiceGetter is a function type for getting the token ID service
// This allows us to avoid circular dependencies
type TokenIDServiceGetter func() TokenIDServiceInterface

type Controller struct {
	service              *Service
	chain                string
	tokenIDServiceGetter TokenIDServiceGetter
}

func NewController(chain string) *Controller {
	return &Controller{
		service:              NewService(),
		chain:                chain,
		tokenIDServiceGetter: nil, // Will be set later
	}
}

// SetTokenIDServiceGetter sets the token ID service getter
func (c *Controller) SetTokenIDServiceGetter(getter TokenIDServiceGetter) {
	c.tokenIDServiceGetter = getter
}

// shouldFilterTokensWithoutID checks if tokens without token_id should be filtered
func shouldFilterTokensWithoutID() bool {
	filterValue := os.Getenv("FILTER_TOKENS_WITHOUT_ID")
	return filterValue == "true" || filterValue == "1"
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

	// Get token ID service
	var tokenIDService TokenIDServiceInterface
	if c.tokenIDServiceGetter != nil {
		tokenIDService = c.tokenIDServiceGetter()
	}

	// Map Moralis token balances to standard format
	var mappedBalances []models.WalletTokenBalance
	for _, token := range balances.Result {
		mappedBalance := MapTokenBalanceToStandard(token, chain, tokenIDService)
		mappedBalances = append(mappedBalances, mappedBalance)
	}

	// Filter balances by token_id if environment variable is set
	mappedBalances = filterBalancesByTokenID(mappedBalances)

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

// GetSolanaWalletTokenBalances retrieves all token balances (including native SOL) for a Solana wallet address
// Uses Moralis Solana Gateway API which has a different endpoint structure
func (c *Controller) GetSolanaWalletTokenBalances(ctx *gin.Context) {
	address := ctx.Param("address")
	if address == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
		return
	}

	// Get optional network parameter (default: mainnet)
	network := ctx.DefaultQuery("network", "mainnet")

	// Get native SOL balance
	nativeBalance, err := c.service.GetSolanaBalance(address, network)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get SPL token balances
	tokenBalances, err := c.service.GetSolanaTokenBalances(address, network)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get token ID service
	var tokenIDService TokenIDServiceInterface
	if c.tokenIDServiceGetter != nil {
		tokenIDService = c.tokenIDServiceGetter()
	}

	// Map native SOL balance to standard format
	var mappedBalances []models.WalletTokenBalance

	// Add native SOL as the first token
	nativeToken := moralis_models.SolanaToken{
		Mint:      "So11111111111111111111111111111111111111112", // Native SOL mint address
		Amount:    nativeBalance.Solana,
		AmountRaw: nativeBalance.Lamports,
		Decimals:  "9",
		Name:      "Solana",
		Symbol:    "SOL",
	}
	mappedNative := MapSolanaTokenToStandard(nativeToken, "solana", true, tokenIDService)
	mappedBalances = append(mappedBalances, mappedNative)

	// Add SPL tokens
	for _, token := range *tokenBalances {
		mappedBalance := MapSolanaTokenToStandard(token, "solana", false, tokenIDService)
		mappedBalances = append(mappedBalances, mappedBalance)
	}

	// Filter balances by token_id if environment variable is set
	mappedBalances = filterBalancesByTokenID(mappedBalances)

	// Prepare response
	response := models.WalletTokenBalancesResponse{
		Success:  true,
		Address:  address,
		Chain:    "solana",
		Balances: mappedBalances,
		Cursor:   "",    // Solana API doesn't use cursor pagination
		HasMore:  false, // No pagination for Solana API
	}

	ctx.JSON(http.StatusOK, response)
}
