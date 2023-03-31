package middleware

import (
	"github.com/gin-gonic/gin"
	"hzer/internal/redis"
	"hzer/internal/response"
	jwtpkg "hzer/pkg/jwt"
)

//JWTAuthCheck 为对Jwt进行过期校验的中间件,此中间件需要在jwt中间件之后注册
//dbs 为mongodb所在数据集
//chackname 为key,一般为用户名
func JWTAuthCheck(chackname string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			token *jwtpkg.JWTLoad //校验的token结构体
			chv   interface{}     //check Member的值（inteface原型）
		)
		token, chv = jwtpkg.GetMember(c, chackname)
		if chv == nil {
			response.FailJson(c, response.TokenWrongful, false)
			c.Abort()
			return
		}
		/*jwts, err := mongodb.GetOneCollection(token.Issuer, bson.D{
			{chackname, chv},
		})
		if jwts == nil || err != nil {
			response.FailJson(c, response.TokenOverdue, false, "")
			c.Abort()
			return
		}*/
		ver, err := redis.GetTokenVersion(token.Issuer, chv.(string))
		if token.Version < ver || err != nil || ver == 0 {
			response.FailJson(c, response.TokenOverdue, false)
			c.Abort()
			return
		}
		c.Next()
	}
}
