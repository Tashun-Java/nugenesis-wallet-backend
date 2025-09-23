package moralis

import (
	"encoding/hex"
	"fmt"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/moralis/moralis_models"
	"github.com/tashunc/nugenesis-wallet-backend/external/models"
	"math/big"
	"strconv"
	"strings"
	"time"
)

// MapHistoryToTransaction converts Moralis history transaction to standard transaction format
func MapHistoryToTransaction(tx moralis_models.HistoryTransaction, walletAddress string) models.Transaction {
	timestamp, _ := time.Parse(time.RFC3339, tx.BlockTimestamp)

	// Determine transaction type and relevant address
	txType := "unknown"
	relevantAddress := ""
	var ethAmount float64
	tokenSymbol := "MATIC" // Default to native token

	// Priority order: Native transfers -> ERC-20 transfers -> NFT transfers -> Category/Summary based -> Legacy parsing

	// 1. Check native transfers first (most reliable)
	if len(tx.NativeTransfers) > 0 {
		ethAmount, tokenSymbol, txType = extractNativeTransferAmountAndType(tx.NativeTransfers, walletAddress)
		if txType == "send" {
			relevantAddress = truncateAddress(tx.ToAddress)
		} else if txType == "receive" {
			relevantAddress = truncateAddress(tx.FromAddress)
		}
	} else if len(tx.ERC20Transfers) > 0 {
		// 2. Check ERC-20 transfers
		ethAmount, tokenSymbol, txType = extractERC20TransferAmountAndType(tx.ERC20Transfers, walletAddress)
		if txType == "send" {
			relevantAddress = truncateAddress(tx.ToAddress)
		} else if txType == "receive" {
			relevantAddress = truncateAddress(tx.FromAddress)
		}
	} else if len(tx.NFTTransfers) > 0 {
		// 3. Check NFT transfers
		ethAmount, tokenSymbol, txType = extractNFTTransferAmountAndType(tx.NFTTransfers, walletAddress)
		if txType == "send" {
			relevantAddress = truncateAddress(tx.ToAddress)
		} else if txType == "receive" {
			relevantAddress = truncateAddress(tx.FromAddress)
		}
	} else if tx.Category != "" {
		// 4. Use Moralis category and summary to determine transaction type
		txType = determineTypeFromCategoryAndSummary(tx.Category, tx.Summary, tx.FromAddress, tx.ToAddress, walletAddress)
		if txType == "send" {
			relevantAddress = truncateAddress(tx.ToAddress)
		} else if txType == "receive" {
			relevantAddress = truncateAddress(tx.FromAddress)
		} else {
			// For contract interactions, use the contract address if available
			if strings.EqualFold(tx.FromAddress, walletAddress) {
				relevantAddress = truncateAddress(tx.ToAddress)
			} else {
				relevantAddress = truncateAddress(tx.FromAddress)
			}
		}
		// For zero value contract interactions, keep amount as 0
		ethAmount = 0
	} else if tx.Value != "" {
		// 5. Legacy: parse native value field directly and determine type manually
		valueWei := new(big.Int)
		if _, ok := valueWei.SetString(tx.Value, 10); ok && valueWei.Cmp(big.NewInt(0)) > 0 {
			valueFlt := new(big.Float).SetInt(valueWei)
			ethDivisor := new(big.Float).SetFloat64(1e18)
			ethValue := new(big.Float).Quo(valueFlt, ethDivisor)
			ethAmount, _ = ethValue.Float64()

			// Manual type determination for legacy parsing
			if strings.EqualFold(tx.FromAddress, walletAddress) {
				txType = "send"
				relevantAddress = truncateAddress(tx.ToAddress)
			} else if strings.EqualFold(tx.ToAddress, walletAddress) {
				txType = "receive"
				relevantAddress = truncateAddress(tx.FromAddress)
			}
		} else {
			// 6. Fallback to checking logs manually
			ethAmount, tokenSymbol = extractTokenTransferAmount(tx.Logs, walletAddress)
			// Manual type determination for log parsing
			if strings.EqualFold(tx.FromAddress, walletAddress) {
				txType = "send"
				relevantAddress = truncateAddress(tx.ToAddress)
			} else if strings.EqualFold(tx.ToAddress, walletAddress) {
				txType = "receive"
				relevantAddress = truncateAddress(tx.FromAddress)
			}
		}
	} else {
		// 7. Final fallback to logs
		ethAmount, tokenSymbol = extractTokenTransferAmount(tx.Logs, walletAddress)
		// Manual type determination for log parsing
		if strings.EqualFold(tx.FromAddress, walletAddress) {
			txType = "send"
			relevantAddress = truncateAddress(tx.ToAddress)
		} else if strings.EqualFold(tx.ToAddress, walletAddress) {
			txType = "receive"
			relevantAddress = truncateAddress(tx.FromAddress)
		}
	}

	// Calculate fee from Moralis transaction_fee field
	var feeAmount float64
	if tx.TransactionFee != "" {
		if fee, err := strconv.ParseFloat(tx.TransactionFee, 64); err == nil {
			feeAmount = fee
		}
	}

	// Determine status
	status := "pending"
	if tx.ReceiptStatus == "1" {
		status = "completed"
	} else if tx.ReceiptStatus == "0" {
		status = "failed"
	}

	// Determine category based on transaction data and Moralis category
	category := "transfer"
	if tx.Category != "" {
		// Use Moralis category if available
		switch strings.ToLower(tx.Category) {
		case "contract interaction":
			category = "contract_interaction"
		case "send", "receive":
			category = "transfer"
		case "nft send", "nft receive":
			category = "nft_transfer"
		default:
			category = strings.ToLower(tx.Category)
		}
	} else if len(tx.InputData) > 2 && tx.InputData != "0x" {
		category = "contract_interaction"
	}

	return models.Transaction{
		ID:        tx.Hash,
		Type:      txType,
		Category:  category,
		Status:    status,
		Token:     tokenSymbol,
		Amount:    fmt.Sprintf("%.6f", ethAmount),
		Value:     fmt.Sprintf("$%.2f", ethAmount*2.0), // Example price, should use real rates
		Address:   relevantAddress,
		ToAddress: truncateAddress(tx.ToAddress),
		Date:      timestamp.Format("2006-01-02"),
		Time:      timestamp.Format("15:04"),
		Fee:       fmt.Sprintf("%.8f", feeAmount),
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

// extractTokenTransferAmount extracts the transfer amount from ERC-20 token transfer logs
func extractTokenTransferAmount(logs []moralis_models.Log, walletAddress string) (float64, string) {
	// ERC-20 Transfer event signature: Transfer(address indexed from, address indexed to, uint256 value)
	// Topic0 for Transfer: 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
	transferEventTopic := "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

	for _, log := range logs {
		// Check if this is a Transfer event
		if log.Topic0 != transferEventTopic {
			continue
		}

		// Extract from and to addresses from topics
		if len(log.Topic1) < 26 || len(log.Topic2) < 26 {
			continue
		}

		fromAddr := "0x" + strings.ToLower(log.Topic1[26:]) // Remove padding
		toAddr := "0x" + strings.ToLower(log.Topic2[26:])   // Remove padding
		walletAddrLower := strings.ToLower(walletAddress)

		// Check if this transfer involves our wallet
		if fromAddr != walletAddrLower && toAddr != walletAddrLower {
			continue
		}

		// Extract amount from data field
		if len(log.Data) < 2 {
			continue
		}

		// Remove 0x prefix and decode hex
		dataHex := log.Data[2:]
		if len(dataHex) == 0 {
			continue
		}

		// Parse the amount (assuming 18 decimals for most tokens)
		amountBytes, err := hex.DecodeString(dataHex)
		if err != nil {
			continue
		}

		amountBig := new(big.Int).SetBytes(amountBytes)
		if amountBig.Cmp(big.NewInt(0)) == 0 {
			continue
		}

		// Convert to float (assuming 18 decimals)
		amountFloat := new(big.Float).SetInt(amountBig)
		divisor := new(big.Float).SetFloat64(1e18)
		result := new(big.Float).Quo(amountFloat, divisor)
		amount, _ := result.Float64()

		// Try to determine token symbol (this is basic - in practice you'd query the contract)
		tokenSymbol := getTokenSymbolFromAddress(log.Address)

		return amount, tokenSymbol
	}

	return 0, "MATIC"
}

// extractNativeTransferAmountAndType extracts amount and transaction type from native transfers (POL/MATIC)
func extractNativeTransferAmountAndType(transfers []moralis_models.NativeTransfer, walletAddress string) (float64, string, string) {
	walletAddrLower := strings.ToLower(walletAddress)

	for _, transfer := range transfers {
		// Check if this transfer involves our wallet
		fromAddr := strings.ToLower(transfer.FromAddress)
		toAddr := strings.ToLower(transfer.ToAddress)

		if fromAddr == walletAddrLower || toAddr == walletAddrLower {
			// Determine transaction type using Direction field from Moralis
			txType := "unknown"
			if transfer.Direction == "send" {
				txType = "send"
			} else if transfer.Direction == "receive" {
				txType = "receive"
			}

			// Use the pre-formatted value from Moralis (much more reliable)
			if transfer.ValueFormatted != "" {
				if amount, err := strconv.ParseFloat(transfer.ValueFormatted, 64); err == nil {
					tokenSymbol := transfer.TokenSymbol
					if tokenSymbol == "" {
						tokenSymbol = "MATIC" // Default fallback
					}
					return amount, tokenSymbol, txType
				}
			}

			// Fallback to parsing raw value
			if transfer.Value != "" {
				amountBig := new(big.Int)
				if _, ok := amountBig.SetString(transfer.Value, 10); ok && amountBig.Cmp(big.NewInt(0)) > 0 {
					amountFloat := new(big.Float).SetInt(amountBig)
					divisor := new(big.Float).SetFloat64(1e18)
					result := new(big.Float).Quo(amountFloat, divisor)
					amount, _ := result.Float64()

					tokenSymbol := transfer.TokenSymbol
					if tokenSymbol == "" {
						tokenSymbol = "MATIC" // Default fallback
					}
					return amount, tokenSymbol, txType
				}
			}
		}
	}

	return 0, "MATIC", "unknown"
}

// extractERC20TransferAmountAndType extracts amount and transaction type from ERC-20 transfers
func extractERC20TransferAmountAndType(transfers []moralis_models.ERC20Transfer, walletAddress string) (float64, string, string) {
	walletAddrLower := strings.ToLower(walletAddress)

	for _, transfer := range transfers {
		// Check if this transfer involves our wallet
		fromAddr := strings.ToLower(transfer.FromAddress)
		toAddr := strings.ToLower(transfer.ToAddress)

		if fromAddr == walletAddrLower || toAddr == walletAddrLower {
			// Determine transaction type using Direction field from Moralis
			txType := "unknown"
			if transfer.Direction == "send" {
				txType = "send"
			} else if transfer.Direction == "receive" {
				txType = "receive"
			}

			// Use the pre-formatted value from Moralis (much more reliable)
			if transfer.ValueFormatted != "" {
				if amount, err := strconv.ParseFloat(transfer.ValueFormatted, 64); err == nil {
					tokenSymbol := transfer.TokenSymbol
					if tokenSymbol == "" {
						tokenSymbol = "TOKEN" // Default fallback
					}
					return amount, tokenSymbol, txType
				}
			}

			// Fallback to parsing raw value with decimals
			if transfer.Value != "" {
				amountBig := new(big.Int)
				if _, ok := amountBig.SetString(transfer.Value, 10); ok && amountBig.Cmp(big.NewInt(0)) > 0 {
					// Use token decimals if available, otherwise default to 18
					decimals := 18
					if transfer.TokenDecimals != "" {
						if d, err := strconv.Atoi(transfer.TokenDecimals); err == nil {
							decimals = d
						}
					}

					amountFloat := new(big.Float).SetInt(amountBig)
					divisor := new(big.Float).SetFloat64(float64(1))
					for i := 0; i < decimals; i++ {
						divisor.Mul(divisor, new(big.Float).SetFloat64(10))
					}
					result := new(big.Float).Quo(amountFloat, divisor)
					amount, _ := result.Float64()

					tokenSymbol := transfer.TokenSymbol
					if tokenSymbol == "" {
						tokenSymbol = "TOKEN" // Default fallback
					}
					return amount, tokenSymbol, txType
				}
			}
		}
	}

	return 0, "TOKEN", "unknown"
}

// extractNFTTransferAmountAndType extracts amount and transaction type from NFT transfers
func extractNFTTransferAmountAndType(transfers []moralis_models.NFTTransfer, walletAddress string) (float64, string, string) {
	walletAddrLower := strings.ToLower(walletAddress)

	for _, transfer := range transfers {
		// Check if this transfer involves our wallet
		fromAddr := strings.ToLower(transfer.FromAddress)
		toAddr := strings.ToLower(transfer.ToAddress)

		if fromAddr == walletAddrLower || toAddr == walletAddrLower {
			// Determine transaction type using Direction field from Moralis
			txType := "unknown"
			if transfer.Direction == "send" {
				txType = "send"
			} else if transfer.Direction == "receive" {
				txType = "receive"
			}

			// For NFTs, we typically show the count/amount rather than value
			amountBig := new(big.Int)
			if _, ok := amountBig.SetString(transfer.Amount, 10); ok {
				amount, _ := amountBig.Float64()

				// Use contract type as token symbol for NFTs
				tokenSymbol := transfer.ContractType
				if tokenSymbol == "" {
					tokenSymbol = "NFT"
				}

				return amount, tokenSymbol, txType
			}
		}
	}

	return 0, "NFT", "unknown"
}

// getTokenSymbolFromAddress returns a token symbol based on contract address
func getTokenSymbolFromAddress(contractAddress string) string {
	// Common Polygon token addresses (simplified mapping)
	knownTokens := map[string]string{
		"0x7d1afa7b718fb893db30a3abc0cfc608aacfebb0": "MATIC",
		"0x2791bca1f2de4661ed88a30c99a7a9449aa84174": "USDC",
		"0xc2132d05d31c914a87c6611c10748aeb04b58e8f": "USDT",
		"0x8f3cf7ad23cd3cadbd9735aff958023239c6a063": "DAI",
		"0x1bfd67037b42cf73acf2047067bd4f2c47d9bfd6": "WBTC",
		"0x7ceb23fd6c6b4db8c56f38e38f33bf1e6df3c3c9": "WETH",
	}

	if symbol, exists := knownTokens[strings.ToLower(contractAddress)]; exists {
		return symbol
	}

	return "TOKEN" // Default fallback
}

// determineTypeFromCategoryAndSummary determines transaction type based on Moralis category and summary fields
func determineTypeFromCategoryAndSummary(category, summary, fromAddress, toAddress, walletAddress string) string {
	walletAddrLower := strings.ToLower(walletAddress)
	fromAddrLower := strings.ToLower(fromAddress)
	toAddrLower := strings.ToLower(toAddress)

	// Handle specific categories from Moralis
	switch strings.ToLower(category) {
	case "send":
		return "send"
	case "receive":
		return "receive"
	case "contract interaction":
		// For contract interactions, determine if it's send/receive based on wallet involvement
		if fromAddrLower == walletAddrLower {
			// Wallet is initiating the interaction - consider it a "send" to the contract
			return "send"
		} else if toAddrLower == walletAddrLower {
			// Wallet is receiving from a contract interaction - consider it a "receive"
			return "receive"
		}
		// If neither from nor to is the wallet, it's still a contract interaction
		return "contract_interaction"
	case "nft send":
		return "send"
	case "nft receive":
		return "receive"
	default:
		// Analyze summary for additional context
		summaryLower := strings.ToLower(summary)

		// Check for common transaction patterns in summary
		if strings.Contains(summaryLower, "sent") || strings.Contains(summaryLower, "transferred") {
			if fromAddrLower == walletAddrLower {
				return "send"
			}
		}

		if strings.Contains(summaryLower, "received") {
			if toAddrLower == walletAddrLower {
				return "receive"
			}
		}

		if strings.Contains(summaryLower, "signed") || strings.Contains(summaryLower, "contract") {
			// Contract interaction initiated by the wallet
			if fromAddrLower == walletAddrLower {
				return "contract_interaction"
			}
		}

		// Fallback: determine by wallet position
		if fromAddrLower == walletAddrLower {
			return "send"
		} else if toAddrLower == walletAddrLower {
			return "receive"
		}
	}

	return "unknown"
}
