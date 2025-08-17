package helius

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Controller struct {
	service *Service
}

func NewController() *Controller {
	return &Controller{
		service: NewService(),
	}
}

func (c *Controller) GetAddressInfo(ctx *gin.Context) {
	address := ctx.Param("address")
	if address == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
		return
	}

	limitStr := ctx.DefaultQuery("limit", "50")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})
		return
	}

	txInfo, err := c.service.GetAddressInfo(address, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//var mappedTxs []models.Transaction
	//for _, tx := range txInfo.Result {
	//	mapped := MapTxToTransaction(tx, address)
	//	mappedTxs = append(mappedTxs, mapped)
	//}

	//ctx.JSON(http.StatusOK, mappedTxs)
	ctx.JSON(http.StatusOK, txInfo)
}
