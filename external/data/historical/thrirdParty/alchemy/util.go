package alchemy

import (
	"fmt"
	"strconv"
	"strings"
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
