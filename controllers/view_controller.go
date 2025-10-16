package controllers

import (
	"e-commerce/config"
	"e-commerce/models"
	"net/http"

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
	var totalUsers int64
	if err := config.DB.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		totalUsers = 0 // fallback if query fails
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title":       "Admin Dashboard",
		"total_users": totalUsers,
		"Active":      "dashboard", // for sidebar highlighting
	})
}

// ---------------- USERS ----------------
func ShowUsersPage(c *gin.Context) {
	c.HTML(http.StatusOK, "users.html", gin.H{
		"title":  "Manage Users",
		"Active": "users",
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
