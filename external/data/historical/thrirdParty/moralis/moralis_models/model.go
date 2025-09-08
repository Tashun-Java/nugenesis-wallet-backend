package moralis_models

type WalletHistoryResponse struct {
	Cursor  string               `json:"cursor"`
	HasMore bool                 `json:"hasMore"`
	Result  []HistoryTransaction `json:"result"`
}

type HistoryTransaction struct {
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	TransactionIndex  string `json:"transaction_index"`
	FromAddress       string `json:"from_address"`
	ToAddress         string `json:"to_address"`
	Value             string `json:"value"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gas_price"`
	GasUsed           string `json:"gas_used"`
	CumulativeGasUsed string `json:"cumulative_gas_used"`
	InputData         string `json:"input"`
	ReceiptStatus     string `json:"receipt_status"`
	BlockTimestamp    string `json:"block_timestamp"`
	BlockNumber       string `json:"block_number"`
	BlockHash         string `json:"block_hash"`
	TransferIndex     []int  `json:"transfer_index,omitempty"`
	Logs              []Log  `json:"logs,omitempty"`
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
