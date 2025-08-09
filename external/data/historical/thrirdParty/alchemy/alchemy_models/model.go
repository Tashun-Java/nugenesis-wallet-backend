package alchemy_models

type TokensByAddressRequest struct {
	Addresses []AddressRequest `json:"addresses"`
	Limit     *int             `json:"limit,omitempty"`
}

type AddressRequest struct {
	Address  string   `json:"address"`
	Networks []string `json:"networks"`
}

type TokensByAddressResponse struct {
	Data TokenData `json:"data"`
}

type TokenData struct {
	Tokens []TokenInfo `json:"tokens"`
}

type TokenInfo struct {
	Address       string        `json:"address"`
	Network       string        `json:"network"`
	TokenAddress  *string       `json:"tokenAddress"`
	TokenBalance  string        `json:"tokenBalance"`
	TokenMetadata TokenMetadata `json:"tokenMetadata"`
	TokenPrices   []TokenPrice  `json:"tokenPrices"`
}

type TokenMetadata struct {
	Symbol   *string `json:"symbol"`
	Decimals *int    `json:"decimals"`
	Name     *string `json:"name"`
	Logo     *string `json:"logo"`
}

type TokenPrice struct {
	Currency      string `json:"currency"`
	Value         string `json:"value"`
	LastUpdatedAt string `json:"lastUpdatedAt"`
}

type TransactionHistoryRequest struct {
	Addresses []AddressRequest `json:"addresses"`
	Before    *string          `json:"before,omitempty"`
	After     *string          `json:"after,omitempty"`
	Limit     *int             `json:"limit,omitempty"`
}

type TransactionHistoryResponse struct {
	Transactions []Transaction `json:"transactions"`
	Before       string        `json:"before"`
	After        string        `json:"after"`
	TotalCount   int           `json:"totalCount"`
}

type Transaction struct {
	Network           string           `json:"network"`
	Hash              string           `json:"hash"`
	TimeStamp         string           `json:"timeStamp"`
	BlockNumber       int              `json:"blockNumber"`
	BlockHash         string           `json:"blockHash"`
	Nonce             string           `json:"nonce"`
	TransactionIndex  int              `json:"transactionIndex"`
	FromAddress       string           `json:"fromAddress"`
	ToAddress         string           `json:"toAddress"`
	ContractAddress   *string          `json:"contractAddress"`
	Value             string           `json:"value"`
	CumulativeGasUsed string           `json:"cumulativeGasUsed"`
	EffectiveGasPrice string           `json:"effectiveGasPrice"`
	GasUsed           string           `json:"gasUsed"`
	Logs              []TransactionLog `json:"logs"`
	InternalTxns      []InternalTxn    `json:"internalTxns"`
}

type TransactionLog struct {
	ContractAddress string   `json:"contractAddress"`
	LogIndex        int      `json:"logIndex"`
	Data            string   `json:"data"`
	Removed         bool     `json:"removed"`
	Topics          []string `json:"topics"`
}

type InternalTxn struct {
	Type         string  `json:"type"`
	FromAddress  string  `json:"fromAddress"`
	ToAddress    string  `json:"toAddress"`
	Value        string  `json:"value"`
	Gas          string  `json:"gas"`
	GasUsed      string  `json:"gasUsed"`
	Input        string  `json:"input"`
	Output       string  `json:"output"`
	Error        *string `json:"error"`
	RevertReason *string `json:"revertReason"`
}
