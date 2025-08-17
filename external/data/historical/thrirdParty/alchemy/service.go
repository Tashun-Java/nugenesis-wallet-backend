package alchemy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/alchemy/alchemy_models"
	"io"
	"net/http"
	"time"
)

const (
	BaseURL = "https://api.g.alchemy.com/data/v1"
	Timeout = 30 * time.Second
)

type Service struct {
	apiKey  *string
	baseURL *string
	client  *http.Client
}

func NewService(apiKey *string, baseURL *string) *Service {
	return &Service{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (s *Service) GetTokensByAddress(addresses []alchemy_models.AddressRequest, limit *int) (*alchemy_models.TokensByAddressResponse, error) {
	url := fmt.Sprintf("%s/%s/assets/tokens/by-address", *s.baseURL, *s.apiKey)

	requestBody := alchemy_models.TokensByAddressRequest{
		Addresses: addresses,
		Limit:     limit,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

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
	fmt.Printf("API response: %s\n", string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response alchemy_models.TokensByAddressResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (s *Service) GetTransactionHistory(addresses []alchemy_models.AddressRequest, before *string, after *string, limit *int) (*alchemy_models.TransactionHistoryResponse, error) {
	url := fmt.Sprintf("https://api.g.alchemy.com/data/v1/%s/transactions/history/by-address", *s.apiKey)

	requestBody := alchemy_models.TransactionHistoryRequest{
		Addresses: addresses,
		Before:    before,
		After:     after,
		Limit:     limit,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

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
	fmt.Printf("API response: %s\n", string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response alchemy_models.TransactionHistoryResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

//func (s *Service) GetSolanaHistory	(addresses []alchemy_models.AddressRequest, limit *int) (*alchemy_models.TokensByAddressResponse, error) {
