package data

import (
	"github.com/gin-gonic/gin"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/alchemy"
	blockchaininfo "github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/blockchain_info"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/blockstream"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/etherscan"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/helius"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/moralis"
	"github.com/tashunc/nugenesis-wallet-backend/external/data/rpc/alchemy/alchemy_general"
	"github.com/tashunc/nugenesis-wallet-backend/external/models/general"
	"os"
	"sync"
)

// Controller pool
type ControllerPool struct {
	bitcoinController          *blockchaininfo.Controller
	blockstreamController      *blockstream.Controller
	ethereumController         *etherscan.Controller
	solanaController           *helius.Controller
	polygonController          *moralis.Controller
	alchemyTokenController     *alchemy.Controller
	alchemyHistoricControllers map[general.CoinType]*alchemy.Controller
	alchemyRPCControllers      map[general.CoinType]*alchemy_general.Controller
	once                       sync.Once
}

var controllerPool *ControllerPool

func (cp *ControllerPool) GetAlchemyRPCController(coinType general.CoinType) *alchemy_general.Controller {
	return cp.alchemyRPCControllers[coinType]
}
func (cp *ControllerPool) GetAlchemyHistoricController(coinType general.CoinType) *alchemy.Controller {
	return cp.alchemyHistoricControllers[coinType]
}

func (cp *ControllerPool) GetBitcoinController() *blockchaininfo.Controller {
	return cp.bitcoinController
}

func (cp *ControllerPool) GetBlockstreamController() *blockstream.Controller {
	return cp.blockstreamController
}

func (cp *ControllerPool) GetEthereumController() *etherscan.Controller {
	return cp.ethereumController
}

func (cp *ControllerPool) GetSolanaController() *helius.Controller {
	return cp.solanaController
}

func (cp *ControllerPool) GetPolygonController() *moralis.Controller {
	return cp.polygonController
}

func (cp *ControllerPool) GetAlchemyTokenController() *alchemy.Controller {
	return cp.alchemyTokenController
}

func initControllers() {
	if controllerPool == nil {
		controllerPool = &ControllerPool{
			alchemyRPCControllers:      make(map[general.CoinType]*alchemy_general.Controller),
			alchemyHistoricControllers: make(map[general.CoinType]*alchemy.Controller),
		}
	}

	controllerPool.once.Do(func() {
		controllerPool.bitcoinController = blockchaininfo.NewController()
		controllerPool.blockstreamController = blockstream.NewController()
		controllerPool.ethereumController = etherscan.NewController()
		controllerPool.solanaController = helius.NewController()
		controllerPool.polygonController = moralis.NewController("polygon")
		//controllerPool.alchemyTokenController = alchemy.NewController()

		envMap := map[general.CoinType]string{
			general.Bitcoin:         "ALCHEMY_BITCOIN_RPC_BASE_URL",
			general.Solana:          "ALCHEMY_SOLANA_RPC_BASE_URL",
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
			// Testnet mappings (mainnet + 160)
			general.EthereumTest:        "ALCHEMY_ETHEREUM_TEST_RPC_BASE_URL",
			general.OptimismTest:        "ALCHEMY_OPTIMISM_TEST_RPC_BASE_URL",
			general.PolygonTest:         "ALCHEMY_POLYGON_TEST_RPC_BASE_URL",
			general.PolygonzkEVMTest:    "ALCHEMY_POLYGONZK_TEST_RPC_BASE_URL",
			general.ArbitrumTest:        "ALCHEMY_ARBITRUM_TEST_RPC_BASE_URL",
			general.ZetaEVMTest:         "ALCHEMY_ZETA_TEST_RPC_BASE_URL",
			general.MantleTest:          "ALCHEMY_MANTLE_TEST_RPC_BASE_URL",
			general.BlastTest:           "ALCHEMY_BLAST_TEST_RPC_BASE_URL",
			general.LineaTest:           "ALCHEMY_LINEA_TEST_RPC_BASE_URL",
			general.RoninTest:           "ALCHEMY_RONIN_TEST_RPC_BASE_URL",
			general.RootstockTest:       "ALCHEMY_ROOTSTOCK_TEST_RPC_BASE_URL",
			general.ArbitrumNovaTest:    "ALCHEMY_ARBITRUMNOVA_TEST_RPC_BASE_URL",
			general.BaseTest:            "ALCHEMY_BASE_TEST_RPC_BASE_URL",
			general.AvalancheCChainTest: "ALCHEMY_AVALANCHE_TEST_RPC_BASE_URL",
			general.BinanceTest:         "ALCHEMY_BINANCE_TEST_RPC_BASE_URL",
			general.CeloTest:            "ALCHEMY_CELO_TEST_RPC_BASE_URL",
			general.SonicTest:           "ALCHEMY_SONIC_TEST_RPC_BASE_URL",
			general.SeiTest:             "ALCHEMY_SEI_TEST_RPC_BASE_URL",
			general.ScrollTest:          "ALCHEMY_SCROLL_TEST_RPC_BASE_URL",
			general.OpBNBTest:           "ALCHEMY_OPBNB_TEST_RPC_BASE_URL",
			general.SepoliaTest:         "ALCHEMY_SEPOLIA_TEST_RPC_BASE_URL",
		}

		for coinType, envVar := range envMap {
			if url := os.Getenv(envVar); url != "" {
				controllerPool.alchemyRPCControllers[coinType] = alchemy_general.NewController(url)
				controllerPool.alchemyHistoricControllers[coinType] = alchemy.NewController(url)
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
		//coinType := general.CoinType(blockchainID)
		//
		//if controller := controllerPool.GetAlchemyHistoricController(coinType); controller != nil {
		//	controller.GetTransactionHistoryByAddress(ctx)
		//} else {
		//	ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		//}
		//
		switch general.CoinType(blockchainID) {
		case general.Bitcoin:
			controllerPool.GetBlockstreamController().GetAddressTransactions(ctx)
		case general.Ethereum:
			controllerPool.GetEthereumController().GetAddressInfo(ctx)
		case general.Solana:
			controllerPool.GetSolanaController().GetAddressInfo(ctx)
		case general.Polygon:
			controllerPool.GetPolygonController().GetWalletHistory(ctx)
		default:
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})

	rg.GET("/tokens/:address", func(ctx *gin.Context) {
		controllerPool.GetAlchemyTokenController().GetTokensByAddress(ctx)
	})

	rg.POST("/tokens", func(ctx *gin.Context) {
		controllerPool.GetAlchemyTokenController().GetTokensByMultipleAddresses(ctx)
	})

	rg.GET("/tokens", func(ctx *gin.Context) {
		controllerPool.GetAlchemyTokenController().GetTokensByAddressQuery(ctx)
	})
}

func RegisterRPCRoutes(rg *gin.RouterGroup) {

	rg.POST("/send", func(ctx *gin.Context) {
		blockchainID := ctx.Param("id")
		coinType := general.CoinType(blockchainID)

		if controller := controllerPool.GetAlchemyRPCController(coinType); controller != nil {
			controller.SendRawTransaction(ctx)
		} else {
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})

	rg.POST("/feeEstimate", func(ctx *gin.Context) {
		blockchainID := ctx.Param("id")
		coinType := general.CoinType(blockchainID)

		if controller := controllerPool.GetAlchemyRPCController(coinType); controller != nil {
			controller.GetEstimateGas(ctx)
		} else {
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})

	rg.GET("/getGasPrice", func(ctx *gin.Context) {
		blockchainID := ctx.Param("id")
		coinType := general.CoinType(blockchainID)

		if controller := controllerPool.GetAlchemyRPCController(coinType); controller != nil {
			controller.GetGasPrice(ctx)
		} else {
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})

	rg.POST("/getCount", func(ctx *gin.Context) {
		blockchainID := ctx.Param("id")
		coinType := general.CoinType(blockchainID)

		if controller := controllerPool.GetAlchemyRPCController(coinType); controller != nil {
			controller.GetTransactionCount(ctx)
		} else {
			ctx.JSON(400, gin.H{"error": "Unsupported blockchain"})
		}
	})

}
