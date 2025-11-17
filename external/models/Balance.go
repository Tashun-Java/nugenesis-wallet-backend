package models

type GetBalanceControllerRequest struct {
	BlockchainID   string `json:"blockchain_id,omitempty"`
	Address        string `json:"address" binding:"required"`
	BlockParameter string `json:"blockParameter,omitempty"`
}

type GetBalanceControllerResponse struct {
	Success bool                     `json:"success"`
	Balance string                   `json:"balance,omitempty"`
	Error   *SendRawTransactionError `json:"error,omitempty"`
	Message string                   `json:"message,omitempty"`
}

type GetTokenBalancesControllerRequest struct {
	BlockchainID      string   `json:"blockchain_id,omitempty"`
	Address           string   `json:"address" binding:"required"`
	ContractAddresses []string `json:"contractAddresses,omitempty"`
}

type TokenBalanceData struct {
	ContractAddress string  `json:"contractAddress"`
	TokenBalance    *string `json:"tokenBalance"`
	Error           *string `json:"error,omitempty"`
}

type GetTokenBalancesControllerResponse struct {
	Success       bool                     `json:"success"`
	Address       string                   `json:"address,omitempty"`
	TokenBalances []TokenBalanceData       `json:"tokenBalances,omitempty"`
	Error         *SendRawTransactionError `json:"error,omitempty"`
	Message       string                   `json:"message,omitempty"`
}

// WalletTokenBalance represents a comprehensive token balance with fiat values
type WalletTokenBalance struct {
	TokenAddress        string  `json:"token_address"`
	Name                string  `json:"name"`
	Symbol              string  `json:"symbol"`
	Logo                string  `json:"logo,omitempty"`
	Thumbnail           string  `json:"thumbnail,omitempty"`
	Decimals            string  `json:"decimals"`
	Balance             string  `json:"balance"`
	BalanceRaw          string  `json:"balance_raw"`
	NativeToken         bool    `json:"native_token"`
	VerifiedContract    bool    `json:"verified_contract"`
	PossibleSpam        bool    `json:"possible_spam"`
	UsdPrice            float64 `json:"usd_price"`
	UsdValue            float64 `json:"usd_value"`
	UsdPrice24hrChange  float64 `json:"usd_price_24hr_change"`
	PortfolioPercentage float64 `json:"portfolio_percentage"`
	SecurityScore       int     `json:"security_score,omitempty"`
	Chain               string  `json:"chain"`
}

// WalletTokenBalancesResponse represents the response for wallet token balances endpoint
type WalletTokenBalancesResponse struct {
	Success  bool                     `json:"success"`
	Address  string                   `json:"address"`
	Chain    string                   `json:"chain"`
	Balances []WalletTokenBalance     `json:"balances"`
	Cursor   string                   `json:"cursor,omitempty"`
	HasMore  bool                     `json:"has_more,omitempty"`
	Error    *SendRawTransactionError `json:"error,omitempty"`
	Message  string                   `json:"message,omitempty"`
}
