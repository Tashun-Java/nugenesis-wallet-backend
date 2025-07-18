package data

import (
	"github.com/gin-gonic/gin"
	blockchaininfo "github.com/tashunc/nugenesis-wallet-backend/internal/data/historical/thrirdParty/blockchain_info"
	"github.com/tashunc/nugenesis-wallet-backend/internal/data/historical/thrirdParty/etherscan"
	"github.com/tashunc/nugenesis-wallet-backend/internal/data/rpc/alchemy/alchemy_general"
	"github.com/tashunc/nugenesis-wallet-backend/internal/models/general"
	"os"
	"sync"
)

// Controller instances (singletons)
var (
	bitcoinController        *blockchaininfo.Controller
	ethereumController       *etherscan.Controller
	alchemySepoliaController *alchemy_general.Controller
	alchemyPolygonController *alchemy_general.Controller
	controllerOnce           sync.Once
)

// Initialize controllers as singletons
func initControllers() {
	alchemySepoliaRpcApiKey := os.Getenv("ALCHEMY_SEPOLIA_RPC_API_KEY")
	alchemyPolygonRpcApiKey := os.Getenv("ALCHEMY_POLYGON_RPC_API_KEY")

	controllerOnce.Do(func() {
		bitcoinController = blockchaininfo.NewController()
		ethereumController = etherscan.NewController()
		alchemySepoliaController = alchemy_general.NewController(&alchemySepoliaRpcApiKey)
		alchemyPolygonController = alchemy_general.NewController(&alchemyPolygonRpcApiKey)
	})
}

func RegisterRoutes(rg *gin.RouterGroup) {
	// Initialize controllers once
	initControllers()

	productGroup := rg.Group("/data")
	productGroup.GET("/", GetProducts)

	blockchainGroup := productGroup.Group("/:id")
	{
		registerHistoricalRoutes(blockchainGroup)
		RegisterRPCRoutes(blockchainGroup)
	}
}

func registerHistoricalRoutes(rg *gin.RouterGroup) {
	rg.GET("/address/:address", func(ctx *gin.Context) {
		blockchainID := ctx.Param("id")

		switch general.CoinType(blockchainID) {
		case general.Bitcoin:
			bitcoinController.GetAddressInfo(ctx)
		case general.Ethereum:
			ethereumController.GetAddressInfo(ctx)
		default:
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})
}

func RegisterRPCRoutes(rg *gin.RouterGroup) {

	rg.POST("/send", func(ctx *gin.Context) {
		blockchainID := ctx.Param("id")

		switch general.CoinType(blockchainID) {
		case general.Polygon:
			alchemyPolygonController.SendRawTransaction(ctx)
		case general.Sepolia:
			alchemySepoliaController.SendRawTransaction(ctx)
		default:
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})

	rg.POST("/feeEstimate", func(ctx *gin.Context) {
		blockchainID := ctx.Param("id")

		switch general.CoinType(blockchainID) {
		case general.Polygon:
			alchemyPolygonController.GetEstimateGas(ctx)
		case general.Sepolia:
			alchemySepoliaController.GetEstimateGas(ctx)
		default:
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})

	rg.GET("/getGasPrice", func(ctx *gin.Context) {
		blockchainID := ctx.Param("id")

		switch general.CoinType(blockchainID) {
		case general.Polygon:
			alchemyPolygonController.GetGasPrice(ctx)
		case general.Sepolia:
			alchemySepoliaController.GetGasPrice(ctx)
		default:
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})

	rg.POST("/getCount", func(ctx *gin.Context) {
		blockchainID := ctx.Param("id")

		switch general.CoinType(blockchainID) {
		case general.Polygon:
			alchemyPolygonController.GetTransactionCount(ctx)
		case general.Sepolia:
			alchemySepoliaController.GetTransactionCount(ctx)
		default:
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})

}
