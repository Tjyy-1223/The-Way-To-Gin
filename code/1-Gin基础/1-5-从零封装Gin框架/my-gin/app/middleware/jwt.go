package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"my-gin/app/common/response"
	"my-gin/app/service"
	"my-gin/global"
	"strconv"
	"time"
)

func JWTAuth(GuardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")
		if tokenStr == "" {
			response.TokenFail(c)
			c.Abort()
			return
		}

		// Token 解析校验
		token, err := jwt.ParseWithClaims(tokenStr, &service.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(global.App.Config.Jwt.Secret), nil
		})

		if err != nil || service.JwtService.IsInBlacklist(tokenStr) {
			global.App.Log.Error("can't pass the jwt token validator ")
			response.TokenFail(c)
			c.Abort()
			return
		}

		// 转换成自定义 service.CustomClaims
		claims := token.Claims.(*service.CustomClaims)
		// Token 发布者校验
		if claims.Issuer != GuardName {
			response.TokenFail(c)
			c.Abort()
			return
		}

		// token 续签处理
		if claims.ExpiresAt-time.Now().Unix() < global.App.Config.Jwt.RefreshGracePeriod {
			lock := global.Lock("refresh_token_local_"+claims.Id, global.App.Config.Jwt.JwtBlacklistGracePeriod)
			if lock.Get() {
				err, user := service.JwtService.GetUserInfo(GuardName, claims.Id)
				if err != nil {
					global.App.Log.Error("service.JwtService.GetUserInfo error!")
					lock.Release()
				} else {
					tokenData, _, _ := service.JwtService.CreateToken(GuardName, user)
					c.Header("new-token", tokenData.AccessToken)
					c.Header("new-expires-in", strconv.Itoa(tokenData.ExpiresIn))
					_ = service.JwtService.JoinBlackList(token)
					lock.Release()
				}
			}
		}

		c.Set("token", token)
		c.Set("id", claims.Id)
	}
}
