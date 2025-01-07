package v2

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddMember(c *gin.Context) {
	// 获取 Get 参数
	name := c.Query("name")
	price := c.DefaultQuery("price", "100")

	c.JSON(http.StatusOK, gin.H{
		"v2":    "AddMember",
		"name":  name,
		"price": price,
	})
}
