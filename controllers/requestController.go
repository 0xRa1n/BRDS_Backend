package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"backend/config"
	"backend/models"

	"github.com/gin-gonic/gin"
)

// generateReferenceNumber creates a unique reference like BR-YYYY-XXXX
func generateReferenceNumber() string {
	year := time.Now().Year()
	bytes := make([]byte, 2)
	rand.Read(bytes)
	randomHex := strings.ToUpper(hex.EncodeToString(bytes))
	return fmt.Sprintf("BR-%d-%s", year, randomHex)
}

// SubmitRequest handles document requests from authenticated users
func SubmitRequest(c *gin.Context) {
	phone, exists := c.Get("phone")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idempotencyKey := c.GetHeader("Idempotency-Key")

	var payload models.SubmitRequestPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Use header or payload for idempotency key
	if idempotencyKey == "" {
		idempotencyKey = payload.IdempotencyKey
	}

	// Check idempotency key first
	if idempotencyKey != "" {
		var existingReq models.DocumentRequest
		if err := config.DB.Where("idempotency_key = ?", idempotencyKey).First(&existingReq).Error; err == nil {
			// Idempotent hit: return the existing reference number
			c.JSON(http.StatusOK, gin.H{
				"message":          "Request already submitted successfully",
				"reference_number": existingReq.ReferenceNumber,
			})
			return
		}
	}

	// Find or Create user automatically
	var user models.User
	if err := config.DB.Where("phone_number = ?", phone).First(&user).Error; err != nil {
		// If user not found, we create them automatically with the phone from JWT
		user = models.User{
			PhoneNumber: phone.(string),
		}
		config.DB.Create(&user)
	}

	// Update user profile with submitted details
	user.FullName = payload.FullName
	user.Address = payload.Address

	parsedDate, err := time.Parse("2006-01-02", payload.DateOfBirth)
	if err == nil {
		user.DateOfBirth = parsedDate
	}
	// Save updated profile
	config.DB.Save(&user)

	// Create document request
	docReq := models.DocumentRequest{
		ReferenceNumber: generateReferenceNumber(),
		UserID:          user.ID,
		DocumentType:    payload.DocumentType,
		Purpose:         payload.Purpose,
		Status:          "Pending",
		IdempotencyKey:  idempotencyKey,
	}

	if err := config.DB.Create(&docReq).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Request submitted successfully",
		"reference_number": docReq.ReferenceNumber,
	})
}
