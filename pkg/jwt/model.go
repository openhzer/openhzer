package jwt

import "github.com/dgrijalva/jwt-go"

type JWTLoad struct {
	UserLoad interface{} `json:"user_load"`
	Version  int64
	jwt.StandardClaims
}
