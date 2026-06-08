package controllers

import (
	"net/http"

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

	token, err := utils.GenerateAdminJWT(admin.Username, admin.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.AdminLoginResponse{
		Token: token,
		Role:  admin.Role,
	})
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
