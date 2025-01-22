package services

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/EgSundqvist/imagebook-imageapi/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware för att hantera autentisering med JWT
func AuthMiddleware() gin.HandlerFunc {
	secretKey := config.AppConfig.JWT.FrontendSecretKey
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Kontrollera att token är signerad med rätt metod och nyckel
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			log.Printf("Token parsing error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Spara användarens information i context för senare användning
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token claims"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

// BackendAuthMiddleware för att hantera autentisering med backend JWT
func BackendAuthMiddleware() gin.HandlerFunc {
	backendSecretKey := config.AppConfig.JWT.BackendSecretKey
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Kontrollera att token är signerad med rätt metod och nyckel
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(backendSecretKey), nil
		})

		if err != nil || !token.Valid {
			log.Printf("Token parsing error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Spara användarens information i context för senare användning
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userID, ok := claims["identity"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token claims"})
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	}
}
