package models

import (
	"fmt"
	"time"
)

type Review struct {
	ID      uint `gorm:"primaryKey"`
	MovieID uint `gorm:"not null"` // Связь с фильмом
	UserID  uint `gorm:"not null"` // Связь с пользователем
	//Content   string    `gorm:"not null"`       // Содержимое отзыва
	Rating    float64   `gorm:"not null"`       // Рейтинг отзыва
	CreatedAt time.Time `gorm:"autoCreateTime"` // Дата создания
	UpdatedAt time.Time `gorm:"autoUpdateTime"` // Дата обновления
}

// Validates if the review has valid content and rating
func (r *Review) Validate() error {
	//if r.Content == "" {
	//	return fmt.Errorf("содержимое отзыва не может быть пустым")
	//}
	if r.Rating < 1 || r.Rating > 5 {
		return fmt.Errorf("рейтинг должен быть между 1 и 5")
	}
	return nil
}
