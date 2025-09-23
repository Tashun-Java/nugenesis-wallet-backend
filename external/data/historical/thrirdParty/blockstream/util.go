package blockstream

import (
	"fmt"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/blockstream/blockstream_models"
	"regexp"
	"strings"
	"time"
)

func ValidateBitcoinAddress(address string) bool {
	if address == "" {
		return false
	}

	p2pkhPattern := regexp.MustCompile(`^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$`)
	bech32Pattern := regexp.MustCompile(`^bc1[a-z0-9]{39,59}$`)

	return p2pkhPattern.MatchString(address) || bech32Pattern.MatchString(address)
}

func FormatSatoshisToBTC(satoshis int64) float64 {
	return float64(satoshis) / 100000000
}

func TruncateAddress(addr string) string {
	if len(addr) <= 20 {
		return addr
	}
	return addr[:10] + "..." + addr[len(addr)-6:]
}

func IsSegwitAddress(address string) bool {
	return strings.HasPrefix(address, "bc1") || strings.HasPrefix(address, "tb1")
}

func MapToStandardizedTransactions(rawTransactions *blockstream_models.AddressTransactionsResponse, targetAddress string) *blockstream_models.StandardizedTransactionsResponse {
	if rawTransactions == nil {
		return &blockstream_models.StandardizedTransactionsResponse{
			Transactions: []blockstream_models.StandardizedTransaction{},
			TotalCount:   0,
		}
	}

	var standardizedTxs []blockstream_models.StandardizedTransaction

	for _, tx := range *rawTransactions {
		standardizedTx := mapSingleTransaction(tx, targetAddress)
		standardizedTxs = append(standardizedTxs, standardizedTx)
	}

	return &blockstream_models.StandardizedTransactionsResponse{
		Transactions: standardizedTxs,
		TotalCount:   len(standardizedTxs),
	}
}

func mapSingleTransaction(tx blockstream_models.TransactionResponse, targetAddress string) blockstream_models.StandardizedTransaction {
	txType := determineTxType(tx, targetAddress)
	status := "completed"
	if !tx.Status.Confirmed {
		status = "pending"
	}

	amount, recipientAddr := calculateAmountAndAddress(tx, targetAddress, txType)

	date, timeStr := formatDateTime(tx.Status.BlockTime)

	fee := fmt.Sprintf("$%.2f", calculateFeeInUSD(tx.Fee))

	value := fmt.Sprintf("$%.2f", amount*getBTCPriceUSD())

	return blockstream_models.StandardizedTransaction{
		ID:      tx.Txid,
		Type:    txType,
		Status:  status,
		Token:   "BTC",
		Amount:  fmt.Sprintf("%.8f", amount),
		Value:   value,
		Address: TruncateAddress(recipientAddr),
		Date:    date,
		Time:    timeStr,
		Fee:     fee,
		Hash:    TruncateAddress(tx.Txid),
	}
}

func determineTxType(tx blockstream_models.TransactionResponse, targetAddress string) string {
	isReceive := false

	for _, output := range tx.Vout {
		if output.ScriptpubkeyAddress == targetAddress {
			isReceive = true
			break
		}
	}

	if isReceive {
		return "receive"
	}
	return "send"
}

func calculateAmountAndAddress(tx blockstream_models.TransactionResponse, targetAddress string, txType string) (float64, string) {
	var amount int64 = 0
	var address string = ""

	if txType == "receive" {
		for _, output := range tx.Vout {
			if output.ScriptpubkeyAddress == targetAddress {
				amount += output.Value
				address = targetAddress
			}
		}
	} else {
		for _, output := range tx.Vout {
			if output.ScriptpubkeyAddress != targetAddress && output.ScriptpubkeyAddress != "" {
				amount += output.Value
				address = output.ScriptpubkeyAddress
				break
			}
		}
	}

	return FormatSatoshisToBTC(amount), address
}

func formatDateTime(blockTime *int64) (string, string) {
	if blockTime == nil {
		now := time.Now()
		return now.Format("2006-01-02"), now.Format("15:04")
	}

	t := time.Unix(*blockTime, 0)
	return t.Format("2006-01-02"), t.Format("15:04")
}

func calculateFeeInUSD(feeInSatoshis int) float64 {
	feeBTC := FormatSatoshisToBTC(int64(feeInSatoshis))
	return feeBTC * getBTCPriceUSD()
}

func getBTCPriceUSD() float64 {
	return 45000.0
}
