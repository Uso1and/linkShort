package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"lnkshrt/internal/domain/config"
	"lnkshrt/internal/domain/models"
)

func AuthMiddleware(cfg *config.ConfigDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Checking authorization...")

		authHeader := c.GetHeader("Authorization")
		log.Printf("Authorization header: %s", authHeader)

		if authHeader == "" {
			log.Println("Authorization header is missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("Token string: %s", tokenString)

		if tokenString == authHeader {
			log.Println("Bearer token not found")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token not found"})
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil {
			log.Printf("Token parsing error: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			log.Printf("Authenticated user: %d (%s)", claims.UserID, claims.Username)
			c.Set("userID", claims.UserID)
			c.Set("username", claims.Username)
		} else {
			log.Println("Invalid token claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		c.Next()
	}
}
