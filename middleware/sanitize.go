package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/gin-gonic/gin"
)

// SQLInjectionPattern is a regex to detect common SQL injection keywords and metacharacters.
// We make it case-insensitive and look for common attack vectors.
var SQLInjectionPattern = regexp.MustCompile(`(?i)(SELECT\s|UPDATE\s|INSERT\s|DELETE\s|DROP\s|ALTER\s|UNION\s|--|;)`)

// SanitizeMiddleware checks URL parameters and JSON body for SQL injection patterns
func SanitizeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Check URL query parameters
		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				if SQLInjectionPattern.MatchString(value) {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Malicious content detected in URL parameters"})
					return
				}
			}
		}

		// 2. Check URL path
		decodedPath, err := url.PathUnescape(c.Request.URL.Path)
		if err == nil && SQLInjectionPattern.MatchString(decodedPath) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Malicious content detected in URL path"})
			return
		}

		// 3. Check Payload (Body)
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
				return
			}

			// Restore the body so downstream handlers can read it
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// Check the raw body string
			if SQLInjectionPattern.MatchString(string(bodyBytes)) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Malicious content detected in payload"})
				return
			}
		}

		c.Next()
	}
}
