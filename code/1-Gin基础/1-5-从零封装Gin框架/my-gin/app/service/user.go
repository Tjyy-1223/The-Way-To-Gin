package service

import (
	"errors"
	"my-gin/app/common/request"
	"my-gin/app/models"
	"my-gin/global"
	"my-gin/utils"
	"strconv"
)

type userService struct {
}

var UserService = new(userService)

// Register 编写用户注册逻辑
func (userService *userService) Register(params request.Register) (err error, user models.User) {
	var result = global.App.DB.Where("mobile = ?", params.Mobile).Select("id").First(&models.User{})
	if result.RowsAffected != 0 {
		err = errors.New("手机号已经存在")
		return
	}

	user = models.User{Name: params.Name, Mobile: params.Mobile, Password: utils.BcryptMake([]byte(params.Password))}
	err = global.App.DB.Create(&user).Error
	return
}

// Login 编写用户登陆逻辑
func (userService *userService) Login(param request.Login) (err error, user *models.User) {
	err = global.App.DB.Where("mobile = ?", param.Mobile).First(&user).Error
	if err != nil || !utils.BcryptMakeCheck([]byte(param.Password), user.Password) {
		err = errors.New("用户名不存在或者密码错误")
	}
	return
}

// GetUserInfo 获取用户信息
func (userService *userService) GetUserInfo(id string) (err error, user models.User) {
	intId, err := strconv.Atoi(id)
	err = global.App.DB.First(&user, intId).Error
	if err != nil {
		err = errors.New("数据不存在")
	}
	return
}
