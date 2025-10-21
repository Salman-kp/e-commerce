package utils

import (
	"crypto/rand"
	"e-commerce/models"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ------------------ Password Hashing ------------------
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ------------------ JWT Functions ------------------
func GenerateJWT(userID int, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"userId": userID,
		"role":   role,
		"exp":    time.Now().Add(time.Minute * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateJWT validates token and returns userId + role
func ValidateJWT(tokenStr string) (int, string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return 0, "", fmt.Errorf("JWT_SECRET not set")
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", fmt.Errorf("invalid token claims")
	}

	userIDFloat, ok := claims["userId"].(float64)
	if !ok {
		return 0, "", fmt.Errorf("invalid userId in token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", fmt.Errorf("invalid role in token")
	}

	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors == jwt.ValidationErrorExpired {
			return int(userIDFloat), role, fmt.Errorf("token expired")
		}
		return 0, "", err
	}

	if !parsedToken.Valid {
		return 0, "", fmt.Errorf("invalid token")
	}

	return int(userIDFloat), role, nil
}

// ------------------ Refresh Token Functions ------------------

// Generate a random refresh token ( hashed)
func GenerateRefreshToken() (string,  error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)
	return token, nil
}

// Save refresh token to DB
func SaveRefreshToken(db *gorm.DB, userID uint, token string, expiresAt time.Time) error {
	rt := models.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	return db.Create(&rt).Error
}

// check user refresh token in db if in .check valid or invalid
func GetRefreshTokenByUserID(db *gorm.DB, userID uint) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := db.
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Order("expires_at DESC").
		First(&rt).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no valid refresh token found")
		}
		return nil, err
	}

	return &rt, nil
}

// Validate refresh token from DB
func ValidateRefreshToken(db *gorm.DB, token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired refresh token")
		}
		return nil, err
	}
	return &rt, nil
}

// Delete refresh token from DB (logout)
func DeleteRefreshToken(db *gorm.DB, token string) error {
	return db.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}