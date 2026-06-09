package controllers

import (
	"fmt"
	"net/http"
	"time"

	"backend/config"
	"backend/models"

	"github.com/gin-gonic/gin"
)

// AdminGetPortalUsers fetches all active registered portal users
func AdminGetPortalUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Order("created_at desc").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch portal users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// AdminUpdatePortalUserPayload represents the payload for updating a portal user
type AdminUpdatePortalUserPayload struct {
	FullName    string `json:"fullName" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Address     string `json:"address" binding:"required"`
	DateOfBirth string `json:"dateOfBirth" binding:"required"`
}

// AdminUpdatePortalUser updates a portal user's details
func AdminUpdatePortalUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var payload AdminUpdatePortalUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload format"})
		return
	}

	dob, err := time.Parse("2006-01-02", payload.DateOfBirth)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date of birth format. Use YYYY-MM-DD"})
		return
	}

	user.FullName = payload.FullName
	user.PhoneNumber = payload.PhoneNumber
	user.Address = payload.Address
	user.DateOfBirth = dob

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update portal user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "data": user})
}

// AdminDeletePortalUser soft-deletes a portal user
func AdminDeletePortalUser(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// AdminGetArchivedAdmins fetches soft-deleted staff users
func AdminGetArchivedAdmins(c *gin.Context) {
	var admins []models.Admin
	if err := config.DB.Unscoped().Where("deleted_at IS NOT NULL").Order("deleted_at desc").Find(&admins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch archived staff"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": admins})
}

// AdminRecoverAdmin restores a soft-deleted staff user
func AdminRecoverAdmin(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Unscoped().Model(&models.Admin{}).Where("id = ?", id).Update("deleted_at", nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to recover staff account"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account recovered successfully"})
}

// AdminGetArchivedPortalUsers fetches soft-deleted portal users
func AdminGetArchivedPortalUsers(c *gin.Context) {
	var users []models.User
	if err := config.DB.Unscoped().Where("deleted_at IS NOT NULL").Order("deleted_at desc").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch archived users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// AdminRecoverPortalUser restores a soft-deleted portal user
func AdminRecoverPortalUser(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Unscoped().Model(&models.User{}).Where("id = ?", id).Update("deleted_at", nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to recover user account"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Account recovered successfully"})
}
