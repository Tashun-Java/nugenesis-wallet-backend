package utils

import (
	"fmt"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/Blockchain_info/blockchain_info_models"
	"github.com/tashunc/nugenesis-wallet-backend/external/models"
	"time"
)

// # Utility functions for the project (e.g., date formatting, error handling)
func MapBlockchainInfoTxToTransaction(tx blockchain_info_models.Tx, userAddress string, txID string) models.Transaction {
	// Convert timestamp to readable date/time
	txTime := time.Unix(tx.Time, 0)

	// Determine if this is a send or receive transaction
	txType := "receive"
	var relevantAddress string
	var amount int64

	// Check if user is sending (address in inputs)
	for _, input := range tx.Inputs {
		if input.PrevOut.Addr == userAddress {
			txType = "send"
			// Find the output that's not the user's address (the recipient)
			for _, output := range tx.Out {
				if output.Addr != userAddress {
					relevantAddress = output.Addr
					amount = output.Value
					break
				}
			}
			break
		}
	}

	// If not sending, check if receiving
	if txType == "receive" {
		for _, output := range tx.Out {
			if output.Addr == userAddress {
				relevantAddress = userAddress
				amount = output.Value
				// Find sender address from inputs
				if len(tx.Inputs) > 0 {
					relevantAddress = tx.Inputs[0].PrevOut.Addr
				}
				break
			}
		}
	}

	// Convert satoshis to BTC (assuming Bitcoin)
	btcAmount := float64(amount) / 100000000.0

	// Calculate approximate USD value (you'd want to use real exchange rates)
	usdValue := btcAmount * 50000.0 // Example rate

	// Convert fee from satoshis to USD
	feeInBTC := float64(tx.Fee) / 100000000.0
	feeInUSD := feeInBTC * 50000.0

	// Determine status based on block confirmation
	status := "pending"
	if tx.BlockHeight > 0 {
		status = "completed"
	}
	if tx.DoubleSpend {
		status = "failed"
	}

	return models.Transaction{
		ID:      txID,
		Type:    txType,
		Status:  status,
		Token:   "BTC", // You might want to make this configurable
		Amount:  fmt.Sprintf("%.8f", btcAmount),
		Value:   fmt.Sprintf("$%.2f", usdValue),
		Address: truncateAddress(relevantAddress),
		Date:    txTime.Format("2006-01-02"),
		Time:    txTime.Format("15:04"),
		Fee:     fmt.Sprintf("$%.2f", feeInUSD),
		Hash:    truncateHash(tx.Hash),
	}
}

// Helper function to truncate address for display
func truncateAddress(addr string) string {
	if len(addr) <= 20 {
		return addr
	}
	return addr[:10] + "..." + addr[len(addr)-6:]
}

// Helper function to truncate hash for display
func truncateHash(hash string) string {
	if len(hash) <= 12 {
		return hash
	}
	return hash[:6] + "..." + hash[len(hash)-4:]
}

// MapTxSliceToTransactionHistory converts multiple Tx to TransactionHistory
//func MapTxSliceToTransactionHistory(txs []Tx, userAddress string) TransactionHistory {
//	transactions := make([]Transaction, len(txs))
//
//	for i, tx := range txs {
//		transactions[i] = MapTxToTransaction(tx, userAddress, i+1)
//	}
//
//	return TransactionHistory{
//		ActiveTab:    "All",
//		SearchQuery:  "",
//		Transactions: transactions,
//	}
//}
