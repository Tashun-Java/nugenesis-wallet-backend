package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func data(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":     "This is middleware data",
		"data":        []string{"item1", "item2", "item3"},
		"server_time": time.Now().Unix(),
	})
}

func submit(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Data submitted successfully",
		"received":     payload,
		"processed_at": time.Now().Unix(),
	})
}

func update(c *gin.Context) {
	id := c.Param("id")
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Resource updated successfully",
		"id":         id,
		"data":       payload,
		"updated_at": time.Now().Unix(),
	})
}

func delete(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message":    "Resource deleted successfully",
		"id":         id,
		"deleted_at": time.Now().Unix(),
	})
}
