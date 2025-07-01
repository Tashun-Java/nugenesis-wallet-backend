package logger

import (
	"github.com/gin-gonic/gin"
	"log"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Started %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
		log.Printf("Completed %s %s with status %d", c.Request.Method, c.Request.URL.Path, c.Writer.Status())
	}
}
