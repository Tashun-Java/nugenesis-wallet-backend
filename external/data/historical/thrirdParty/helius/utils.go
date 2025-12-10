package helius

import (
	"fmt"
	"sync"
	"time"

	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/helius/helius_models"
	"github.com/tashunc/nugenesis-wallet-backend/external/models"
	"github.com/tashunc/nugenesis-wallet-backend/static/staticServices"
)

var programMapping = map[string]string{
	"11111111111111111111111111111111":            "native_transfer", // System Program
	"TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA": "token_transfer",  // SPL Token Program
	"Stake11111111111111111111111111111111111111": "stake",           // Stake Program
	"Vote111111111111111111111111111111111111111": "vote",            // Vote Program
	"ComputeBudget111111111111111111111111111111": "contract_call",   // Compute Budget Program
	"MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr": "memo",            // Memo Program
}

var (
	// AssetService instance for accessing Solana token symbols
	assetService     *staticServices.AssetService
	assetServiceOnce sync.Once
)

// initAssetService initializes the AssetService singleton
func initAssetService() {
	assetService = staticServices.NewAssetService()
}

// getTokenSymbol returns the token symbol for a given mint address
// Uses the static AssetService cache from the static folder
func getTokenSymbol(mint string) string {
	// Ensure AssetService is initialized (happens only once)
	assetServiceOnce.Do(initAssetService)

	// Use AssetService to get the token symbol
	return assetService.GetTokenSymbolByMint(mint)
}

func MapTxToTransaction(tx helius_models.Transaction, address string) []models.Transaction {
	// Extract the necessary fields from the transaction data
	signature := tx.Signature
	fee := float64(tx.Fee) / 1e9 // Convert lamports to SOL (1 SOL = 1e9 lamports)
	feePayer := tx.FeePayer
	timestamp := tx.Timestamp

	// Convert timestamp to human-readable date and time
	date, timeFormatted := formatTimestamp(timestamp)

	// Default values
	transactionType := "unknown"
	token := "SOL"
	amount := "0" // Since no native transfer or token transfer
	value := "0"  // Same as amount
	//address := feePayer

	// Check if there are instructions and map the programId to type
	for _, instruction := range tx.Instructions {
		programID := instruction.ProgramID

		// Set the transaction type based on programID
		switch programID {
		case "11111111111111111111111111111111":
			transactionType = "native_transfer" // System Program
		case "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA":
			transactionType = "token_transfer" // SPL Token Program
		case "Stake11111111111111111111111111111111111111":
			transactionType = "stake" // Stake Program
		case "Vote111111111111111111111111111111111111111":
			transactionType = "vote" // Vote Program
		case "ComputeBudget111111111111111111111111111111":
			transactionType = "contract_call" // Compute Budget Program
		default:
			transactionType = "program_interaction" // Custom Program
		}
		//}
		//
		//// Convert timestamp to human-readable date and time
		//t := time.Unix(timestamp, 0)
		//date := t.Format("2006-01-02")
		//timeFormatted := t.Format("15:04:05")
		//
		//// Default values
		//transactionType := "unknown"
		//token := "SOL"
		//amount := "0" // Since no native transfer or token transfer
		//value := "0"  // Same as amount
		//address := feePayer

		// Check if there are instructions and map the programId to type
		//instructions := tx["instructions"].([]interface{})
		//for _, instruction := range instructions {
		//	programID := instruction.(map[string]interface{})["programId"].(string)
		//
		//	// Set the transaction type based on programID
		//	if txType, exists := programMapping[programID]; exists {
		//		transactionType = txType
		//	}
	}

	// Return the mapped transaction
	var mappedTransactions []models.Transaction
	if len(tx.NativeTransfers) == 0 {
		mappedTransactions = append(mappedTransactions, models.Transaction{
			ID:       signature,
			Type:     "send",
			Category: transactionType,
			Status:   "success", // assuming no error in the transaction
			Token:    token,
			Amount:   amount,
			Value:    value,
			Address:  feePayer,
			Date:     date,
			Time:     timeFormatted,
			Fee:      fmt.Sprintf("%.9f", fee),
			Hash:     signature,
		})
	} else {

		for _, transfer := range tx.NativeTransfers {
			if address == transfer.FromUserAccount {
				transactionType = "send"
			} else if address == transfer.ToUserAccount {
				transactionType = "receive"
			} else {
				continue
			}
			solAmount := float64(transfer.Amount) / 1e9
			mappedTransactions = append(mappedTransactions, models.Transaction{
				ID:        signature,
				Type:      transactionType,
				Category:  "native_transfer",
				Status:    "success", // assuming no error in the transaction
				Token:     token,
				Amount:    fmt.Sprintf("%.9f", solAmount),
				Value:     fmt.Sprintf("%.9f", solAmount),
				Address:   transfer.FromUserAccount, // Fixed: use actual sender, not feePayer
				ToAddress: transfer.ToUserAccount,
				Date:      date,
				Time:      timeFormatted,
				Fee:       fmt.Sprintf("%.9f", fee),
				Hash:      signature,
			})
		}

	}

	if len(tx.TokenTransfers) > 0 {
		for _, transfer := range tx.TokenTransfers {
			// Skip if address is not involved
			if address != transfer.FromUserAccount && address != transfer.ToUserAccount {
				continue
			}

			var transactionType string
			if address == transfer.FromUserAccount {
				transactionType = "send"
			} else {
				transactionType = "receive"
			}

			tokenMint := transfer.Mint
			tokenSymbol := getTokenSymbol(tokenMint)
			tokenAmount := transfer.TokenAmount

			mappedTransactions = append(mappedTransactions, models.Transaction{
				ID:        signature,
				Type:      transactionType,
				Category:  "token_transfer",
				Status:    "success",
				Token:     tokenSymbol, // Now shows "USDC" instead of mint
				Amount:    fmt.Sprintf("%.9f", tokenAmount),
				Value:     fmt.Sprintf("%.9f", tokenAmount),
				Address:   transfer.FromUserAccount,
				ToAddress: transfer.ToUserAccount,
				Date:      date,
				Time:      timeFormatted,
				Fee:       fmt.Sprintf("%.9f", fee),
				Hash:      signature,
			})
		}
	}

	return mappedTransactions

}

func formatTimestamp(timestamp int64) (string, string) {
	t := time.Unix(timestamp, 0)
	date := t.Format("2006-01-02")        // Format date as YYYY-MM-DD
	timeFormatted := t.Format("15:04:05") // Format time as HH:MM:SS
	return date, timeFormatted
}
