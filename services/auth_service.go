package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/smtp"
	"os"
	"time"

	"e-commerce/models"
	"e-commerce/utils"

	"gorm.io/gorm"
)

// ---------- helper: send email using SMTP (Gmail) ----------
func sendEmail(to, subject, body string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	if from == "" || password == "" {
		fmt.Printf("SMTP not configured - email to %s: %s\n%s\n", to, subject, body)
		return nil
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from, to, subject, body)

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
}

// ---------- Signup ----------
func SignupService(db *gorm.DB, fullName, email, password string) error {
	var existing models.User
	if err := db.Where("email = ?", email).First(&existing).Error; err == nil {
		return errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user := models.User{
		FullName:     fullName,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         "user",
		IsVerified:   false,
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	// send signup OTP
	return SendOTPService(db, user.ID, user.Email, "signup")
}

// ---------- Login ----------
func LoginService(db *gorm.DB, email, password string) (string, string, error) {
	// cleanup expired refresh tokens (optional housekeeping)
	db.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{})

	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", "", errors.New("user not found")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", "", errors.New("invalid credentials")
	}

	if !user.IsVerified {
		return "", "", errors.New("email not verified")
	}

	accessToken, err := utils.GenerateJWT(int(user.ID), user.Role)
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}

	refreshPlain, hashedToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	if err := utils.SaveRefreshToken(db, user.ID, hashedToken, time.Now().Add(7*24*time.Hour)); err != nil {
		return "", "", errors.New("failed to save refresh token")
	}

	return accessToken, refreshPlain, nil
}

// ---------- Forgot Password ----------
func ForgotPasswordService(db *gorm.DB, email string) error {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	// create and send OTP for reset_password
	return SendOTPService(db, user.ID, user.Email, "reset_password")
}

// ---------- Reset Password ----------
func ResetPasswordService(db *gorm.DB, email, otpCode, newPassword string) error {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	var otp models.OTP
	if err := db.Where("user_id = ? AND otp_code = ? AND purpose = ? AND is_used = ?", user.ID, otpCode, "reset_password", false).
		First(&otp).Error; err != nil {
		return errors.New("invalid otp")
	}

	if time.Now().After(otp.ExpiresAt) {
		return errors.New("otp expired")
	}

	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.PasswordHash = hashed
	if err := db.Save(&user).Error; err != nil {
		return errors.New("failed to update password")
	}

	otp.IsUsed = true
	db.Save(&otp)

	return nil
}

// ---------- Refresh ----------
func RefreshService(db *gorm.DB, refreshToken string) (string, error) {
	rt, err := utils.ValidateRefreshToken(db, refreshToken)
	if err != nil {
		return "", err
	}

	var user models.User
	if err := db.First(&user, rt.UserID).Error; err != nil {
		return "", errors.New("user not found")
	}

	newToken, err := utils.GenerateJWT(int(user.ID), user.Role)
	return newToken, err
}

// ---------- Logout ----------
func LogoutService(db *gorm.DB, refreshToken string) error {
	return utils.DeleteRefreshToken(db, refreshToken)
}

// ---------- OTP helpers ----------
func generateSecureOTP() (string, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", nBig.Int64()), nil
}

// ---------- OTP functions ----------
func SendOTPService(db *gorm.DB, userID uint, email, purpose string) error {
	db.Where("user_id = ? AND purpose = ? AND is_used = ?", userID, purpose, false).Delete(&models.OTP{})

	code, err := generateSecureOTP()
	if err != nil {
		return errors.New("failed to generate otp")
	}

	otp := models.OTP{
		UserID:    userID,
		OTPCode:   code,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Purpose:   purpose,
		IsUsed:    false,
	}

	if err := db.Create(&otp).Error; err != nil {
		return err
	}

	body := fmt.Sprintf("Your OTP code is: %s\nIt expires in 10 minutes.", code)
	subject := "Your OTP Code"
	return sendEmail(email, subject, body)
}

func VerifyOTPService(db *gorm.DB, email, otpCode string) error {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	var otp models.OTP
	if err := db.Where("user_id = ? AND otp_code = ? AND is_used = ?", user.ID, otpCode, false).First(&otp).Error; err != nil {
		return errors.New("invalid otp")
	}

	if time.Now().After(otp.ExpiresAt) {
		return errors.New("otp expired")
	}

	otp.IsUsed = true
	db.Save(&otp)

	user.IsVerified = true
	db.Save(&user)

	return nil
}

func ResendOTPService(db *gorm.DB, email string) error {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	return SendOTPService(db, user.ID, user.Email, "signup")
}