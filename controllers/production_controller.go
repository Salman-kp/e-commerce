package controllers

import (
	"net/http"
	"strconv"
	"time"

	"e-commerce/config"
	"e-commerce/models"

	"github.com/gin-gonic/gin"
)

// Allowed production status values
var allowedStatuses = map[string]bool{
	"started":     true,
	"in_progress": true,
	"completed":   true,
}

// ---------------- START PRODUCTION ----------------
func StartProductionHandler(c *gin.Context) {
	idParam := c.Param("id")
	productID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, uint(productID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Product not found"})
		return
	}

	// Prevent multiple active productions
	var existing models.ProductProduction
	if err := config.DB.Where("product_id = ? AND status != ?", uint(productID), "completed").First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Production already in progress for this product"})
		return
	}

	production := models.ProductProduction{
		ProductID: uint(productID),
		Status:    "started",
	}

	if err := config.DB.Create(&production).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to start production"})
		return
	}

	// Preload Product
	if err := config.DB.Preload("Product").First(&production, production.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch production with product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Production started", "data": production})
}

// ---------------- UPDATE PRODUCTION STATUS ----------------
func UpdateProductionStatusHandler(c *gin.Context) {
	idParam := c.Param("id")
	productionID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid production ID"})
		return
	}

	var production models.ProductProduction
	if err := config.DB.First(&production, uint(productionID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Production not found"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if !allowedStatuses[input.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid status value"})
		return
	}

	production.Status = input.Status
	if input.Status == "completed" {
		now := time.Now()
		production.CompletedAt = &now
	}

	if err := config.DB.Save(&production).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update production status"})
		return
	}

	// Preload Product
	if err := config.DB.Preload("Product").First(&production, production.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch production with product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Production status updated", "data": production})
}

// ---------------- GET PRODUCTION DETAILS ----------------
func GetProductionDetailsHandler(c *gin.Context) {
	idParam := c.Param("id")
	productionID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid production ID"})
		return
	}

	var production models.ProductProduction
	if err := config.DB.Preload("Product").First(&production, uint(productionID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Production not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": production})
}

// ---------------- GET ALL PRODUCTIONS ----------------
func GetAllProductionsHandler(c *gin.Context) {
	var productions []models.ProductProduction
	if err := config.DB.Preload("Product").Find(&productions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch productions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": productions})
}