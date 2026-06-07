package controllers

import (
	"net/http"
	"time"

	"backend/config"
	"backend/models"

	"github.com/gin-gonic/gin"
)

// GetProfile fetches the current user's profile based on the JWT phone number
func GetProfile(c *gin.Context) {
	phone, exists := c.Get("phone")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := config.DB.Where("phone_number = ?", phone).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, models.User{})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile updates the user's personal information
func UpdateProfile(c *gin.Context) {
	phone, exists := c.Get("phone")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var user models.User
	if err := config.DB.Where("phone_number = ?", phone).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.FullName = req.FullName
	user.Address = req.Address

	parsedDate, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err == nil {
		user.DateOfBirth = parsedDate
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
