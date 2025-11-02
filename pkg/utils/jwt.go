package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type Claims struct {
	UserID   string `json:"sub"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func SignToken(userId, username string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpires := os.Getenv("JWT_EXPIRES_IN")

	claims := jwt.MapClaims{
		"sub":  userId,
		"name": username,
		"iat":  time.Now().Unix(),
	}
	// срок жизни
	if jwtExpires != "" {
		duration, err := time.ParseDuration(jwtExpires)
		if err != nil {
			return "", fmt.Errorf("invalid JWT_EXPIRES_IN: %v", err)
		}
		claims["exp"] = time.Now().Add(duration).Unix()
	} else {
		claims["exp"] = time.Now().Add(15 * time.Minute).Unix()
	}

	// создаём токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// подписываем
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}
	return signedToken, nil
}

func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32) // 32 байта = 256 бит
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
