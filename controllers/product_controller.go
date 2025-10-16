package controllers

import (
	"e-commerce/config"
	"e-commerce/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ---------------- CREATE PRODUCT ----------------
func CreateProductHandler(c *gin.Context) {
	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product := models.Product{
		Name:          input.Name,
		Description:   input.Description,
		Price:         input.Price,
		StockQuantity: input.StockQuantity,
		Category:      input.Category,
		ImageURL:      input.ImageURL,
	}
	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}
	c.JSON(http.StatusCreated,gin.H{"message":"Product Created Successfully","product":product})
}

// ---------------- UPDATE PRODUCT ----------------
// func UpdateProductHandler(c *gin.Context){
// idParam :=c.Param("id")
// id,err:=strconv.ParseUint(idParam,10,64)

// }