package auth

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWT(tokenString string) (*Claims, error) {
	secretKey := os.Getenv("SECRET")

	keyFunc := func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, err
}
