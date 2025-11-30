package coingecko

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controller handles HTTP requests for CoinGecko API
type Controller struct {
	service *Service
}

// NewController creates a new CoinGecko controller instance
func NewController() *Controller {
	return &Controller{
		service: NewService(),
	}
}

// GetTokenPrices fetches prices for specified token IDs
func (c *Controller) GetTokenPrices(ctx *gin.Context) {
	ids := ctx.Query("ids")
	vsCurrencies := ctx.DefaultQuery("vs_currencies", "usd")

	if ids == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ids parameter is required"})
		return
	}

	prices, err := c.service.GetPrices(ids, vsCurrencies)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, prices)
}
