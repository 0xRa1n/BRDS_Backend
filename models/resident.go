package models

import (
	"time"

	"gorm.io/gorm"
)

// Resident represents the resident entity in the database
type Resident struct {
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
	FullName    string `json:"full_name"`
	Address     string `json:"address"`
	DateOfBirth string `json:"date_of_birth"` // Can be parsed later
}

// AuthResponse represents the response when an OTP is verified successfully
type AuthResponse struct {
	Token string `json:"token"`
}
