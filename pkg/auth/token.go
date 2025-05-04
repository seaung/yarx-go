package auth

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	identityKey string
	key         string
}

var (
	ErrMissingHeader = errors.New("Authorization头部长度为零")
	config           = Config{"key", "identityKey"}
	once             sync.Once
)

func InitIdentityKey(key, identity string) {
	once.Do(func() {
		if key != "" {
			config.key = key
		}

		if identity != "" {
			config.identityKey = identity
		}
	})
}

func ParseToken(tstring, key string) (identity string, err error) {
	token, err := jwt.Parse(tstring, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(key), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		identity = claims[config.identityKey].(string)
	}

	return
}

func ParseRequest(c *gin.Context) (string, error) {
	header := c.Request.Header.Get("Authorization")
	if len(header) == 0 {
		return "", ErrMissingHeader
	}

	var t string

	fmt.Sscanf(header, "Bearer %s", &t)

	return ParseToken(t, config.key)
}

func GenerateToken(identityKey string) (token string, err error) {
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		config.identityKey: identityKey,
		"nbf":              time.Now().Unix(),
		"iat":              time.Now().Unix(),
		"exp":              time.Now().Unix(),
	})

	token, err = tokenString.SignedString([]byte(config.key))

	return
}
