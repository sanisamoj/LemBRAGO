package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"lembrago.com/lembrago/internal/config"
	"lembrago.com/lembrago/models"
)

type CustomClaims struct {
	UserID string           `json:"id"`
	OrgID  string           `json:"orgId"`
	Role   *models.UserRole `json:"role"`
	Email  *string          `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID, orgID string, role models.UserRole) (string, error) {
	appConfig := config.GetAppConfig()
	expirationTime := time.Now().Add(29 * time.Hour)
	claims := &CustomClaims{
		UserID: userID,
		OrgID:  orgID,
		Role:   &role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    appConfig.JWTIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(appConfig.JWTSecret)
	if err != nil {
		return "", fmt.Errorf("Error ao assinar o Token")
	}

	return tokenStr, nil
}

func GenerateJWTUserCreation(userID, orgID string, role models.UserRole, email string) (string, error) {
	appConfig := config.GetAppConfig()
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &CustomClaims{
		UserID: userID,
		OrgID:  orgID,
		Role:   &role,
		Email:  &email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    appConfig.JWTIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(appConfig.JWTSecretUserCreation)
	if err != nil {
		return "", fmt.Errorf("Error ao assinar o Token")
	}

	return tokenStr, nil
}
