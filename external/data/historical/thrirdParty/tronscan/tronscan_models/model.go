package tronscan_models

import (
	"encoding/json"
	"strconv"
)

// FlexibleFloat is a custom type that can unmarshal both string and numeric JSON values
type FlexibleFloat float64

// UnmarshalJSON implements custom unmarshaling for FlexibleFloat
func (ff *FlexibleFloat) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as float64 first
	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		*ff = FlexibleFloat(f)
		return nil
	}

	// If that fails, try to unmarshal as string and convert
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// Convert string to float64
	if s == "" {
		*ff = 0
		return nil
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*ff = FlexibleFloat(f)
	return nil
}

// Float64 returns the float64 value
func (ff FlexibleFloat) Float64() float64 {
	return float64(ff)
}

// TronScanResponse represents the API response from TronScan
type TronScanResponse struct {
	Data       []Transaction `json:"data"`
	Total      int           `json:"total"`
	RangeTotal int           `json:"rangeTotal"`
}

// Transaction represents a single Tron transaction
type Transaction struct {
	BlockNumber   int64        `json:"block"`
	Hash          string       `json:"hash"`
	Timestamp     int64        `json:"timestamp"`
	OwnerAddress  string       `json:"ownerAddress"`
	ToAddress     string       `json:"toAddress"`
	ToAddressList []string     `json:"toAddressList"`
	ContractType  int          `json:"contractType"`
	Confirmed     bool         `json:"confirmed"`
	Revert        bool         `json:"revert"`
	ContractData  ContractData `json:"contractData"`
	TokenInfo     *TokenInfo   `json:"tokenInfo,omitempty"`
	TokenType     string       `json:"tokenType,omitempty"`
	Cost          *Cost        `json:"cost,omitempty"`
	TriggerInfo   *TriggerInfo `json:"trigger_info,omitempty"`
}

// ContractData contains transaction details
type ContractData struct {
	Amount          int64  `json:"amount,omitempty"`
	AssetName       string `json:"asset_name,omitempty"`
	OwnerAddress    string `json:"owner_address"`
	ToAddress       string `json:"to_address"`
	Data            string `json:"data,omitempty"`
	CallValue       int64  `json:"call_value,omitempty"`
	ContractAddress string `json:"contract_address,omitempty"`
}

// TokenInfo contains information about tokens involved in the transaction
type TokenInfo struct {
	TokenID      string `json:"tokenId"`
	TokenAbbr    string `json:"tokenAbbr"`
	TokenName    string `json:"tokenName"`
	TokenDecimal int    `json:"tokenDecimal"`
	TokenCanShow int    `json:"tokenCanShow"`
	TokenType    string `json:"tokenType"`
	TokenLogo    string `json:"tokenLogo"`
	TokenLevel   string `json:"tokenLevel"`
	Vip          bool   `json:"vip"`
}

// Cost represents transaction costs
type Cost struct {
	NetFee            int64 `json:"net_fee"`
	EnergyUsage       int64 `json:"energy_usage"`
	EnergyFee         int64 `json:"energy_fee"`
	EnergyUsageTotal  int64 `json:"energy_usage_total"`
	OriginEnergyUsage int64 `json:"origin_energy_usage"`
	NetUsage          int64 `json:"net_usage"`
}

// TriggerInfo contains smart contract trigger information
type TriggerInfo struct {
	Method     string                 `json:"method"`
	Data       string                 `json:"data"`
	Parameter  map[string]interface{} `json:"parameter"`
	MethodID   string                 `json:"methodId"`
	MethodName string                 `json:"methodName"`
	Contract   string                 `json:"contract"`
	CallValue  int64                  `json:"call_value"`
}

// TronGridResponse represents the API response from TronGrid
type TronGridResponse struct {
	Data    []TronGridTransaction `json:"data"`
	Success bool                  `json:"success"`
	Meta    TronGridMeta          `json:"meta"`
}

// TronGridTransaction represents a transaction from TronGrid API
type TronGridTransaction struct {
	TxID           string          `json:"txID"`
	BlockNumber    int64           `json:"blockNumber"`
	BlockTimestamp int64           `json:"block_timestamp"`
	RawData        TronGridRawData `json:"raw_data"`
	RawDataHex     string          `json:"raw_data_hex"`
	Ret            []TronGridRet   `json:"ret"`
	Signature      []string        `json:"signature"`
	ContractResult []string        `json:"contract_result,omitempty"`
}

// TronGridRawData contains raw transaction data
type TronGridRawData struct {
	Contract      []TronGridContract `json:"contract"`
	RefBlockBytes string             `json:"ref_block_bytes"`
	RefBlockHash  string             `json:"ref_block_hash"`
	Expiration    int64              `json:"expiration"`
	Timestamp     int64              `json:"timestamp"`
	FeeLimit      int64              `json:"fee_limit,omitempty"`
}

// TronGridContract contains contract information
type TronGridContract struct {
	Type      string                 `json:"type"`
	Parameter map[string]interface{} `json:"parameter"`
}

// TronGridRet contains transaction result
type TronGridRet struct {
	ContractRet string `json:"contractRet"`
	Fee         int64  `json:"fee,omitempty"`
}

// TronGridMeta contains pagination metadata
type TronGridMeta struct {
	At          int64  `json:"at"`
	Fingerprint string `json:"fingerprint"`
	PageSize    int    `json:"page_size"`
}

// AccountTokensResponse represents the response from TronScan account tokens API
type AccountTokensResponse struct {
	Data  []TokenBalance `json:"data"`
	Total int            `json:"total"`
}

// TokenBalance represents a token balance for a Tron account
type TokenBalance struct {
	TokenID          string        `json:"tokenId"`
	TokenName        string        `json:"tokenName"`
	TokenAbbr        string        `json:"tokenAbbr"`
	TokenDecimal     int           `json:"tokenDecimal"`
	TokenCanShow     int           `json:"tokenCanShow"`
	TokenType        string        `json:"tokenType"`
	TokenLogo        string        `json:"tokenLogo"`
	TokenPriceInTrx  float64       `json:"tokenPriceInTrx"`
	Amount           FlexibleFloat `json:"amount,omitempty"`
	Balance          string        `json:"balance"`
	TokenPriceInUsd  float64       `json:"tokenPriceInUsd,omitempty"`
	AmountInUsd      float64       `json:"amountInUsd,omitempty"`
	NrOfTokenHolders int           `json:"nrOfTokenHolders,omitempty"`
	TransferCount    int           `json:"transferCount,omitempty"`
	Vip              bool          `json:"vip,omitempty"`
	OwnerAddress     string        `json:"owner_address,omitempty"`
}

// AccountInfoResponse represents account info including TRX balance
type AccountInfoResponse struct {
	TronPower             int64          `json:"tronPower"`
	Balance               int64          `json:"balance"`
	TotalFrozen           int64          `json:"totalFrozen"`
	Bandwidth             Bandwidth      `json:"bandwidth,omitempty"`
	TRC20TokenBalances    []TokenBalance `json:"trc20token_balances,omitempty"`
	TokenBalances         []TokenBalance `json:"tokenBalances,omitempty"`
	Address               string         `json:"address"`
	Name                  string         `json:"name,omitempty"`
	TotalTransactionCount int            `json:"totalTransactionCount"`
}

// Bandwidth represents bandwidth information
type Bandwidth struct {
	EnergyRemaining   int64 `json:"energyRemaining"`
	TotalEnergyLimit  int64 `json:"totalEnergyLimit"`
	TotalEnergyWeight int64 `json:"totalEnergyWeight"`
	NetUsed           int64 `json:"netUsed"`
	NetLimit          int64 `json:"netLimit"`
}
