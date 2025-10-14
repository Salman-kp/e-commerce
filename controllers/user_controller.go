package controllers

import (
	"e-commerce/config"
	"e-commerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfileHandler(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDInterface.(int)

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          user.ID,
		"full_name":   user.FullName,
		"email":       user.Email,
		"role":        user.Role,
		"is_blocked":  user.IsBlocked,
		"is_verified": user.IsVerified,
		"avatar_url":  user.AvatarURL,
		"address":     user.Address,
	})
}

func UpdateProfileHandler(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDInterface.(int)

	var body struct {
		FullName  *string `json:"full_name"`
		Address   *string `json:"address"`
		AvatarURL *string `json:"avatar_url"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if body.FullName != nil {
		user.FullName = *body.FullName
	}
	if body.Address != nil {
		user.Address = *body.Address
	}
	if body.AvatarURL != nil {
		user.AvatarURL = body.AvatarURL
	}
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"profile": gin.H{
			"id":         user.ID,
			"full_name":  user.FullName,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
			"address":    user.Address,
		},
	})
}
