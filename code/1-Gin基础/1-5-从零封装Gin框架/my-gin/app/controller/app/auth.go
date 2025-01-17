package app

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"my-gin/app/common/request"
	"my-gin/app/common/response"
	"my-gin/app/service"
)

// Login 用户登陆
func Login(c *gin.Context) {
	var form request.Login
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if err, user := service.UserService.Login(form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		tokenData, err, _ := service.JwtService.CreateToken(service.AppGuardName, user)
		if err != nil {
			response.BusinessFail(c, err.Error())
		}
		response.Success(c, tokenData)
	}
}

func Info(c *gin.Context) {
	err, user := service.UserService.GetUserInfo(c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, user)
}

func LogOut(c *gin.Context) {
	err := service.JwtService.JoinBlackList(c.Keys["token"].(*jwt.Token))
	if err != nil {
		response.BusinessFail(c, "登出失败")
		return
	}
	response.Success(c, nil)
}
