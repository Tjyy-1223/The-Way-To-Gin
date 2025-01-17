package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"my-gin/app/common/response"
	"my-gin/app/service"
	"my-gin/global"
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

		c.Set("token", token)
		c.Set("id", claims.Id)
	}
}
