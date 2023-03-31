package middleware

import (
	"github.com/gin-gonic/gin"
	"hzer/internal/response"
	jwtpkg "hzer/pkg/jwt"
)

//JWTAuth 为JWT中间件，客户端下需要在header带上Authorization: Bearer <token>
//issure 为可选验证签名，支持多参选择
func JWTAuth(issure ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			response.FailJson(c, response.NoAccess, false)
			c.Abort()
			return
		}
		if len(token) < 8 {
			response.FailJson(c, response.TokenWrongful, false)
			c.Abort()
			return
		}
		token = token[7:]
		chk, jwt, err := jwtpkg.CheckToken(token, "")
		if err != nil {
			errmsg := response.TokenWrongful
			errmsg.Msg = "令牌验证错误"
			response.FailJson(c, errmsg, true, err.Error())
			c.Abort()
			return
		}
		chs := true
		for _, v := range issure {
			if v != jwt.Issuer {
				chs = false
				break
			}
		}
		if !chs {
			errmsg := response.TokenWrongful
			errmsg.Msg = "签名错误"
			response.FailJson(c, errmsg, true, err.Error())
			c.Abort()
			return
		}
		if !chk {
			response.FailJson(c, response.NoAccess, true, "破损令牌")
			c.Abort()
			return
		}
		c.Set("token", jwt)
		c.Next()
	}
}
