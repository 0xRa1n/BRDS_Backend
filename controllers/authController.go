package controllers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"backend/config"
	"backend/models"
	"backend/services"
	"backend/utils"

	"log"

	"github.com/gin-gonic/gin"
)

// generateOTP creates a 6-digit OTP
func generateOTP() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%06d", n.Int64())
}

// SendOTP handles sending the OTP
func SendOTP(c *gin.Context) {
	// payload: { phone_number: "" }
	var req models.OTPSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Main logic
	if !config.AppCache.CheckRateLimit(req.PhoneNumber) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests. Please try again later."})
		return
	}

	code := generateOTP()
	config.AppCache.SaveOTP(req.PhoneNumber, code)

	err := services.SendSMS(req.PhoneNumber, code)
	if err != nil {
		// print the error
		log.Println("Send SMS: " + err.Error());
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP SMS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// VerifyOTP handles OTP validation and resident upsert
func VerifyOTP(c *gin.Context) {
	// payload: { phoneNumber: "", code: "", fullName: "", address: "", dateOfBirth: "" }
	var req models.OTPVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Main logic
	if !config.AppCache.CheckRateLimit(req.PhoneNumber + "_verify") {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests. Please try again later."})
		return
	}

	if !config.AppCache.VerifyOTP(req.PhoneNumber, req.Code) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	var user models.User
	res := config.DB.Where("phone_number = ?", req.PhoneNumber).First(&user)

	// If user doesn't exist, create it. Otherwise, update it.
	user.PhoneNumber = req.PhoneNumber

	if res.Error != nil {
		config.DB.Create(&user)
	} else {
		config.DB.Save(&user)
	}

	token, err := utils.GenerateJWT(user.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{Token: token})
}

// Logout handles JWT blocklisting
func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid Authorization header"})
		return
	}

	// Main logic
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.ExtractJWTClaims(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	jti, _ := claims["jti"].(string)
	expF, _ := claims["exp"].(float64)
	expiresAt := time.Unix(int64(expF), 0)

	config.AppCache.BlockJWT(jti, expiresAt)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
