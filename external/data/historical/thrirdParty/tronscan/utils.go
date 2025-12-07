package tronscan

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/tronscan/tronscan_models"
	"github.com/tashunc/nugenesis-wallet-backend/external/models"
)

const (
	// Tron uses 6 decimal places for TRX (1 TRX = 1,000,000 SUN)
	TRXDecimals = 6
	TRXDivisor  = 1000000
)

// MapToStandardTransaction converts a TronScan transaction to standard format
func MapToStandardTransaction(tx tronscan_models.Transaction, walletAddress string) models.Transaction {
	// Determine transaction type and relevant address
	txType := "unknown"
	relevantAddress := ""
	var amount float64
	tokenSymbol := "TRX"

	// Normalize addresses for comparison (Tron addresses are case-sensitive but we normalize)
	walletAddrLower := strings.ToLower(walletAddress)
	fromAddrLower := strings.ToLower(tx.OwnerAddress)
	toAddrLower := ""

	if tx.ToAddress != "" {
		toAddrLower = strings.ToLower(tx.ToAddress)
	} else if tx.ContractData.ToAddress != "" {
		toAddrLower = strings.ToLower(tx.ContractData.ToAddress)
	}

	// Determine transaction type based on addresses
	if fromAddrLower == walletAddrLower && toAddrLower == walletAddrLower {
		txType = "self"
		relevantAddress = truncateAddress(tx.ToAddress)
	} else if fromAddrLower == walletAddrLower {
		txType = "send"
		relevantAddress = truncateAddress(tx.ToAddress)
		if relevantAddress == "" {
			relevantAddress = truncateAddress(tx.ContractData.ToAddress)
		}
	} else if toAddrLower == walletAddrLower {
		txType = "receive"
		relevantAddress = truncateAddress(tx.OwnerAddress)
	} else {
		// Contract interaction where wallet is not directly involved
		txType = "contract_interaction"
		if tx.ContractData.ContractAddress != "" {
			relevantAddress = truncateAddress(tx.ContractData.ContractAddress)
		}
	}

	// Get amount and token information
	if tx.TokenInfo != nil {
		// Token transfer (TRC10 or TRC20)
		tokenSymbol = tx.TokenInfo.TokenAbbr

		// For TRC20 tokens, try to get amount from trigger_info parameter
		if tx.TriggerInfo != nil && tx.TriggerInfo.Parameter != nil {
			if valueParam, ok := tx.TriggerInfo.Parameter["value"].(string); ok {
				// Parse the hex/decimal value from trigger info
				if parsedValue, err := strconv.ParseFloat(valueParam, 64); err == nil {
					decimals := tx.TokenInfo.TokenDecimal
					if decimals == 0 {
						decimals = 6 // Default to 6 decimals
					}
					divisor := math.Pow10(decimals)
					amount = parsedValue / divisor
				}
			}
		} else if tx.ContractData.Amount > 0 {
			// TRC10 token or fallback for TRC20
			decimals := tx.TokenInfo.TokenDecimal
			if decimals == 0 {
				decimals = 6 // Default to 6 decimals
			}
			divisor := math.Pow10(decimals)
			amount = float64(tx.ContractData.Amount) / divisor
		}
	} else if tx.ContractData.Amount > 0 {
		// Native TRX transfer or TRC10 without token info
		amount = float64(tx.ContractData.Amount) / TRXDivisor
		tokenSymbol = "TRX"
	} else if tx.ContractData.CallValue > 0 {
		// Smart contract call with TRX value
		amount = float64(tx.ContractData.CallValue) / TRXDivisor
		tokenSymbol = "TRX"
	}

	// Calculate fee
	var feeAmount float64
	if tx.Cost != nil {
		// Total fee = energy fee + net fee (in SUN, convert to TRX)
		totalFeeSun := tx.Cost.EnergyFee + tx.Cost.NetFee
		feeAmount = float64(totalFeeSun) / TRXDivisor
	}

	// Determine status
	status := "completed"
	if !tx.Confirmed {
		status = "pending"
	}
	if tx.Revert {
		status = "failed"
	}

	// Determine category
	category := "transfer"
	switch tx.ContractType {
	case 1:
		category = "transfer" // TransferContract
	case 2:
		category = "transfer" // TransferAssetContract (TRC10)
	case 31:
		category = "contract_interaction" // TriggerSmartContract
	default:
		if tx.ContractData.ContractAddress != "" || tx.TriggerInfo != nil {
			category = "contract_interaction"
		}
	}

	// Convert timestamp from milliseconds to time.Time
	timestamp := time.Unix(tx.Timestamp/1000, 0)

	return models.Transaction{
		ID:        tx.Hash,
		Type:      txType,
		Category:  category,
		Status:    status,
		Token:     tokenSymbol,
		Amount:    fmt.Sprintf("%.6f", amount),
		Value:     fmt.Sprintf("$%.2f", amount*0.10), // Placeholder price, should use real rates
		Address:   relevantAddress,
		ToAddress: truncateAddress(tx.ToAddress),
		Date:      timestamp.Format("2006-01-02"),
		Time:      timestamp.Format("15:04"),
		Fee:       fmt.Sprintf("%.6f", feeAmount),
		Hash:      truncateHash(tx.Hash),
	}
}

// Helper function to truncate Tron address for display
func truncateAddress(addr string) string {
	if len(addr) <= 20 {
		return addr
	}
	// Tron addresses are typically 34 characters (starting with T)
	if len(addr) >= 34 {
		return addr[:8] + "..." + addr[len(addr)-6:]
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

// ValidateTronAddress validates a Tron address format
func ValidateTronAddress(address string) bool {
	// Tron addresses start with 'T' and are 34 characters long (base58)
	if len(address) != 34 {
		return false
	}
	if !strings.HasPrefix(address, "T") {
		return false
	}
	// Basic validation - could be enhanced with base58 checksum validation
	return true
}

// MapTokenBalanceToStandard converts a TronScan token balance to standard wallet token balance format
func MapTokenBalanceToStandard(token tronscan_models.TokenBalance, address string, tokenIDService TokenIDServiceInterface) models.WalletTokenBalance {
	// Parse balance
	balance := token.Balance
	if balance == "" {
		// If balance is empty, use amount if available
		if token.Amount.Float64() > 0 {
			balance = fmt.Sprintf("%.6f", token.Amount.Float64())
		} else {
			balance = "0"
		}
	}

	// Determine if this is native TRX
	isNativeToken := false
	tokenAddress := token.TokenID
	symbol := token.TokenAbbr
	name := token.TokenName

	// Parse USD price and value
	usdPrice := token.TokenPriceInUsd
	usdValue := token.AmountInUsd

	// Calculate portfolio percentage if needed (would need total portfolio value)
	portfolioPercentage := 0.0

	// Lookup token ID from service
	tokenID := ""
	if tokenIDService != nil {
		tokenID = tokenIDService.GetTokenID("tron", tokenAddress)
	}

	return models.WalletTokenBalance{
		TokenAddress:        tokenAddress,
		TokenID:             tokenID,
		Name:                name,
		Symbol:              symbol,
		Logo:                token.TokenLogo,
		Thumbnail:           token.TokenLogo,
		Decimals:            fmt.Sprintf("%d", token.TokenDecimal),
		Balance:             balance,
		BalanceRaw:          token.Balance,
		NativeToken:         isNativeToken,
		VerifiedContract:    token.TokenCanShow == 1,
		PossibleSpam:        false,
		UsdPrice:            usdPrice,
		UsdValue:            usdValue,
		UsdPrice24hrChange:  0, // TronScan doesn't provide 24h change
		PortfolioPercentage: portfolioPercentage,
		SecurityScore:       0,
		Chain:               "tron",
	}
}

// CreateNativeTRXBalance creates a native TRX balance entry
func CreateNativeTRXBalance(trxBalance int64, usdPrice float64, tokenIDService TokenIDServiceInterface) models.WalletTokenBalance {
	// Convert from SUN to TRX (1 TRX = 1,000,000 SUN)
	balanceInTRX := float64(trxBalance) / TRXDivisor
	balanceStr := fmt.Sprintf("%.6f", balanceInTRX)

	// Calculate USD value
	usdValue := balanceInTRX * usdPrice

	// Lookup token ID for native TRX
	tokenID := ""
	if tokenIDService != nil {
		tokenID = tokenIDService.GetTokenIDForNative("tron", "TRX")
	}

	return models.WalletTokenBalance{
		TokenAddress:        "_", // Special marker for native token
		TokenID:             tokenID,
		Name:                "Tronix",
		Symbol:              "TRX",
		Logo:                "",
		Thumbnail:           "",
		Decimals:            "6",
		Balance:             balanceStr,
		BalanceRaw:          fmt.Sprintf("%d", trxBalance),
		NativeToken:         true,
		VerifiedContract:    true,
		PossibleSpam:        false,
		UsdPrice:            usdPrice,
		UsdValue:            usdValue,
		UsdPrice24hrChange:  0,
		PortfolioPercentage: 0,
		SecurityScore:       100,
		Chain:               "tron",
	}
}
