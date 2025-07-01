package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.GET("/login", Login)
	auth.GET("/callback", Callback)

	rg.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})
}
