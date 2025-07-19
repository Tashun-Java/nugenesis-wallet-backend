package alchemy

import (
	"github.com/gin-gonic/gin"
	"github.com/tashunc/nugenesis-wallet-backend/internal/data/historical/thrirdParty/alchemy/alchemy_models"
	"net/http"
	"strings"
)

type Controller struct {
	service *Service
}

func NewController() *Controller {
	return &Controller{
		service: NewService(),
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

	response, err := c.service.GetTokensByAddress(addresses)
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

	response, err := c.service.GetTokensByAddress(request.Addresses)
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

	response, err := c.service.GetTokensByAddress(addressRequests)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
