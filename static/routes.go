package static

import (
	"github.com/gin-gonic/gin"
	"github.com/tashunc/nugenesis-wallet-backend/static/staticControllers"
	"github.com/tashunc/nugenesis-wallet-backend/static/staticServices"
)

var assetController *staticControllers.AssetController

func initControllers() {
	if assetController == nil {
		assetService := staticServices.NewAssetService()
		assetController = staticControllers.NewAssetController(assetService)
	}
}

func RegisterRoutes(rg *gin.RouterGroup) {
	initControllers()

	staticGroup := rg.Group("/static")
	{
		assetGroup := staticGroup.Group("/assets")
		{
			assetGroup.GET("/symbols", assetController.GetAllSymbols)
			assetGroup.GET("/symbol/:symbol", assetController.GetByCoinSymbol)
			assetGroup.POST("/refresh", assetController.ForceRefresh)
			assetGroup.GET("/cache/stats", assetController.GetCacheStats)
		}
	}
}
