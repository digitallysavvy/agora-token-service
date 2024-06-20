package token_service

import (
	"log"
	"net/http"
	"os"

	"github.com/AgoraIO-Community/agora-backend-service/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// TokenService represents the main application TokenService.
type TokenService struct {
	// Server is the HTTP server for the application.
	Server *http.Server

	// Sigint is a channel to handle OS signals, such as Ctrl+C.
	Sigint chan os.Signal

	// appID is the identifier for the application.
	appID string

	// appCertificate is the certificate used by the application.
	appCertificate string

	// allowOrigin specifies the allowed origin for Cross-Origin Resource Sharing (CORS).
	allowOrigin string

	// middleware holds the CORS and NoCache configurations for the application
	middleware *middleware.Middleware
}

// NewTokenService returns a TokenService pointer with all configurations set.
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
func (s *TokenService) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/token")
	api.Use(s.middleware.NoCache())
	api.Use(s.middleware.CORSMiddleware())
	api.GET("/ping", s.Ping)
	api.POST("/getNew", s.GetToken)
}

// Ping is a simple handler for the /ping route.
func (s *TokenService) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}