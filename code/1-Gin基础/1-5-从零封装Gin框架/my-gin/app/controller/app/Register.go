package app

import (
	"github.com/gin-gonic/gin"
	"my-gin/app/common/request"
	"my-gin/app/common/response"
	"my-gin/app/service"
)

// Register 用户注册
func Register(c *gin.Context) {
	var form request.Register
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if err, user := service.UserService.Register(form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, user)
	}
}
