package models

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"regexp"
	"time"
)

var jwtKey = []byte("your_secret_key") // Секретный ключ для подписи токена

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Role     string `json:"role" gorm:"default:'user'"` // поле для роли
}

// Функция для валидации данных пользователя
func (user *User) Validate() error {
	// Проверка формата email
	if !isValidEmail(user.Email) {
		return errors.New("некорректный формат email")
	}
	// Проверка длины имени пользователя
	if len(user.Username) < 3 || len(user.Username) > 20 {
		return errors.New("имя пользователя должно быть от 3 до 20 символов")
	}
	// Проверка длины пароля
	if len(user.Password) < 6 {
		return errors.New("пароль должен быть не менее 6 символов")
	}
	return nil
}

// Функция для проверки корректности email
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// Хеширование пароля
func (user *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

// Проверка пароля при входе
func (user *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

// Структура для хранения данных токена
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"` //  поле для роли
	jwt.StandardClaims
}

// Генерация токена
func GenerateJWT(username string, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		Role:     role, // Добавьте роль здесь
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
func ParseJWT(tokenString string, claims *Claims) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	return claims.Username, nil
}
func RefreshJWT(oldToken string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(oldToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("старый токен недействителен")
	}

	// Устанавливаем новый срок действия токена
	expirationTime := time.Now().Add(24 * time.Hour)
	claims.ExpiresAt = expirationTime.Unix()
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return newToken.SignedString(jwtKey)
}
