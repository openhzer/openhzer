package jwt_test

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"hzer/internal/middleware"
	"hzer/internal/response"
	jwtpkg "hzer/pkg/jwt"
)

// 业务模型
type testUseModel struct {
	ID int `json:"ID"`
}

// 定义一个负载模型
type testPayloadModel struct {
	Data testUseModel
	jwt.StandardClaims
}

func testHandler(c *gin.Context) {
	// 获取jwt数据并处理错误
	// 可以在中间件完成这一步,根据业务自行扩展
	token, err := jwtpkg.GetJwtProto(c, &testPayloadModel{})
	if err != nil {
		response.FailJson(c, response.FailStruct{
			Code: 401,
			Msg:  err.Error(),
		}, false)
	}
	fmt.Println(token)
}

func main() {
	r := gin.Default()

	//使用中间件时绑定业务jwt模型,并且可以设置验证的签发人
	r.Use(middleware.JWTAuth(&testPayloadModel{}, "is"))
	r.GET("/test", testHandler)
	r.Run()
}
