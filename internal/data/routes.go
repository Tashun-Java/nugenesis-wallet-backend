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

// Controller pool
type ControllerPool struct {
	bitcoinController  *blockchaininfo.Controller
	ethereumController *etherscan.Controller
	alchemyControllers map[general.CoinType]*alchemy_general.Controller
	once               sync.Once
}

var controllerPool *ControllerPool

func (cp *ControllerPool) GetAlchemyController(coinType general.CoinType) *alchemy_general.Controller {
	return cp.alchemyControllers[coinType]
}

func (cp *ControllerPool) GetBitcoinController() *blockchaininfo.Controller {
	return cp.bitcoinController
}

func (cp *ControllerPool) GetEthereumController() *etherscan.Controller {
	return cp.ethereumController
}

func initControllers() {
	if controllerPool == nil {
		controllerPool = &ControllerPool{
			alchemyControllers: make(map[general.CoinType]*alchemy_general.Controller),
		}
	}

	controllerPool.once.Do(func() {
		controllerPool.bitcoinController = blockchaininfo.NewController()
		controllerPool.ethereumController = etherscan.NewController()

		envMap := map[general.CoinType]string{
			general.Ethereum:        "ALCHEMY_ETHEREUM_RPC_BASE_URL",
			general.Optimism:        "ALCHEMY_OPTIMISM_RPC_BASE_URL",
			general.Polygon:         "ALCHEMY_POLYGON_RPC_BASE_URL",
			general.PolygonzkEVM:    "ALCHEMY_POLYGONZK_RPC_BASE_URL",
			general.Arbitrum:        "ALCHEMY_ARBITRUM_RPC_BASE_URL",
			general.ZetaEVM:         "ALCHEMY_ZETA_RPC_BASE_URL",
			general.Fantom:          "ALCHEMY_FANTOM_RPC_BASE_URL",
			general.Mantle:          "ALCHEMY_MANTLE_RPC_BASE_URL",
			general.Blast:           "ALCHEMY_BLAST_RPC_BASE_URL",
			general.Linea:           "ALCHEMY_LINEA_RPC_BASE_URL",
			general.Ronin:           "ALCHEMY_RONIN_RPC_BASE_URL",
			general.Rootstock:       "ALCHEMY_ROOTSTOCK_RPC_BASE_URL",
			general.ArbitrumNova:    "ALCHEMY_ARBITRUMNOVA_RPC_BASE_URL",
			general.Base:            "ALCHEMY_BASE_RPC_BASE_URL",
			general.AvalancheCChain: "ALCHEMY_AVALANCHE_RPC_BASE_URL",
			general.Binance:         "ALCHEMY_BINANCE_RPC_BASE_URL",
			general.Celo:            "ALCHEMY_CELO_RPC_BASE_URL",
			general.Metis:           "ALCHEMY_METIS_RPC_BASE_URL",
			general.Sonic:           "ALCHEMY_SONIC_RPC_BASE_URL",
			general.Sei:             "ALCHEMY_SEI_RPC_BASE_URL",
			general.Scroll:          "ALCHEMY_SCROLL_RPC_BASE_URL",
			general.OpBNB:           "ALCHEMY_OPBNB_RPC_BASE_URL",
			general.Sepolia:         "ALCHEMY_SEPOLIA_RPC_BASE_URL",
		}

		for coinType, envVar := range envMap {
			if url := os.Getenv(envVar); url != "" {
				controllerPool.alchemyControllers[coinType] = alchemy_general.NewController(url)
			}
		}
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
			controllerPool.GetBitcoinController().GetAddressInfo(ctx)
		case general.Ethereum:
			controllerPool.GetEthereumController().GetAddressInfo(ctx)
		default:
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})
}

func RegisterRPCRoutes(rg *gin.RouterGroup) {

	rg.POST("/send", func(ctx *gin.Context) {
		blockchainID := ctx.Param("id")
		coinType := general.CoinType(blockchainID)

		if controller := controllerPool.GetAlchemyController(coinType); controller != nil {
			controller.SendRawTransaction(ctx)
		} else {
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})

	rg.POST("/feeEstimate", func(ctx *gin.Context) {
		blockchainID := ctx.Param("id")
		coinType := general.CoinType(blockchainID)

		if controller := controllerPool.GetAlchemyController(coinType); controller != nil {
			controller.GetEstimateGas(ctx)
		} else {
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})

	rg.GET("/getGasPrice", func(ctx *gin.Context) {
		blockchainID := ctx.Param("id")
		coinType := general.CoinType(blockchainID)

		if controller := controllerPool.GetAlchemyController(coinType); controller != nil {
			controller.GetGasPrice(ctx)
		} else {
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})

	rg.POST("/getCount", func(ctx *gin.Context) {
		blockchainID := ctx.Param("id")
		coinType := general.CoinType(blockchainID)

		if controller := controllerPool.GetAlchemyController(coinType); controller != nil {
			controller.GetTransactionCount(ctx)
		} else {
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})

}
