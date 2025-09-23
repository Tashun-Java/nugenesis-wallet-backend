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
