package routes

import (
	"Project/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Получение всех пользователей
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := models.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении пользователей"})
		return
	}
	c.JSON(http.StatusOK, users)
}
