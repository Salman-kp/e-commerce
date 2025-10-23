package controllers

import (
	"net/http"

	"e-commerce/config"
	"e-commerce/models"
	"e-commerce/services"
	"e-commerce/utils"

	"github.com/gin-gonic/gin"
)

// ------------------ SIGNUP ------------------
func SignupHandler(c *gin.Context) {
	var body struct {
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.SignupService(config.DB, body.FullName, body.Email, body.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Signup successful. Please verify your email using OTP."})
}

// ------------------ LOGIN ------------------
func LoginHandler(c *gin.Context) {
	// var body struct {
	// 	Email    string `json:"email" binding:"required,email"`
	// 	Password string `json:"password" binding:"required"`
	// }
	// if err := c.ShouldBindJSON(&body); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	// accessToken, refreshToken, role, err := services.LoginService(config.DB, body.Email, body.Password)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	// 	return
	// }

	email := c.PostForm("email")
	password := c.PostForm("password")

	if email == "" || password == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"title": "Login Page",
			"error": "⚠️ Email and password cannot be empty",
		})
		return
	}

	accessToken, role, err := services.LoginService(config.DB, email, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"title": "Login Page",
			"error": "❌ Invalid email or password",
		})
		return
	}
	//Set accessToken in cookie
	c.SetCookie("access_token", accessToken, 30*70, "/", "localhost", false, true) // 30 minutes
	
	if role == "admin" {
		c.Redirect(http.StatusSeeOther, "/view/dashboard")
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login Page",
		"error": "❌ Only admin users are allowed",
	})
}

// ------------------ OTP ------------------
func SendOTPHandler(c *gin.Context) {
	var body struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if err := config.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := services.SendOTPService(config.DB, user.ID, user.Email, "generic"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func VerifyOTPHandler(c *gin.Context) {
	var body struct {
		Email string `json:"email" binding:"required,email"`
		OTP   string `json:"otp" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.VerifyOTPService(config.DB, body.Email, body.OTP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func ResendOTPHandler(c *gin.Context) {
	var body struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ResendOTPService(config.DB, body.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OTP resent successfully"})
}

// -------------------- Forgot Password --------------------
func ForgotPasswordHandler(c *gin.Context) {
	var body struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ForgotPasswordService(config.DB, body.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to your email"})
}

// -------------------- Reset Password --------------------
func ResetPasswordHandler(c *gin.Context) {
	var body struct {
		Email       string `json:"email" binding:"required,email"`
		OTP         string `json:"otp" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ResetPasswordService(config.DB, body.Email, body.OTP, body.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

// ------------------ REFRESH TOKEN ------------------
func RefreshTokenHandler(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil || accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access token"})
		return
	}

	userID, _, err := utils.ValidateJWT(accessToken)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
		return
	}

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login again"})
		return
	}

	newToken, err := services.RefreshService(config.DB, uint(userID))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login again"})
		return
	}

	c.SetCookie("access_token", newToken, 30*60, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"access_token": newToken})
}

// ------------------ LOGOUT ------------------
func LogoutHandler(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil || accessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing access token"})
		return
	}

	userID, _, err := utils.ValidateJWT(accessToken)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid access token"})
		return
	}

	if err := services.LogoutService(config.DB, uint(userID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}