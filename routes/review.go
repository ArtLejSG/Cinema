package routes

import (
	"Project/db"
	"Project/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Добавление отзыва
func AddReview(c *gin.Context) {
	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
		return
	}
	review.UserID = userID.(uint)

	// Проверка существования фильма
	if err := db.DB.First(&models.Movie{}, review.MovieID).Error; err != nil {
		log.Println("Фильм не найден:", err) // Логирование ошибки
		c.JSON(http.StatusNotFound, gin.H{"error": "Фильм не найден"})
		return
	}

	// Валидация отзыва
	if err := review.Validate(); err != nil {
		log.Println("Ошибка валидации отзыва:", err) // Логирование ошибки валидации
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Сохранение отзыва в базу данных
	if err := db.DB.Create(&review).Error; err != nil {
		log.Println("Ошибка при добавлении отзыва:", err) // Логирование ошибки
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления отзыва"})
		return
	}

	c.JSON(http.StatusCreated, review) // Возвращаем созданный отзыв
}

// Получение отзывов для фильма
func GetReviewsByMovie(c *gin.Context) {
	var reviews []models.Review
	movieID := c.Param("id") // Изменяем на movieID для соответствия

	if err := db.DB.Where("movie_id = ?", movieID).Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении отзывов"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// Обновление отзыва
func UpdateReview(c *gin.Context) {
	var review models.Review
	id := c.Param("id")

	if err := db.DB.First(&review, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Отзыв не найден"})
		return
	}

	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Валидация перед обновлением
	if err := review.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Save(&review)
	c.JSON(http.StatusOK, review)
}

// Удаление отзыва
func DeleteReview(c *gin.Context) {
	id := c.Param("id")
	var review models.Review

	if err := db.DB.First(&review, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Отзыв не найден"})
		return
	}

	if err := db.DB.Delete(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении отзыва"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
