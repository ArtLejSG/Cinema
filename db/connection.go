package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// Переменная для хранения экземпляра базы данных
var DB *gorm.DB

// Функция для подключения к базе данных
func Connect() {
	dsn := "host=localhost user=movie_user password=Kiri567 dbname=movies_db port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}
}
