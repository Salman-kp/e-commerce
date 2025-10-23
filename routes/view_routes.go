package routes

import (
	"e-commerce/controllers"
	"e-commerce/middlewares"

	"github.com/gin-gonic/gin"
)

func AdminViewRoutes(r *gin.Engine) {
	r.GET("/login", controllers.ShowLoginPage)
	view := r.Group("/view")
	view.Use(middlewares.AdminAuthMiddleware())
	{
		view.GET("/dashboard", controllers.ShowDashboard)
		view.GET("/users", controllers.ShowUsersPage)
		view.GET("/products", controllers.ShowProductsPage)
		view.GET("/orders", controllers.ShowOrdersPage)

		//---------USER EDTITE
		view.GET("/users/edit/:id", controllers.ShowEditUserPage)
		//---------- PRODUCT CREATE & UPDATE
		view.GET("/products/create", controllers.ShowCreateProductPage)
		view.GET("/products/edit/:id", controllers.ShowEditProductPage)
	}
}
