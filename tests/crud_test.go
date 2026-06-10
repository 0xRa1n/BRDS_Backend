package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/config"
	"backend/middleware"
	"backend/models"
	"backend/routes"
	"backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.DocumentRequest{})
	config.DB = db
	config.InitCache()
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.SanitizeMiddleware())
	routes.RegisterAuthRoutes(router)
	routes.RegisterUserRoutes(router)
	routes.RegisterRequestRoutes(router)
	return router
}

func TestCRUDOperations(t *testing.T) {
	setupTestDB()
	router := setupTestRouter()

	testPhone := "09123456789"

	// Create user manually for testing since auth is bypassed by JWT
	config.DB.Create(&models.User{PhoneNumber: testPhone})

	token, _ := utils.GenerateJWT(testPhone)

	// 1. Test Update Profile
	t.Run("Update Profile", func(t *testing.T) {
		payload := map[string]string{
			"full_name":     "John Doe",
			"address":       "123 Street",
			"date_of_birth": "1990-01-01",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("PUT", "/api/v1/users/profile", bytes.NewBuffer(body))
		req.AddCookie(&http.Cookie{Name: "jwt_token", Value: token})
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}

		var user models.User
		config.DB.Where("phone_number = ?", testPhone).First(&user)
		if user.FullName != "John Doe" {
			t.Errorf("expected name John Doe, got %s", user.FullName)
		}
	})

	// 2. Test Get Profile
	t.Run("Get User Profile", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users/profile", nil)
		req.AddCookie(&http.Cookie{Name: "jwt_token", Value: token})

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	// 3. Test Submit Document Request
	t.Run("Submit Request", func(t *testing.T) {
		payload := map[string]string{
			"full_name":       "John Doe",
			"address":         "123 Street",
			"date_of_birth":   "1990-01-01",
			"document_type":   "Barangay Clearance",
			"purpose":         "Employment",
			"idempotency_key": "test_key_1",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/v1/requests", bytes.NewBuffer(body))
		req.AddCookie(&http.Cookie{Name: "jwt_token", Value: token})
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}

		var docReq models.DocumentRequest
		config.DB.Where("idempotency_key = ?", "test_key_1").First(&docReq)
		if docReq.DocumentType != "Barangay Clearance" {
			t.Errorf("expected Barangay Clearance, got %s", docReq.DocumentType)
		}
	})

	// 4. Test Get Document Requests
	t.Run("Get User Requests", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/requests", nil)
		req.AddCookie(&http.Cookie{Name: "jwt_token", Value: token})

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})
}
