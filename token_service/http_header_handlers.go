package token_service

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *TokenService) nocache() gin.HandlerFunc {
	return func(c *gin.Context) {
		// set headers
		c.Header("Cache-Control", "private, no-cache, no-store, must-revalidate")
		c.Header("Expires", "-1")
		c.Header("Pragma", "no-cache")
	}
}

// Add CORSMiddleware to handle CORS requests and set the necessary headers
func (s *TokenService) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if !s.isOriginAllowed(origin) {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Origin not allowed",
			})
			c.Abort()
			return
		}
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func (s *TokenService) isOriginAllowed(origin string) bool {
	if s.allowOrigin == "*" {
		return true
	}

	allowedOrigins := strings.Split(s.allowOrigin, ",")
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}

	return false
}
