package helius

import (
	"encoding/json"
	"fmt"
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
		apiKey:  os.Getenv("HELIUS_API_KEY"),
		baseURL: "https://api.helius.xyz/v0",
		client:  &http.Client{},
	}
}

// GetAddressInfo retrieves information about a specific address
func (s *Service) GetAddressInfo(address string, limit, offset int) (interface{}, error) {
	url := fmt.Sprintf("%s/addresses/%s/transactions?api-key=%s&limit=%d&offset=%d", s.baseURL, address, s.apiKey, limit, offset)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	body, err := io.ReadAll(resp.Body)

	fmt.Println("Helius response:", string(body))

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, err
}
