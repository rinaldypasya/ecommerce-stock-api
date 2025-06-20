package auth

import (
	"ecommerce-stock-api/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTSecret() ([]byte, error) {
	secret := config.AppConfig.JWTSecret
	if secret == "" {
		return nil, errors.New("JWT_SECRET environment variable not set")
	}
	return []byte(secret), nil
}

func GenerateJWT(userID uint) (string, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
