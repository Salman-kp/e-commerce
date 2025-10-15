package routes

import (
	"e-commerce/controllers"
	"github.com/gin-gonic/gin"
)

func AdminViewRoutes(r *gin.Engine) {
	view := r.Group("/view")
	{
		view.GET("/login", controllers.ShowLoginPage)
		view.GET("/dashboard", controllers.ShowDashboard)
		view.GET("/users", controllers.ShowUsersPage)
		view.GET("/products", controllers.ShowProductsPage)
		view.GET("/orders", controllers.ShowOrdersPage)
	}
}