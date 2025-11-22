package moralis

import (
	"encoding/json"
	"fmt"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/moralis/moralis_models"
	"io"
	"net/http"
	"os"
)

type Service struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewService() *Service {
	return &Service{
		apiKey:  os.Getenv("MORALIS_API_KEY"),
		baseURL: "https://deep-index.moralis.io/api/v2.2",
		client:  &http.Client{},
	}
}

func (s *Service) GetWalletHistory(address string, chain string, cursor string, limit int) (*moralis_models.WalletHistoryResponse, error) {
	url := fmt.Sprintf(
		"%s/wallets/%s/history?chain=%s&order=DESC",
		s.baseURL, address, chain,
	)

	if cursor != "" {
		url += fmt.Sprintf("&cursor=%s", cursor)
	}

	if limit > 0 {
		url += fmt.Sprintf("&limit=%d", limit)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", s.apiKey)
	req.Header.Set("accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request Moralis API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Moralis API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Moralis response: %w", err)
	}

	var apiResp moralis_models.WalletHistoryResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Moralis response: %w", err)
	}

	return &apiResp, nil
}

// GetWalletTokenBalances retrieves token balances (including native tokens) for a wallet address
// Supports multiple chains and includes USD values for tokens
func (s *Service) GetWalletTokenBalances(address string, chain string, cursor string, limit int, excludeSpam bool) (*moralis_models.WalletTokenBalancesResponse, error) {
	// Build the URL with required parameters
	url := fmt.Sprintf(
		"%s/wallets/%s/tokens?chain=%s",
		s.baseURL, address, chain,
	)

	// Add optional parameters
	if cursor != "" {
		url += fmt.Sprintf("&cursor=%s", cursor)
	}

	if limit > 0 {
		url += fmt.Sprintf("&limit=%d", limit)
	}

	if excludeSpam {
		url += "&exclude_spam=true"
	}

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-API-Key", s.apiKey)
	req.Header.Set("accept", "application/json")

	// Execute request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request Moralis API: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Moralis API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Moralis response: %w", err)
	}

	// Unmarshal JSON response
	var apiResp moralis_models.WalletTokenBalancesResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Moralis response: %w", err)
	}

	return &apiResp, nil
}

// GetSolanaTokenBalances retrieves SPL token balances for a Solana wallet address
// Uses the Solana-specific Moralis Gateway API
func (s *Service) GetSolanaTokenBalances(address string, network string) (*moralis_models.SolanaTokenBalancesResponse, error) {
	// Default to mainnet if network is not specified
	if network == "" {
		network = "mainnet"
	}

	// Build the Solana Gateway API URL
	url := fmt.Sprintf(
		"https://solana-gateway.moralis.io/account/%s/%s/tokens",
		network, address,
	)

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-API-Key", s.apiKey)
	req.Header.Set("accept", "application/json")

	// Execute request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request Moralis Solana API: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Moralis Solana API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Moralis Solana response: %w", err)
	}

	// Unmarshal JSON response
	var apiResp moralis_models.SolanaTokenBalancesResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Moralis Solana response: %w", err)
	}

	return &apiResp, nil
}

// GetSolanaBalance retrieves native SOL balance for a Solana wallet address
func (s *Service) GetSolanaBalance(address string, network string) (*moralis_models.SolanaBalanceResponse, error) {
	// Default to mainnet if network is not specified
	if network == "" {
		network = "mainnet"
	}

	// Build the Solana Gateway API URL
	url := fmt.Sprintf(
		"https://solana-gateway.moralis.io/account/%s/%s/balance",
		network, address,
	)

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-API-Key", s.apiKey)
	req.Header.Set("accept", "application/json")

	// Execute request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request Moralis Solana API: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Moralis Solana API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Moralis Solana response: %w", err)
	}

	// Unmarshal JSON response
	var apiResp moralis_models.SolanaBalanceResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Moralis Solana response: %w", err)
	}

	return &apiResp, nil
}
