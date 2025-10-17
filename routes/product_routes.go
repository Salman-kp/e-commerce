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
	    admin.POST("/products/:id/production", controllers.StartProductionHandler)              
		admin.PUT("/products/:id/production/status", controllers.UpdateProductionStatusHandler) 
		admin.GET("/products/:id/production", controllers.GetProductionDetailsHandler)          
		admin.GET("/products/production", controllers.GetAllProductionsHandler)               
	}
	public := r.Group("/products")
	{
		public.GET("", controllers.GetProductsHandler)
		public.GET("/:id", controllers.GetProductByIDHandler)
	}
}
