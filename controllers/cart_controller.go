package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"e-commerce/config"
	"e-commerce/models"
)

// ---------------- RESPONSE STRUCTS ----------------
type CartItemResponse struct {
	ID        uint           `json:"id"`
	Product   ProductSummary `json:"product"`
	Quantity  int            `json:"quantity"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type ProductSummary struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stock_quantity"`
	ImageURL      string  `json:"image_url"`
}

// ---------------- ADD TO CART ----------------
func AddToCart(c *gin.Context) {
	userIDInt, exists := c.Get("userID")
	id, ok := userIDInt.(int)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(id)

	var input struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if product exists
	var product models.Product
	if err := config.DB.First(&product, input.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Check stock
	if input.Quantity > product.StockQuantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock available"})
		return
	}

	// Check if already in cart
	var existing models.CartItem
	if err := config.DB.Where("user_id = ? AND product_id = ?", userID, input.ProductID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Product already in cart. Use PUT to update quantity."})
		return
	}

	// Create cart item
	cartItem := models.CartItem{
		UserID:    userID,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
	}
	if err := config.DB.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
		return
	}

	// Preload product
	if err := config.DB.Preload("Product").First(&cartItem, cartItem.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart item"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product added to cart", "cart_item": mapCartItem(cartItem)})
}

// ---------------- GET CART ITEMS ----------------
func GetCartItems(c *gin.Context) {
	userIDInt, exists := c.Get("userID")
	id, ok := userIDInt.(int)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(id)

	var cartItems []models.CartItem
	config.DB.Preload("Product").Where("user_id = ?", userID).Find(&cartItems)

	var resp []CartItemResponse
	for _, item := range cartItems {
		resp = append(resp, mapCartItem(item))
	}

	c.JSON(http.StatusOK, gin.H{"cart_items": resp})
}

// ---------------- UPDATE CART ITEM ----------------
func UpdateCartItem(c *gin.Context) {
	userIDInt, exists := c.Get("userID")
	id, ok := userIDInt.(int)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(id)

	cartID := c.Param("id")
	cartIDUint, err := strconv.ParseUint(cartID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart item ID"})
		return
	}

	var input struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cartItem models.CartItem
	if err := config.DB.Preload("Product").Where("id = ? AND user_id = ?", cartIDUint, userID).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	if input.Quantity > cartItem.Product.StockQuantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock available"})
		return
	}

	cartItem.Quantity = input.Quantity
	config.DB.Save(&cartItem)

	c.JSON(http.StatusOK, gin.H{"message": "Cart item updated", "cart_item": mapCartItem(cartItem)})
}

// ---------------- DELETE CART ITEM ----------------
func DeleteCartItem(c *gin.Context) {
	userIDInt, exists := c.Get("userID")
	id, ok := userIDInt.(int)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(id)

	cartID := c.Param("id")
	cartIDUint, err := strconv.ParseUint(cartID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart item ID"})
		return
	}

	result := config.DB.Where("id = ? AND user_id = ?", cartIDUint, userID).Delete(&models.CartItem{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove cart item"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item removed"})
}

// ---------------- HELPER ----------------
func mapCartItem(item models.CartItem) CartItemResponse {
	return CartItemResponse{
		ID: item.ID,
		Product: ProductSummary{
			ID:            item.Product.ID,
			Name:          item.Product.Name,
			Description:   item.Product.Description,
			Price:         item.Product.Price,
			StockQuantity: item.Product.StockQuantity,
			ImageURL:      item.Product.ImageURL,
		},
		Quantity:  item.Quantity,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}