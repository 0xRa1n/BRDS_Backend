package controllers

import (
	"fmt"
	"html"
	"net/http"
	"time"

	"backend/config"
	"backend/models"
	"backend/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AdminLogin handles staff/admin authentication
func AdminLogin(c *gin.Context) {
	var req models.AdminLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var admin models.Admin
	if err := config.DB.Where("username = ?", req.Username).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Update LoginHistory
	admin.LoginHistory = append(admin.LoginHistory, time.Now())
	config.DB.Save(&admin)

	token, err := utils.GenerateAdminJWT(admin.Username, admin.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.AdminLoginResponse{
		Token:    token,
		Role:     admin.Role,
		FullName: admin.FullName,
	})
}

// AdminGetUsers fetches all admin/staff users
func AdminGetUsers(c *gin.Context) {
	var admins []models.Admin
	if err := config.DB.Order("created_at desc").Find(&admins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": admins})
}

// AdminCreateUserPayload represents the payload for creating a new admin user
type AdminCreateUserPayload struct {
	FullName string `json:"fullName" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
	Status   string `json:"status" binding:"required"`
}

// AdminCreateUser creates a new admin or staff user
func AdminCreateUser(c *gin.Context) {
	var payload AdminCreateUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	admin := models.Admin{
		UniqueID:     "USR-" + fmt.Sprintf("%d", time.Now().Unix()), // Simple unique ID generator
		FullName:     html.EscapeString(payload.FullName),
		Username:     html.EscapeString(payload.Username),
		PasswordHash: string(passwordHash),
		Role:         html.EscapeString(payload.Role),
		Status:       html.EscapeString(payload.Status),
		LoginHistory: []time.Time{},
	}

	if err := config.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user. Username might be taken."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "data": admin})
}

// AdminUpdateUserPayload represents the payload for updating an admin user
type AdminUpdateUserPayload struct {
	FullName string `json:"fullName" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password"` // Optional, only update if provided
	Role     string `json:"role" binding:"required"`
	Status   string `json:"status" binding:"required"`
}

// AdminUpdateUser updates an existing admin or staff user
func AdminUpdateUser(c *gin.Context) {
	id := c.Param("id")
	var admin models.Admin

	if err := config.DB.First(&admin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var payload AdminUpdateUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	admin.FullName = html.EscapeString(payload.FullName)
	admin.Username = html.EscapeString(payload.Username)
	admin.Role = html.EscapeString(payload.Role)
	admin.Status = html.EscapeString(payload.Status)

	if payload.Password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		admin.PasswordHash = string(passwordHash)
	}

	if err := config.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "data": admin})
}

// AdminDeleteUser deletes an admin user
func AdminDeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Admin{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// AdminGetAllRequests fetches all document requests, with optional status filter
func AdminGetAllRequests(c *gin.Context) {
	statusFilter := c.Query("status")
	var requests []models.DocumentRequest

	query := config.DB.Preload("User").Order("created_at desc")
	if statusFilter != "" && statusFilter != "All" {
		query = query.Where("status = ?", statusFilter)
	}

	if err := query.Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": requests})
}

// AdminCreateRequest payload for staff creating a request
type AdminCreateRequestPayload struct {
	FullName     string `json:"full_name" binding:"required"`
	ContactNumber string `json:"contact_number" binding:"required"`
	DocumentType string `json:"document_type" binding:"required"`
	Purpose      string `json:"purpose" binding:"required"`
}

// AdminCreateRequest allows a staff member to create a request for a resident
func AdminCreateRequest(c *gin.Context) {
	var payload AdminCreateRequestPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// 1. Find or create User by phone number
	var user models.User
	if err := config.DB.Where("phone_number = ?", payload.ContactNumber).First(&user).Error; err != nil {
		// User not found, create new
		user = models.User{
			PhoneNumber: payload.ContactNumber,
			FullName:    payload.FullName,
		}
		if err := config.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	} else {
		// Update name if different
		if user.FullName == "" || user.FullName != payload.FullName {
			user.FullName = payload.FullName
			config.DB.Save(&user)
		}
	}

	// 2. Create the DocumentRequest
	refNum := generateReferenceNumber()
	
	req := models.DocumentRequest{
		ReferenceNumber: refNum,
		UserID:          user.ID,
		DocumentType:    payload.DocumentType,
		Purpose:         payload.Purpose,
		Status:          "Pending",
	}

	if err := config.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document request"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":          "Request created successfully",
		"reference_number": req.ReferenceNumber,
	})
}
