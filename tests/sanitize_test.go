package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func TestSanitizeMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.SanitizeMiddleware())
	router.POST("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	router.POST("/test/*path", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	tests := []struct {
		name           string
		method         string
		url            string
		body           string
		expectedStatus int
	}{
		{
			name:           "Clean Request",
			method:         "POST",
			url:            "/test",
			body:           `{"phone": "12345"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Clean Query",
			method:         "POST",
			url:            "/test?name=john",
			body:           `{"phone": "12345"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "SQLi in URL Query",
			method:         "POST",
			url:            "/test?name=john'%3B%20DROP%20TABLE%20users%3B--",
			body:           `{"phone": "12345"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "SQLi in JSON Body",
			method:         "POST",
			url:            "/test",
			body:           `{"phone": "12345", "query": "SELECT * FROM users"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "SQLi in URL Path",
			method:         "POST",
			url:            "/test/SELECT%20*%20FROM%20users",
			body:           `{"phone": "12345"}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.url, bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d for %s", tt.expectedStatus, w.Code, tt.name)
			}
		})
	}
}
