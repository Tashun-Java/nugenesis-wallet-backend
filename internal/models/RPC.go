package models

type SendRawTransactionError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SendRawTransactionControllerRequest struct {
	SignedTransactions []string `json:"params" binding:"required"`
}

type SendRawTransactionControllerResponse struct {
	Success         bool                     `json:"success"`
	TransactionHash string                   `json:"transactionHash,omitempty"`
	Error           *SendRawTransactionError `json:"error,omitempty"`
	Message         string                   `json:"message,omitempty"`
}

type EstimateGasControllerRequest struct {
	From     string `json:"from" binding:"required"`
	To       string `json:"to,omitempty"`
	Gas      string `json:"gas,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	Value    string `json:"value,omitempty"`
	Data     string `json:"data,omitempty"`
	Nonce    string `json:"nonce,omitempty"`
}

type EstimateGasControllerResponse struct {
	Success      bool                     `json:"success"`
	EstimatedGas string                   `json:"estimatedGas,omitempty"`
	Error        *SendRawTransactionError `json:"error,omitempty"`
	Message      string                   `json:"message,omitempty"`
}

type GetTransactionCountControllerRequest struct {
	Address        string `json:"address" binding:"required"`
	BlockParameter string `json:"blockParameter,omitempty"`
}

type GetTransactionCountControllerResponse struct {
	Success          bool                     `json:"success"`
	TransactionCount string                   `json:"transactionCount,omitempty"`
	Error            *SendRawTransactionError `json:"error,omitempty"`
	Message          string                   `json:"message,omitempty"`
}
