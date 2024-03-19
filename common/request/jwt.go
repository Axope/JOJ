package request

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type CustomClaims struct {
	UUID  string
	ID    uint
	Admin int
	jwt.RegisteredClaims
}

func NewCustomClaims(uuid string, id uint, admin int, expire time.Duration) *CustomClaims {
	return &CustomClaims{
		UUID:  uuid,
		ID:    id,
		Admin: admin,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
		},
	}
}
