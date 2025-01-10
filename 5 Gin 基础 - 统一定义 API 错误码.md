## 第五章 Gin 基础 - 统一定义 API 错误码

### 5.1 API 错误码意义

在 Gin 或任何其他 Web 框架中，统一定义 API 错误码（即 HTTP 状态码和自定义的错误码）是一个良好的开发实践，具有多个重要的好处，主要体现在以下几个方面：

+ **一致的响应结构**：通过统一的错误码，客户端可以更容易地解析和处理错误。如果错误响应总是遵循一个固定的格式，客户端无需针对每一个不同的错误做不同的判断和处理逻辑。
+ **易于维护的 API**：无论是前端开发人员还是后端开发人员，了解所有 API 错误码及其含义有助于协作，并且在出现问题时更容易追踪。

通过统一的错误码，可以清晰地区分不同类别的错误，比如：

- **用户错误**（如请求参数缺失、格式错误）通常使用 4xx 状态码。
- **服务器错误**（如数据库连接失败、内部异常）通常使用 5xx 状态码。
- 特定业务逻辑错误可以通过自定义错误码进一步细分，例如 `1001` 表示“用户名已存在”，`1002` 表示“密码不符合要求”等。

**以下是一个简单的统一错误码设计示例：**

```go
// 定义常见的 HTTP 状态码及业务错误码
const (
    // HTTP 状态码
    HTTPStatusOK                  = 200
    HTTPStatusBadRequest          = 400
    HTTPStatusUnauthorized        = 401
    HTTPStatusForbidden           = 403
    HTTPStatusNotFound            = 404
    HTTPStatusInternalServerError = 500

    // 自定义业务错误码
    ErrCodeUserAlreadyExists      = 1001
    ErrCodeInvalidPassword        = 1002
    ErrCodeInvalidParams          = 1003
    ErrCodeDatabaseError          = 2001
)
```

在 API 错误响应中，你可以返回类似以下的结构：

```json
{
    "status": "error",
    "message": "用户名已存在",
    "code": 1001
}
```



### 5.2 修改实践

参考文章：[统一定义API错误码](https://github.com/xinliangnote/Go/blob/master/01-Gin%E6%A1%86%E6%9E%B6/06-%E7%BB%9F%E4%B8%80%E5%AE%9A%E4%B9%89%20API%20%E9%94%99%E8%AF%AF%E7%A0%81.md)

#### 5.2.1 修改之前

在使用 `gin` 开发接口的时候，返回接口数据是这样写的。

```go
type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// always return http.StatusOK
c.JSON(http.StatusOK, response{
	Code: 20101,
	Msg:  "用户手机号不合法",
	Data: nil,
})
```

这种写法 `code`、`msg` 都是**在哪需要返回在哪定义** ，没有进行统一管理。

#### 5.2.2 修改之后

```go
// 比如，返回“用户手机号不合法”错误
c.JSON(http.StatusOK, errno.ErrUserPhone.WithID(c.GetString("trace-id")))

// 正确返回
c.JSON(http.StatusOK, errno.OK.WithData(data).WithID(c.GetString("trace-id")))
```

`errno.ErrUserPhone`、`errno.OK` 表示自定义的错误码，下面会看到定义的地方。

+ `.WithID()` 设置当前请求的唯一ID，也可以理解为链路ID，忽略也可以。

+ `.WithData()` 设置成功时返回的数据。

下面分享下编写的 `errno` 包源码，非常简单，希望大家不要介意。

> **Gin Context**: Gin 的 `Context` 是处理 HTTP 请求的核心结构，它不仅包含了请求和响应的详细信息，还提供了一个 key-value 存储机制来存储和传递数据。你可以将数据存储在 `Context` 中，在整个请求处理周期中进行传递。
>
> **`context.GetString()`**:
>
> - `GetString(key string)` 方法用于从 `Context` 中获取一个与给定 `key` 相关联的值。如果存在该 `key`，并且其值是字符串类型，`GetString()` 将返回该字符串。如果该 `key` 不存在，则返回空字符串 `""` 和 `false`（表示未找到）。
> - 在本例中，`context.GetString("trace-id")` 试图从 `Context` 中获取名为 `"trace-id"` 的字符串值。
>
> 在实际应用中，通常会通过中间件为每个请求生成一个唯一的 `trace-id`，并将其存储在 `Context` 中。一个常见的做法是使用 UUID 或其他方法生成一个唯一的 `trace-id`，然后将它存入 `Context`。
>
> ```go
> package main
> 
> import (
>     "github.com/gin-gonic/gin"
>     "github.com/google/uuid"
>     "log"
> )
> 
> // 中间件：生成并存储 trace-id
> func TraceIDMiddleware() gin.HandlerFunc {
>     return func(c *gin.Context) {
>         // 生成一个新的 trace-id
>         traceID := uuid.New().String()
> 
>         // 将 trace-id 存入 Context 中
>         c.Set("trace-id", traceID)
> 
>         // 将 trace-id 放到响应头中（可选）
>         c.Header("X-Trace-ID", traceID)
> 
>         // 继续请求处理
>         c.Next()
>     }
> }
> 
> func main() {
>     router := gin.Default()
> 
>     // 使用中间件
>     router.Use(TraceIDMiddleware())
> 
>     // 处理请求并获取 trace-id
>     router.GET("/hello", func(c *gin.Context) {
>         // 从 Context 获取 trace-id
>         traceID := c.GetString("trace-id")
>         log.Printf("Trace ID: %s", traceID)
> 
>         // 返回 trace-id
>         c.JSON(200, gin.H{
>             "message": "Hello, World!",
>             "trace-id": traceID,
>         })
>     })
> 
>     router.Run(":8080")
> }
> 
> ```



#### 5.2.3 errno 源码学习

```go
// errno/errno.go

package errno

import (
	"encoding/json"
)

var _ Error = (*err)(nil)

type Error interface {
	// i 为了避免被其他包实现
	i()
	// WithData 设置成功时返回的数据
	WithData(data interface{}) Error
	// WithID 设置当前请求的唯一ID
	WithID(id string) Error
	// ToString 返回 JSON 格式的错误详情
	ToString() string
}

type err struct {
	Code int         `json:"code"`         // 业务编码
	Msg  string      `json:"msg"`          // 错误描述
	Data interface{} `json:"data"`         // 成功时返回的数据
	ID   string      `json:"id,omitempty"` // 当前请求的唯一ID，便于问题定位，忽略也可以
}

func NewError(code int, msg string) Error {
	return &err{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

func (e *err) i() {}

func (e *err) WithData(data interface{}) Error {
	e.Data = data
	return e
}

func (e *err) WithID(id string) Error {
	e.ID = id
	return e
}

// ToString 返回 JSON 格式的错误详情
func (e *err) ToString() string {
	err := &struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
		ID   string      `json:"id,omitempty"`
	}{
		Code: e.Code,
		Msg:  e.Msg,
		Data: e.Data,
		ID:   e.ID,
	}

	raw, _ := json.Marshal(err)
	return string(raw)
}

// errno/code.go

package errno

var (
	// OK
	OK = NewError(0, "OK")

	// 服务级错误码
	ErrServer    = NewError(10001, "服务异常，请联系管理员")
	ErrParam     = NewError(10002, "参数有误")
	ErrSignParam = NewError(10003, "签名参数有误")

	// 模块级错误码 - 用户模块
	ErrUserPhone   = NewError(20101, "用户手机号不合法")
	ErrUserCaptcha = NewError(20102, "用户验证码有误")

	// ...
)
```

错误码规则：

- 错误码需在 `code.go` 文件中定义。
- 错误码需为 > 0 的数，反之表示正确。

错误码为 5 位数

| 1            | 01           | 01         |
| ------------ | ------------ | ---------- |
| 服务级错误码 | 模块级错误码 | 具体错误码 |

- 服务级别错误码：1 位数进行表示，比如 1 为系统级错误；2 为普通错误，通常是由用户非法操作引起。
- 模块级错误码：2 位数进行表示，比如 01 为用户模块；02 为订单模块。
- 具体错误码：2 位数进行表示，比如 01 为手机号不合法；02 为验证码输入错误。