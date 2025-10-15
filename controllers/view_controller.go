package controllers

import (
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
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title": "Admin Dashboard",
	})
}

// ---------------- USERS ----------------
func ShowUsersPage(c *gin.Context) {
	c.HTML(http.StatusOK, "users.html", gin.H{
		"title": "Manage Users",
	})
}

// ---------------- PRODUCTS ----------------
func ShowProductsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "products.html", gin.H{
		"title": "Products Page",
	})
}

// ---------------- ORDERS ----------------
func ShowOrdersPage(c *gin.Context) {
	c.HTML(http.StatusOK, "orders.html", gin.H{
		"title": "Orders Page",
	})
}