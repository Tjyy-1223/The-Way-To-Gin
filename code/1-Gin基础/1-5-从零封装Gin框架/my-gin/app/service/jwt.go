package service

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"my-gin/global"
	"my-gin/utils"
	"strconv"
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

// 获取黑名单缓存 key
func (jwtService *jwtService) getBlackListKey(tokenStr string) string {
	return "jwt_black_list:" + utils.MD5([]byte(tokenStr))
}

// JoinBlackList token 加入黑名单
func (jwtService *jwtService) JoinBlackList(token *jwt.Token) (err error) {
	nowUnix := time.Now().Unix()
	timer := time.Duration(token.Claims.(*CustomClaims).ExpiresAt-nowUnix) * time.Second
	// 将 token 剩余时间设置为缓存有效期，并将当前时间作为缓存 value 值
	err = global.App.Redis.SetNX(context.Background(), jwtService.getBlackListKey(token.Raw), strconv.FormatInt(nowUnix, 10), timer).Err()
	return
}

// IsInBlacklist token 是否在黑名单中
func (jwtService *jwtService) IsInBlacklist(tokenStr string) bool {
	blackedUnixStr, err := global.App.Redis.Get(context.Background(), jwtService.getBlackListKey(tokenStr)).Result()
	if err != nil {
		return false
	}

	blackedUnix, err := strconv.ParseInt(blackedUnixStr, 10, 64)
	if blackedUnixStr == "" || err != nil {
		return false
	}

	// JwtBlacklistGracePeriod 为黑名单宽限时间，避免并发请求失效
	if time.Now().Unix()-blackedUnix < global.App.Config.Jwt.JwtBlacklistGracePeriod {
		return false
	}
	return true
}

// GetUserInfo 根据不同客户端 token ，查询不同用户表数据
func (jwtService *jwtService) GetUserInfo(GuardName string, id string) (err error, user JwtUser) {
	switch GuardName {
	case AppGuardName:
		return UserService.GetUserInfo(id)
	default:
		err = errors.New("guard " + GuardName + " dose not exist")
	}
	return
}
