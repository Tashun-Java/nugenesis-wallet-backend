package alchemy_sepolia

import (
	"fmt"
	"github.com/tashunc/nugenesis-wallet-backend/internal/data/historical/thrirdParty/etherscan/etherscan_models"
	"github.com/tashunc/nugenesis-wallet-backend/internal/models"
	"strconv"
	"time"
)

func MapTxToTransaction(tx etherscan_models.TxEntry, userAddress string) models.Transaction {
	timestamp, _ := strconv.ParseInt(tx.TimeStamp, 10, 64)
	txTime := time.Unix(timestamp, 0)

	valueWei, _ := strconv.ParseFloat(tx.Value, 64)
	ethAmount := valueWei / 1e18
	usdValue := ethAmount * 3500 // Example ETH to USD rate

	feeWei, _ := strconv.ParseFloat(tx.GasPrice, 64)
	gasUsed, _ := strconv.ParseFloat(tx.GasUsed, 64)
	feeETH := (feeWei * gasUsed) / 1e18
	feeUSD := feeETH * 3500

	txType := "receive"
	relevantAddress := tx.From
	if tx.From == userAddress {
		txType = "send"
		relevantAddress = tx.To
	}

	status := "completed"
	if tx.IsError == "1" {
		status = "failed"
	}

	return models.Transaction{
		ID:      tx.Hash,
		Type:    txType,
		Status:  status,
		Token:   "ETH",
		Amount:  fmt.Sprintf("%.6f", ethAmount),
		Value:   fmt.Sprintf("$%.2f", usdValue),
		Address: truncateAddress(relevantAddress),
		Date:    txTime.Format("2006-01-02"),
		Time:    txTime.Format("15:04"),
		Fee:     fmt.Sprintf("$%.2f", feeUSD),
		Hash:    truncateHash(tx.Hash),
	}
}

func truncateAddress(address string) string {
	if len(address) <= 10 {
		return address
	}
	return address[:6] + "..." + address[len(address)-4:]
}

func truncateHash(hash string) string {
	if len(hash) <= 12 {
		return hash
	}
	return hash[:8] + "..." + hash[len(hash)-4:]
}
