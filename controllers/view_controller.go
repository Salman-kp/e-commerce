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
	if err := config.DB.Find(&users).Error; err != nil {
		users = []models.User{}
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
	c.HTML(http.StatusOK, "products.html", gin.H{
		"title":  "Products Page",
		"Active": "products",
	})
}

// ---------------- ORDERS ----------------
func ShowOrdersPage(c *gin.Context) {
	c.HTML(http.StatusOK, "orders.html", gin.H{
		"title":  "Orders Page",
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
