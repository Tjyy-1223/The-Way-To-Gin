package service

import (
	"errors"
	"my-gin/app/common/request"
	"my-gin/app/models"
	"my-gin/global"
	"my-gin/utils"
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
