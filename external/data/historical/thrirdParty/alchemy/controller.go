package alchemy

import (
	"github.com/gin-gonic/gin"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/alchemy/alchemy_models"
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
