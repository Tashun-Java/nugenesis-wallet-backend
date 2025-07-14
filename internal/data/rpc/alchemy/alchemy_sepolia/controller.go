package alchemy_sepolia

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tashunc/nugenesis-wallet-backend/internal/data/rpc/alchemy/alchemy_models"
	"github.com/tashunc/nugenesis-wallet-backend/internal/models"
)

type Controller struct {
	service *Service
}

func NewController() *Controller {
	return &Controller{
		service: NewService(os.Getenv("ALCHEMY_SEPOLIA_API_KEY")),
	}
}

func (c *Controller) SendRawTransaction(ctx *gin.Context) {
	var request models.SendRawTransactionControllerRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.SendRawTransactionControllerResponse{
			Success: false,
			Message: "Invalid request format",
			Error: &models.SendRawTransactionError{
				Code:    400,
				Message: err.Error(),
			},
		})
		return
	}

	response, err := c.service.PostSendRawTransaction(request.SignedTransactions)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.SendRawTransactionControllerResponse{
			Success: false,
			Message: "Failed to send transaction",
			Error: &models.SendRawTransactionError{
				Code:    500,
				Message: err.Error(),
			},
		})
		return
	}

	if response.Error != nil {
		ctx.JSON(http.StatusBadRequest, models.SendRawTransactionControllerResponse{
			Success: false,
			Message: "Transaction failed",
			Error:   response.Error,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.SendRawTransactionControllerResponse{
		Success:         true,
		TransactionHash: response.Result,
		Message:         "Transaction sent successfully",
	})
}

func (c *Controller) GetEstimateGas(ctx *gin.Context) {
	var request models.EstimateGasControllerRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.EstimateGasControllerResponse{
			Success: false,
			Message: "Invalid request format",
			Error: &models.SendRawTransactionError{
				Code:    400,
				Message: err.Error(),
			},
		})
		return
	}

	transactionObject := alchemy_models.TransactionObject{
		From:     request.From,
		To:       request.To,
		Gas:      request.Gas,
		GasPrice: request.GasPrice,
		Value:    request.Value,
		Data:     request.Data,
		Nonce:    request.Nonce,
	}

	response, err := c.service.GetEstimateGas(transactionObject)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.EstimateGasControllerResponse{
			Success: false,
			Message: "Failed to estimate gas",
			Error: &models.SendRawTransactionError{
				Code:    500,
				Message: err.Error(),
			},
		})
		return
	}

	if response.Error != nil {
		ctx.JSON(http.StatusBadRequest, models.EstimateGasControllerResponse{
			Success: false,
			Message: "Gas estimation failed",
			Error:   response.Error,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.EstimateGasControllerResponse{
		Success:      true,
		EstimatedGas: response.Result,
		Message:      "Gas estimated successfully",
	})
}

func (c *Controller) GetTransactionCount(ctx *gin.Context) {
	var request models.GetTransactionCountControllerRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.GetTransactionCountControllerResponse{
			Success: false,
			Message: "Invalid request format",
			Error: &models.SendRawTransactionError{
				Code:    400,
				Message: err.Error(),
			},
		})
		return
	}

	blockParameter := request.BlockParameter
	if blockParameter == "" {
		blockParameter = "pending"
	}

	response, err := c.service.GetTransactionCount(request.Address, blockParameter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.GetTransactionCountControllerResponse{
			Success: false,
			Message: "Failed to get transaction count",
			Error: &models.SendRawTransactionError{
				Code:    500,
				Message: err.Error(),
			},
		})
		return
	}

	if response.Error != nil {
		ctx.JSON(http.StatusBadRequest, models.GetTransactionCountControllerResponse{
			Success: false,
			Message: "Transaction count retrieval failed",
			Error:   response.Error,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.GetTransactionCountControllerResponse{
		Success:          true,
		TransactionCount: response.Result,
		Message:          "Transaction count retrieved successfully",
	})
}
