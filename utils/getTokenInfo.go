package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func GetTokenInfo(jwtSecret []byte, tokenStr string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	}

	_, err := jwt.ParseWithClaims(tokenStr, claims, keyFunc)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
