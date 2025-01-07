package main

import (
	"GinDemo/config"
	"GinDemo/router"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode) // 默认为 debug 模式，设置为发布模式
	r := gin.Default()
	router.InitRouter(r)
	r.Run(config.PORT)
}
