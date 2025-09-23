package static

import (
	"github.com/gin-gonic/gin"
	"github.com/tashunc/nugenesis-wallet-backend/static/staticControllers"
	"github.com/tashunc/nugenesis-wallet-backend/static/staticServices"
)

var assetController *staticControllers.AssetController
var blockchainController *staticControllers.BlockchainController

func initControllers() {
	if assetController == nil {
		assetService := staticServices.NewAssetService()
		blockchainService := staticServices.NewBlockchainService(assetService)
		assetController = staticControllers.NewAssetController(assetService)
		blockchainController = staticControllers.NewBlockchainController(blockchainService)
	}
}

func RegisterRoutes(rg *gin.RouterGroup) {
	initControllers()

	staticGroup := rg.Group("/static")
	{
		// Serve static asset files (logos, etc.)
		staticGroup.Static("/assetsLogo", "./assets")

		// API routes for asset data
		assetGroup := staticGroup.Group("/assets")
		{
			assetGroup.GET("/symbols", assetController.GetAllSymbols)
			assetGroup.GET("/symbol/:symbol", assetController.GetByCoinSymbol)
			assetGroup.POST("/refresh", assetController.ForceRefresh)
			assetGroup.POST("/generate-ids", assetController.GenerateAllIDs)
			assetGroup.GET("/cache/stats", assetController.GetCacheStats)
			assetGroup.GET("/health", assetController.HealthCheck)
		}

		// API routes for blockchain data
		blockchainGroup := staticGroup.Group("/blockchains")
		{
			blockchainGroup.GET("/", blockchainController.GetAllBlockchains)
			blockchainGroup.GET("/id/:id", blockchainController.GetBlockchainByID)
			blockchainGroup.GET("/name/:name/id", blockchainController.GetBlockchainID)
		}
	}
}
