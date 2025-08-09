package blockstream

import (
	"encoding/json"
	"fmt"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/blockstream/blockstream_models"
	"io"
	"net/http"
	"time"
)

const (
	BaseURL = "https://blockstream.info/api"
	Timeout = 30 * time.Second
)

type Service struct {
	baseURL string
	client  *http.Client
}

func NewService() *Service {
	return &Service{
		baseURL: BaseURL,
		client: &http.Client{
			Timeout: Timeout,
		},
	}
}

func (s *Service) GetAddressTransactions(address string) (*blockstream_models.AddressTransactionsResponse, error) {
	url := fmt.Sprintf("%s/address/%s/txs", s.baseURL, address)

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
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response blockstream_models.AddressTransactionsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (s *Service) GetAddressTransactionsStandardized(address string) (*blockstream_models.StandardizedTransactionsResponse, error) {
	rawTransactions, err := s.GetAddressTransactions(address)
	if err != nil {
		return nil, err
	}

	standardizedResponse := MapToStandardizedTransactions(rawTransactions, address)
	return standardizedResponse, nil
}
