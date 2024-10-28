package main

import (
	"Project/db"
	"Project/middleware"
	"Project/models"
	"Project/routes"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		routes.HandleError(c, err, http.StatusBadRequest)
		return
	}
	// Валидация данных
	if err := user.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка существующего пользователя
	var existingUser models.User
	if err := db.DB.Where("username = ? OR email = ?", user.Username, user.Email).First(&existingUser).Error; err == nil {
		// Если пользователь найден, возвращаем ошибку
		c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким именем или email уже существует"})
		return
	}

	// Хеширование пароля
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хешировании пароля"})
		log.Println("Ошибка при хешировании пароля:", err)
		return
	}

	// Сохранение пользователя в базу данных
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось зарегистрировать пользователя"})
		log.Println("Ошибка при сохранении пользователя в базу данных:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Регистрация прошла успешно"})
}
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Incoming request: %s %s", c.Request.Method, c.Request.URL)
		c.Next()
		log.Printf("Response status: %d", c.Writer.Status())
	}
}
func main() {
	// Инициализация базы данных
	db.Connect()
	models.Migrate(db.DB)

	// Инициализация сервера Gin
	r := gin.Default()
	r.Use(Logger())
	// Группа маршрутов для аутентификации
	auth := r.Group("/auth")
	{
		auth.POST("/register", RegisterUser) // Регистрация пользователя
		auth.POST("/login", routes.LoginUser)
		auth.POST("/refresh", middleware.RefreshToken)
	}
	// Защищённые маршруты
	protected := r.Group("/")
	protected.Use(middleware.Auth())
	{
		protected.GET("/movies", routes.GetMovies)              // Получение всех фильмов
		protected.GET("/users", routes.GetUsers)                // Получение всех пользователей
		protected.GET("/movies-filter", routes.GetMoviesFilter) // фильмы с фильтрацией
		// Добавьте другие защищённые маршруты, если нужно
	}

	authRoutes := r.Group("/auth")
	authRoutes.POST("/logout", routes.LogoutUser) // Маршрут для логаута

	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.Auth(), middleware.RoleCheck("admin"))
	{
		adminRoutes.POST("/movies", routes.AddMovie)
		adminRoutes.PUT("/movies/:id", routes.UpdateMovie)
		adminRoutes.DELETE("/movies/:id", routes.DeleteMovie)
	}
	reviewRoutes := r.Group("/reviews")
	reviewRoutes.Use(middleware.Auth()) // Только для авторизованных пользователей
	{
		reviewRoutes.POST("/", routes.AddReview) // Добавление отзыва
		reviewRoutes.GET("/:id", routes.GetReviewsByMovie)
		reviewRoutes.DELETE("/:id", routes.DeleteReview)
		reviewRoutes.PUT("/:id", routes.UpdateReview)
		// Получение отзывов по фильму
	}
	// Пример маршрута для проверки работы сервера
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Запуск сервера на порту 8081
	if err := r.Run(":8081"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}

}
