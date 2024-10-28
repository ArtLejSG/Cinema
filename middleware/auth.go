package middleware

import (
	"Project/db"
	"Project/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Middleware для проверки токена и аутентификации пользователя
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен не предоставлен"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &models.Claims{}
		_, err := models.ParseJWT(token, claims) // Реализуйте эту функцию
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
			c.Abort()
			return
		}

		blacklisted, err := models.IsTokenBlacklisted(db.DB, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки токена"})
			c.Abort()
			return
		}

		if blacklisted {
			c.JSON(http.StatusForbidden, gin.H{"error": "Токен отозван"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// Middleware для проверки роли пользователя
func RoleCheck(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Middleware для обновления токена
func RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен не предоставлен"})
		return
	}

	tokenParts := strings.Split(authHeader, "Bearer ")
	if len(tokenParts) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат токена"})
		return
	}

	oldToken := tokenParts[1]

	newToken, err := models.RefreshJWT(oldToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка обновления токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken})
}
