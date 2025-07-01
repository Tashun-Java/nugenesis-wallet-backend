package blockchaininfo

import (
	"github.com/gin-gonic/gin"
	"github.com/tashunc/nugenesis-wallet-backend/internal/models"
	"log"
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
	log.Printf("[START] %s %s?%s", address, limit, limit)

	addressInfo, err := c.service.GetAddressInfo(address, limit, offset)

	var mappedTransactions []models.Transaction
	log.Println(addressInfo)
	if addressInfo != nil {
		log.Printf("[START] %d ", len(addressInfo.Txs))
		for _, tx := range addressInfo.Txs {
			mappedTx := MapTxToTransaction(tx, address, addressInfo.Hash160)
			mappedTransactions = append(mappedTransactions, mappedTx)
		}
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//reponse := utils.MapTxToTransaction(addressInfo, address, addressInfo.Hash160)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//return mappedTransactions, nil
	ctx.JSON(http.StatusOK, mappedTransactions)
}

//func (c *Controller) GetMultiAddress(ctx *gin.Context) {
//	addressesParam := ctx.Query("addresses")
//	if addressesParam == "" {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "addresses parameter is required"})
//		return
//	}
//
//	addresses := strings.Split(addressesParam, ",")
//	limitStr := ctx.DefaultQuery("limit", "50")
//	offsetStr := ctx.DefaultQuery("offset", "0")
//
//	limit, err := strconv.Atoi(limitStr)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
//		return
//	}
//
//	offset, err := strconv.Atoi(offsetStr)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})
//		return
//	}
//
//	multiAddr, err := c.service.GetMultiAddress(addresses, limit, offset)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	ctx.JSON(http.StatusOK, multiAddr)
//}
//
//func (c *Controller) GetBalance(ctx *gin.Context) {
//	address := ctx.Param("address")
//	if address == "" {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
//		return
//	}
//
//	balance, err := c.service.GetBalance(address)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	ctx.JSON(http.StatusOK, balance)
//}
//
//func (c *Controller) GetUnspentOutputs(ctx *gin.Context) {
//	address := ctx.Param("address")
//	if address == "" {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "address parameter is required"})
//		return
//	}
//
//	unspentOutputs, err := c.service.GetUnspentOutputs(address)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	ctx.JSON(http.StatusOK, unspentOutputs)
//}
//
//func (c *Controller) PushTransaction(ctx *gin.Context) {
//	var request blockchain_info_models.PushTxRequest
//	if err := ctx.ShouldBindJSON(&request); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
//		return
//	}
//
//	if request.Tx == "" {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "tx field is required"})
//		return
//	}
//
//	response, err := c.service.PushTransaction(request.Tx)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	if response.Error != "" {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": response.Error})
//		return
//	}
//
//	ctx.JSON(http.StatusOK, gin.H{"message": response.Notice})
//}
//
//func (c *Controller) GetTransaction(ctx *gin.Context) {
//	txHash := ctx.Param("hash")
//	if txHash == "" {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "hash parameter is required"})
//		return
//	}
//
//	tx, err := c.service.GetTransaction(txHash)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	ctx.JSON(http.StatusOK, tx)
//}
