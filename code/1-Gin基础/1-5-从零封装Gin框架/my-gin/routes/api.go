package routes

import (
	"github.com/gin-gonic/gin"
	"my-gin/app/controller/app"
	"my-gin/app/middleware"
	"my-gin/app/service"
	"net/http"
	"time"
)

// SetApiGroupRoutes 定义 api 分组路由
func SetApiGroupRoutes(router *gin.RouterGroup) {
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.GET("/test", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "success")
	})

	router.POST("/auth/register", app.Register)
	router.POST("/auth/login", app.Login)
	authRouter := router.Group("").Use(middleware.JWTAuth(service.AppGuardName))
	{
		authRouter.POST("/auth/info", app.Info)
	}
}
