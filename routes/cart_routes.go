package routes

import (
	"e-commerce/controllers"
	"e-commerce/middlewares"

	"github.com/gin-gonic/gin"
)

func CartRoutes(r *gin.Engine){
	cart:=r.Group("/cart")
	cart.Use(middlewares.UserAuthMiddleware())
	{
		cart.POST("",controllers.AddToCart)
		cart.GET("",controllers.GetCartItems)
		cart.PUT("/:id",controllers.UpdateCartItem)
		cart.DELETE("/:id",controllers.DeleteCartItem)
	}
}