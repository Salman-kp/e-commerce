package routes

import (
	"e-commerce/controllers"
	"e-commerce/middlewares"

	"github.com/gin-gonic/gin"
)

func WishlistRoutes(r *gin.Engine) {
	wishlist := r.Group("/wishlist")
	wishlist.Use(middlewares.UserAuthMiddleware())
	{
		wishlist.POST("", controllers.AddToWishlist)
		wishlist.GET("", controllers.GetWishlist)
		wishlist.DELETE("/:product_id", controllers.RemoveFromWishlist)
	}
}