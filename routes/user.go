package routes

import (
	"Project/db"
	"Project/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func LoginUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		HandleError(c, err, http.StatusBadRequest)
		return
	}
	// Поиск пользователя в базе данных
	var existingUser models.User
	if err := db.DB.Where("username = ?", user.Username).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}
	// Проверка пароля
	if err := existingUser.CheckPassword(user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	// Генерация JWT, передайте роль
	token, err := models.GenerateJWT(existingUser.Username, existingUser.Role) // Предполагая, что у вас есть поле Role в модели User
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
func LogoutUser(c *gin.Context) {
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

	// Добавьте токен в чёрный список или выполните другую логику для его аннулирования
	err := models.BlacklistToken(db.DB, oldToken)
	if err != nil {
		log.Printf("Ошибка при добавлении токена в чёрный список: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка логаута"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Успешный логаут"})
}
