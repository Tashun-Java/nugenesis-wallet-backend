package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tashunc/nugenesis-wallet-backend/config"
	"github.com/tashunc/nugenesis-wallet-backend/external/auth"
	"github.com/tashunc/nugenesis-wallet-backend/external/data"
	"github.com/tashunc/nugenesis-wallet-backend/static"

	"github.com/tashunc/nugenesis-wallet-backend/external/user"
	"github.com/tashunc/nugenesis-wallet-backend/pkg/logger"
	"net/http"
	"time"
)

func main() {
	cfg := config.LoadConfig()

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-Nonce", "x-nonce-timestamp", "x-nonce-hash"}
	router.Use(cors.New(corsConfig))
	router.Use(logger.GinLogger())

	//nonceStore := staticServices.NewNonceStore()
	// API group
	api := router.Group("/api")
	{
		user.RegisterRoutes(api)
		data.RegisterRoutes(api)
		auth.RegisterRoutes(api)
		static.RegisterRoutes(api)
		//middleware.RegisterRoutes(api, nonceStore)

	}
	router.GET("/api/nonce/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"nonce_length_min": 32,
			"ttl_seconds":      300, // 5 minutes
			"hash_algorithm":   "SHA256",
			"format":           "nonce + timestamp -> SHA256",
		})
	})

	router.POST("/api/nonce/register", func(c *gin.Context) {
		var request struct {
			Nonce     string `json:"nonce" binding:"required"`
			Timestamp int64  `json:"timestamp" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		timestamp := time.Unix(request.Timestamp, 0)

		// Check if timestamp is reasonable (not too old, not too far in future)
		now := time.Now()
		if timestamp.Before(now.Add(-1*time.Minute)) || timestamp.After(now.Add(1*time.Minute)) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp"})
			return
		}

		//hash := nonceStore.RegisterNonce(request.Nonce, timestamp)

		//c.JSON(http.StatusOK, gin.H{
		//	"hash":      hash,
		//	"timestamp": request.Timestamp,
		//	"expires":   timestamp.Add(5 * time.Minute).Unix(),
		//})
	})

	err := router.Run(":" + cfg.Port)
	if err != nil {
		fmt.Print("Failed to start server:", err)
	}
}
