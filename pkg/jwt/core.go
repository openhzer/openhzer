package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

var SecretKey = "6oFChDc2tQhXlM8XSTqAnqJik4kkpU2v"

//JWT jwt对象
type JWT struct {
	// 声明签名信息
	SigningKey []byte
}

//NewJWT 初始化jwt对象
func NewJWT() *JWT {
	return &JWT{
		[]byte(SecretKey),
	}
}

//CreateToken 调用jwt-go库生成token,编码的算法为jwt.SigningMethodHS256
func (j *JWT) CreateToken(payload JWTLoad) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString(j.SigningKey)
}

//ParserToken 将token解码并验证
func (j *JWT) ParserToken(tokenString string) (*JWTLoad, error) {
	// https://gowalker.org/github.com/dgrijalva/jwt-go#ParseWithClaims
	// 输入用户自定义的Claims结构体对象,token,以及自定义函数来解析token字符串为jwt的Token结构体指针
	// Keyfunc是匿名函数类型: type Keyfunc func(*Token) (interface{}, error)
	// func ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {}
	token, err := jwt.ParseWithClaims(tokenString, &JWTLoad{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		// https://gowalker.org/github.com/dgrijalva/jwt-go#ValidationError
		// jwt.ValidationError 是一个无效token的错误结构
		if ve, ok := err.(*jwt.ValidationError); ok {
			// ValidationErrorMalformed是一个uint常量，表示token不可用
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, fmt.Errorf("token不可用")
				// ValidationErrorExpired表示Token过期
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, fmt.Errorf("token过期")
				// ValidationErrorNotValidYet表示无效token
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, fmt.Errorf("无效的token")
			} else {
				return nil, fmt.Errorf("token不可用")
			}

		}
	}
	// 将token中的claims信息解析出来并断言成用户自定义的有效载荷结构
	if claims, ok := token.Claims.(*JWTLoad); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("token无效")
}

//CreateToken 用于快速生成一个Token
//UserLoad 用户负载
//ExpiresAt 多少秒之后过期，单位：秒
//Issuer 签名颁发着
func CreateToken(UserLoad interface{}, ExpiresAt, Version int64, Issuer string) (Token string, err error) {
	// 构造SignKey: 签名和解签名需要使用一个值
	j := NewJWT()

	// 构造用户claims信息(负荷)
	claims := JWTLoad{
		UserLoad: UserLoad,
		Version:  Version,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,      // 签名生效时间
			ExpiresAt: time.Now().Unix() + ExpiresAt, // 签名过期时间
			Issuer:    Issuer,                        // 签名颁发者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

//CheckToken 用于检查Token是否有效
//issuer为可选参数
func CheckToken(token string, issuer string) (bool, *JWTLoad, error) {
	j := NewJWT()
	load, err := j.ParserToken(token)
	if err != nil {
		return false, load, err
	}
	if load.ExpiresAt < time.Now().Unix() {
		return false, load, errors.New("token已过期，请重新登录")
	}
	if issuer == "" {
		return true, load, nil
	}
	if issuer == load.Issuer {
		return true, load, nil
	}
	return false, load, errors.New("token签名错误，请重新登录")
}

//GetMember 用于从GinContext中的JWTUserStruct中取回数据，返回interface{}，数据可能为nil，需要自己断言转换
func GetMember(c *gin.Context, Member string) (*JWTLoad, interface{}) {
	token, _ := c.Get("token")
	if token == nil {
		return nil, nil
	}
	load := token.(*JWTLoad)
	if load.UserLoad == nil {
		return nil, nil
	}
	ujwt := load.UserLoad.(map[string]interface{})
	v, ok := ujwt[Member]
	if !ok {
		return nil, nil
	}
	return load, v
}
