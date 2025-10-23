package controllers

import (
	"net/http"
	"strconv"

	"e-commerce/config"
	"e-commerce/models"

	"github.com/gin-gonic/gin"
)

// ---------------- INPUT STRUCT ----------------
type ProductInput struct {
	Name          string  `json:"name" binding:"required"`
	Description   string  `json:"description"`
	Price         float64 `json:"price" binding:"required"`
	StockQuantity int     `json:"stock_quantity" binding:"required"`
	Category      string  `json:"category"`
	ImageURL      string  `json:"image_url"`
}

// ---------------- CREATE PRODUCT ----------------
func CreateProductHandler(c *gin.Context) {
	name := c.PostForm("Name")
	description := c.PostForm("Description")
	priceStr := c.PostForm("Price")
	stockStr := c.PostForm("StockQuantity")
	category := c.PostForm("Category")
	imageURL := c.PostForm("ImageURL")

	// Validate required fields
	if name == "" || description == "" || priceStr == "" || stockStr == "" || category == "" || imageURL == "" {
		c.String(http.StatusBadRequest, "All fields are required")
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid price")
		return
	}

	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid stock quantity")
		return
	}

	product := models.Product{
		Name:          name,
		Description:   description,
		Price:         price,
		StockQuantity: stock,
		Category:      category,
		ImageURL:      imageURL,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to create product")
		return
	}

	c.Redirect(http.StatusSeeOther, "/view/products")
}

// ---------------- UPDATE PRODUCT ----------------
func UpdateProductHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid product ID")
		return
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.String(http.StatusNotFound, "Product not found")
		return
	}

	// Read form values
	name := c.PostForm("Name")
	description := c.PostForm("Description")
	priceStr := c.PostForm("Price")
	stockStr := c.PostForm("StockQuantity")
	category := c.PostForm("Category")
	imageURL := c.PostForm("ImageURL")

	// Validate required fields
	if name == "" || description == "" || priceStr == "" || stockStr == "" || category == "" || imageURL == "" {
		c.String(http.StatusBadRequest, "All fields are required")
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid price")
		return
	}

	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid stock quantity")
		return
	}

	product.Name = name
	product.Description = description
	product.Price = price
	product.StockQuantity = stock
	product.Category = category
	product.ImageURL = imageURL

	if err := config.DB.Save(&product).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to update product")
		return
	}

	c.Redirect(http.StatusSeeOther, "/view/products")
}

// ---------------- DELETE PRODUCT ----------------
func DeleteProductHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid product ID")
		return
	}

	if err := config.DB.Delete(&models.Product{}, id).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to delete product")
		return
	}

	c.Redirect(http.StatusSeeOther, "/view/products")
}


/*---------------------------------------------------- JSON BASED-------------------------------------------------*/

// ---------------- CREATE PRODUCT ----------------
// func CreateProductHandler(c *gin.Context) {
// 	var input ProductInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	product := models.Product{
// 		Name:          input.Name,
// 		Description:   input.Description,
// 		Price:         input.Price,
// 		StockQuantity: input.StockQuantity,
// 		Category:      input.Category,
// 		ImageURL:      input.ImageURL,
// 	}
// 	if err := config.DB.Create(&product).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "product": product})
// }

// ---------------- UPDATE PRODUCT ----------------
// func UpdateProductHandler(c *gin.Context) {
// 	idParam := c.Param("id")
// 	id, err := strconv.ParseUint(idParam, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
// 		return
// 	}
// 	var product models.Product
// 	if err := config.DB.First(&product, uint(id)).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
// 		return
// 	}
// 	var input ProductInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if input.Name != "" {
// 		product.Name = input.Name
// 	}
// 	if input.Description != "" {
// 		product.Description = input.Description
// 	}
// 	if input.Price != 0 {
// 		product.Price = input.Price
// 	}
// 	if input.StockQuantity != 0 {
// 		product.StockQuantity = input.StockQuantity
// 	}
// 	if input.Category != "" {
// 		product.Category = input.Category
// 	}
// 	if input.ImageURL != "" {
// 		product.ImageURL = input.ImageURL
// 	}
// 	if err := config.DB.Save(&product).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "product": product})
// }

// ---------------- DELETE PRODUCT ----------------

// func DeleteProductHandler(c *gin.Context) {
// 	idParam := c.Param("id")
// 	id, err := strconv.ParseUint(idParam, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
// 		return
// 	}
// 	if err := config.DB.Delete(&models.Product{}, uint(id)).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
// }

// ---------------- GET ALL PRODUCTS (PUBLIC) ----------------
func GetProductsHandler(c *gin.Context) {
	var products []models.Product
	if err := config.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

// ---------------- GET PRODUCT BY ID (PUBLIC) ----------------
func GetProductByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}
