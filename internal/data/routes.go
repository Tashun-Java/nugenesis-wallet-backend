package data

import (
	"github.com/gin-gonic/gin"
	blockchaininfo "github.com/tashunc/nugenesis-wallet-backend/internal/data/historical/thrirdParty/blockchain_info"

	"github.com/tashunc/nugenesis-wallet-backend/internal/data/historical/thrirdParty/etherscan"
	"github.com/tashunc/nugenesis-wallet-backend/internal/data/rpc/alchemy/alchemy_sepolia"
	"github.com/tashunc/nugenesis-wallet-backend/internal/models"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	productGroup := rg.Group("/data")
	productGroup.GET("/", GetProducts)

	// Create routes based on blockchain decimal values
	blockchainGroup := productGroup.Group("/:id")
	{
		registerHistoricalRoutes(blockchainGroup)
		RegisterRPCRoutes(blockchainGroup)
	}

	// Legacy routes for backward compatibility
	//btcGroup := productGroup.Group("/btc")
	//{
	//	c := blockchaininfo.NewController()
	//	btcGroup.GET("/address/:address", c.GetAddressInfo)
	//}
	//
	//ethGroup := productGroup.Group("/eth")
	//{
	//	c := etherscan.NewController()
	//	ethGroup.GET("/address/:address", c.GetAddressInfo)
	//}

	// Ethereum RPC routes
}

func registerHistoricalRoutes(rg *gin.RouterGroup) {
	rg.GET("/address/:address", func(ctx *gin.Context) {
		blockchainID := ctx.Param("id")

		switch models.BlockchainDecimal(blockchainID) {
		case models.Bitcoin:
			c := blockchaininfo.NewController()
			c.GetAddressInfo(ctx)
		case models.Ethereum:
			c := etherscan.NewController()
			c.GetAddressInfo(ctx)
		default:
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})
}

func RegisterRPCRoutes(rg *gin.RouterGroup) {

	rg.POST("/send", func(ctx *gin.Context) {
		blockchainID := ctx.Param("id")

		switch models.BlockchainDecimal(blockchainID) {
		case models.Sepolia:
			c := alchemy_sepolia.NewController()
			c.SendRawTransaction(ctx)
		default:
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})

	rg.POST("/feeEstimate", func(ctx *gin.Context) {
		blockchainID := ctx.Param("id")

		switch models.BlockchainDecimal(blockchainID) {
		case models.Sepolia:
			c := alchemy_sepolia.NewController()
			c.GetEstimateGas(ctx)
		default:
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})

	rg.GET("/getGasPrice", func(ctx *gin.Context) {
		blockchainID := ctx.Param("id")

		switch models.BlockchainDecimal(blockchainID) {
		case models.Sepolia:
			c := alchemy_sepolia.NewController()
			c.GetGasPrice(ctx)
		default:
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})

	rg.POST("/getCount", func(ctx *gin.Context) {
		blockchainID := ctx.Param("id")

		switch models.BlockchainDecimal(blockchainID) {
		case models.Sepolia:
			c := alchemy_sepolia.NewController()
			c.GetTransactionCount(ctx)
		default:
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})

}
