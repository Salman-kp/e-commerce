package controllers

import (
	"e-commerce/config"
	"e-commerce/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ---------------- LOGIN ----------------
func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login Page",
	})
}
// ---------------- DASHBOARD ----------------
func ShowDashboard(c *gin.Context) {
	var totalUsers, totalProducts, totalOrders int64

	// count total users
	if err := config.DB.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		totalUsers = 0
	}
	// Count total products
	if err := config.DB.Model(&models.Product{}).Count(&totalProducts).Error; err != nil {
		totalProducts = 0
	}
	// count total orders
	if err := config.DB.Model(&models.Order{}).Count(&totalOrders).Error; err != nil {
		totalOrders = 0
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title":          "Admin Dashboard",
		"total_users":    totalUsers,
		"total_products": totalProducts,
		"total_orders":   totalOrders,
		"Active":         "dashboard", // for sidebar highlighting
	})
}




// ---------------- USERS ----------------
func ShowUsersPage(c *gin.Context) {
	var users []models.User
	if err := config.DB.Order("id ASC").Find(&users).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "users.html", gin.H{"error": "Failed to fetch users"})
		return
	}

	c.HTML(http.StatusOK, "users.html", gin.H{
		"title":  "Manage Users",
		"users":  users,
		"Active": "users",
	})
}
// ---------------- SHOW EDIT USER ----------------
func ShowEditUserPage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.String(http.StatusNotFound, "User not found")
		return
	}

	c.HTML(http.StatusOK, "edit_user.html", gin.H{
		"title": "Edit User",
		"user":  user,
	})
}




// ---------------- PRODUCTS ----------------
func ShowProductsPage(c *gin.Context) {
	var products []models.Product
	if err := config.DB.Order("id ASC").Find(&products).Error; err != nil {
		products = []models.Product{}
	}

	c.HTML(http.StatusOK, "products.html", gin.H{
		"title":    "Manage Products",
		"products": products,
		"Active":   "products",
	})
}
// ---------------- CREATE PRODUCT PAGE ----------------
func ShowCreateProductPage(c *gin.Context) {
	c.HTML(http.StatusOK, "create_product.html", gin.H{
		"title": "Add Product",
	})
}
// ---------------- EDIT PRODUCT PAGE ----------------
func ShowEditProductPage(c *gin.Context) {
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

	c.HTML(http.StatusOK, "edit_product.html", gin.H{
		"title":   "Edit Product",
		"product": product,
	})
}



// ---------------- ORDERS ----------------
func ShowOrdersPage(c *gin.Context) {
	var orders []models.Order
	if err := config.DB.Preload("User").Preload("OrderItems").Order("id ASC").Find(&orders).Error; err != nil {
		orders = []models.Order{}
	}
	c.HTML(http.StatusOK, "orders.html", gin.H{
		"title":  "Manage Orders",
		"orders": orders,
		"Active": "orders",
	})
}

// -------MIDDLEWARE
func MethodOverride() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			if method := c.PostForm("_method"); method != "" {
				c.Request.Method = method
			}
		}
		c.Next()
	}
}
