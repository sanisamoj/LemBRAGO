package utils

import (
	"fmt"
	"strings"
)

func GetTokenInH(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("Header is empty")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("Header is empty")
	}
	tokenString := parts[1]
	return tokenString, nil
}
