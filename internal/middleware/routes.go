package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/tashunc/nugenesis-wallet-backend/services"
)

func RegisterRoutes(rg *gin.RouterGroup, nonceStore *services.NonceStore) {
	protected := rg.Group("/middleware")
	protected.Use(services.ClientNonceAuthMiddleware(nonceStore))

	protected.GET("/data", data)
	protected.POST("/submit", submit)
	protected.PUT("/update", update)
	protected.DELETE("/delete", delete)

}
