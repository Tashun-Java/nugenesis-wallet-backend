package data

import (
	"github.com/gin-gonic/gin"
	"github.com/tashunc/nugenesis-wallet-backend/internal/data/historical/thrirdParty/blockchain_info"
	"github.com/tashunc/nugenesis-wallet-backend/internal/data/historical/thrirdParty/etherscan"
	"github.com/tashunc/nugenesis-wallet-backend/internal/data/rpc/alchemy/alchemy_sepolia"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	productGroup := rg.Group("/data")
	productGroup.GET("/", GetProducts)

	btcGroup := productGroup.Group("/btc")
	{
		c := blockchaininfo.NewController()
		btcGroup.GET("/address/:address", c.GetAddressInfo)
	}

	ethGroup := productGroup.Group("/eth")
	{
		c := etherscan.NewController()
		ethGroup.GET("/address/:address", c.GetAddressInfo)
	}
	//ethAlchemyGroup := productGroup.Group("/eth")
	{
		c := alchemy_sepolia.NewController()
		ethGroup.POST("/send", c.SendRawTransaction)
		ethGroup.POST("/feeEstimate", c.GetEstimateGas)
		ethGroup.POST("/fees", c.GetEstimateGas)
		ethGroup.POST("/getCount", c.GetTransactionCount)
		ethGroup.GET("/getGasPrice", c.GetGasPrice)

	}
}
