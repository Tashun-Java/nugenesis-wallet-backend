package staticModels

// BlockchainInfo represents the native token information for a blockchain
type BlockchainInfo struct {
	Name        string   `json:"name"`
	Website     string   `json:"website"`
	Description string   `json:"description"`
	Explorer    string   `json:"explorer"`
	Symbol      string   `json:"symbol"`
	Type        string   `json:"type"`
	Decimals    int      `json:"decimals"`
	Status      string   `json:"status"`
	RPCURL      string   `json:"rpc_url,omitempty"`
	Links       []Link   `json:"links"`
	Tags        []string `json:"tags,omitempty"`
}

// AssetInfo represents token/asset information
type AssetInfo struct {
	Name        string   `json:"name"`
	Website     string   `json:"website"`
	Description string   `json:"description"`
	Explorer    string   `json:"explorer"`
	Symbol      string   `json:"symbol"`
	Type        string   `json:"type"`
	Decimals    int      `json:"decimals"`
	Status      string   `json:"status"`
	ID          string   `json:"id"`
	Links       []Link   `json:"links"`
	Tags        []string `json:"tags,omitempty"`
}

// Link represents external links
type Link struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// AssetResponse represents the response structure for GetByCoinSymbol
type AssetResponse struct {
	ID           string `json:"id"`
	Symbol       string `json:"symbol"`
	Name         string `json:"name"`
	Blockchain   string `json:"blockchain"`
	BlockchainID string `json:"blockchain_id,omitempty"`
	Address      string `json:"address,omitempty"`
	LogoPath     string `json:"logo_path"`
	Type         string `json:"type"`
	Decimals     int    `json:"decimals"`
	Status       string `json:"status"`
}

// AssetsBySymbolResponse represents the response for GetByCoinSymbol endpoint
type AssetsBySymbolResponse struct {
	Symbol string          `json:"symbol"`
	Assets []AssetResponse `json:"assets"`
	Count  int             `json:"count"`
}

// AllSymbolsResponse represents the response for GetAllSymbols endpoint
type AllSymbolsResponse struct {
	Symbols []string `json:"symbols"`
	Count   int      `json:"count"`
}

// AllAssetsResponse represents the response for GetAllAssets endpoint
type AllAssetsResponse struct {
	Assets []AssetResponse `json:"assets"`
	Count  int             `json:"count"`
}

// RefreshResponse represents the response for ForceRefresh endpoint
type RefreshResponse struct {
	Message string `json:"message"`
}

// CacheStatsResponse represents the response for GetCacheStats endpoint
type CacheStatsResponse struct {
	CacheStats map[string]interface{} `json:"cache_stats"`
}

// ErrorResponse represents error responses
type ErrorResponse struct {
	Error string `json:"error"`
}

// BlockchainResponse represents a blockchain with its ID
type BlockchainResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Symbol     string `json:"symbol,omitempty"`
	IsTestnet  bool   `json:"is_testnet"`
	HasAssets  bool   `json:"has_assets"`
	AssetCount int    `json:"asset_count"`
}

// AllBlockchainsResponse represents the response for GetAllBlockchains endpoint
type AllBlockchainsResponse struct {
	Blockchains []BlockchainResponse `json:"blockchains"`
	Count       int                  `json:"count"`
}

// BlockchainByIDResponse represents the response for GetBlockchainByID endpoint
type BlockchainByIDResponse struct {
	Blockchain BlockchainResponse `json:"blockchain"`
}
