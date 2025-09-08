package moralis

import (
	"fmt"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/moralis/moralis_models"
	"github.com/tashunc/nugenesis-wallet-backend/external/models"
	"math/big"
	"time"
)

// MapHistoryToTransaction converts Moralis history transaction to standard transaction format
func MapHistoryToTransaction(tx moralis_models.HistoryTransaction, walletAddress string) models.Transaction {
	timestamp, _ := time.Parse(time.RFC3339, tx.BlockTimestamp)

	// Determine transaction type and relevant address
	txType := "unknown"
	relevantAddress := ""
	if tx.FromAddress == walletAddress {
		txType = "send"
		relevantAddress = truncateAddress(tx.ToAddress)
	} else if tx.ToAddress == walletAddress {
		txType = "receive"
		relevantAddress = truncateAddress(tx.FromAddress)
	}

	// Convert Wei to ETH (assuming 18 decimals for native token)
	valueWei := new(big.Int)
	valueWei.SetString(tx.Value, 10)
	valueFlt := new(big.Float).SetInt(valueWei)
	ethDivisor := new(big.Float).SetFloat64(1e18)
	ethValue := new(big.Float).Quo(valueFlt, ethDivisor)
	ethAmount, _ := ethValue.Float64()

	// Calculate fee in ETH
	//gasPrice, _ := strconv.ParseInt(tx.GasPrice, 10, 64)
	//gasUsed, _ := strconv.ParseInt(tx.GasUsed, 10, 64)
	//feeWei := big.NewInt(gasPrice * gasUsed)
	//feeFlt := new(big.Float).SetInt(feeWei)
	//feeAmount, _ := feeEth.Float64()

	// Determine status
	status := "pending"
	if tx.ReceiptStatus == "1" {
		status = "completed"
	} else if tx.ReceiptStatus == "0" {
		status = "failed"
	}

	// Determine category based on transaction data
	category := "transfer"
	if len(tx.InputData) > 2 && tx.InputData != "0x" {
		category = "contract_interaction"
	}

	return models.Transaction{
		ID:        tx.Hash,
		Type:      txType,
		Category:  category,
		Status:    status,
		Token:     "MATIC",
		Amount:    fmt.Sprintf("%.6f", ethAmount),
		Value:     fmt.Sprintf("$%.2f", ethAmount*2.0), // Example price, should use real rates
		Address:   relevantAddress,
		ToAddress: truncateAddress(tx.ToAddress),
		Date:      timestamp.Format("2006-01-02"),
		Time:      timestamp.Format("15:04"),
		Hash:      truncateHash(tx.Hash),
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
