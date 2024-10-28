package models

import "gorm.io/gorm"

// Структура для черного списка токенов
type BlacklistedToken struct {
	ID    uint   `gorm:"primaryKey"`
	Token string `gorm:"unique;not null"` // Сделайте токен уникальным, чтобы избежать дублирования
}

// Функция для добавления токена в черный список
func BlacklistToken(db *gorm.DB, token string) error {
	return db.Create(&BlacklistedToken{Token: token}).Error
}

// Функция для проверки, находится ли токен в черном списке
func IsTokenBlacklisted(db *gorm.DB, token string) (bool, error) {
	var blacklistedToken BlacklistedToken
	err := db.Where("token = ?", token).First(&blacklistedToken).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Токен не найден в черном списке
			return false, nil
		}
		// Произошла ошибка, возвращаем её
		return false, err
	}
	// Токен найден в черном списке
	return true, nil
}
