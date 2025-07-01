package data

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//func GetProducts(c *gin.Context) {
//	c.JSON(http.StatusOK, gin.H{"products": []string{"Laptop", "Phone"}})
//}

func GetProducts(c *gin.Context) {
	data, _ := GetUniPricesForToken("bitcoin,0chain", "usd")
	c.JSON(http.StatusOK, gin.H{"products": data})
}
