package service

import (
	"github.com/dgrijalva/jwt-go"
	"my-gin/global"
	"time"
)

type jwtService struct {
}

var JwtService = new(jwtService)

// JwtUser 所有需要颁发 token 的用户模型必须实现这个接口
type JwtUser interface {
	GetUid() string
}

// CustomClaims 自定义 Claims
type CustomClaims struct {
	jwt.StandardClaims
}

const (
	TokenType    = "bearer"
	AppGuardName = "app"
)

type TokenOutPut struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// CreateToken 生成 Token
func (jwtService *jwtService) CreateToken(GuardName string, user JwtUser) (tokenData TokenOutPut, err error, token *jwt.Token) {
	token = jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		CustomClaims{ // 自定义声明
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Unix() + global.App.Config.Jwt.JwtTtl, // 设置 token 的过期时间，使用 Unix 时间戳 + 配置的有效期（JwtTtl）
				Id:        user.GetUid(),                                    // 设置 token 的唯一 ID，通常是用户 ID
				Issuer:    GuardName,                                        // 设置 token 的颁发者，用于区分不同客户端的 token
				NotBefore: time.Now().Unix() - 1000,                         // 设置 token 的生效时间，这里设置为当前时间减去 1000 秒（可能是为了提前几秒使用）
			},
		},
	)

	tokenStr, err := token.SignedString([]byte(global.App.Config.Jwt.Secret))
	tokenData = TokenOutPut{
		tokenStr,
		int(global.App.Config.Jwt.JwtTtl),
		TokenType,
	}
	return
}
