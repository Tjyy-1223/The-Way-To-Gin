package router

import (
	"GinDemo/middleware/logger"
	"GinDemo/middleware/recover"
	"GinDemo/middleware/sign"
	v1 "GinDemo/router/v1"
	v2 "GinDemo/router/v2"
	"GinDemo/validator/member"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitRouter(r *gin.Engine) {
	r.Use(logger.LoggerToFile(), recover.Recover())
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
