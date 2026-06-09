package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

var DB *gorm.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		log.Fatal("Missing required database environment variables")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Manila",
		host, user, password, dbname, port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto Migrate the models
	err = DB.AutoMigrate(&models.User{}, &models.DocumentRequest{}, &models.Admin{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}

	seedAdmin()

	// Setup Connection Pooling
	sqlDB, err := DB.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	log.Println("Database connection established and models migrated.")
}

func seedAdmin() {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash default admin password: %v", err)
	}

	var adminCount int64
	DB.Model(&models.Admin{}).Where("username = ?", "admin").Count(&adminCount)
	if adminCount == 0 {
		adminAcc := models.Admin{
			UniqueID:     "ADM-01",
			FullName:     "System Administrator",
			Username:     "admin",
			PasswordHash: string(passwordHash),
			Role:         "Admin",
			LoginHistory: []time.Time{},
		}
		
		if err := DB.Create(&adminAcc).Error; err != nil {
			log.Fatalf("Failed to seed default admin: %v", err)
		}
		log.Println("Default admin account seeded.")
	}

	var staffCount int64
	DB.Model(&models.Admin{}).Where("username = ?", "staff").Count(&staffCount)
	if staffCount == 0 {
		staffAcc := models.Admin{
			UniqueID:     "DEF-STAFF-01",
			FullName:     "Default Staff",
			Username:     "staff",
			PasswordHash: string(passwordHash),
			Role:         "Staff",
			LoginHistory: []time.Time{},
		}
		
		if err := DB.Create(&staffAcc).Error; err != nil {
			log.Fatalf("Failed to seed default staff: %v", err)
		}
		log.Println("Default staff account seeded.")
	} else {
		// Update existing staff account that might have null fields from previous schema
		DB.Model(&models.Admin{}).Where("username = ?", "staff").Updates(map[string]interface{}{
			"unique_id": "DEF-STAFF-01",
			"full_name": "Default Staff",
			"role":      "Staff",
			"status":    "Active",
		})
	}
}
