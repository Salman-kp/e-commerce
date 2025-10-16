package routes

import (
	"e-commerce/controllers"
	"e-commerce/middlewares"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	admin.Use(middlewares.AdminAuthMiddleware())
	{
		admin.POST("/products", controllers.CreateProductHandler)
	    admin.PUT("/products/:id", controllers.UpdateProductHandler)
		admin.DELETE("/products/:id", controllers.DeleteProductHandler)
	    admin.POST("/products/:id/production", controllers.StartProductionHandler)              // start production route
		admin.PUT("/products/:id/production/status", controllers.UpdateProductionStatusHandler) // update production status route
		admin.GET("/products/:id/production", controllers.GetProductionDetailsHandler)          // get production details route
		admin.GET("/products/production", controllers.GetAllProductionsHandler)                 // get all productions route
	}
}
