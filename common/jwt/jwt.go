package jwt

import (
	"errors"
	"github.com/Axope/JOJ/common/request"
	"github.com/Axope/JOJ/configs"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/sync/singleflight"
	"net/http"
	"time"
)

var (
	control             = &singleflight.Group{}
	ErrTokenExpired     = errors.New("token is expired")
	ErrTokenNotValidYet = errors.New("token not active yet")
	ErrTokenMalformed   = errors.New("that's not a token")
	ErrTokenInvalid     = errors.New("couldn't handle this token")
)

// header 中需要携带的默认字段
const TokenField = "Authorization"

type JWT struct {
	SigningKey []byte
	Expire     time.Duration
}

func (j *JWT) CreateToken(claims request.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 续签 使用归并回源避免并发问题
func (j *JWT) CreateTokenByOldToken(oldToken string, claims request.CustomClaims) (string, error) {
	v, err, _ := control.Do("JWT"+oldToken, func() (interface{}, error) {
		return j.CreateToken(claims)
	})
	return v.(string), err
}

func (j *JWT) ParseToken(tokenString string) (*request.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, ErrTokenInvalid

	} else {
		return nil, ErrTokenInvalid
	}
}

func (j *JWT) GetToken(c *gin.Context) (string, error) {
	// 检查 cookie 和 header 中的字段
	token, err := c.Cookie(TokenField)
	if err != nil {
		if err == http.ErrNoCookie {
			return c.Request.Header.Get(TokenField), nil
		}
		return "", err
	}
	return token, nil
}

func (j *JWT) GetClaims(c *gin.Context) (*request.CustomClaims, error) {
	token, err := j.GetToken(c)
	if err != nil {
		return nil, err
	}
	claims, err := j.ParseToken(token)
	if err != nil {
		return nil, err
	}
	return claims, err
}

func (j *JWT) GetUserID(c *gin.Context) (uint, error) {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := j.GetClaims(c); err != nil {
			return 0, err
		} else {
			return cl.ID, nil
		}
	} else {
		c, ok := claims.(*request.CustomClaims)
		if !ok {
			return 0, errors.New("conversion failed")
		}
		return c.ID, nil
	}
}

var j *JWT

func InitJWT() error {
	JWTConfig := configs.GetJWTConfig()
	expire, err := time.ParseDuration(JWTConfig.Expire)
	if err != nil {
		// panic(fmt.Sprintf("expire(%v) parse error, Valid time units are 'ns', 'us' (or 'µs'), 'ms', 's', 'm', 'h'.", JWTConfig.Expire))
		return err
	}
	j = &JWT{
		SigningKey: []byte(JWTConfig.SigningKey),
		Expire:     expire,
	}
	return nil
}
func GetJWT() *JWT {
	return j
}
