package routes

import (
	"e-commerce/controllers"
	"e-commerce/middlewares"

	"github.com/gin-gonic/gin"
)

func OrdeRoutes(r *gin.Engine) {
	order := r.Group("/order")
	order.Use(middlewares.UserAuthMiddleware())
	{
		order.POST("", controllers.PlaceOrder)
		order.GET("", controllers.GetUserOrders)
		order.GET("/:id", controllers.GetOrder)
		order.DELETE("/:id", controllers.DeleteOrder)
	}

	adminOrders := r.Group("/admin/orders")
	adminOrders.Use(middlewares.AdminAuthMiddleware())
	{
		adminOrders.GET("", controllers.GetAllOrders) 
		adminOrders.PUT("/:id", controllers.UpdateOrderStatusAdmin)
	}
}
