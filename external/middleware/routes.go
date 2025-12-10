package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/tashunc/nugenesis-wallet-backend/static/staticServices"
)

func RegisterRoutes(rg *gin.RouterGroup, nonceStore *staticServices.NonceStore) {
	protected := rg.Group("/middleware")
	protected.Use(staticServices.ClientNonceAuthMiddleware(nonceStore))

	protected.GET("/data", data)
	protected.POST("/submit", submit)
	protected.PUT("/update", update)
	protected.DELETE("/delete", delete)

}
