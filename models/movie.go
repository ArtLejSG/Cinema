package models

import (
	"gorm.io/gorm"
)

var DB *gorm.DB

// Структура для таблицы фильмов
type Movie struct {
	ID          uint    `gorm:"primaryKey"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Year        int     `json:"year"`
	Genre       string  `json:"genre"`
	Rating      float64 `json:"rating"`
}

// Функция для миграции модели
func Migrate(db *gorm.DB) {
	DB = db
	db.AutoMigrate(&User{}, &Movie{}, &BlacklistedToken{}, &Review{})
}
