package data

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// Struct for reusable client
var client = resty.New().
	SetHostURL("https://api.coingecko.com/api/v3").
	SetHeader("Accept", "application/json")

// Fetches crypto token prices in specified fiat currencies
func GetUniPricesForToken(ids, vsCurrencies string) (map[string]map[string]float64, error) {
	var result map[string]map[string]float64

	resp, err := client.R().
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
