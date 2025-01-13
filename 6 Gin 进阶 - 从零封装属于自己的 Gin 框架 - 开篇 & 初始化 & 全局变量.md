# 第六章 Gin 进阶 - 从零封装属于自己的 Gin 框架 - 开篇 & 初始化 & 全局变量

学习资料参考：[手把手，带你从零封装Gin框架](https://juejin.cn/post/7016742808560074783)

## 1 总体结构

整体项目的结构以及对应含义如下：

```
my-gin
	- app
		- common  			// 公共模块（请求、响应结构体等）
		- controllers   // 业务调度器
		- middleware 		// 中间件
		- models 				// 数据库结构体
		- services 			// 业务层
  - bootstrap 			// 项目启动初始化
  - config 					// 配置结构体
  - global 					// 全局变量
  - routes 					// 路由定义
  - static 					// 静态资源（允许外部访问）
  - storage 				// 系统日志、文件等静态资源
  - utils 					// 工具函数
  - config.yaml 		// 配置文件
  - go.mod  				// 项目启动文件
  main.go
```



## 2 项目初始化

1. 先在目录下创建一个目录 `my-gin` 用来存放项目代码

```
mkdir my-gin
```

2. 在项目根目录下，初始化 `go.mod` 文件

```
go mod init my-gin
```

3. 安装 Gin

```
go get -u github.com/gin-gonic/gin
```

4. 在项目根目录下编写 `main.go` 文件

```go
package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    r := gin.Default()

    // 测试路由
    r.GET("/ping", func(c *gin.Context) {
        c.String(http.StatusOK, "pong")
    })

    // 启动服务器
    r.Run(":8080")
}
```

执行 `go run main.go` 启动应用

可以看到 HTTP 服务器启动成功，打开 http://127.0.0.1:8080/ping 测试路由



## 3 配置初始化 & 全局变量

配置文件是每个项目必不可少的部分，用来保存应用基本数据、数据库配置等信息，避免要修改一个配置项需要到处找的尴尬。

这里使用 [viper](https://link.juejin.cn/?target=https%3A%2F%2Fgithub.com%2Fspf13%2Fviper) 作为配置管理方案，它支持 JSON、TOML、YAML、HCL、envfile、Java properties 等多种格式的配置文件，并且能够监听配置文件的修改，进行热重载，详细介绍大家可以去官方文档查看

安装：

```
go get -u github.com/spf13/viper 
```



#### **3.1 编辑配置文件**

在项目根目录下新建一个文件 `config.yaml` ，初期先将项目的基本配置放入，后续我们会添加更多配置信息

```yaml
app: # 应用基本配置
  env: local # 环境名称
  port: 8888 # 服务监听端口号
  app_name: gin-app # 应用名称
  app_url: http://localhost # 应用域名
```



#### **3.2 编写配置结构体**

在项目根目录下新建文件夹 `config`，用于存放所有配置对应的结构体

新建 `config.go` 文件，定义 `Configuration` 结构体，其 `App` 属性对应 `config.yaml` 中的 `app`

```go
package config

type Configuration struct {
    App App `mapstructure:"app" json:"app" yaml:"app"`
}
```

新建 `app.go` 文件，定义 `App` 结构体，**其所有属性分别对应 `config.yaml` 中 `app` 下的所有配置**

```go
package config

type App struct {
    Env string `mapstructure:"env" json:"env" yaml:"env"`
    Port string `mapstructure:"port" json:"port" yaml:"port"`
    AppName string `mapstructure:"app_name" json:"app_name" yaml:"app_name"`
    AppUrl string `mapstructure:"app_url" json:"app_url" yaml:"app_url"`
}
```

注意：配置结构体中 `mapstructure` 标签需对应 `config.ymal` 中的配置名称， `viper` 会根据标签 value 值把 `config.yaml` 的数据赋予给结构体



#### **3.3 全局变量**

新建 `global/app.go` 文件，定义 `Application` 结构体，用来存放一些项目启动时的变量，便于调用，目前先将 `viper` 结构体和 `Configuration` 结构体放入，后续会添加其他成员属性

```go
package global

import (
    "github.com/spf13/viper"
    "jassue-gin/config"
)

type Application struct {
    ConfigViper *viper.Viper
    Config config.Configuration
}

var App = new(Application)
```



#### 3.4 使用 viper 载入配置 - 重点

新建 `bootstrap/config.go` 文件，编写代码：

```go
package bootstrap

import (
    "fmt"
    "github.com/fsnotify/fsnotify"
    "github.com/spf13/viper"
    "jassue-gin/global"
    "os"
)

func InitializeConfig() *viper.Viper {
    // 设置配置文件路径
    config := "config.yaml"
    // 生产环境可以通过设置环境变量来改变配置文件路径
    if configEnv := os.Getenv("VIPER_CONFIG"); configEnv != "" {
        config = configEnv
    }

    // 初始化 viper
    v := viper.New()
    v.SetConfigFile(config)
    v.SetConfigType("yaml")
    if err := v.ReadInConfig(); err != nil {
        panic(fmt.Errorf("read config failed: %s \n", err))
    }

    // 监听配置文件
    v.WatchConfig()
    v.OnConfigChange(func(in fsnotify.Event) {
        fmt.Println("config file changed:", in.Name)
        // 重载配置
        if err := v.Unmarshal(&global.App.Config); err != nil {
            fmt.Println(err)
        }
    })
    // 将配置赋值给全局变量
    if err := v.Unmarshal(&global.App.Config); err != nil {
       fmt.Println(err)
    }

    return v
}
```



#### 3.5 初始化配置 

修改 `main.go` 文件

```go
package main

import (
    "github.com/gin-gonic/gin"
    "jassue-gin/bootstrap"
    "jassue-gin/global"
    "net/http"
)

func main() {
    // 初始化配置
    bootstrap.InitializeConfig()

    r := gin.Default()

    // 测试路由
    r.GET("/ping", func(c *gin.Context) {
        c.String(http.StatusOK, "pong")
    })

    // 启动服务器
    r.Run(":" + global.App.Config.App.Port)
}
```

执行 `go run main.go` ，启动应用，服务器监听的端口是已经是配置文件里的端口号了.

尝试访问：http://127.0.0.1:8888/ping 测试路由