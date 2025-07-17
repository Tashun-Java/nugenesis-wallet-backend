package alchemy_models

import "github.com/tashunc/nugenesis-wallet-backend/internal/models"

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

type SendRawTransactionRequest struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      int      `json:"id"`
}

type SendRawTransactionResponse struct {
	Jsonrpc string                          `json:"jsonrpc"`
	Id      int                             `json:"id"`
	Result  string                          `json:"result,omitempty"`
	Error   *models.SendRawTransactionError `json:"error,omitempty"`
}

type TransactionObject struct {
	From  string `json:"from"`
	To    string `json:"to,omitempty"`
	Value string `json:"value,omitempty"`
	Data  string `json:"data,omitempty"`
	Nonce string `json:"nonce,omitempty"`
}

type EstimateGasRequest struct {
	Jsonrpc string              `json:"jsonrpc"`
	Method  string              `json:"method"`
	Params  []TransactionObject `json:"params"`
	Id      int                 `json:"id"`
}

type EstimateGasResponse struct {
	Jsonrpc string                          `json:"jsonrpc"`
	Id      int                             `json:"id"`
	Result  string                          `json:"result,omitempty"`
	Error   *models.SendRawTransactionError `json:"error,omitempty"`
}

type GetTransactionCountRequest struct {
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	Id      int      `json:"id"`
}

type GetTransactionCountResponse struct {
	Jsonrpc string                          `json:"jsonrpc"`
	Id      int                             `json:"id"`
	Result  string                          `json:"result,omitempty"`
	Error   *models.SendRawTransactionError `json:"error,omitempty"`
}
