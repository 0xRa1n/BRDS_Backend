package models

import (
	"time"

	"gorm.io/gorm"
)

// Admin represents an administrator or staff member
type Admin struct {
	gorm.Model
	UniqueID     string      `json:"uniqueId" gorm:"uniqueIndex"`
	FullName     string      `json:"fullName"`
	Username     string      `json:"username" gorm:"uniqueIndex;not null"`
	PasswordHash string      `json:"-" gorm:"not null"`
	Role         string      `json:"role" gorm:"not null;default:'Staff'"` // Admin or Staff
	Status       string      `json:"status" gorm:"not null;default:'Active'"`
	LoginHistory []time.Time `json:"loginHistory" gorm:"serializer:json"`
}

// AdminLoginRequest represents the payload for admin login
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminLoginResponse represents the payload returned upon successful admin login
type AdminLoginResponse struct {
	Token    string `json:"token"`
	Role     string `json:"role"`
	FullName string `json:"fullName"`
}
