package helius_models

//import (
//	"fmt"
//	"time"
//)

// Define the AccountData struct
type AccountData struct {
	Account             string               `json:"account"`
	NativeBalanceChange int64                `json:"nativeBalanceChange"`
	TokenBalanceChanges []TokenBalanceChange `json:"tokenBalanceChanges"`
}

// Define the TokenBalanceChange struct for token transfers
type TokenBalanceChange struct {
	Mint        string `json:"mint"`
	TokenAmount struct {
		Decimals    int64  `json:"decimals"`
		TokenAmount string `json:"tokenAmount"`
	} `json:"rawTokenAmount"`
	TokenAccount string `json:"tokenAccount"`
	UserAccount  string `json:"userAccount"`
}

// Define the InnerInstruction struct
type InnerInstruction struct {
	Accounts  []string `json:"accounts"`
	Data      string   `json:"data"`
	ProgramID string   `json:"programId"`
}

// Define the Instruction struct
type Instruction struct {
	Accounts          []string           `json:"accounts"`
	Data              string             `json:"data"`
	InnerInstructions []InnerInstruction `json:"innerInstructions"`
	ProgramID         string             `json:"programId"`
}

// Define the NativeTransfer struct for native SOL transfers
type NativeTransfer struct {
	Amount          int64  `json:"amount"`
	FromUserAccount string `json:"fromUserAccount"`
	ToUserAccount   string `json:"toUserAccount"`
}

// Define the Transaction struct
type Transaction struct {
	AccountData      []AccountData          `json:"accountData"`
	Description      string                 `json:"description"`
	Events           map[string]interface{} `json:"events"`
	Fee              int64                  `json:"fee"`
	FeePayer         string                 `json:"feePayer"`
	Instructions     []Instruction          `json:"instructions"`
	NativeTransfers  []NativeTransfer       `json:"nativeTransfers"`
	Signature        string                 `json:"signature"`
	Slot             int64                  `json:"slot"`
	Source           string                 `json:"source"`
	Timestamp        int64                  `json:"timestamp"`
	TokenTransfers   []interface{}          `json:"tokenTransfers"`
	TransactionError interface{}            `json:"transactionError"`
	Type             string                 `json:"type"`
}
