package alchemy

import (
	"fmt"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/alchemy/alchemy_models"
	"github.com/tashunc/nugenesis-wallet-backend/external/models"
	"strconv"
	"strings"
	"time"
)

//func MapTokensToPortfolio(response *alchemy_models.TokensByAddressResponse) []models.Token {
//	var tokens []models.Token
//
//	for _, data := range response.Data {
//		for _, tokenInfo := range data.Tokens {
//			token := models.Token{
//				Symbol:    tokenInfo.Symbol,
//				Name:      tokenInfo.Name,
//				Balance:   formatBalance(tokenInfo.Balance, tokenInfo.Decimals),
//				Value:     fmt.Sprintf("$%.2f", tokenInfo.BalanceUSD),
//				Logo:      tokenInfo.Logo,
//				Address:   tokenInfo.ContractAddress,
//				Network:   data.Network,
//				TokenType: tokenInfo.TokenType,
//				Decimals:  tokenInfo.Decimals,
//			}
//			tokens = append(tokens, token)
//		}
//	}
//
//	return tokens
//}

func formatBalance(balance string, decimals int) string {
	if balance == "" || balance == "0" {
		return "0"
	}

	balanceFloat, err := strconv.ParseFloat(balance, 64)
	if err != nil {
		return balance
	}

	divisor := 1.0
	for i := 0; i < decimals; i++ {
		divisor *= 10
	}

	actualBalance := balanceFloat / divisor

	if actualBalance < 0.0001 {
		return fmt.Sprintf("%.8f", actualBalance)
	} else if actualBalance < 1 {
		return fmt.Sprintf("%.6f", actualBalance)
	} else {
		return fmt.Sprintf("%.4f", actualBalance)
	}
}

func truncateAddress(addr string) string {
	if len(addr) <= 20 {
		return addr
	}
	return addr[:10] + "..." + addr[len(addr)-6:]
}

func ValidateEthereumAddress(address string) bool {
	if !strings.HasPrefix(address, "0x") {
		return false
	}

	if len(address) != 42 {
		return false
	}

	hexPart := address[2:]
	for _, char := range hexPart {
		if !isHexChar(char) {
			return false
		}
	}

	return true
}

func isHexChar(char rune) bool {
	return (char >= '0' && char <= '9') ||
		(char >= 'a' && char <= 'f') ||
		(char >= 'A' && char <= 'F')
}

func ValidateNetwork(network string) bool {
	validNetworks := []string{
		"eth-mainnet",
		"eth-goerli",
		"eth-sepolia",
		"polygon-mainnet",
		"polygon-mumbai",
		"arbitrum-mainnet",
		"arbitrum-goerli",
		"optimism-mainnet",
		"optimism-goerli",
	}

	for _, validNetwork := range validNetworks {
		if network == validNetwork {
			return true
		}
	}

	return false
}

func truncateHash(hash string) string {
	if len(hash) <= 12 {
		return hash
	}
	return hash[:8] + "..." + hash[len(hash)-4:]
}

func pow10(n int) int64 {
	result := int64(1)
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}

func MapAssetTransferToTransaction(transfer alchemy_models.AssetTransfer, userAddress string) models.Transaction {
	// Parse block number and metadata
	var txTime time.Time
	if transfer.Metadata != nil && transfer.Metadata.BlockTimestamp != "" {
		txTime, _ = time.Parse(time.RFC3339, transfer.Metadata.BlockTimestamp)
	} else {
		txTime = time.Now()
	}

	// Determine transaction type (send or receive)
	txType := "receive"
	relevantAddress := transfer.From
	if strings.EqualFold(transfer.From, userAddress) {
		txType = "send"
		if transfer.To != nil {
			relevantAddress = *transfer.To
		}
	}

	// Determine token and amount
	token := "ETH"
	var amount float64
	category := transfer.Category

	if transfer.Asset != nil && *transfer.Asset != "" {
		token = *transfer.Asset
	}

	if transfer.Value != nil {
		amount = *transfer.Value
	} else if transfer.RawContract.Value != nil {
		// Parse hex value
		if valueStr := *transfer.RawContract.Value; strings.HasPrefix(valueStr, "0x") {
			if val, err := strconv.ParseInt(valueStr[2:], 16, 64); err == nil {
				// Get decimals
				decimals := 18
				if transfer.RawContract.Decimal != nil {
					if d, err := strconv.Atoi(*transfer.RawContract.Decimal); err == nil {
						decimals = d
					} else if strings.HasPrefix(*transfer.RawContract.Decimal, "0x") {
						if d, err := strconv.ParseInt((*transfer.RawContract.Decimal)[2:], 16, 64); err == nil {
							decimals = int(d)
						}
					}
				}
				amount = float64(val) / float64(pow10(decimals))
			}
		}
	}

	// Handle NFTs
	if transfer.Category == "erc721" || transfer.Category == "erc1155" {
		if transfer.Erc721TokenId != nil {
			token = fmt.Sprintf("%s #%s", token, *transfer.Erc721TokenId)
		}
		amount = 1
	}

	// Estimate USD value (simplified, in production you'd use real-time prices)
	usdValue := amount * 3500 // Example conversion rate

	// Status is always completed for confirmed transfers
	status := "completed"

	// Fee information is not provided in asset transfers
	fee := "$0.00"

	return models.Transaction{
		ID:        transfer.Hash,
		Type:      txType,
		Category:  category,
		Status:    status,
		Token:     token,
		Amount:    fmt.Sprintf("%.6f", amount),
		Value:     fmt.Sprintf("$%.2f", usdValue),
		Address:   truncateAddress(relevantAddress),
		ToAddress: truncateAddress(getToAddress(transfer)),
		Date:      txTime.Format("2006-01-02"),
		Time:      txTime.Format("15:04"),
		Fee:       fee,
		Hash:      truncateHash(transfer.Hash),
	}
}

func getToAddress(transfer alchemy_models.AssetTransfer) string {
	if transfer.To != nil {
		return *transfer.To
	}
	return ""
}
