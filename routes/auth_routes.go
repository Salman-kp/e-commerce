package routes

import (
	"e-commerce/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {

	auth := r.Group("/auth")
	{
		auth.POST("/signup", controllers.SignupHandler)
		auth.POST("/login", controllers.LoginHandler)
		auth.POST("/send-otp", controllers.SendOTPHandler)
		auth.POST("/verify-otp", controllers.VerifyOTPHandler)
		auth.POST("/forgot-password", controllers.ForgotPasswordHandler)
		auth.POST("/reset-password", controllers.ResetPasswordHandler)
		auth.POST("/resend-otp", controllers.ResendOTPHandler)

		// New refresh token endpoints
		auth.POST("/refresh", controllers.RefreshTokenHandler)
		auth.POST("/logout", controllers.LogoutHandler)
	}
}
