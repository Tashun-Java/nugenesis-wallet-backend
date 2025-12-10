package coingecko

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/tashunc/nugenesis-wallet-backend/external/models"
)

// Service handles CoinGecko API interactions
type Service struct {
	client *resty.Client
}

// NewService creates a new CoinGecko service instance
func NewService() *Service {
	return &Service{
		client: resty.New().
			SetHostURL("https://api.coingecko.com/api/v3").
			SetHeader("Accept", "application/json"),
	}
}

// GetPrices fetches crypto token prices from CoinGecko API
func (s *Service) GetPrices(ids, vsCurrencies string) (PriceResponse, error) {
	var result PriceResponse

	resp, err := s.client.R().
		SetQueryParams(map[string]string{
			"ids":           ids,
			"vs_currencies": vsCurrencies,
		}).
		SetResult(&result).
		Get("/simple/price")

	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return result, nil
}

// GetCoinGeckoID returns the CoinGecko ID for a given token symbol
func GetCoinGeckoID(symbol string) string {
	upperSymbol := strings.ToUpper(symbol)
	if id, exists := SymbolToCoinGeckoID[upperSymbol]; exists {
		return id
	}
	// Return lowercase symbol as fallback
	return strings.ToLower(symbol)
}

// EnrichBalancesWithPrices fetches missing prices from CoinGecko and enriches balances
func (s *Service) EnrichBalancesWithPrices(balances []models.WalletTokenBalance) []models.WalletTokenBalance {
	// Collect symbols that need price data
	symbolToIndex := make(map[string][]int) // Map symbol to balance indices

	for i, balance := range balances {
		// Check if price is missing (0 or not set)
		if balance.UsdPrice == 0 && balance.Balance != "0" && balance.Balance != "" {
			symbol := strings.ToUpper(balance.Symbol)
			symbolToIndex[symbol] = append(symbolToIndex[symbol], i)
		}
	}

	// If no prices need to be fetched, return as-is
	if len(symbolToIndex) == 0 {
		return balances
	}

	// Build CoinGecko IDs string
	var coinGeckoIDs []string
	for symbol := range symbolToIndex {
		coinGeckoID := GetCoinGeckoID(symbol)
		coinGeckoIDs = append(coinGeckoIDs, coinGeckoID)
	}
	idsString := strings.Join(coinGeckoIDs, ",")

	// Fetch prices from CoinGecko
	prices, err := s.GetPrices(idsString, "usd")
	if err != nil {
		// Log error but don't fail - just return balances without enrichment
		return balances
	}

	// Update balances with fetched prices
	for symbol, indices := range symbolToIndex {
		coinGeckoID := GetCoinGeckoID(symbol)
		if priceData, exists := prices[coinGeckoID]; exists {
			if usdPrice, hasPriceUSD := priceData["usd"]; hasPriceUSD {
				for _, idx := range indices {
					// Update USD price
					balances[idx].UsdPrice = usdPrice

					// Calculate USD value from balance and price
					if balanceFloat, err := strconv.ParseFloat(balances[idx].Balance, 64); err == nil {
						balances[idx].UsdValue = balanceFloat * usdPrice
					}
				}
			}
		}
	}

	return balances
}
