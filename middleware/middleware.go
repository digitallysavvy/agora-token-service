package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	AllowOrigin string
}

func NewMiddleware(allowOrigin string) *Middleware {
	return &Middleware{AllowOrigin: allowOrigin}
}

func (m *Middleware) NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		// set headers
		c.Header("Cache-Control", "private, no-cache, no-store, must-revalidate")
		c.Header("Expires", "-1")
		c.Header("Pragma", "no-cache")
	}
}

// Add CORSMiddleware to handle CORS requests and set the necessary headers
func (m *Middleware) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if !m.isOriginAllowed(origin) {
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

func (m *Middleware) isOriginAllowed(origin string) bool {
	if m.AllowOrigin == "*" {
		return true
	}

	allowedOrigins := strings.Split(m.AllowOrigin, ",")
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}

	return false
}

func (m *Middleware) TimestampMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Proceed to the next middleware/handler
		c.Next()

		// Add the current timestamp to the response header
		timestamp := time.Now().Format(time.RFC3339)
		c.Writer.Header().Set("X-Timestamp", timestamp)
	}
}

// TimestampMiddleware adds a timestamp to each response
// func (m *Middleware) TimestampMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Create a custom ResponseWriter to capture the response body
// 		writer := &ResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
// 		c.Writer = writer

// 		// Proceed to the next middleware/handler
// 		c.Next()

// 		// Add the current timestamp to the response header
// 		timestamp := time.Now().Format(time.RFC3339)
// 		c.Writer.Header().Set("X-Timestamp", timestamp)

// 		// Add the current timestamp to the response body
// 		var response map[string]interface{}
// 		if err := json.Unmarshal(writer.body.Bytes(), &response); err != nil {
// 			response = make(map[string]interface{})
// 		}
// 		response["timestamp"] = timestamp

// 		// Write the modified response body back to the client
// 		if err := json.NewEncoder(c.Writer).Encode(response); err != nil {
// 			c.String(500, "Failed to encode response")
// 		}
// 	}
// }

// // ResponseWriter is a custom writer to capture the response body
// type ResponseWriter struct {
// 	gin.ResponseWriter
// 	body *bytes.Buffer
// }

// func (w ResponseWriter) Write(b []byte) (int, error) {
// 	w.body.Write(b)
// 	return w.ResponseWriter.Write(b)
// }
