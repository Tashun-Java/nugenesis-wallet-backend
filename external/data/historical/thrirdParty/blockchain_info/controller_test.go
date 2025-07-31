package blockchaininfo

//
//import (
//	"bytes"
//	"encoding/json"
//	"github.com/tashunc/nugenesis-wallet-backend/external/data/rpc/thrirdParty/blockchain_info/blockchain_info_models"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/gin-gonic/gin"
//)
//
//func setupTestRouter() *gin.Engine {
//	gin.SetMode(gin.TestMode)
//	router := gin.New()
//	controller := NewController()
//	controller.RegisterRoutes(router)
//	return router
//}
//
//func TestController_GetAddressInfo(t *testing.T) {
//	mockService := &MockService{
//		GetAddressInfoFunc: func(address string, limit int, offset int) (*blockchain_info_models.AddressInfo, error) {
//			return &blockchain_info_models.AddressInfo{
//				Address:      address,
//				FinalBalance: 1000000000,
//				NTx:          5,
//			}, nil
//		},
//	}
//
//	controller := &Controller{service: mockService}
//	router := gin.New()
//	router.GET("/api/blockchain/address/:address", controller.GetAddressInfo)
//
//	req, _ := http.NewRequest("GET", "/api/blockchain/address/1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa?limit=10&offset=5", nil)
//	w := httptest.NewRecorder()
//	router.ServeHTTP(w, req)
//
//	if w.Code != http.StatusOK {
//		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
//	}
//
//	var response blockchain_info_models.AddressInfo
//	err := json.Unmarshal(w.Body.Bytes(), &response)
//	if err != nil {
//		t.Fatalf("Failed to unmarshal response: %v", err)
//	}
//
//	if response.Address != "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa" {
//		t.Errorf("Expected address 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa, got %s", response.Address)
//	}
//
//	if response.FinalBalance != 1000000000 {
//		t.Errorf("Expected balance 1000000000, got %d", response.FinalBalance)
//	}
//}
//
//func TestController_GetAddressInfo_MissingAddress(t *testing.T) {
//	router := setupTestRouter()
//
//	req, _ := http.NewRequest("GET", "/api/blockchain/address/", nil)
//	w := httptest.NewRecorder()
//	router.ServeHTTP(w, req)
//
//	if w.Code != http.StatusNotFound {
//		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
//	}
//}
//
//func TestController_GetAddressInfo_InvalidLimit(t *testing.T) {
//	router := setupTestRouter()
//
//	req, _ := http.NewRequest("GET", "/api/blockchain/address/1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa?limit=invalid", nil)
//	w := httptest.NewRecorder()
//	router.ServeHTTP(w, req)
//
//	if w.Code != http.StatusBadRequest {
//		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
//	}
//
//	var response map[string]string
//	json.Unmarshal(w.Body.Bytes(), &response)
//	if response["error"] != "invalid limit parameter" {
//		t.Errorf("Expected error message 'invalid limit parameter', got %s", response["error"])
//	}
//}
//
//func TestController_GetMultiAddress(t *testing.T) {
//	mockService := &MockService{
//		GetMultiAddressFunc: func(addresses []string, limit int, offset int) (*blockchain_info_models.MultiAddress, error) {
//			return &blockchain_info_models.MultiAddress{
//				Addresses: []blockchain_info_models.AddressInfo{
//					{Address: addresses[0], FinalBalance: 1000000000},
//					{Address: addresses[1], FinalBalance: 2000000000},
//				},
//				Wallet: struct {
//					NTx           int   `json:"n_tx"`
//					NTxFiltered   int   `json:"n_tx_filtered"`
//					TotalReceived int64 `json:"total_received"`
//					TotalSent     int64 `json:"total_sent"`
//					FinalBalance  int64 `json:"final_balance"`
//				}{
//					FinalBalance: 3000000000,
//				},
//			}, nil
//		},
//	}
//
//	controller := &Controller{service: mockService}
//	router := gin.New()
//	router.GET("/api/blockchain/multiaddr", controller.GetMultiAddress)
//
//	req, _ := http.NewRequest("GET", "/api/blockchain/multiaddr?addresses=addr1,addr2&limit=10", nil)
//	w := httptest.NewRecorder()
//	router.ServeHTTP(w, req)
//
//	if w.Code != http.StatusOK {
//		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
//	}
//
//	var response blockchain_info_models.MultiAddress
//	err := json.Unmarshal(w.Body.Bytes(), &response)
//	if err != nil {
//		t.Fatalf("Failed to unmarshal response: %v", err)
//	}
//
//	if len(response.Addresses) != 2 {
//		t.Errorf("Expected 2 addresses, got %d", len(response.Addresses))
//	}
//
//	if response.Wallet.FinalBalance != 3000000000 {
//		t.Errorf("Expected wallet balance 3000000000, got %d", response.Wallet.FinalBalance)
//	}
//}
//
//func TestController_GetMultiAddress_MissingAddresses(t *testing.T) {
//	router := setupTestRouter()
//
//	req, _ := http.NewRequest("GET", "/api/blockchain/multiaddr", nil)
//	w := httptest.NewRecorder()
//	router.ServeHTTP(w, req)
//
//	if w.Code != http.StatusBadRequest {
//		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
//	}
//
//	var response map[string]string
//	json.Unmarshal(w.Body.Bytes(), &response)
//	if response["error"] != "addresses parameter is required" {
//		t.Errorf("Expected error message 'addresses parameter is required', got %s", response["error"])
//	}
//}
//
//func TestController_GetBalance(t *testing.T) {
//	mockService := &MockService{
//		GetBalanceFunc: func(address string) (*blockchain_info_models.Balance, error) {
//			return &blockchain_info_models.Balance{
//				Address:       address,
//				Balance:       1500000000,
//				TotalReceived: 2000000000,
//				TotalSent:     500000000,
//			}, nil
//		},
//	}
//
//	controller := &Controller{service: mockService}
//	router := gin.New()
//	router.GET("/api/blockchain/balance/:address", controller.GetBalance)
//
//	req, _ := http.NewRequest("GET", "/api/blockchain/balance/1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", nil)
//	w := httptest.NewRecorder()
//	router.ServeHTTP(w, req)
//
//	if w.Code != http.StatusOK {
//		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
//	}
//
//	var response blockchain_info_models.Balance
//	err := json.Unmarshal(w.Body.Bytes(), &response)
//	if err != nil {
//		t.Fatalf("Failed to unmarshal response: %v", err)
//	}
//
//	if response.Balance != 1500000000 {
//		t.Errorf("Expected balance 1500000000, got %d", response.Balance)
//	}
//}
//
//func TestController_GetUnspentOutputs(t *testing.T) {
//	mockService := &MockService{
//		GetUnspentOutputsFunc: func(address string) (*blockchain_info_models.UnspentOutputs, error) {
//			return &blockchain_info_models.UnspentOutputs{
//				UnspentOutputs: []blockchain_info_models.UnspentOutput{
//					{
//						TxHash:        "abc123",
//						TxOutputN:     0,
//						Value:         1000000000,
//						Confirmations: 6,
//					},
//				},
//			}, nil
//		},
//	}
//
//	controller := &Controller{service: mockService}
//	router := gin.New()
//	router.GET("/api/blockchain/unspent/:address", controller.GetUnspentOutputs)
//
//	req, _ := http.NewRequest("GET", "/api/blockchain/unspent/1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", nil)
//	w := httptest.NewRecorder()
//	router.ServeHTTP(w, req)
//
//	if w.Code != http.StatusOK {
//		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
//	}
//
//	var response blockchain_info_models.UnspentOutputs
//	err := json.Unmarshal(w.Body.Bytes(), &response)
//	if err != nil {
//		t.Fatalf("Failed to unmarshal response: %v", err)
//	}
//
//	if len(response.UnspentOutputs) != 1 {
//		t.Errorf("Expected 1 unspent output, got %d", len(response.UnspentOutputs))
//	}
//
//	if response.UnspentOutputs[0].Value != 1000000000 {
//		t.Errorf("Expected value 1000000000, got %d", response.UnspentOutputs[0].Value)
//	}
//}
//
//func TestController_PushTransaction_Success(t *testing.T) {
//	mockService := &MockService{
//		PushTransactionFunc: func(rawTx string) (*blockchain_info_models.PushTxResponse, error) {
//			return &blockchain_info_models.PushTxResponse{
//				Notice: "Transaction Submitted",
//				Error:  "",
//			}, nil
//		},
//	}
//
//	controller := &Controller{service: mockService}
//	router := gin.New()
//	router.POST("/api/blockchain/pushtx", controller.PushTransaction)
//
//	requestBody := blockchain_info_models.PushTxRequest{
//		Tx: "0100000001...",
//	}
//	bodyBytes, _ := json.Marshal(requestBody)
//
//	req, _ := http.NewRequest("POST", "/api/blockchain/pushtx", bytes.NewBuffer(bodyBytes))
//	req.Header.Set("Content-Type", "application/json")
//	w := httptest.NewRecorder()
//	router.ServeHTTP(w, req)
//
//	if w.Code != http.StatusOK {
//		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
//	}
//
//	var response map[string]string
//	err := json.Unmarshal(w.Body.Bytes(), &response)
//	if err != nil {
//		t.Fatalf("Failed to unmarshal response: %v", err)
//	}
//
//	if response["message"] != "Transaction Submitted" {
//		t.Errorf("Expected message 'Transaction Submitted', got %s", response["message"])
//	}
//}
//
//func TestController_PushTransaction_Error(t *testing.T) {
//	mockService := &MockService{
//		PushTransactionFunc: func(rawTx string) (*blockchain_info_models.PushTxResponse, error) {
//			return &blockchain_info_models.PushTxResponse{
//				Notice: "",
//				Error:  "Transaction decode failed",
//			}, nil
//		},
//	}
//
//	controller := &Controller{service: mockService}
//	router := gin.New()
//	router.POST("/api/blockchain/pushtx", controller.PushTransaction)
//
//	requestBody := blockchain_info_models.PushTxRequest{
//		Tx: "invalid-tx",
//	}
//	bodyBytes, _ := json.Marshal(requestBody)
//
//	req, _ := http.NewRequest("POST", "/api/blockchain/pushtx", bytes.NewBuffer(bodyBytes))
//	req.Header.Set("Content-Type", "application/json")
//	w := httptest.NewRecorder()
//	router.ServeHTTP(w, req)
//
//	if w.Code != http.StatusBadRequest {
//		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
//	}
//
//	var response map[string]string
//	err := json.Unmarshal(w.Body.Bytes(), &response)
//	if err != nil {
//		t.Fatalf("Failed to unmarshal response: %v", err)
//	}
//
//	if response["error"] != "Transaction decode failed" {
//		t.Errorf("Expected error 'Transaction decode failed', got %s", response["error"])
//	}
//}
//
//func TestController_PushTransaction_InvalidJSON(t *testing.T) {
//	router := setupTestRouter()
//
//	req, _ := http.NewRequest("POST", "/api/blockchain/pushtx", bytes.NewBuffer([]byte("invalid json")))
//	req.Header.Set("Content-Type", "application/json")
//	w := httptest.NewRecorder()
//	router.ServeHTTP(w, req)
//
//	if w.Code != http.StatusBadRequest {
//		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
//	}
//
//	var response map[string]string
//	json.Unmarshal(w.Body.Bytes(), &response)
//	if response["error"] != "invalid request body" {
//		t.Errorf("Expected error 'invalid request body', got %s", response["error"])
//	}
//}
//
//func TestController_PushTransaction_EmptyTx(t *testing.T) {
//	router := setupTestRouter()
//
//	requestBody := blockchain_info_models.PushTxRequest{
//		Tx: "",
//	}
//	bodyBytes, _ := json.Marshal(requestBody)
//
//	req, _ := http.NewRequest("POST", "/api/blockchain/pushtx", bytes.NewBuffer(bodyBytes))
//	req.Header.Set("Content-Type", "application/json")
//	w := httptest.NewRecorder()
//	router.ServeHTTP(w, req)
//
//	if w.Code != http.StatusBadRequest {
//		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
//	}
//
//	var response map[string]string
//	json.Unmarshal(w.Body.Bytes(), &response)
//	if response["error"] != "tx field is required" {
//		t.Errorf("Expected error 'tx field is required', got %s", response["error"])
//	}
//}
