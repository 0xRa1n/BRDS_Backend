package models

import (
	"time"
	"gorm.io/gorm"
)

// DocumentRequest represents a user's request for a barangay document
type DocumentRequest struct {
	gorm.Model
	ReferenceNumber string `json:"reference_number" gorm:"uniqueIndex;not null"`
	UserID          uint   `json:"user_id" gorm:"index;not null"`
	User            User   `json:"user" gorm:"foreignKey:UserID"`
	Address         string `json:"address" gorm:"not null;default:''"`
	DocumentType    string `json:"document_type" gorm:"not null"`
	Purpose         string `json:"purpose" gorm:"not null"`
	Status          string     `json:"status" gorm:"not null;default:'Pending'"`
	Remarks         string     `json:"remarks"`
	IdempotencyKey  string     `json:"idempotency_key" gorm:"uniqueIndex"`
	AppointmentDate *time.Time `json:"appointment_date"`
}

// SubmitRequestPayload represents the payload from the frontend for submitting a request
type SubmitRequestPayload struct {
	FullName       string `json:"full_name" binding:"required"`
	Address        string `json:"address" binding:"required"`
	DateOfBirth    string `json:"date_of_birth" binding:"required"`
	DocumentType   string `json:"document_type" binding:"required"`
	Purpose        string `json:"purpose" binding:"required"`
	IdempotencyKey string `json:"idempotency_key"` // Optional in body, but we'll extract from header or body
}
