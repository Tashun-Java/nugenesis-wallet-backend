package data

import (
	"github.com/gin-gonic/gin"
	"github.com/tashunc/nugenesis-wallet-backend/internal/data/BlockchainServices/thrirdParty/etherscan"
)
import "github.com/tashunc/nugenesis-wallet-backend/internal/data/BlockchainServices/thrirdParty/blockchain_info"

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
}
