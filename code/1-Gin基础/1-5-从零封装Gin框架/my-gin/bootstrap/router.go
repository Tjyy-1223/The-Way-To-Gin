package bootstrap

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"my-gin/global"
	"my-gin/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	// 前端项目静态资源
	router.StaticFile("/", "./static/dist/index.html")
	router.Static("/assets", "./static/dist/assets")
	router.StaticFile("/vite.svg", "./static/dist/vite.svg")

	// 其他静态资源
	router.Static("/public", "./static")
	router.Static("/storage", "./storage/app/public")

	// 注册 api 分组路由
	apiGroup := router.Group("/api")
	routes.SetApiGroupRoutes(apiGroup)

	return router
}

// RunServer 启动服务器
func RunServer() {
	r := setupRouter()

	srv := &http.Server{
		Addr:    ":" + global.App.Config.App.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器，设置 5 s的超时时间
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
