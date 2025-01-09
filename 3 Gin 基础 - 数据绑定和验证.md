## 第三章 Gin 基础 - 数据绑定和验证

在 **Gin** 中，数据绑定和验证是指将 HTTP 请求中的数据（如 URL 参数、查询参数、请求体等）解析为 Go 结构体（或其他数据类型），并对其进行验证以确保数据符合预期格式或规则。

- **数据绑定**：Gin 会将请求中的数据（如 JSON、表单、查询参数等）解析并绑定到结构体上。
- **数据验证**：绑定数据后，Gin 会根据结构体中的 `binding` 标签验证数据是否合法（如必填、范围限制等）。
- **结合使用**：Gin 提供了 `ShouldBind` 系列方法，结合数据绑定与验证，确保请求的数据符合预期。
- **自定义验证**：通过 `validator` 库，开发者可以创建自定义的验证规则。

这种机制能够大大简化请求数据的解析和验证过程，提高 API 的健壮性和安全性。

### 3.1 数据绑定 - Binding

**数据绑定** 是指将客户端发来的请求数据（如 `JSON`, `Form`, `Query` 等）自动解析并绑定到 Go 结构体中的字段。Gin 提供了多种绑定方式，如绑定 JSON 数据、表单数据、URL 查询参数等。

**常见的绑定方式：**

- **JSON 绑定**：用于将请求体中的 JSON 数据绑定到 Go 结构体。
- **表单绑定**：用于将 `application/x-www-form-urlencoded` 类型的请求数据绑定到结构体。
- **查询参数绑定**：用于将 URL 查询字符串中的参数绑定到结构体。
- **路径参数绑定**：用于将 URL 路径中的参数绑定到结构体。
- **XML 绑定**：用于将 XML 格式的数据绑定到结构体。

**示例：JSON 绑定**

假设客户端发送一个携带 JSON 信息的请求：

```go
{
  "username": "john",
  "password": "secret"
}

```

可以通过 `ShouldBindJSON` 方法将其绑定到结构体：

```go
type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, gin.H{"message": "Login successful", "user": req.Username})
}

```

`c.ShouldBindJSON(&req)` 会将请求体中的 JSON 数据绑定到 `req` 结构体中。

如果请求数据不符合结构体字段的要求（比如缺少 `username` 或 `password` 字段），`ShouldBindJSON` 会返回错误。



### 3.2 数据验证 - Validation

**数据验证** 是指在数据绑定后，检查数据是否满足某些条件或规则。Gin 中的验证通常是通过结构体标签和 **`binding`** 标签来实现的。

Gin 内部基于 `github.com/go-playground/validator/v10` 实现数据验证。

**示例：验证数据**

在上面的 `LoginRequest` 结构体中，`binding` 标签就用于定义字段验证规则：

```go
type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}
```

+ `binding:"required"` 表示 `Username` 和 `Password` 字段不能为空。

Gin 会在绑定数据后自动执行这些验证规则。如果验证失败，Gin 会返回一个错误。

**示例：使用更多验证规则**

```go
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=30"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

func register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, gin.H{"message": "Registration successful", "user": req.Username})
}
```

在这个例子中，`RegisterRequest` 结构体包含了更多的验证规则：

- `Username`：必须是 3 到 30 个字符。
- `Email`：必须是有效的电子邮件地址。
- `Password`：必须至少 6 个字符。

#### 3.3.1 基本验证规则

- `required`：字段必须存在且不能为零值（如空字符串、零、空数组、空切片等）
- `omitempty`：字段为空时不做校验，适用于选择性验证
- `min`：字段值必须大于或等于指定的最小值（通常用于数字或字符串）
- `max`：字段值必须小于或等于指定的最大值（通常用于数字或字符串）
- `len`：字段长度必须等于指定值（用于字符串、切片、数组）

#### 3.3.2 字符串相关验证规则

- `min`：字符串最小长度
- `max`：字符串最大长度
- `len`：字符串固定长度
- `email`：字段必须是有效的电子邮件格式
- `url`：字段必须是有效的 URL 格式
- `uuid`：字段必须是有效的 UUID 格式
- `alpha`：字段只能包含字母（大小写均可）
- `alphanumeric`：字段只能包含字母和数字
- `numeric`：字段只能包含数字
- `hexadecimal`：字段必须是十六进制字符
- `lowercase`：字段必须是小写字母
- `uppercase`：字段必须是大写字母
- `contains`：字符串中必须包含指定的子串

#### 3.3.3 数值相关验证规则

- `min`：数字最小值
- `max`：数字最大值
- `gt`：大于指定值
- `lt`：小于指定值
- `eq`：等于指定值
- `ne`：不等于指定值
- `in`：值必须在指定的范围内（支持数组或范围）
- `not_in`：值不能在指定的范围内

#### 3.3.4 时间验证规则

- `datetime`：字段必须是有效的时间日期格式
- `before`：日期/时间必须早于指定时间
- `after`：日期/时间必须晚于指定时间

#### 3.3.5 数组和切片相关验证规则

- `dive`：对数组或切片的每个元素应用验证规则
- `unique`：数组或切片中的元素必须唯一
- `min`：切片或数组的最小长度
- `max`：切片或数组的最大长度

#### 3.3.6 结构体相关规则

- `valid`：校验结构体中的所有字段，适用于嵌套结构体
- `required`：结构体字段必须包含且不为零值
- `omitempty`：字段为空时不验证，适用于可选字段
- `dive`：递归验证结构体内部的字段

#### 3.3.7 其他常见规则

- `default`：设置默认值（通常结合 `binding` 使用）
- `exclusive`：排除某些字段的验证
- `exclude`：排除字段验证



### 3.3 绑定与验证结合使用

Gin 支持在数据绑定时自动执行验证，并返回错误信息。如果绑定或验证失败，Gin 会返回相应的错误响应。

**示例：绑定和验证结合**

```go
type User struct {
    Name  string `json:"name" binding:"required"`
    Age   int    `json:"age" binding:"gte=0,lte=130"`
}

func createUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        // 如果绑定或验证失败，返回错误信息
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, gin.H{"message": "User created", "user": user})
}
```

在这个例子中：

- `Name` 字段是必需的，不能为空（`required`）。
- `Age` 字段必须在 0 到 130 之间（`gte=0,lte=130`）。
- `c.ShouldBindJSON(&user)` 会先尝试绑定请求体中的 JSON 数据到 `user` 结构体。如果请求的数据不符合要求，或者字段不符合验证规则（如 `Age` 超过了最大值 130），会返回 400 错误，且错误信息会被包含在响应中。



### 3.4 自定义规则验证

Gin 还支持自定义验证规则。通过 `validator` 库，你可以创建自定义的验证函数，并将其应用到结构体字段上。

**示例：自定义验证规则**

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
)

func main() {
    // 获取 gin 默认的 validator 实例
    v := validator.New()

    // 注册自定义验证规则
    v.RegisterValidation("isadult", func(fl validator.FieldLevel) bool {
        return fl.Field().Int() >= 18
    })

    // 创建结构体
    type User struct {
        Name string `json:"name" binding:"required"`
        Age  int    `json:"age" binding:"isadult"`
    }

    // 创建 Gin 路由
    r := gin.Default()

    r.POST("/create", func(c *gin.Context) {
        var user User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, gin.H{"message": "User created", "user": user})
    })

    r.Run(":8080")
}
```

在这个例子中，`isadult` 是一个自定义的验证规则，它验证 `Age` 是否大于等于 18。如果 `Age` 小于 18，会返回错误。



### 3.5 自定义验证方法实践

本节参考：[04-数据绑定和验证.md](https://github.com/xinliangnote/Go/blob/master/01-Gin%E6%A1%86%E6%9E%B6/04-%E6%95%B0%E6%8D%AE%E7%BB%91%E5%AE%9A%E5%92%8C%E9%AA%8C%E8%AF%81.md)

首先，写一个验证方法： **validator/member/member.go**

```go
package member

import (
	"github.com/go-playground/validator/v10"
)

func NameValid(f1 validator.FieldLevel) bool {
	s := f1.Field().String()
	if s == "admin" {
		return false
	}
	return true
}
```

接下来，在路由中绑定：

**router/router.go**

```go
package router

import (
	"ginDemo/middleware/logger"
	"ginDemo/middleware/sign"
	"ginDemo/router/v1"
	"ginDemo/router/v2"
	"ginDemo/validator/member"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
)

func InitRouter(r *gin.Engine)  {

	r.Use(logger.LoggerToFile())

	// v1 版本
	GroupV1 := r.Group("/v1")
	{
		GroupV1.Any("/product/add", v1.AddProduct)
		GroupV1.Any("/member/add", v1.AddMember)
	}

	// v2 版本
	GroupV2 := r.Group("/v2").Use(sign.Sign())
	{
		GroupV2.Any("/product/add", v2.AddProduct)
		GroupV2.Any("/member/add", v2.AddMember)
	}

	// 绑定验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("NameValid", member.NameValid)
	}
}
```

最后，看一下调用的代码。

**router/v1/member.go**

```go
package v1

import (
	"ginDemo/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddMember(c *gin.Context) {

	res := entity.Result{}
	mem := entity.Member{}

	if err := c.ShouldBind(&mem); err != nil {
		res.SetCode(entity.CODE_ERROR)
		res.SetMessage(err.Error())
		c.JSON(http.StatusForbidden, res)
		c.Abort()
		return
	}

	// 处理业务(下次再分享)

	data := map[string]interface{}{
		"name" : mem.Name,
		"age"  : mem.Age,
	}
	res.SetCode(entity.CODE_ERROR)
	res.SetData(data)
	c.JSON(http.StatusOK, res)
}
```

访问看看效果吧。

访问：`http://localhost:8080/v1/member/add`

```
{
    "code": -1,
    "msg": "Key: 'Member.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'Member.Age' Error:Field validation for 'Age' failed on the 'required' tag",
    "data": null
}
```

访问：`http://localhost:8080/v1/member/add?name=1`

```
{"code":-1,"msg":"Key: 'Member.Age' Error:Field validation for 'Age' failed on the 'required' tag","data":null}
```

访问：`http://localhost:8080/v1/member/add?age=1`

```
{
    "code": -1,
    "msg": "Key: 'Member.Age' Error:Field validation for 'Age' failed on the 'required' tag",
    "data": null
}
```

访问：`http://localhost:8080/v1/member/add?name=admin&age=1`

```
{
    "code": -1,
    "msg": "Key: 'Member.Name' Error:Field validation for 'Name' failed on the 'NameValid' tag",
    "data": null
}
```

访问：`http://localhost:8080/v1/member/add?name=1&age=1`

```
{
    "code": -1,
    "msg": "Key: 'Member.Age' Error:Field validation for 'Age' failed on the 'gt' tag",
    "data": null
}
```

访问：`http://localhost:8080/v1/member/add?name=1&age=121`

```
{
    "code": -1,
    "msg": "Key: 'Member.Age' Error:Field validation for 'Age' failed on the 'lt' tag",
    "data": null
}
```

访问：`http://localhost:8080/v1/member/add?name=Tom&age=30`

```
{
    "code": 1,
    "msg": "",
    "data": {
        "age": 30,
        "name": "Tom"
    }
}
```

未避免返回信息过多，错误提示咱们也可以统一。

```
if err := c.ShouldBind(&mem); err != nil {
	res.SetCode(entity.CODE_ERROR)
	res.SetMessage("参数验证错误")
	c.JSON(http.StatusForbidden, res)
	c.Abort()
	return
}
```

这一次目录结构调整了一些，在这里说一下：

```
├─ ginDemo
│  ├─ common        //公共方法
│     ├── common.go
│  ├─ config        //配置文件
│     ├── config.go
│  ├─ entity        //实体
│     ├── ...
│  ├─ middleware    //中间件
│     ├── logger
│         ├── ...
│     ├── sign
│         ├── ...
│  ├─ router        //路由
│     ├── ...
│  ├─ validator     //验证器
│     ├── ...
│  ├─ vendor        //扩展包
│     ├── github.com
│         ├── ...
│     ├── golang.org
│         ├── ...
│     ├── gopkg.in
│         ├── ...
│  ├─ Gopkg.toml
│  ├─ Gopkg.lock
│  ├─ main.go
```

将 `sign` 和 `logger` 调整为中间件，并放到 `middleware` 中间件 目录。

新增了 `common` 公共方法目录。

新增了 `validator` 验证器目录。

新增了 `entity` 实体目录。