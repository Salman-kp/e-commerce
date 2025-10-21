package controllers

import (
	"e-commerce/config"
	"e-commerce/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ---------------- ADD TO WISHLIST ----------------
func AddToWishlist(c *gin.Context) {
	userIDInt, exists := c.Get("userID")
	id, ok := userIDInt.(int)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(id)

	var body struct {
		ProductID uint `json:"product_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var existing models.WishlistItem
	if err := config.DB.Where("user_id = ? AND product_id = ?", userID, body.ProductID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Product already in wishlist"})
		return
	}

	wishlist := models.WishlistItem{
		UserID:    userID,
		ProductID: body.ProductID,
	}
	if err := config.DB.Create(&wishlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product to wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added to wishlist successfully"})
}

// ---------------- GET WISHLIST ----------------
func GetWishlist(c *gin.Context) {
	userIDInt, exists := c.Get("userID")
	id, ok := userIDInt.(int)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(id)

	var wishlist []models.WishlistItem
	if err := config.DB.Preload("Product").Where("user_id = ?", userID).Order("created_at desc").Find(&wishlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wishlist": wishlist})
}

// ---------------- REMOVE FROM WISHLIST ----------------
func RemoveFromWishlist(c *gin.Context) {
	userIDInt, exists := c.Get("userID")
	id, ok := userIDInt.(int)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(id)

	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	pid := uint(productID)

	result := config.DB.Where("user_id = ? AND product_id = ?", userID, pid).Delete(&models.WishlistItem{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove product from wishlist"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found in wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product removed from wishlist successfully"})
}