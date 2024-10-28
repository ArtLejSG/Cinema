package tests

import (
	"bytes"
	"encoding/json"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"Project/models" // Импортируйте ваш пакет для моделей
	"Project/routes" // Импортируйте ваши маршруты, если они нужны для тестирования
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var db *gorm.DB

func setup() {
	var err error
	dsn := "host=localhost user=movie_user password=Kiri567 dbname=movies_db port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	// Дополнительные действия, такие как миграции, можно выполнить здесь
	models.Migrate(db) // Добавьте миграцию моделей
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.POST("/reviews", routes.AddReview)
	r.GET("/reviews/:movieID", routes.GetReviewsByMovie)
	r.PUT("/reviews/:id", routes.UpdateReview)
	r.DELETE("/reviews/:id", routes.DeleteReview)
	return r
}

func TestMain(m *testing.M) {
	setup()         // Вызываем setup перед запуском тестов
	code := m.Run() // Запускаем тесты
	os.Exit(code)   // Завершаем программу с кодом завершения
}

func TestAddReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	// Подготовка данных для теста
	review := models.Review{
		MovieID: 1,  // Убедитесь, что этот фильм существует
		UserID:  10, // Убедитесь, что пользователь с ID 10 существует
		Content: "Отличный фильм!",
		Rating:  5.0,
	}

	reviewData, err := json.Marshal(review)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/reviews", bytes.NewBuffer(reviewData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("TEST_TOKEN"))

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)

	var createdReview models.Review
	err = json.Unmarshal(res.Body.Bytes(), &createdReview)
	assert.NoError(t, err)
	assert.NotZero(t, createdReview.ID) // Проверка, что отзыв был создан
}

func TestGetReviewsByMovie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/reviews/1", nil) // Предполагается, что фильм с ID 1 существует
	req.Header.Set("Authorization", "Bearer "+os.Getenv("TEST_TOKEN"))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var reviews []models.Review
	err := json.Unmarshal(res.Body.Bytes(), &reviews)
	assert.NoError(t, err)
	assert.NotEmpty(t, reviews) // Проверка, что отзывы были получены
}

func TestUpdateReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	// Подготовка данных для теста
	updatedReview := models.Review{
		Content: "Изменённый отзыв",
		Rating:  4.0,
	}

	// Тестируем обновление существующего отзыва
	updatedReviewData, err := json.Marshal(updatedReview)
	assert.NoError(t, err)

	req, _ := http.NewRequest("PUT", "/reviews/1", bytes.NewBuffer(updatedReviewData)) // Предполагается, что отзыв с ID 1 существует
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("TEST_TOKEN"))

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var review models.Review
	err = json.Unmarshal(res.Body.Bytes(), &review)
	assert.NoError(t, err)
	assert.Equal(t, updatedReview.Content, review.Content) // Проверка, что отзыв был обновлён

	// Тестируем обновление несуществующего отзыва
	req, _ = http.NewRequest("PUT", "/reviews/999", nil) // Предполагается, что отзыв с ID 999 не существует
	req.Header.Set("Authorization", "Bearer "+os.Getenv("TEST_TOKEN"))
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestDeleteReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	req, _ := http.NewRequest("DELETE", "/reviews/1", nil) // Предполагается, что отзыв с ID 1 существует
	req.Header.Set("Authorization", "Bearer "+os.Getenv("TEST_TOKEN"))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNoContent, res.Code)
}
