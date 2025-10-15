package routes

import (
	"github.com/gin-gonic/gin"
	"e-commerce/controllers"
	"e-commerce/middlewares"
)

func AdminRoutes(r *gin.Engine) {

	admin := r.Group("/admin")
	admin.Use(middlewares.AdminAuthMiddleware())
	{
		admin.PUT("/users/:id", controllers.UpdateUserHandler)
		admin.POST("/users/:id/block", controllers.BlockUserHandler)
		admin.POST("/users/:id/unblock", controllers.UnblockUserHandler)
		admin.GET("/users", controllers.GetAllUsersHandler)
		admin.GET("/users/:id", controllers.GetUserByIDHandler)
		admin.DELETE("/users/:id", controllers.DeleteUserHandler)
	}
}