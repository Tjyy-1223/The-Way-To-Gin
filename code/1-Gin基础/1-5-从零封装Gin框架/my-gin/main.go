package main

import (
	"my-gin/bootstrap"
	"my-gin/global"
)

func main() {
	// 初始化配置
	bootstrap.InitializeConfig()

	// 初始化日志
	global.App.Log = bootstrap.InitializeLog()
	global.App.Log.Info("log init success!")

	// 初始化数据库
	global.App.DB = bootstrap.InitializeDB()

	// 初始化验证器
	bootstrap.InitializeValidator()

	// 初始化 Redis
	global.App.Redis = bootstrap.InitializeRedis()

	// 初始化文件系统
	bootstrap.InitializeStorage()

	// 程序关闭前，释放数据库连接
	defer func() {
		if global.App.DB != nil {
			db, _ := global.App.DB.DB()
			err := db.Close()
			if err != nil {
				global.App.Log.Error("failed to close the database in the process.")
			}
		}
	}()

	// 启动服务器
	bootstrap.RunServer()
}
