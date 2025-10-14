package routes

import (
	"e-commerce/controllers"
	"e-commerce/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	user := r.Group("/user")
	user.Use(middlewares.UserAuthMiddleware())
	{
		user.GET("/profile", controllers.GetProfileHandler)
		user.PUT("/profile", controllers.UpdateProfileHandler)
	}
}
