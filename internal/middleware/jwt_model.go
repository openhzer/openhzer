package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	jwtPkg "hzer/pkg/jwt"
)

type AdminLoad struct {
	Username string `json:"username"`
}

type WechatLoad struct {
	UUID string `json:"uuid"`
}

type UserLoad struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

func NewUserLoad(ID uint, ExpiresAt int64, Issuer string) *UserLoad {
	return &UserLoad{
		ID:               ID,
		RegisteredClaims: jwtPkg.CreateStandardClaims(ExpiresAt, Issuer),
	}
}
