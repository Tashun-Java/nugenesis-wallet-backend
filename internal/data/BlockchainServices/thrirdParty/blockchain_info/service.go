package blockchaininfo

import (
	"encoding/json"
	"fmt"
	"github.com/tashunc/nugenesis-wallet-backend/internal/data/BlockchainServices/thrirdParty/blockchain_info/blockchain_info_models"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	BaseURL = "https://blockchain.info"
	Timeout = 30 * time.Second
)

type Service struct {
	client  *http.Client
	baseURL string
}

func NewService() *Service {
	return &Service{
		client: &http.Client{
			Timeout: Timeout,
		},
		baseURL: BaseURL,
	}
}

// https://www.blockchain.com/explorer/api/blockchain_api
func (s *Service) GetAddressInfo(address string, limit int, offset int) (*blockchain_info_models.AddressInfo, error) {
	url := fmt.Sprintf("%s/rawaddr/%s?limit=%d&offset=%d", s.baseURL, address, limit, offset)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get address info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var addressInfo blockchain_info_models.AddressInfo
	if err := json.Unmarshal(body, &addressInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &addressInfo, nil
}

func (s *Service) GetMultiAddress(addresses []string, limit int, offset int) (*blockchain_info_models.MultiAddress, error) {
	addressesStr := strings.Join(addresses, "|")
	url := fmt.Sprintf("%s/multiaddr?active=%s&limit=%d&offset=%d", s.baseURL, addressesStr, limit, offset)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get multi address info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var multiAddr blockchain_info_models.MultiAddress
	if err := json.Unmarshal(body, &multiAddr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &multiAddr, nil
}

func (s *Service) GetBalance(address string) (*blockchain_info_models.Balance, error) {
	url := fmt.Sprintf("%s/balance?active=%s", s.baseURL, address)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var balanceMap map[string]blockchain_info_models.Balance
	if err := json.Unmarshal(body, &balanceMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	balance, exists := balanceMap[address]
	if !exists {
		return nil, fmt.Errorf("address not found in response")
	}

	return &balance, nil
}

func (s *Service) GetUnspentOutputs(address string) (*blockchain_info_models.UnspentOutputs, error) {
	url := fmt.Sprintf("%s/unspent?active=%s", s.baseURL, address)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get unspent outputs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var unspentOutputs blockchain_info_models.UnspentOutputs
	if err := json.Unmarshal(body, &unspentOutputs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &unspentOutputs, nil
}

func (s *Service) PushTransaction(rawTx string) (*blockchain_info_models.PushTxResponse, error) {
	url := fmt.Sprintf("%s/pushtx", s.baseURL)

	payload := fmt.Sprintf("tx=%s", rawTx)
	resp, err := s.client.Post(url, "application/x-www-form-urlencoded", strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to push transaction: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var pushResp blockchain_info_models.PushTxResponse
	if resp.StatusCode == http.StatusOK {
		pushResp.Notice = string(body)
	} else {
		pushResp.Error = string(body)
	}

	return &pushResp, nil
}

func (s *Service) GetTransaction(txHash string) (*blockchain_info_models.Tx, error) {
	url := fmt.Sprintf("%s/rawtx/%s", s.baseURL, txHash)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var tx blockchain_info_models.Tx
	if err := json.Unmarshal(body, &tx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &tx, nil
}
