package user

//import (
//	"github.com/gin-gonic/gin"
//	"net/http"
//)
//
//func GetUser(c *gin.Context) {
//	id := c.Param("id")
//	user, err := GetUserByID(id)
//	if err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
//		return
//	}
//	c.JSON(http.StatusOK, user)
//}
