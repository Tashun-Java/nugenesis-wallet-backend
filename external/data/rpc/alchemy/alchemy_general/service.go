package alchemy_general

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/rpc/alchemy/alchemy_models"
	"io"
	"log"
	"net/http"
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

func (s *Service) PostSendRawTransaction(signedTxs []string) (*alchemy_models.SendRawTransactionResponse, error) {

	url := fmt.Sprintf("%s%s", *s.baseURL, *s.apiKey)

	request := &alchemy_models.SendRawTransactionRequest{
		Jsonrpc: "2.0",
		Method:  "eth_sendRawTransaction",
		Params:  signedTxs,
		Id:      1,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("Alchemy API response: %s", string(respBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("alchemy API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var response alchemy_models.SendRawTransactionResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (s *Service) GetEstimateGas(transactionObject alchemy_models.TransactionObject) (*alchemy_models.EstimateGasResponse, error) {
	url := fmt.Sprintf("%s%s", *s.baseURL, *s.apiKey)

	request := &alchemy_models.EstimateGasRequest{
		Jsonrpc: "2.0",
		Method:  "eth_estimateGas",
		Params:  []alchemy_models.TransactionObject{transactionObject},
		Id:      1,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("Alchemy API response: %s", string(respBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("alchemy API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var response alchemy_models.EstimateGasResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (s *Service) GetTransactionCount(address string, blockParameter string) (*alchemy_models.GetTransactionCountResponse, error) {
	url := fmt.Sprintf("%s%s", *s.baseURL, *s.apiKey)

	request := &alchemy_models.GetTransactionCountRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getTransactionCount",
		Params:  []string{address, blockParameter},
		Id:      1,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("Alchemy API response: %s", string(respBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("alchemy API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var response alchemy_models.GetTransactionCountResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to  unmarshal response: %w", err)
	}

	return &response, nil
}

func (s *Service) GetGasPrice() (*alchemy_models.GetGasPriceResponse, error) {
	url := fmt.Sprintf("%s%s", *s.baseURL, *s.apiKey)

	request := &alchemy_models.GetGasPriceRequest{
		Jsonrpc: "2.0",
		Method:  "eth_gasPrice",
		Params:  []interface{}{},
		Id:      1,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("Alchemy API response: %s", string(respBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("alchemy API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var response alchemy_models.GetGasPriceResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}
