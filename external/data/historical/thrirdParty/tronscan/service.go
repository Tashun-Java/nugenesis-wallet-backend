package tronscan

import (
	"encoding/json"
	"fmt"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/tronscan/tronscan_models"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	// TronScan API is more reliable and provides better transaction data
	DefaultBaseURL = "https://apilist.tronscan.org/api"
	// Alternative: TronGrid API (official Tron Foundation API)
	TronGridURL = "https://api.trongrid.io"
	Timeout     = 30 * time.Second
)

type Service struct {
	baseURL     string
	tronGridURL string
	apiKey      string
	client      *http.Client
}

func NewService() *Service {
	apiKey := os.Getenv("TRONGRID_API_KEY")

	return &Service{
		baseURL:     DefaultBaseURL,
		tronGridURL: TronGridURL,
		apiKey:      apiKey,
		client: &http.Client{
			Timeout: Timeout,
		},
	}
}

// GetAddressTransactions fetches transactions for an address from TronScan API
func (s *Service) GetAddressTransactions(address string, limit int, start int) (*tronscan_models.TronScanResponse, error) {
	if limit == 0 {
		limit = 50 // Default limit
	}

	url := fmt.Sprintf("%s/transaction?sort=-timestamp&count=true&limit=%d&start=%d&address=%s",
		s.baseURL, limit, start, address)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("TronScan API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response tronscan_models.TronScanResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetAddressTransactionsFromTronGrid fetches transactions from TronGrid API (alternative)
func (s *Service) GetAddressTransactionsFromTronGrid(address string, limit int, fingerprint string) (*tronscan_models.TronGridResponse, error) {
	if limit == 0 {
		limit = 20
	}

	url := fmt.Sprintf("%s/v1/accounts/%s/transactions?limit=%d&order_by=block_timestamp,desc",
		s.tronGridURL, address, limit)

	if fingerprint != "" {
		url += fmt.Sprintf("&fingerprint=%s", fingerprint)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add API key if available
	if s.apiKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", s.apiKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("TronGrid API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response tronscan_models.TronGridResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetAccountTokens fetches all token balances for an address from TronScan API
func (s *Service) GetAccountTokens(address string, start int, limit int) (*tronscan_models.AccountTokensResponse, error) {
	if limit == 0 {
		limit = 20
	}

	url := fmt.Sprintf("%s/account/tokens?address=%s&start=%d&limit=%d",
		s.baseURL, address, start, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("TronScan API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response tronscan_models.AccountTokensResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetAccountInfo fetches account information including TRX balance
func (s *Service) GetAccountInfo(address string) (*tronscan_models.AccountInfoResponse, error) {
	url := fmt.Sprintf("%s/account?address=%s", s.baseURL, address)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("TronScan API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response tronscan_models.AccountInfoResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}
