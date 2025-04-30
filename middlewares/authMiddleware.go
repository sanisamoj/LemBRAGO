package middlewares

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"lembrago.com/lembrago/cache"
	"lembrago.com/lembrago/models"
	"lembrago.com/lembrago/utils"
)

func AuthMiddleware(jwtSecretParam []byte, requiredRoles []models.UserRole) gin.HandlerFunc {
	if len(jwtSecretParam) == 0 {
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Configuração de segurança interna inválida"})
		}
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}
		tokenString := parts[1]
		tExist, err := cache.Get(tokenString)

		if tExist != "" || err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
			return
		}

		claims := &utils.CustomClaims{}

		keyFunc := func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecretParam, nil
		}

		token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)

		if err != nil {
			log.Printf("Token parsing/validation error: %v", err)
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			} else if errors.Is(err, jwt.ErrSignatureInvalid) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token signature is invalid"})
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("orgID", claims.OrgID)

		var currentUserRole models.UserRole
		if claims.Role != nil && *claims.Role != "" {
			c.Set("role", *claims.Role)
			currentUserRole = *claims.Role
		}

		if claims.Email != nil && *claims.Email != "" {
			c.Set("email", *claims.Email)
		}

		if len(requiredRoles) > 0 {
			roleFound := slices.Contains(requiredRoles, models.UserRole(currentUserRole))

			if !roleFound {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied: insufficient role"})
				return
			}
		}

		c.Next()
	}
}
