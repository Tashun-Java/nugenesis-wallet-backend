package etherscan_models

type AddressResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Result  []TxEntry `json:"result"`
}

type TxEntry struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             string `json:"value"`
	GasPrice          string `json:"gasPrice"`
	GasUsed           string `json:"gasUsed"`
	IsError           string `json:"isError"`
	Confirmations     string `json:"confirmations"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	TxReceiptStatus   string `json:"txreceipt_status"`
	Input             string `json:"input"`
}
