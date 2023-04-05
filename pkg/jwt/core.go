package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"hzer/pkg/util"
	"os"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

const ExpiresAtFiveDays = 3600 * 24 * 5

// var SecretKey = "6oFChDc2tQhXlM8XSTqAnqJik4kkpU2v"
var SecretKey = initKey()

func initKey() string {
	sk := os.Getenv("SecretKey")
	return util.Ifs(sk == "", "6oFChDc2tQhXlM8XSTqAnqJik4kkpU2v", sk)
}

type LoadModel interface {
	Valid() error
}

// JWT jwt对象
type JWT struct {
	// 声明签名信息
	SigningKey []byte
}

// NewJWT 初始化jwt对象
func NewJWT() *JWT {
	return &JWT{
		[]byte(SecretKey),
	}
}

// CreateToken 调用jwt-go库生成token,编码的算法为jwt.SigningMethodHS256
func (j *JWT) CreateToken(payload jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString(j.SigningKey)
}

// ParserToken 将token解码并验证
func (j *JWT) ParserToken(tokenString string, model LoadModel) (jwt.Claims, error) {
	// https://gowalker.org/github.com/dgrijalva/jwt-go#ParseWithClaims
	// 输入用户自定义的Claims结构体对象,token,以及自定义函数来解析token字符串为jwt的Token结构体指针
	// Keyfunc是匿名函数类型: type Keyfunc func(*Token) (interface{}, error)
	// func ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {}

	token, err := jwt.ParseWithClaims(tokenString, model, func(token *jwt.Token) (interface{}, error) {
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
	if token.Valid {
		return token.Claims, nil
	}
	/*// 将token中的claims信息解析出来并断言成用户自定义的有效载荷结构
	if claims, ok := token.Claims.(*JWTLoad); ok && token.Valid {
		return claims, nil
	}*/
	return nil, fmt.Errorf("token无效")
}

// CreateToken 用于快速生成一个Token
// UserLoad 用户负载
// ExpiresAt 多少秒之后过期，单位：秒
// Issuer 签名颁发着
func CreateToken(JWTLoad jwt.Claims) (Token string, err error) {
	// 构造SignKey: 签名和解签名需要使用一个值
	j := NewJWT()
	// 构造用户claims信息(负荷)
	token, err := j.CreateToken(JWTLoad)
	if err != nil {
		return "", err
	}
	return token, nil
}

// CreateStandardClaims 快速创建签名数据
func CreateStandardClaims(ExpiresAt int64, Issuer string) jwt.StandardClaims {
	return jwt.StandardClaims{
		NotBefore: time.Now().Unix() - 1000,      // 签名生效时间
		ExpiresAt: time.Now().Unix() + ExpiresAt, // 签名过期时间
		Issuer:    Issuer,                        // 签名颁发者
	}
}

// CheckToken 用于检查Token是否有效
// issuer为可选参数
func CheckToken(token string, model LoadModel, issuer string) (bool, interface{}, *jwt.StandardClaims, error) {
	j := NewJWT()
	load, err := j.ParserToken(token, model)
	if err != nil {
		return false, nil, nil, err
	}
	retv := reflect.ValueOf(model)
	claimser := retv.Elem().FieldByName("StandardClaims")
	if !claimser.CanSet() {
		return false, nil, nil, errors.New("JWTLoad Can Not Set")
	}
	claims := claimser.Interface().(jwt.StandardClaims)
	if claims.ExpiresAt < time.Now().Unix() {
		return false, load, &claims, errors.New("token已过期，请重新登录")
	}
	if issuer == "" {
		return true, load, &claims, nil
	}
	if issuer == claims.Issuer {
		return true, load, &claims, nil
	}
	return false, load, &claims, errors.New("token签名错误，请重新登录")
}

// GetJwtProto 从Gin中获取Jwt原型体
// model为业务内jwt模型
func GetJwtProto[T any](c *gin.Context, model T) (T, error) {
	tokener, _ := c.Get("token")
	if tokener == nil {
		return model, errors.New("JWT Is NULL")
	}
	token, ok := tokener.(T)
	if !ok {
		return model, errors.New("JWTLoad Error")
	}
	return token, nil
}

// GetMember 用于从GinContext中的JWTUserStruct中取回数据，返回interface{}，数据可能为nil，需要自己断言转换
/**
 * $param {gin.context}
 */
/*func GetMember(c *gin.Context, Member string) (*JWTLoad, interface{}) {
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
}*/

// GetTokenLoad 用于从GinContext中取回JWTUsermapm
func GetTokenLoad(c *gin.Context) (*JWTLoad, map[string]interface{}) {
	token, _ := c.Get("token")
	if token == nil {
		return nil, nil
	}
	load := token.(*JWTLoad)
	if load.UserLoad == nil {
		return nil, nil
	}
	return load, load.UserLoad.(map[string]interface{})
}

func ShouldBindTokenLoad(c *gin.Context, obj any) error {
	var e error
	token, _ := c.Get("token")
	if token == nil {
		e = errors.New("token is invalid")
	}
	load := token.(*JWTLoad)
	if load.UserLoad == nil {
		e = errors.New("token is illegal")
	}
	for k, v := range load.UserLoad.(map[string]interface{}) {
		fmt.Println(v, k)
		e = SetField(obj, k, v)
		if e != nil {
			return e
		}
	}
	return nil
}
func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)
	fmt.Println(structFieldValue)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}
