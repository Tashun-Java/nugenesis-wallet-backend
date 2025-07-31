package blockchaininfo

import (
	"encoding/json"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/blockchain_info/blockchain_info_models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestService_GetAddressInfo(t *testing.T) {
	mockResponse := blockchain_info_models.AddressInfo{
		Hash160:       "test-hash160",
		Address:       "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
		NTx:           10,
		TotalReceived: 5000000000,
		TotalSent:     2000000000,
		FinalBalance:  3000000000,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rawaddr/1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa" {
			t.Errorf("Expected path /rawaddr/1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa, got %s", r.URL.Path)
		}

		if r.URL.Query().Get("limit") != "50" {
			t.Errorf("Expected limit=50, got %s", r.URL.Query().Get("limit"))
		}

		if r.URL.Query().Get("offset") != "0" {
			t.Errorf("Expected offset=0, got %s", r.URL.Query().Get("offset"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	service := &Service{
		client:  &http.Client{},
		baseURL: server.URL,
	}

	result, err := service.GetAddressInfo("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", 50, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Address != mockResponse.Address {
		t.Errorf("Expected address %s, got %s", mockResponse.Address, result.Address)
	}

	if result.FinalBalance != mockResponse.FinalBalance {
		t.Errorf("Expected balance %d, got %d", mockResponse.FinalBalance, result.FinalBalance)
	}
}

func TestService_GetAddressInfo_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	service := &Service{
		client:  &http.Client{},
		baseURL: server.URL,
	}

	_, err := service.GetAddressInfo("invalid-address", 50, 0)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}

	expectedError := "API returned status 400"
	if err.Error() != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err.Error())
	}
}

func TestService_GetBalance(t *testing.T) {
	address := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
	mockResponse := map[string]blockchain_info_models.Balance{
		address: {
			Address:       address,
			Balance:       3000000000,
			TotalReceived: 5000000000,
			TotalSent:     2000000000,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/balance" {
			t.Errorf("Expected path /balance, got %s", r.URL.Path)
		}

		if r.URL.Query().Get("active") != address {
			t.Errorf("Expected active=%s, got %s", address, r.URL.Query().Get("active"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	service := &Service{
		client:  &http.Client{},
		baseURL: server.URL,
	}

	result, err := service.GetBalance(address)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedBalance := mockResponse[address]
	if result.Balance != expectedBalance.Balance {
		t.Errorf("Expected balance %d, got %d", expectedBalance.Balance, result.Balance)
	}
}

func TestService_GetUnspentOutputs(t *testing.T) {
	address := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
	mockResponse := blockchain_info_models.UnspentOutputs{
		UnspentOutputs: []blockchain_info_models.UnspentOutput{
			{
				TxHash:        "abc123",
				TxOutputN:     0,
				Value:         1000000000,
				Confirmations: 6,
			},
			{
				TxHash:        "def456",
				TxOutputN:     1,
				Value:         2000000000,
				Confirmations: 10,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/unspent" {
			t.Errorf("Expected path /unspent, got %s", r.URL.Path)
		}

		if r.URL.Query().Get("active") != address {
			t.Errorf("Expected active=%s, got %s", address, r.URL.Query().Get("active"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	service := &Service{
		client:  &http.Client{},
		baseURL: server.URL,
	}

	result, err := service.GetUnspentOutputs(address)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result.UnspentOutputs) != 2 {
		t.Errorf("Expected 2 unspent outputs, got %d", len(result.UnspentOutputs))
	}

	if result.UnspentOutputs[0].Value != 1000000000 {
		t.Errorf("Expected first output value 1000000000, got %d", result.UnspentOutputs[0].Value)
	}
}

func TestService_PushTransaction_Success(t *testing.T) {
	rawTx := "0100000001..."
	successMessage := "Transaction Submitted"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pushtx" {
			t.Errorf("Expected path /pushtx, got %s", r.URL.Path)
		}

		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		r.ParseForm()
		if r.FormValue("tx") != rawTx {
			t.Errorf("Expected tx=%s, got %s", rawTx, r.FormValue("tx"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(successMessage))
	}))
	defer server.Close()

	service := &Service{
		client:  &http.Client{},
		baseURL: server.URL,
	}

	result, err := service.PushTransaction(rawTx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Notice != successMessage {
		t.Errorf("Expected notice %s, got %s", successMessage, result.Notice)
	}

	if result.Error != "" {
		t.Errorf("Expected no error message, got %s", result.Error)
	}
}

func TestService_PushTransaction_Error(t *testing.T) {
	rawTx := "invalid-tx"
	errorMessage := "Transaction decode failed"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorMessage))
	}))
	defer server.Close()

	service := &Service{
		client:  &http.Client{},
		baseURL: server.URL,
	}

	result, err := service.PushTransaction(rawTx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Error != errorMessage {
		t.Errorf("Expected error %s, got %s", errorMessage, result.Error)
	}

	if result.Notice != "" {
		t.Errorf("Expected no notice, got %s", result.Notice)
	}
}

func TestService_GetTransaction(t *testing.T) {
	txHash := "abc123def456"
	mockResponse := blockchain_info_models.Tx{
		Hash:        txHash,
		Ver:         1,
		VinSz:       1,
		VoutSz:      2,
		Size:        250,
		Fee:         10000,
		BlockHeight: 700000,
		Time:        1640995200,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/rawtx/" + txHash
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	service := &Service{
		client:  &http.Client{},
		baseURL: server.URL,
	}

	result, err := service.GetTransaction(txHash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Hash != mockResponse.Hash {
		t.Errorf("Expected hash %s, got %s", mockResponse.Hash, result.Hash)
	}

	if result.Fee != mockResponse.Fee {
		t.Errorf("Expected fee %d, got %d", mockResponse.Fee, result.Fee)
	}
}

func TestService_GetMultiAddress(t *testing.T) {
	addresses := []string{"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2"}
	mockResponse := blockchain_info_models.MultiAddress{
		Addresses: []blockchain_info_models.AddressInfo{
			{
				Address:      addresses[0],
				FinalBalance: 1000000000,
			},
			{
				Address:      addresses[1],
				FinalBalance: 2000000000,
			},
		},
		Wallet: struct {
			NTx           int   `json:"n_tx"`
			NTxFiltered   int   `json:"n_tx_filtered"`
			TotalReceived int64 `json:"total_received"`
			TotalSent     int64 `json:"total_sent"`
			FinalBalance  int64 `json:"final_balance"`
		}{
			FinalBalance: 3000000000,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/multiaddr" {
			t.Errorf("Expected path /multiaddr, got %s", r.URL.Path)
		}

		expectedActive := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa|1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2"
		if r.URL.Query().Get("active") != expectedActive {
			t.Errorf("Expected active=%s, got %s", expectedActive, r.URL.Query().Get("active"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	service := &Service{
		client:  &http.Client{},
		baseURL: server.URL,
	}

	result, err := service.GetMultiAddress(addresses, 50, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result.Addresses) != 2 {
		t.Errorf("Expected 2 addresses, got %d", len(result.Addresses))
	}

	if result.Wallet.FinalBalance != 3000000000 {
		t.Errorf("Expected wallet balance 3000000000, got %d", result.Wallet.FinalBalance)
	}
}
