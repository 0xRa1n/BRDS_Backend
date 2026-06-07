package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents the user entity in the database
type User struct {
	gorm.Model
	PhoneNumber string    `json:"phone_number" gorm:"uniqueIndex;not null"`
	FullName    string    `json:"full_name"`
	Address     string    `json:"address"`
	DateOfBirth time.Time `json:"date_of_birth" gorm:"type:date"`
}

// OTPSendRequest represents the payload for sending an OTP
type OTPSendRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

// OTPVerifyRequest represents the payload for verifying an OTP
type OTPVerifyRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Code        string `json:"code" binding:"required,len=6"`
}

// AuthResponse represents the response when an OTP is verified successfully
type AuthResponse struct {
	Token string `json:"token"`
}

// ProfileUpdateRequest represents the payload for updating resident details
type ProfileUpdateRequest struct {
	FullName    string `json:"full_name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	DateOfBirth string `json:"date_of_birth"`
}
