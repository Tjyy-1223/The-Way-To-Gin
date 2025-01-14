# 第七章 Gin 进阶 - 从零封装属于自己的 Gin 框架 - 日志初始化

学习资料参考：[手把手，带你从零封装Gin框架](https://juejin.cn/post/7016742808560074783)

本篇来讲一下怎么将日志服务集成到项目中，它也是框架中必不可少的，**平时代码调试，线上 Bug 分析都离不开它。**

这里将使用 [zap](https://link.juejin.cn?target=https%3A%2F%2Fgithub.com%2Fuber-go%2Fzap) 作为日志库，一般来说，日志都是需要写入到文件保存的，这也是 zap 唯一缺少的部分，所以我将结合 [lumberjack](https://link.juejin.cn?target=https%3A%2F%2Fgithub.com%2Fnatefinch%2Flumberjack) 来使用，实现日志切割归档的功能。



## 7.1 zap

Zap 是 Go 语言中的一个高性能、结构化的日志库，由 Uber 开发。它的特点是高效、易用，并且支持结构化日志输出。与传统的日志库（如 `log`）相比，Zap 具有更高的性能和更多的功能。

1. **高性能**：Zap 的设计目标之一就是性能。它通过预先分配内存和避免不必要的内存分配来优化性能，非常适合高频调用的场景。
2. **结构化日志**：支持输出结构化日志，允许以键值对的方式记录日志数据。这使得日志可以更方便地进行查询、分析和处理，尤其是在使用像 Elasticsearch、Prometheus 等日志收集系统时。
3. **级别控制**：支持不同级别的日志（如 `Debug`、`Info`、`Warn`、`Error`、`DPanic`、`Panic`、`Fatal`），并且可以灵活配置。
4. **JSON 格式输出**：Zap 默认支持 JSON 格式输出，也可以输出常见的日志格式（例如普通文本），根据应用的需求来选择。
5. **异步日志**：通过 `zapcore`，Zap 支持异步日志，进一步提升性能。



## 7.2 Lumberjack

Lumberjack 是 Go 语言中的一个日志文件轮换库，用于管理日志文件的滚动、压缩、删除等操作。它特别适合需要长时间运行的服务（例如 Web 服务器），可以确保日志文件不会无限制增长而占用大量磁盘空间。

1. **日志文件滚动**：当日志文件达到指定大小时，Lumberjack 会自动将当前日志文件重命名并创建一个新的日志文件，继续写入日志。
2. **压缩旧日志**：可以设置旧的日志文件（如 `log.1`、`log.2` 等）进行压缩，通常压缩为 `.gz` 格式，以节省磁盘空间。
3. **保留数量控制**：可以指定最多保留多少个备份日志文件。例如，可以配置最多保留 7 个文件，超过的旧日志文件会被删除。
4. **时间限制**：可以设置基于时间的滚动（例如每天滚动），而不仅仅是基于文件大小。
5. **支持并发**：Lumberjack 适合多线程或高并发环境下使用，可以确保日志写入的可靠性。



## 7.3 使用实践

#### 7.3.1 安装

```
go get -u go.uber.org/zap
go get -u gopkg.in/natefinch/lumberjack.v2
```



#### 7.3.2 定义日志配置项

新建 `config/log.go` 文件，定义 `zap` 和 `lumberjack` 初始化需要使用的配置项，大家可以根据自己的喜好去定制

```go
package config

type Log struct {
    Level string `mapstructure:"level" json:"level" yaml:"level"`
    RootDir string `mapstructure:"root_dir" json:"root_dir" yaml:"root_dir"`
    Filename string `mapstructure:"filename" json:"filename" yaml:"filename"`
    Format string `mapstructure:"format" json:"format" yaml:"format"`
    ShowLine bool `mapstructure:"show_line" json:"show_line" yaml:"show_line"`
    MaxBackups int `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
    MaxSize int `mapstructure:"max_size" json:"max_size" yaml:"max_size"` // MB
    MaxAge int `mapstructure:"max_age" json:"max_age" yaml:"max_age"` // day
    Compress bool `mapstructure:"compress" json:"compress" yaml:"compress"`
}
```

`config/config.go` 添加 `Log` 成员属性

```go
package config

type Configuration struct {
    App App `mapstructure:"app" json:"app" yaml:"app"`
    Log Log `mapstructure:"log" json:"log" yaml:"log"`
}
```

`config.yaml` 增加对应配置项

```go
log:
  level: info # 日志等级
  root_dir: ./storage/logs # 日志根目录
  filename: app.log # 日志文件名称
  format: # 写入格式 可选json
  show_line: true # 是否显示调用行
  max_backups: 3 # 旧文件的最大个数
  max_size: 500 # 日志文件最大大小（MB）
  max_age: 28 # 旧文件的最大保留天数
  compress: true # 是否压缩
```



#### 7.3.3 定义 utils 工具函数

新建 `utils/directory.go` 文件，编写 `PathExists` 函数，用于判断路径是否存在

```go
package utils

import "os"

func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}
```



#### 7.3.4 初始化 zap

`zap` 的具体使用说明可查看[官方文档](https://link.juejin.cn/?target=https%3A%2F%2Fpkg.go.dev%2Fgo.uber.org%2Fzap)

新建 `bootstrap/log.go` 文件，编写：

```go
package bootstrap

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
    "jassue-gin/global"
    "jassue-gin/utils"
    "os"
    "time"
)

var (
    level zapcore.Level // zap 日志等级
    options []zap.Option // zap 配置项
)

func InitializeLog() *zap.Logger {
    // 创建根目录
    createRootDir()

    // 设置日志等级
    setLogLevel()

    if global.App.Config.Log.ShowLine {
        options = append(options, zap.AddCaller())
    }

    // 初始化 zap
    return zap.New(getZapCore(), options...)
}

func createRootDir() {
    if ok, _ := utils.PathExists(global.App.Config.Log.RootDir); !ok {
        _ = os.Mkdir(global.App.Config.Log.RootDir, os.ModePerm)
    }
}

func setLogLevel() {
    switch global.App.Config.Log.Level {
    case "debug":
        level = zap.DebugLevel
        options = append(options, zap.AddStacktrace(level))
    case "info":
        level = zap.InfoLevel
    case "warn":
        level = zap.WarnLevel
    case "error":
        level = zap.ErrorLevel
        options = append(options, zap.AddStacktrace(level))
    case "dpanic":
        level = zap.DPanicLevel
    case "panic":
        level = zap.PanicLevel
    case "fatal":
        level = zap.FatalLevel
    default:
        level = zap.InfoLevel
    }
}

// 扩展 Zap
func getZapCore() zapcore.Core {
    var encoder zapcore.Encoder

    // 调整编码器默认配置
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
        encoder.AppendString(time.Format("[" + "2006-01-02 15:04:05.000" + "]"))
    }
    encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
        encoder.AppendString(global.App.Config.App.Env + "." + l.String())
    }

    // 设置编码器
    if global.App.Config.Log.Format == "json" {
        encoder = zapcore.NewJSONEncoder(encoderConfig)
    } else {
        encoder = zapcore.NewConsoleEncoder(encoderConfig)
    }

    return zapcore.NewCore(encoder, getLogWriter(), level)
}

// 使用 lumberjack 作为日志写入器
func getLogWriter() zapcore.WriteSyncer {
    file := &lumberjack.Logger{
        Filename:   global.App.Config.Log.RootDir + "/" + global.App.Config.Log.Filename,
        MaxSize:    global.App.Config.Log.MaxSize,
        MaxBackups: global.App.Config.Log.MaxBackups,
        MaxAge:     global.App.Config.Log.MaxAge,
        Compress:   global.App.Config.Log.Compress,
    }

    return zapcore.AddSync(file)
}
```



#### 7.3.5 定义全局变量 Log

在 `global/app.go` 中，添加 `Log` 成员属性

```go
package global

import (
    "github.com/spf13/viper"
    "go.uber.org/zap"
    "jassue-gin/config"
)

type Application struct {
    ConfigViper *viper.Viper
    Config config.Configuration
    Log *zap.Logger
}

var App = new(Application)
```



#### 7.3.6 测试

在 `main.go` 中调用日志初始化函数，并尝试写入日志

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

    // 初始化日志
    global.App.Log = bootstrap.InitializeLog()
    global.App.Log.Info("log init success!")

    r := gin.Default()

    // 测试路由
    r.GET("/ping", func(c *gin.Context) {
        c.String(http.StatusOK, "pong")
    })

    // 启动服务器
    r.Run(":" + global.App.Config.App.Port)
}
```

启动 `main.go` ，生成 `storage/logs/app.log` 文件，表示日志初始化成功，文件内容显示如下：

```
[2021-10-12 19:17:46.997]	local.info	jassue-gin/main.go:16	log init success!
```



#### 7.3.7 补充

把日志即写入文件又打印到控制台呢，开发过程中主要还是看控制台日志

使用NewMultiWriteSyncer 方法加入多个日志写入器

```go
var writes = []zapcore.WriteSyncer{getLogWriter(), zapcore.AddSync(os.Stdout)}
return zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writes...), level)
```

