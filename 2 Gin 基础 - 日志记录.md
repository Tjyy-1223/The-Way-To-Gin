## 第二章 Gin 基础 - 日志记录

### 1 Gin 日志使用

在 **Gin** 中输出日志有多种方式，包括使用 Gin 自带的日志功能、通过自定义中间件来输出日志、或使用第三方日志库。下面是几种常见的输出日志的方法。

#### 1.1 使用 Gin 自带的日志功能

Gin 默认会记录所有请求的日志，并输出到标准输出（通常是控制台）。每个请求的日志会包括请求方法、路径、响应状态码、响应时间等信息。

```go
import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default() // gin.Default() 会自动使用默认的日志中间件
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    r.Run(":8080")
}
```

`gin.Default()` 使用了 Gin 内建的中间件 `Logger`，它会自动记录所有请求的信息：

```csharp
[GIN] 2025/01/07 - 15:04:32 | 200 |   10.123ms |   127.0.0.1:8080 | GET     "/ping"
```

日志包含了以下信息：

- 请求时间（例如：`2025/01/07 - 15:04:32`）
- HTTP 响应状态码（例如：`200`）
- 请求处理耗时（例如：`10.123ms`）
- 请求来源的 IP 地址（例如：`127.0.0.1`）
- 请求的 HTTP 方法和路径（例如：`GET /ping`）



#### 1.2 自定义中间件输出日志

如果你想自定义日志格式或更细粒度地控制日志输出（例如记录请求的参数或处理过程中的其他信息），你可以创建自定义中间件。示例：

```go
import (
    "fmt"
    "github.com/gin-gonic/gin"
)

func CustomLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 请求开始时的日志
        fmt.Printf("Start Request: %s %s\n", c.Request.Method, c.Request.URL.Path)
        
        // 继续处理请求
        c.Next()

        // 请求结束后的日志
        fmt.Printf("End Request: %s %s - Status: %d\n", c.Request.Method, c.Request.URL.Path, c.Writer.Status())
    }
}

func main() {
    r := gin.Default()

    // 使用自定义日志中间件
    r.Use(CustomLogger())

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    r.Run(":8080")
}
```

**输出示例**：

```
Start Request: GET /ping
End Request: GET /ping - Status: 200
```



#### 1.3 使用 `log` 标准库输出日志

Gin 的上下文对象（`*gin.Context`）提供了对请求和响应的访问，你可以通过 `log` 标准库输出日志。示例：

```go
import (
    "log"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // 使用自定义日志
    r.Use(func(c *gin.Context) {
        log.Printf("Received %s request for %s\n", c.Request.Method, c.Request.URL.Path)
        c.Next()
        log.Printf("Response status: %d\n", c.Writer.Status())
    })

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    r.Run(":8080")
}
```

**控制台输出**：

```
2025/01/07 15:04:32 Received GET request for /ping
2025/01/07 15:04:32 Response status: 200
```



#### 1.4 使用第三方日志库（如 `logrus` 或 `zap`）

如果你需要更复杂的日志功能（如日志级别、输出到文件、结构化日志等），可以集成第三方日志库，比如 `logrus` 或 `zap`。

**示例：使用 `logrus` 输出日志**

```go
import (
    "github.com/gin-gonic/gin"
    log "github.com/sirupsen/logrus"
)

func main() {
    // 配置 logrus
    log.SetFormatter(&log.TextFormatter{
        FullTimestamp: true,
    })

    r := gin.Default()

    // 使用 logrus 作为日志输出
    r.Use(func(c *gin.Context) {
        log.WithFields(log.Fields{
            "method": c.Request.Method,
            "path":   c.Request.URL.Path,
        }).Info("Request received")
        c.Next()
        log.WithFields(log.Fields{
            "status": c.Writer.Status(),
        }).Info("Response sent")
    })

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    r.Run(":8080")
}
```

**控制台输出示例**：

```
time="2025-01-07T15:04:32+08:00" level=info msg="Request received" method=GET path=/ping
time="2025-01-07T15:04:32+08:00" level=info msg="Response sent" status=200
```



#### 1.5 日志配置和输出到文件

你还可以通过配置日志库将日志输出到文件而不是标准输出（控制台）。

**示例：将日志输出到文件（使用 `logrus`）**

```go
import (
    "github.com/gin-gonic/gin"
    log "github.com/sirupsen/logrus"
    "os"
)

func main() {
    // 打开文件进行日志输出
    file, err := os.OpenFile("gin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal(err)
    }
    log.SetOutput(file)
    log.SetFormatter(&log.TextFormatter{
        FullTimestamp: true,
    })

    r := gin.Default()

    // 使用 logrus 作为日志输出
    r.Use(func(c *gin.Context) {
        log.WithFields(log.Fields{
            "method": c.Request.Method,
            "path":   c.Request.URL.Path,
        }).Info("Request received")
        c.Next()
        log.WithFields(log.Fields{
            "status": c.Writer.Status(),
        }).Info("Response sent")
    })

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    r.Run(":8080")
}
```

此时，所有日志信息会输出到 `gin.log` 文件中。



#### 1.6 自定义日志格式

如果你想进一步自定义日志的格式，可以通过 `gin.LoggerWithFormatter()` 方法来指定自定义格式的日志输出。

**示例：自定义日志格式**

```go
import (
    "github.com/gin-gonic/gin"
    "log"
    "strings"
)

func main() {
    // 自定义日志格式
    gin.DefaultWriter = log.Writer()
    r := gin.New()

    // 自定义 Logger 中间件
    r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        // 自定义日志格式：请求方法、路径、状态码、响应时间
        return strings.Join([]string{
            "[Custom Log]",
            param.TimeStamp.Format("2006/01/02 - 15:04:05"),
            param.Method,
            param.Path,
            "Status:", param.StatusCode,
            "Duration:", param.Latency.String(),
        }, " ") + "\n"
    }))

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    r.Run(":8080")
}
```

**自定义日志输出示例：**

```
[Custom Log] 2025/01/07 - 15:04:32 GET /ping Status: 200 Duration: 10.123ms
```



#### 1.7 总结

- **Gin 自带的日志功能**：通过 `gin.Default()` 自动启用，记录基本请求日志。
- **自定义中间件日志**：通过 `r.Use(...)` 定义自定义中间件，可以根据需要输出更多的日志信息。
- **标准库 `log`**：可以直接使用 Go 的 `log` 包输出日志。
- **第三方日志库（如 `logrus`、`zap` 等）**：可以集成这些库，提供更强大、更灵活的日志功能，包括日志级别、输出格式、输出到文件等。

根据实际需求，可以选择合适的方式来记录和管理日志。



### 2 Gin 日志实践 - logrus

查了很多资料，Go 的日志记录用的最多的还是 `github.com/sirupsen/logrus`。

> Logrus is a structured logger for Go (golang), completely API compatible with the standard library logger.

Gin 框架的日志默认只会在控制台输出，咱们利用 `Logrus` 封装一个中间件，将日志记录到文件中。

这篇文章就是学习和使用 `Logrus` 。



#### 2.1 日志格式

比如，我们约定日志格式为 Text，包含字段如下：

`请求时间`、`日志级别`、`状态码`、`执行时间`、`请求IP`、`请求方式`、`请求路由`。

接下来，咱们利用 `Logrus` 实现它。



#### 2.2 Logrus 使用

首先需要对 logrus 进行安装

```
import "github.com/sirupsen/logrus"
```

为了方便的使用 logrus，我们将其提供的日志功能设置成一个中间件，如项目中的 logger.go

+ 日志可以记录到 File 中，定义一个 `LoggerToFile` 方法。
+ 日志可以记录到 MongoDB 中，定义一个 `LoggerToMongo` 方法。
+ 日志可以记录到 ES 中，定义一个 `LoggerToES` 方法。
+ 日志可以记录到 MQ 中，定义一个 `LoggerToMQ` 方法。

**实现 `LoggerToFile` 方法，其他的可以根据自己的需求进行实现。这个 `logger` 中间件，创建好了，可以任意在其他项目中进行迁移使用。代码框架如下：**

middleware/logger.go

```go
package middleware

import (
	"fmt"
	"ginDemo/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {

	logFilePath := config.Log_FILE_PATH
	logFileName := config.LOG_FILE_NAME

	//日志文件
	fileName := path.Join(logFilePath, logFileName)

	//写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	//实例化
	logger := logrus.New()

	//设置输出
	logger.Out = src

	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{})

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

// 日志记录到 MongoDB
func LoggerToMongo() gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}

// 日志记录到 ES
func LoggerToES() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 MQ
func LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
```

日志中间件写好了，怎么调用呢？

**只需在 main.go 中新增：**

```go
engine := gin.Default() //在这行后新增
engine.Use(middleware.LoggerToFile())
```

你也可以在 `engine.Use()` 中添加多个不同的中间件，并且它们会按顺序执行。

```go
package main

import (
    "github.com/gin-gonic/gin"
    "yourapp/middleware"
)

func main() {
    engine := gin.Default()

    engine.Use(middleware.LoggerToFile("log1.txt"))
    engine.Use(middleware.LoggerToFile("log2.txt"))

    engine.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    engine.Run(":8080")
}
```

假设 `LoggerToFile` 中间件将日志写入文件 `log1.txt` 和 `log2.txt`，执行的顺序是：

1. 第一个 `LoggerToFile("log1.txt")` 中间件会写日志到 `log1.txt`。
2. 第二个 `LoggerToFile("log2.txt")` 中间件会写日志到 `log2.txt`。

需要注意的是，如果多个中间件都对相同的资源（例如日志文件）进行操作，可能会引发一些并发问题或日志重复问题。为避免这种情况，可以使用线程安全的方式（如互斥锁）来管理对日志文件的访问。

**运行一下，看看日志：**

```go
time="2019-07-17T22:10:45+08:00" level=info msg="| 200 |      27.698µs |             ::1 | GET | /v1/product/add?name=a&price=10 |"
time="2019-07-17T22:10:46+08:00" level=info msg="| 200 |      27.239µs |             ::1 | GET | /v1/product/add?name=a&price=10 |"
```

**这个 `time="2019-07-17T22:10:45+08:00"` ，这个时间格式不是咱们想要的，怎么办？**

时间需要格式化一下，修改 `logger.SetFormatter`

```go
//设置日志格式
logger.SetFormatter(&logrus.TextFormatter{
	TimestampFormat:"2006-01-02 15:04:05",
})
```

时间变得正常了。

**我不喜欢文本格式，喜欢 JSON 格式，怎么办？**

```go
//设置日志格式
logger.SetFormatter(&logrus.JSONFormatter{
	TimestampFormat:"2006-01-02 15:04:05",
})
```

执行以下，再看日志：

```
{"level":"info","msg":"| 200 |       24.78µs |             ::1 | GET | /v1/product/add?name=a\u0026price=10 |","time":"2019-07-17 22:23:55"}
{"level":"info","msg":"| 200 |      26.946µs |             ::1 | GET | /v1/product/add?name=a\u0026price=10 |","time":"2019-07-17 22:23:56"}
```

**msg 信息太多，不方便看，怎么办？**

```go
// 日志格式
logger.WithFields(logrus.Fields{
	"status_code"  : statusCode,
	"latency_time" : latencyTime,
	"client_ip"    : clientIP,
	"req_method"   : reqMethod,
	"req_uri"      : reqUri,
}).Info()
```

执行以下，再看日志：

```
{"client_ip":"::1","latency_time":26681,"level":"info","msg":"","req_method":"GET","req_uri":"/v1/product/add?name=a\u0026price=10","status_code":200,"time":"2019-07-17 22:37:54"}
{"client_ip":"::1","latency_time":24315,"level":"info","msg":"","req_method":"GET","req_uri":"/v1/product/add?name=a\u0026price=10","status_code":200,"time":"2019-07-17 22:37:55"}
```

说明一下：`time`、`msg`、`level` 这些参数是 logrus 自动加上的。



**logrus 支持输出文件名和行号吗？**

不支持，作者的回复是太耗性能。

不过网上也有人通过 Hook 的方式实现了，选择在生产环境使用的时候，记得做性能测试。

**logrus 支持日志分割吗？**

不支持，但有办法实现它。

1、可以利用 `Linux logrotate`，统一由运维进行处理。

2、可以利用 `file-rotatelogs` 实现

需要导入包：

```0
github.com/lestrrat-go/file-rotatelogs
github.com/rifflock/lfshook
```



#### 2.3 核心代码

##### 2.3.1 middleware/logger.go

```go
package middleware

import (
	"GinDemo/config"
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	logrus "github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

// LoggerToFile log the records into the file
func LoggerToFile() gin.HandlerFunc {
	logFilePath := config.LOG_FILE_PATH
	logFileName := config.LOG_FILE_NAME

	// log file
	fileName := path.Join(logFilePath, logFileName)

	// write into the file
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	logger := logrus.New()
	logger.Out = src                   // output file
	logger.SetLevel(logrus.DebugLevel) // log lever is DEBUG

	// set rotate logs for segment
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d.log",
		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// add a hook
	logger.AddHook(lfHook)
  
  // 返回的函数中使用了 logger 才是关键
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		reqMethod := c.Request.Method
		reqUrl := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIp := c.ClientIP()
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIp,
			"req_method":   reqMethod,
			"req_uri":      reqUrl,
		}).Info()
	}
}
```

##### 2.3.2 main.go

```go
package main

import (
	"GinDemo/config"
	"GinDemo/middleware"
	"GinDemo/router"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode) // 默认为 debug 模式，设置为发布模式
	engine := gin.Default()
	engine.Use(middleware.LoggerToFile())
	router.InitRouter(engine)
	engine.Run(config.PORT)
}
```

##### 2.3.3 输出测试

这时会新生成一个文件 `system.log.20250107.log`，日志内容与上面的格式一致。