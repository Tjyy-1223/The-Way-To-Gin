## 第四章 Gin 基础 - 自定义错误处理

项目代码 [ Gin 基础 - 自定义错误处理](https://github.com/Tjyy-1223/The-Way-To-Gin/tree/main/code/1-Gin%E5%9F%BA%E7%A1%80/1-4-%E8%87%AA%E5%AE%9A%E4%B9%89%E9%94%99%E8%AF%AF%E5%A4%84%E7%90%86/GinDemo)

### 1 错误处理的意义

在 Gin 框架中，自定义错误处理的作用主要是对 HTTP 请求中的异常情况进行管理和响应。通过自定义错误处理，开发者可以灵活地控制错误的响应方式、格式以及错误信息的输出，增强应用的稳定性和可维护性。

具体来说，自定义错误处理可以实现以下几个作用：

#### 1.1 统一的错误响应格式

自定义错误处理可以确保所有的错误响应具有一致的结构，这有助于客户端处理错误。例如，所有的错误响应可以统一为以下格式：

```json
{
  "status": "error",
  "message": "具体错误信息",
  "code": 错误代码
}
```

这样客户端可以通过统一的格式进行错误处理。

#### 1.2 增强可读性和可维护性

自定义错误处理帮助开发者集中管理错误响应的逻辑。比如，可以在一个地方定义所有错误的处理方式，避免了在多个地方编写重复的错误处理代码。

#### 1.3 错误分类和处理

通过自定义错误处理，开发者可以根据不同的错误类型（如验证错误、数据库错误、第三方 API 错误等）进行不同的处理。例如：

- 业务逻辑错误返回 400（Bad Request）
- 身份验证错误返回 401（Unauthorized）
- 服务器内部错误返回 500（Internal Server Error）

#### 1.4 添加自定义的错误信息

Gin 默认的错误处理可能比较简单，只返回 HTTP 状态码和一些简短的信息。通过自定义错误处理，可以添加更多的上下文信息（如调试信息、错误发生的源头等），帮助开发者快速定位问题。

#### 1.5 错误恢复与请求的继续处理

在某些情况下，可以通过自定义错误处理进行错误恢复，使得服务器能够在某些异常发生时继续处理其他请求。例如，某个请求可能会发生错误，但应用仍然可以继续响应其他的请求而不会完全崩溃。

#### 1.6 示例

在 Gin 中，通常可以通过中间件来实现自定义的错误处理。以下是一个简单的示例：



### 2 自定义错误处理实践

参考文章：[自定义错误处理](https://github.com/xinliangnote/Go/blob/master/01-Gin%E6%A1%86%E6%9E%B6/05-%E8%87%AA%E5%AE%9A%E4%B9%89%E9%94%99%E8%AF%AF%E5%A4%84%E7%90%86.md)

#### 2.1 默认错误处理

默认的错误处理是 `errors.New("错误信息")`，这个信息通过 error 类型的返回值进行返回。

举个简单的例子：

```go
func hello(name string) (str string, err error) {
	if name == "" {
		err = errors.New("name 不能为空")
		return
	}
	str = fmt.Sprintf("hello: %s", name)
	return
}
```

当调用这个方法时：

```go
var name = ""
str, err :=  hello(name)
if err != nil {
	fmt.Println(err.Error())
	return
}
```

这就是默认的错误处理，下面还会用这个例子进行说。这个默认的错误处理，只是得到了一个错误信息的字符串。

然而，对于一些时间场景中

+ 可能还想得到发生错误时的 `时间`、`文件名`、`方法名`、`行号` 等信息。
+ 又或者还想得到错误时进行告警，比如 `短信告警`、`邮件告警`、`微信告警` 等。

最后，我还想调用的时候，不那么复杂，就和默认错误处理类似，比如：

```go
alarm.WeChat("错误信息")
return
```

这样，我们就得到了我们想要的信息（`时间`、`文件名`、`方法名`、`行号`），并通过 `微信` 的方式进行告警通知我们。

同理，`alarm.Email("错误信息")`、`alarm.Sms("错误信息")` 我们得到的信息是一样的，只是告警方式不同而已。

还要保证，我们业务逻辑中，获取错误的时候，只获取错误信息即可。

**上面这些想出来的，就是今天要实现的，自定义错误处理。**



#### 2.2 自定义错误处理

##### 2.2.1 errors.go

**首先看下面这段代码：**

```go
package main

import (
	"errors"
	"fmt"
)

func hello(name string) (str string, err error) {
	if name == "" {
		err = errors.New("name 不能为空")
		return
	}
	str = fmt.Sprintf("hello: %s", name)
	return
}

func main() {
	var name = ""
	fmt.Println("param:", name)

	str, err := hello(name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(str)
}
```

输出：

```
param: Tom
hello: Tom
```

当 name = "" 时，输出：

```
param:
name 不能为空
```

**建议每个函数都要有错误处理，error 应该为最后一个返回值。**

咱们一起看下官方 errors.go

```go
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package errors implements functions to manipulate errors.
package errors

// New returns an error that formats as the given text.
func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
```

上面的代码，并不复杂，参照上面的，咱们进行写一个自定义错误处理。



##### 2.2.2 自定义error

咱们定义一个 alarm.go，用于处理告警。

```go
package alarm

import (
	"encoding/json"
	"fmt"
	"ginDemo/common/function"
	"path/filepath"
	"runtime"
	"strings"
)

type errorString struct {
	s string
}

type errorInfo struct {
	Time     string `json:"time"`
	Alarm    string `json:"alarm"`
	Message  string `json:"message"`
	Filename string `json:"filename"`
	Line     int    `json:"line"`
	Funcname string `json:"funcname"`
}

func (e *errorString) Error() string {
	return e.s
}

func New (text string) error {
	alarm("INFO", text)
	return &errorString{text}
}

// 发邮件
func Email (text string) error {
	alarm("EMAIL", text)
	return &errorString{text}
}

// 发短信
func Sms (text string) error {
	alarm("SMS", text)
	return &errorString{text}
}

// 发微信
func WeChat (text string) error {
	alarm("WX", text)
	return &errorString{text}
}

// 告警方法
func  alarm(level string, str string) {
	// 当前时间
	currentTime := function.GetTimeStr()

	// 定义 文件名、行号、方法名
	fileName, line, functionName := "?", 0 , "?"

	pc, fileName, line, ok := runtime.Caller(2)
	if ok {
		functionName = runtime.FuncForPC(pc).Name()
		functionName = filepath.Ext(functionName)
		functionName = strings.TrimPrefix(functionName, ".")
	}

	var msg = errorInfo {
		Time     : currentTime,
		Alarm    : level,
		Message  : str,
		Filename : fileName,
		Line     : line,
		Funcname : functionName,
	}

	jsons, errs := json.Marshal(msg)

	if errs != nil {
		fmt.Println("json marshal error:", errs)
	}

	errorJsonInfo := string(jsons)

	fmt.Println(errorJsonInfo)

	if level == "EMAIL" {
		// 执行发邮件

	} else if level == "SMS" {
		// 执行发短信

	} else if level == "WX" {
		// 执行发微信

	} else if level == "INFO" {
		// 执行记日志
	}
}
```

**看下如何调用：**

```go
package v1

import (
	"fmt"
	"ginDemo/common/alarm"
	"ginDemo/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddProduct(c *gin.Context)  {
	// 获取 Get 参数
	name := c.Query("name")

	var res = entity.Result{}

	str, err := hello(name)
	if err != nil {
		res.SetCode(entity.CODE_ERROR)
		res.SetMessage(err.Error())
		c.JSON(http.StatusOK, res)
		c.Abort()
		return
	}

	res.SetCode(entity.CODE_SUCCESS)
	res.SetMessage(str)
	c.JSON(http.StatusOK, res)
}

func hello(name string) (str string, err error) {
	if name == "" {
		err = alarm.WeChat("name 不能为空")
		return
	}
	str = fmt.Sprintf("hello: %s", name)
	return
}
```

访问：`http://localhost:8080/v1/product/add?name=a`

```
{
    "code": 1,
    "msg": "hello: a",
    "data": null
}
```

未抛出错误，不会输出信息。

访问：`http://localhost:8080/v1/product/add`

```
{
    "code": -1,
    "msg": "name 不能为空",
    "data": null
}
```

抛出了错误，输出信息如下：

```
{"time":"2019-07-23 22:19:17","alarm":"WX","message":"name 不能为空","filename":"绝对路径/ginDemo/router/v1/product.go","line":33,"funcname":"hello"}
```

到这里，报错时我们收到了 `时间`、`错误信息`、`文件名`、`行号`、`方法名` 了。调用起来，也比较简单。

**如何继续进行告警通知：**

+ 在这里存储数据到队列中，再执行异步任务具体去消耗，这块就不实现了，大家可以去完善。

读取 `文件名`、`方法名`、`行号` 使用的是 `runtime.Caller()`。



##### 2.2.3 panic 和 recover

我们还知道，Go 有 `panic` 和 `recover`，它们是干什么的呢，接下来咱们就说说。

+ **当程序不能继续运行的时候，才应该使用 panic 抛出错误。**
+ 当程序发生 panic 后，在 defer(延迟函数) 内部可以调用 recover 进行控制，不过有个前提条件，只有在相同的 Go 协程中才可以。

panic 分两个，一种是有意抛出的，一种是无意的写程序马虎造成的，咱们一个个说。

**有意抛出的 panic：**

```go
package main

import (
	"fmt"
)

func main() {

	fmt.Println("-- 1 --")

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panic: %s\n", r)
		}
		fmt.Println("-- 2 --")
	}()
	
	panic("i am panic")
}
```

输出：

```
-- 1 --
panic: i am panic
-- 2 --
```

**无意抛出的 panic：**

```go
package main

import (
	"fmt"
)

func main() {

	fmt.Println("-- 1 --")

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panic: %s\n", r)
		}
		fmt.Println("-- 2 --")
	}()


	var slice = [] int {1, 2, 3, 4, 5}

	slice[6] = 6
}
```

输出：

```
-- 1 --
panic: runtime error: index out of range
-- 2 --
```

**上面的两个我们都通过 `recover` 捕获到了，那我们如何在 Gin 框架中使用呢？如果收到 `panic` 时，也想进行告警怎么实现呢？**

既然想实现告警，先在 ararm.go 中定义一个 `Panic()` 方法，当项目发生 `panic` 异常时，调用这个方法，这样就实现告警了。

```go
// Panic 异常
func Panic (text string) error {
	alarm("PANIC", text)
	return &errorString{text}
}
```

那我们怎么捕获到呢？

使用中间件进行捕获，写一个 `recover` 中间件。

```go
package recover

import (
	"fmt"
	"ginDemo/common/alarm"
	"github.com/gin-gonic/gin"
)

func Recover()  gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				alarm.Panic(fmt.Sprintf("%s", r))
			}
		}()
		c.Next()
	}
}
```

路由调用中间件：

```go
r.Use(logger.LoggerToFile(), recover.Recover())

//Use 可以传递多个中间件。
```

**验证下吧，咱们先抛出两个异常，看看能否捕获到？**

**还是修改 product.go 这个文件吧。**

**有意抛出 panic：**

```go
package v1

import (
	"fmt"
	"ginDemo/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddProduct(c *gin.Context)  {
	// 获取 Get 参数
	name := c.Query("name")

	var res = entity.Result{}

	str, err := hello(name)
	if err != nil {
		res.SetCode(entity.CODE_ERROR)
		res.SetMessage(err.Error())
		c.JSON(http.StatusOK, res)
		c.Abort()
		return
	}

	res.SetCode(entity.CODE_SUCCESS)
	res.SetMessage(str)
	c.JSON(http.StatusOK, res)
}

func hello(name string) (str string, err error) {
	if name == "" {
		// 有意抛出 panic
		panic("i am panic")
		return
	}
	str = fmt.Sprintf("hello: %s", name)
	return
}
```

访问：`http://localhost:8080/v1/product/add`

界面是空白的。

抛出了异常，输出信息如下：

```
{"time":"2019-07-23 22:42:37","alarm":"PANIC","message":"i am panic","filename":"绝对路径/ginDemo/middleware/recover/recover.go","line":13,"funcname":"1"}
```

很显然，定位的文件名、方法名、行号不是我们想要的。

需要调整 `runtime.Caller(2)`，这个代码在 `alarm.go 的 alarm` 方法中。

**将 2 调整成 4 ，看下输出信息：**

```
{"time":"2019-07-23 22:45:24","alarm":"PANIC","message":"i am panic","filename":"绝对路径/ginDemo/router/v1/product.go","line":33,"funcname":"hello"}
```

这就对了。

> 调用栈的深度（skip 值）：
>
> - `runtime.Caller(2)` 表示你要获取调用栈中**第 2 层调用者的信息**。
> - `runtime.Caller(4)` 表示你要获取调用栈中**第 4 层调用者的信息**。

**无意抛出 panic：**

```go
// 上面代码不变

func hello(name string) (str string, err error) {
	if name == "" {
		// 无意抛出 panic
		var slice = [] int {1, 2, 3, 4, 5}
		slice[6] = 6
		return
	}
	str = fmt.Sprintf("hello: %s", name)
	return
}
```

界面是空白的。

抛出了异常，输出信息如下：

```
{"time":"2019-07-23 22:50:06","alarm":"PANIC","message":"runtime error: index out of range","filename":"绝对路径/runtime/panic.go","line":44,"funcname":"panicindex"}
```

很显然，定位的文件名、方法名、行号也不是我们想要的。

将 4 调整成 5 ，看下输出信息：

```
{"time":"2019-07-23 22:55:27","alarm":"PANIC","message":"runtime error: index out of range","filename":"绝对路径/ginDemo/router/v1/product.go","line":34,"funcname":"hello"}
```

这就对了。