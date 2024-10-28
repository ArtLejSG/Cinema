package routes

import (
	"Project/db"
	"Project/models"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Добавление нового фильма
func AddMovie(c *gin.Context) {
	var newMovie models.Movie
	if err := c.ShouldBindJSON(&newMovie); err != nil {
		HandleError(c, err, http.StatusBadRequest)
		return
	}
	//
	if newMovie.Title == "" || newMovie.Year <= 0 {
		HandleError(c, errors.New("необходимо указать корректное название и год"), http.StatusBadRequest)
		return
	}
	// Проверка роли пользователя
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещён"})
		return
	}
	db.DB.Create(&newMovie)
	c.JSON(http.StatusOK, newMovie)
}

// Получение всех фильмов
func GetMovies(c *gin.Context) {
	var movies []models.Movie
	if err := db.DB.Find(&movies).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, movies)
}

// Получение фильма по ID
func GetMovieByID(c *gin.Context) {
	var movie models.Movie
	if err := db.DB.First(&movie, c.Param("id")).Error; err != nil {
		HandleError(c, err, http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, movie)
}

// Обновление фильма
func UpdateMovie(c *gin.Context) {
	var movie models.Movie
	if err := db.DB.First(&movie, c.Param("id")).Error; err != nil {
		HandleError(c, err, http.StatusBadRequest)
		return
	}
	if err := c.ShouldBindJSON(&movie); err != nil {
		HandleError(c, err, http.StatusBadRequest)
		return
	}
	// Проверка роли пользователя
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещён"})
		return
	}
	db.DB.Save(&movie)
	c.JSON(http.StatusOK, movie)
}

// Удаление фильма
func DeleteMovie(c *gin.Context) {
	var movie models.Movie
	if err := db.DB.First(&movie, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Фильм не найден"})
		return
	}
	// Проверка роли пользователя
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещён"})
		return
	}
	db.DB.Delete(&movie)
	c.JSON(http.StatusOK, gin.H{"message": "Фильм удален"})
}
func HandleError(c *gin.Context, err error, status int) {
	log.Println(err) // Логируем ошибку для отладки
	c.JSON(status, gin.H{"error": err.Error()})
}
func GetMoviesFilter(c *gin.Context) {
	var movies []models.Movie

	// Получение параметров запроса для фильтрации и сортировки
	title := c.Query("title")
	genre := c.Query("genre")
	Year := c.Query("year")
	sortBy := c.Query("sort_by") // например: "title", "release_year", "rating"
	order := c.Query("order")    // например: "ASC" или "DESC"

	query := db.DB

	// Фильтрация
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}
	if genre != "" {
		query = query.Where("genre ILIKE ?", "%"+genre+"%") // Можно использовать ILIKE для нечувствительной к регистру фильтрации
	}
	if Year != "" {
		query = query.Where("year = ?", Year)
	}

	// Сортировка
	if sortBy != "" {
		switch sortBy {
		case "title", "year", "rating":
			if order != "ASC" && order != "DESC" {
				order = "ASC" // Устанавливаем порядок по умолчанию
			}
			query = query.Order(sortBy + " " + order)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный параметр сортировки"})
			return
		}
	}

	if err := query.Find(&movies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении фильмов"})
		return
	}

	c.JSON(http.StatusOK, movies)
}
