package cloud_recording_service

import (
	"log"
	"os"

	"github.com/AgoraIO-Community/agora-backend-service/middleware"
	"github.com/AgoraIO-Community/agora-backend-service/token_service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// CloudRecordingService represents the cloud recording service.
type CloudRecordingService struct {
	appID               string
	appCertificate      string
	customerID          string
	customerCertificate string
	allowOrigin         string
	middleware          *middleware.Middleware
	tokenService        *token_service.TokenService
	baseURL            string
}

// NewCloudRecordingService returns a CloudRecordingService pointer with all configurations set.
func NewCloudRecordingService(tokenService *token_service.TokenService) *CloudRecordingService {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	appIDEnv, appIDExists := os.LookupEnv("APP_ID")
	appCertEnv, appCertExists := os.LookupEnv("APP_CERTIFICATE")
	customerIDEnv, customerIDExists := os.LookupEnv("CUSTOMER_ID")
	customerCertEnv, customerCertExists := os.LookupEnv("CUSTOMER_CERTIFICATE")
	corsAllowOrigin, _ := os.LookupEnv("CORS_ALLOW_ORIGIN")
	baseURLEnv, baseURLExists := os.LookupEnv("BASE_URL")

	if !appIDExists || !appCertExists || !customerIDExists || !customerCertExists || !baseURLExists {
		log.Fatal("FATAL ERROR: ENV not properly configured, check .env file for BASE_URL, APP_ID, APP_CERTIFICATE, CUSTOMER_ID, and CUSTOMER_CERTIFICATE")
	}

	return &CloudRecordingService{
		appID:               appIDEnv,
		appCertificate:      appCertEnv,
		customerID:          customerIDEnv,
		customerCertificate: customerCertEnv,
		allowOrigin:         corsAllowOrigin,
		middleware:          middleware.NewMiddleware(corsAllowOrigin),
		tokenService:        tokenService,
		baseURL:             baseURLEnv,
	}
}

// RegisterRoutes registers the routes for the CloudRecordingService.
func (s *CloudRecordingService) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/cloud_recording")
	api.Use(s.middleware.NoCache())
	api.Use(s.middleware.CORSMiddleware())
	api.GET("/ping", s.Ping)
	api.POST("/acquireResource", s.AcquireResource)
	api.POST("/startRecording", s.StartRecording)
	api.POST("/stopRecording", s.StopRecording)
	api.GET("/getStatus", s.GetStatus)

	updateAPI := api.Group("/update")
	updateAPI.POST("/subscriber-list", s.UpdateSubscriptionList)
	updateAPI.POST("/layout", s.UpdateLayout)
}

// Ping is a simple handler for the /ping route.
func (s *CloudRecordingService) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}