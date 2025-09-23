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
