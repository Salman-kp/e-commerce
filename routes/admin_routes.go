package routes

import (
	"e-commerce/controllers"
	"e-commerce/middlewares"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {

	admin := r.Group("/admin")
	admin.Use(middlewares.AdminAuthMiddleware())
	{
		admin.GET("/users", controllers.GetAllUsersHandler)
		admin.GET("/users/:id", controllers.GetUserByIDHandler)
		admin.PUT("/users/:id", controllers.UpdateUserHandler)
		admin.DELETE("/users/:id", controllers.DeleteUserHandler)
		admin.POST("/users/:id/block", controllers.BlockUserHandler)
		admin.POST("/users/:id/unblock", controllers.UnblockUserHandler)
	}
}