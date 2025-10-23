package routes

import (
	"e-commerce/controllers"
	"e-commerce/middlewares"

	"github.com/gin-gonic/gin"
)

func PaymentRoutes(r *gin.Engine) {
	payments := r.Group("/payments")
	payments.Use(middlewares.UserAuthMiddleware())
	{
		payments.POST("/create", controllers.CreatePaymentIntent)
	}
	// Admin routes
	adminPayments := r.Group("/admin/payments")
	adminPayments.Use(middlewares.AdminAuthMiddleware())
	{
		adminPayments.PUT("/:payment_id/update", controllers.UpdatePaymentStatus)
	}

}
