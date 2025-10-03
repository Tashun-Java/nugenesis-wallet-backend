package staticControllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tashunc/nugenesis-wallet-backend/static/staticModels"
	"github.com/tashunc/nugenesis-wallet-backend/static/staticServices"
)

// AssetController handles HTTP requests for asset operations
type AssetController struct {
	assetService *staticServices.AssetService
}

// NewAssetController creates a new AssetController instance
func NewAssetController(assetService *staticServices.AssetService) *AssetController {
	return &AssetController{
		assetService: assetService,
	}
}

// GetByCoinSymbol handles GET requests to get assets by coin symbol
func (c *AssetController) GetByCoinSymbol(ctx *gin.Context) {
	symbol := ctx.Param("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, staticModels.ErrorResponse{
			Error: "Symbol parameter is required",
		})
		return
	}

	assets, err := c.assetService.GetByCoinSymbol(symbol)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, staticModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	response := staticModels.AssetsBySymbolResponse{
		Symbol: symbol,
		Assets: assets,
		Count:  len(assets),
	}
	ctx.JSON(http.StatusOK, response)
}

// GetAllSymbols handles GET requests to get all available coin symbols
func (c *AssetController) GetAllSymbols(ctx *gin.Context) {
	symbols, err := c.assetService.GetAllSymbols()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, staticModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	response := staticModels.AllSymbolsResponse{
		Symbols: symbols,
		Count:   len(symbols),
	}
	ctx.JSON(http.StatusOK, response)
}

// GetAllAssets handles GET requests to get all assets with full details
func (c *AssetController) GetAllAssets(ctx *gin.Context) {
	assets, err := c.assetService.GetAllAssets()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, staticModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	response := staticModels.AllAssetsResponse{
		Assets: assets,
		Count:  len(assets),
	}
	ctx.JSON(http.StatusOK, response)
}

// ForceRefresh handles POST requests to force cache refresh
func (c *AssetController) ForceRefresh(ctx *gin.Context) {
	err := c.assetService.ForceRefresh()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, staticModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	response := staticModels.RefreshResponse{
		Message: "Cache refreshed successfully",
	}
	ctx.JSON(http.StatusOK, response)
}

// GetCacheStats handles GET requests to get cache statistics
func (c *AssetController) GetCacheStats(ctx *gin.Context) {
	stats := c.assetService.GetCacheStats()
	response := staticModels.CacheStatsResponse{
		CacheStats: stats,
	}
	ctx.JSON(http.StatusOK, response)
}

// GenerateAllIDs handles POST requests to generate IDs for all assets
func (c *AssetController) GenerateAllIDs(ctx *gin.Context) {
	err := c.assetService.GenerateAllIDs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, staticModels.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	stats := c.assetService.GetCacheStats()
	ctx.JSON(http.StatusOK, gin.H{
		"message":        "All asset IDs generated successfully",
		"total_mappings": stats["total_id_mappings"],
		"next_id":        stats["next_id"],
	})
}

// HealthCheck handles GET requests to check service health
func (c *AssetController) HealthCheck(ctx *gin.Context) {
	stats := c.assetService.GetCacheStats()
	ctx.JSON(http.StatusOK, gin.H{
		"status":      "healthy",
		"cache_stats": stats,
		"timestamp":   stats["last_update"],
	})
}
