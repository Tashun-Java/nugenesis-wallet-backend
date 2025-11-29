package alchemy

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/alchemy/alchemy_models"
	"github.com/tashunc/nugenesis-wallet-backend/external/models"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Controller struct {
	service *Service
}

func NewController(baseUrl string) *Controller {
	AlchemyApiKey := os.Getenv("ALCHEMY_API_KEY")
	controllerBaseURL := baseUrl
	return &Controller{
		service: NewService(&AlchemyApiKey, &controllerBaseURL),
	}
}

func (c *Controller) GetTokensByAddress(ctx *gin.Context) {
	address := ctx.Param("address")
	if address == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
		return
	}

	network := ctx.DefaultQuery("network", "eth-mainnet")

	addresses := []alchemy_models.AddressRequest{
		{
			Address:  address,
			Networks: []string{network},
		},
	}

	// Handle pagination parameters
	limit := ctx.Query("limit")

	var limitPtr *int

	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			limitPtr = &l
		}
	}

	response, err := c.service.GetTokensByAddress(addresses, limitPtr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) GetTokensByMultipleAddresses(ctx *gin.Context) {
	var request alchemy_models.TokensByAddressRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if len(request.Addresses) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "addresses array is required"})
		return
	}

	response, err := c.service.GetTokensByAddress(request.Addresses, request.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) GetTokensByAddressQuery(ctx *gin.Context) {
	addressesParam := ctx.Query("addresses")
	if addressesParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "addresses parameter is required"})
		return
	}

	addresses := strings.Split(addressesParam, ",")
	network := ctx.DefaultQuery("network", "eth-mainnet")

	var addressRequests []alchemy_models.AddressRequest
	for _, addr := range addresses {
		addressRequests = append(addressRequests, alchemy_models.AddressRequest{
			Address:  strings.TrimSpace(addr),
			Networks: []string{network},
		})
	}

	// Handle pagination parameters
	limit := ctx.Query("limit")

	var limitPtr *int

	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			limitPtr = &l
		}
	}

	response, err := c.service.GetTokensByAddress(addressRequests, limitPtr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) GetTransactionHistoryByAddress(ctx *gin.Context) {
	address := ctx.Param("address")
	if address == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
		return
	}

	network := ctx.DefaultQuery("network", "eth-mainnet")

	addresses := []alchemy_models.AddressRequest{
		{
			Address:  address,
			Networks: []string{network},
		},
	}

	// Handle pagination parameters
	limit := ctx.Query("limit")
	before := ctx.Query("before")
	after := ctx.Query("after")

	var limitPtr *int
	var beforePtr *string
	var afterPtr *string

	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			limitPtr = &l
		}
	}

	if before != "" {
		beforePtr = &before
	}

	if after != "" {
		afterPtr = &after
	}

	response, err := c.service.GetTransactionHistory(addresses, beforePtr, afterPtr, limitPtr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) GetTransactionHistoryByMultipleAddresses(ctx *gin.Context) {
	var request alchemy_models.TransactionHistoryRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if len(request.Addresses) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "addresses array is required"})
		return
	}

	response, err := c.service.GetTransactionHistory(request.Addresses, request.Before, request.After, request.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) GetAssetTransfers(ctx *gin.Context) {
	address := ctx.Param("address")
	if address == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
		return
	}

	// Parse pagination parameters
	limitStr := ctx.DefaultQuery("limit", "100")
	pageKey := ctx.Query("pageKey")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	// Convert limit to hex string for Alchemy API (divide by 2 since we're making 2 requests)
	maxCount := fmt.Sprintf("0x%x", limit/2)

	categories := []string{"external", "internal", "erc20", "erc721", "erc1155"}

	// Request 1: Get transactions sent FROM the address
	fromRequest := &alchemy_models.AssetTransfersRequest{
		FromAddress:  &address,
		Category:     categories,
		MaxCount:     &maxCount,
		WithMetadata: func() *bool { b := true; return &b }(),
	}

	if pageKey != "" {
		fromRequest.PageKey = &pageKey
	}

	fromResponse, err := c.service.GetAssetTransfers(fromRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch sent transactions: " + err.Error()})
		return
	}

	// Request 2: Get transactions sent TO the address
	toRequest := &alchemy_models.AssetTransfersRequest{
		ToAddress:    &address,
		Category:     categories,
		MaxCount:     &maxCount,
		WithMetadata: func() *bool { b := true; return &b }(),
	}

	if pageKey != "" {
		toRequest.PageKey = &pageKey
	}

	toResponse, err := c.service.GetAssetTransfers(toRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch received transactions: " + err.Error()})
		return
	}

	// Combine both responses
	allTransfers := append(fromResponse.Transfers, toResponse.Transfers...)

	// Map to standardized Transaction model
	var mappedTxs []models.Transaction
	for _, transfer := range allTransfers {
		mapped := MapAssetTransferToTransaction(transfer, address)
		mappedTxs = append(mappedTxs, mapped)
	}

	// Sort by block number (descending - most recent first)
	// Simple bubble sort for demonstration - you might want to use a better sorting algorithm
	for i := 0; i < len(mappedTxs)-1; i++ {
		for j := 0; j < len(mappedTxs)-i-1; j++ {
			if mappedTxs[j].Date < mappedTxs[j+1].Date {
				mappedTxs[j], mappedTxs[j+1] = mappedTxs[j+1], mappedTxs[j]
			}
		}
	}

	ctx.JSON(http.StatusOK, mappedTxs)
}
