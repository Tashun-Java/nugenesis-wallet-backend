package blockstream

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	service *Service
}

func NewController() *Controller {
	return &Controller{
		service: NewService(),
	}
}

func (c *Controller) GetAddressTransactions(ctx *gin.Context) {
	address := ctx.Param("address")
	if address == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
		return
	}

	if !ValidateBitcoinAddress(address) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid Bitcoin address format"})
		return
	}

	response, err := c.service.GetAddressTransactionsStandardized(address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
