package token

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	key         string
	identifyKey string
}

var (
	MissingHeaderError = errors.New("`Authorization` 头部为空")
	config             = Config{"Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgqBBCWQ", "identifySecureKey"}
	once               sync.Once
)

func Init(key, identityKey string) {
	once.Do(func() {
		if key != "" {
			config.key = key
		}
		if identityKey != "" {
			config.identifyKey = identityKey
		}
	})
}

// 根据指定的密钥解析token，成功返回token上下文，否则报错
func ParseToken(tokenString, key string) (string, error) {
	// 1. 解析token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// 2. 确保加密算法与预期的一致
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(key), nil
	})

	// 解析token失败
	if err != nil {
		return "", err
	}

	var identifyKey string

	// 如果解析成功则将token中的主题提取出来
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		identifyKey = claims[config.identifyKey].(string)
	}

	return identifyKey, nil
}

// 解析http请求
func ParseRequest(c *gin.Context) (token string, err error) {
	header := c.Request.Header.Get("Authorization")

	if len(header) == 0 {
		return "", MissingHeaderError
	}

	var t string

	fmt.Sscanf(header, "Bearer %s", &t)

	return ParseToken(t, config.key)
}

// 签发token
func SignToken(identifyKey string) (tokenString string, err error) {
	// 1. 包装token内容
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		config.identifyKey: identifyKey,
		"nbf":              time.Now().Unix(),
		"iat":              time.Now().Unix(),
		"exp":              time.Now().Add(100000 * time.Hour).Unix(),
	})

	tokenString, err = token.SignedString([]byte(config.key))

	return
}
