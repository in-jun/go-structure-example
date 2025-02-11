package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/in-jun/go-structure-example/internal/pkg/config"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

var (
	secretKey        = []byte(config.AppConfig.JWTSecret)
	accessExpiry, _  = time.ParseDuration(config.AppConfig.JWTAccessExpiry)
	refreshExpiry, _ = time.ParseDuration(config.AppConfig.JWTRefreshExpiry)
)

func GenerateAccessToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
