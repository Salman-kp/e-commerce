package middlewares

import (
	"net/http"
	"strings"

	"e-commerce/config"
	"e-commerce/utils"

	"github.com/gin-gonic/gin"
)

// ---------------- AdminAuthMiddleware
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var accessToken string
		authHeader := c.GetHeader("Authorization")

		if authHeader != "" {
			accessToken = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			cookieToken, err := c.Cookie("access_token")
			if err != nil || cookieToken == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access token"})
				c.Abort()
				return
			}
			accessToken = cookieToken
		}

		userID, role, err := utils.ValidateJWT(accessToken)
		if err != nil {
			if userID != 0 {
				rt, err := utils.GetRefreshTokenByUserID(config.DB, uint(userID))
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login again"})
					c.Abort()
					return
				}

				_, err = utils.ValidateRefreshToken(config.DB, rt.Token)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login again"})
					c.Abort()
					return
				}

				newAccessToken, err := utils.GenerateJWT(userID, role)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login again"})
					c.Abort()
					return
				}

				c.SetCookie("access_token", newAccessToken, 30*60, "/", "localhost", false, true)
				accessToken = newAccessToken
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login again"})
				c.Abort()
				return
			}
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Admins only"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	}
}


//---------------- UserAuthMiddleware
func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var accessToken string
		authHeader := c.GetHeader("Authorization")

		if authHeader != "" {
			accessToken = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			cookieToken, err := c.Cookie("access_token")
			if err != nil || cookieToken == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access token"})
				c.Abort()
				return
			}
			accessToken = cookieToken
		}

		userID, role, err := utils.ValidateJWT(accessToken)
		if err != nil {
			if userID != 0 {
				rt, err := utils.GetRefreshTokenByUserID(config.DB, uint(userID))
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login again"})
					c.Abort()
					return
				}

				_, err = utils.ValidateRefreshToken(config.DB, rt.Token)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login again"})
					c.Abort()
					return
				}

				newAccessToken, err := utils.GenerateJWT(userID, role)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login again"})
					c.Abort()
					return
				}

				c.SetCookie("access_token", newAccessToken, 30*60, "/", "localhost", false, true)
				accessToken = newAccessToken
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login again"})
				c.Abort()
				return
			}
		}

		if role != "user" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Users only"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	}
}