package moralis_models

type WalletHistoryResponse struct {
	Cursor  string               `json:"cursor"`
	HasMore bool                 `json:"hasMore"`
	Result  []HistoryTransaction `json:"result"`
}

type HistoryTransaction struct {
	Hash              string           `json:"hash"`
	Nonce             string           `json:"nonce"`
	TransactionIndex  string           `json:"transaction_index"`
	FromAddress       string           `json:"from_address"`
	ToAddress         string           `json:"to_address"`
	Value             string           `json:"value"`
	Gas               string           `json:"gas"`
	GasPrice          string           `json:"gas_price"`
	GasUsed           string           `json:"gas_used"`
	CumulativeGasUsed string           `json:"cumulative_gas_used"`
	InputData         string           `json:"input"`
	ReceiptStatus     string           `json:"receipt_status"`
	BlockTimestamp    string           `json:"block_timestamp"`
	BlockNumber       string           `json:"block_number"`
	BlockHash         string           `json:"block_hash"`
	TransactionFee    string           `json:"transaction_fee,omitempty"`
	TransferIndex     []int            `json:"transfer_index,omitempty"`
	Logs              []Log            `json:"logs,omitempty"`
	NFTTransfers      []NFTTransfer    `json:"nft_transfers,omitempty"`
	ERC20Transfers    []ERC20Transfer  `json:"erc20_transfers,omitempty"`
	NativeTransfers   []NativeTransfer `json:"native_transfers,omitempty"`
	Category          string           `json:"category,omitempty"`
	Summary           string           `json:"summary,omitempty"`
}

type NFTTransfer struct {
	LogIndex           int    `json:"log_index"`
	Value              string `json:"value"`
	ContractType       string `json:"contract_type"`
	TransactionType    string `json:"transaction_type"`
	TokenAddress       string `json:"token_address"`
	TokenID            string `json:"token_id"`
	FromAddress        string `json:"from_address"`
	ToAddress          string `json:"to_address"`
	Amount             string `json:"amount"`
	Operator           string `json:"operator"`
	PossibleSpam       bool   `json:"possible_spam"`
	VerifiedCollection bool   `json:"verified_collection"`
	Direction          string `json:"direction"`
}

type ERC20Transfer struct {
	TokenName        string `json:"token_name"`
	TokenSymbol      string `json:"token_symbol"`
	TokenLogo        string `json:"token_logo"`
	TokenDecimals    string `json:"token_decimals"`
	FromAddress      string `json:"from_address"`
	ToAddress        string `json:"to_address"`
	Address          string `json:"address"`
	LogIndex         int    `json:"log_index"`
	Value            string `json:"value"`
	ValueFormatted   string `json:"value_formatted"`
	PossibleSpam     bool   `json:"possible_spam"`
	VerifiedContract bool   `json:"verified_contract"`
	SecurityScore    int    `json:"security_score"`
	Direction        string `json:"direction"`
}

type NativeTransfer struct {
	FromAddress         string `json:"from_address"`
	ToAddress           string `json:"to_address"`
	Value               string `json:"value"`
	ValueFormatted      string `json:"value_formatted"`
	Direction           string `json:"direction"`
	InternalTransaction bool   `json:"internal_transaction"`
	TokenSymbol         string `json:"token_symbol"`
	TokenLogo           string `json:"token_logo"`
}

type Log struct {
	LogIndex         string `json:"log_index"`
	TransactionHash  string `json:"transaction_hash"`
	TransactionIndex string `json:"transaction_index"`
	Address          string `json:"address"`
	Data             string `json:"data"`
	Topic0           string `json:"topic0,omitempty"`
	Topic1           string `json:"topic1,omitempty"`
	Topic2           string `json:"topic2,omitempty"`
	Topic3           string `json:"topic3,omitempty"`
}

// WalletTokenBalancesResponse represents the response from Moralis /wallets/{address}/tokens endpoint
type WalletTokenBalancesResponse struct {
	Cursor  string         `json:"cursor,omitempty"`
	Page    int            `json:"page,omitempty"`
	Result  []TokenBalance `json:"result"`
	HasMore bool           `json:"hasMore,omitempty"`
}

// TokenBalance represents a single token balance in the wallet
type TokenBalance struct {
	TokenAddress                    string   `json:"token_address"`
	Name                            string   `json:"name"`
	Symbol                          string   `json:"symbol"`
	Logo                            string   `json:"logo,omitempty"`
	Thumbnail                       string   `json:"thumbnail,omitempty"`
	Decimals                        int      `json:"decimals"`
	Balance                         string   `json:"balance"`
	BalanceFormatted                string   `json:"balance_formatted"`
	PossibleSpam                    bool     `json:"possible_spam"`
	VerifiedContract                bool     `json:"verified_contract"`
	TotalSupply                     string   `json:"total_supply,omitempty"`
	TotalSupplyFormatted            string   `json:"total_supply_formatted,omitempty"`
	PercentageRelativeToTotalSupply *float64 `json:"percentage_relative_to_total_supply,omitempty"`
	SecurityScore                   *int     `json:"security_score,omitempty"`
	UsdPrice                        *float64 `json:"usd_price,omitempty"`
	UsdPrice24hrPercentChange       *float64 `json:"usd_price_24hr_percent_change,omitempty"`
	UsdPrice24hrUsdChange           *float64 `json:"usd_price_24hr_usd_change,omitempty"`
	UsdValue                        *float64 `json:"usd_value,omitempty"`
	UsdValue24hrUsdChange           *float64 `json:"usd_value_24hr_usd_change,omitempty"`
	NativeToken                     *bool    `json:"native_token,omitempty"`
	PortfolioPercentage             *float64 `json:"portfolio_percentage,omitempty"`
}

// Solana-specific models for Moralis Solana Gateway API

// SolanaTokenBalancesResponse represents the response from Moralis Solana getSPL endpoint
// The API returns an array of tokens directly
type SolanaTokenBalancesResponse []SolanaToken

// SolanaToken represents a single SPL token balance
type SolanaToken struct {
	AssociatedTokenAddress string  `json:"associatedTokenAddress"`
	Mint                   string  `json:"mint"`
	AmountRaw              string  `json:"amountRaw"`
	Amount                 string  `json:"amount"`
	Decimals               string  `json:"decimals"`
	Name                   string  `json:"name,omitempty"`
	Symbol                 string  `json:"symbol,omitempty"`
	Logo                   string  `json:"logo,omitempty"`
	UsdPrice               float64 `json:"usdPrice,omitempty"`
	UsdValue               float64 `json:"usdValue,omitempty"`
}

// SolanaBalanceResponse represents the response from Moralis Solana balance endpoint
type SolanaBalanceResponse struct {
	Lamports string `json:"lamports"`
	Solana   string `json:"solana"`
}

// SolanaPortfolioResponse represents the response from Moralis Solana portfolio endpoint
type SolanaPortfolioResponse struct {
	NativeBalance struct {
		Lamports string `json:"lamports"`
		Solana   string `json:"solana"`
	} `json:"nativeBalance"`
	Tokens    []SolanaToken `json:"tokens"`
	NftTokens []SolanaToken `json:"nftTokens,omitempty"`
}
