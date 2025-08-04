package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup) {
	userGroup := rg.Group("/users")
	userGroup.GET("/", GetUsers)
}
