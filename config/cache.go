package config

import (
	"sync"
	"time"
)

// PendingOTP holds the OTP and its expiration time
type PendingOTP struct {
	Code      string
	ExpiresAt time.Time
}

// RateLimitEntry holds attempts and window
type RateLimitEntry struct {
	Attempts  int
	ExpiresAt time.Time
}

// Cache handles in-memory storage for mock rate-limiting, OTPs, and JWT blocklist
type Cache struct {
	mu           sync.RWMutex
	otps         map[string]PendingOTP     // phone -> PendingOTP
	rateLimits   map[string]RateLimitEntry // phone -> RateLimitEntry
	jwtBlocklist map[string]time.Time      // jti -> expiration time
}

// Global cache instance
var AppCache *Cache

func InitCache() {
	AppCache = &Cache{
		otps:         make(map[string]PendingOTP),
		rateLimits:   make(map[string]RateLimitEntry),
		jwtBlocklist: make(map[string]time.Time),
	}
	
	// Start cleanup routine
	go AppCache.cleanupRoutine()
}

// CheckRateLimit returns true if allowed, false if rate limited
func (c *Cache) CheckRateLimit(phone string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.rateLimits[phone]
	now := time.Now()

	// If entry expired, reset it
	if exists && now.After(entry.ExpiresAt) {
		exists = false
	}

	if !exists {
		c.rateLimits[phone] = RateLimitEntry{
			Attempts:  1,
			ExpiresAt: now.Add(10 * time.Minute), // 10 min window
		}
		return true
	}

	if entry.Attempts >= 3 {
		return false
	}

	entry.Attempts++
	c.rateLimits[phone] = entry
	return true
}

// SaveOTP saves an OTP for a phone number
func (c *Cache) SaveOTP(phone, code string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.otps[phone] = PendingOTP{
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute), // 5 min TTL
	}
}

// VerifyOTP verifies the OTP. Returns true if valid, false otherwise.
func (c *Cache) VerifyOTP(phone, code string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	otp, exists := c.otps[phone]
	if !exists {
		return false
	}

	if time.Now().After(otp.ExpiresAt) {
		delete(c.otps, phone) // clean up expired
		return false
	}

	if otp.Code == code {
		// Invalidate immediately after successful use
		delete(c.otps, phone)
		return true
	}

	return false
}

// BlockJWT adds a JWT jti to the blocklist
func (c *Cache) BlockJWT(jti string, expiresAt time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.jwtBlocklist[jti] = expiresAt
}

// IsJWTBlocked checks if a JWT is blocklisted
func (c *Cache) IsJWTBlocked(jti string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	exp, exists := c.jwtBlocklist[jti]
	if !exists {
		return false
	}

	if time.Now().After(exp) {
		return false // Expiration passed, not blocked conceptually (though we'd clean it up)
	}

	return true
}

// cleanupRoutine periodically removes expired items
func (c *Cache) cleanupRoutine() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()

		for phone, otp := range c.otps {
			if now.After(otp.ExpiresAt) {
				delete(c.otps, phone)
			}
		}

		for phone, rl := range c.rateLimits {
			if now.After(rl.ExpiresAt) {
				delete(c.rateLimits, phone)
			}
		}

		for jti, exp := range c.jwtBlocklist {
			if now.After(exp) {
				delete(c.jwtBlocklist, jti)
			}
		}

		c.mu.Unlock()
	}
}
