package token_service

import (
	"log"
	"net/http"
	"os"

	"github.com/AgoraIO-Community/agora-backend-service/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// TokenService represents the main application token service.
// It holds the necessary configurations and dependencies for managing tokens.
type TokenService struct {
	Server         *http.Server           // The HTTP server for the application
	Sigint         chan os.Signal         // Channel to handle OS signals, such as Ctrl+C
	appID          string                 // The Agora app ID
	appCertificate string                 // The Agora app certificate
	allowOrigin    string                 // The allowed origin for CORS
	middleware     *middleware.Middleware // Middleware for handling requests
}

// NewTokenService returns a TokenService pointer with all configurations set.
// It loads environment variables, validates their presence, and initializes the TokenService struct.
//
// Returns:
//   - *TokenService: The initialized TokenService struct.
//
// Behavior:
//   - Loads environment variables from the .env file.
//   - Retrieves and validates necessary environment variables.
//   - Initializes and returns a TokenService struct with the loaded configurations.
//
// Notes:
//   - Logs a fatal error and exits if any required environment variables are missing.
func NewTokenService() *TokenService {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	appIDEnv, appIDExists := os.LookupEnv("APP_ID")
	appCertEnv, appCertExists := os.LookupEnv("APP_CERTIFICATE")
	corsAllowOrigin, _ := os.LookupEnv("CORS_ALLOW_ORIGIN")

	if !appIDExists || !appCertExists || len(appIDEnv) == 0 || len(appCertEnv) == 0 {
		log.Fatal("FATAL ERROR: ENV not properly configured, check .env file or APP_ID and APP_CERTIFICATE")
	}

	return &TokenService{
		appID:          appIDEnv,
		appCertificate: appCertEnv,
		allowOrigin:    corsAllowOrigin,
		middleware:     middleware.NewMiddleware(corsAllowOrigin),
	}
}

// RegisterRoutes registers the routes for the TokenService.
// It sets up the API endpoints and applies necessary middleware for request handling.
//
// Parameters:
//   - r: *gin.Engine - The Gin engine instance to register the routes with.
//
// Behavior:
//   - Creates an API group for token routes.
//   - Applies middleware for NoCache and CORS.
//   - Registers routes for ping and getNew token.
//
// Notes:
//   - This function organizes the API routes and ensures that requests are handled with appropriate middleware.
func (s *TokenService) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/token")
	api.Use(s.middleware.NoCache())
	api.Use(s.middleware.CORSMiddleware())
	api.GET("/ping", s.Ping)
	api.POST("/getNew", s.GetToken)
}

// Ping is a simple handler for the /ping route.
// It responds with a "pong" message to indicate that the service is running.
//
// Parameters:
//   - c: *gin.Context - The Gin context representing the HTTP request and response.
//
// Behavior:
//   - Sends a JSON response with a "pong" message.
//
// Notes:
//   - This function is useful for health checks and ensuring that the service is up and running.
func (s *TokenService) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}