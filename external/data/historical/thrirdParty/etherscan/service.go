package etherscan

import (
	"encoding/json"
	"fmt"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/etherscan/etherscan_models"
	"io"
	"net/http"
)

type Service struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewService(apiKey string) *Service {
	return &Service{
		apiKey:  apiKey,
		baseURL: "https://api.etherscan.io/api",
		client:  &http.Client{},
	}
}

func (s *Service) GetAddressInfo(address string, limit, offset int) (*etherscan_models.AddressResponse, error) {
	url := fmt.Sprintf(
		"%s?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&page=%d&offset=%d&sort=desc&apikey=%s",
		s.baseURL, address, (offset/limit)+1, limit, s.apiKey,
	)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to request Etherscan: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Etherscan API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Etherscan response: %w", err)
	}

	var apiResp etherscan_models.AddressResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Etherscan response: %w", err)
	}

	if apiResp.Status != "1" {
		return nil, fmt.Errorf("Etherscan API error: %s", apiResp.Message)
	}

	return &apiResp, nil
}
