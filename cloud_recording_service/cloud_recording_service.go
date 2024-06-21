package cloud_recording_service

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/AgoraIO-Community/agora-backend-service/middleware"
	"github.com/AgoraIO-Community/agora-backend-service/token_service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// CloudRecordingService represents the cloud recording service.
// It holds the necessary configurations and dependencies for managing cloud recordings.
type CloudRecordingService struct {
	appID               string                      // The Agora app ID
	appCertificate      string                      // The Agora app certificate
	customerID          string                      // The customer ID for authentication
	customerCertificate string                      // The customer certificate for authentication
	allowOrigin         string                      // The allowed origin for CORS
	middleware          *middleware.Middleware      // Middleware for handling requests
	tokenService        *token_service.TokenService // Token service for generating tokens
	baseURL             string                      // The base URL for the Agora cloud recording API
	storageConfig       StorageConfig
}

// NewCloudRecordingService returns a CloudRecordingService pointer with all configurations set.
// It loads environment variables, validates their presence, and initializes the CloudRecordingService struct.
//
// Parameters:
//   - tokenService: *token_service.TokenService - The token service for generating tokens.
//
// Returns:
//   - *CloudRecordingService: The initialized CloudRecordingService struct.
//
// Behavior:
//   - Loads environment variables from the .env file.
//   - Retrieves and validates necessary environment variables.
//   - Initializes and returns a CloudRecordingService struct with the loaded configurations.
//
// Notes:
//   - Logs a fatal error and exits if any required environment variables are missing.
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

	secretKey, secretKeyExists := os.LookupEnv("STORAGE_SECRET_KEY")
	vendorStr, vendorExists := os.LookupEnv("STORAGE_VENDOR")
	regionStr, regionExists := os.LookupEnv("STORAGE_REGION")
	bucket, bucketExists := os.LookupEnv("STORAGE_BUCKET")
	accessKey, accessKeyExists := os.LookupEnv("STORAGE_ACCESS_KEY")

	if !appIDExists || !appCertExists || !customerIDExists || !customerCertExists || !baseURLExists ||
	!secretKeyExists || !vendorExists || !regionExists || !bucketExists || !accessKeyExists {
		log.Fatal("FATAL ERROR: ENV not properly configured, check .env file for all required variables")
	}

	vendor, err := strconv.Atoi(vendorStr)
	if err != nil {
		log.Fatal("FATAL ERROR: Invalid STORAGE_VENDOR value")
	}
	region, err := strconv.Atoi(regionStr)
	if err != nil {
		log.Fatal("FATAL ERROR: Invalid STORAGE_REGION value")
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
		storageConfig: StorageConfig{
			SecretKey:  secretKey,
			Vendor:     vendor,
			Region:     region,
			Bucket:     bucket,
			AccessKey:  accessKey,
		},
	}
}

// RegisterRoutes registers the routes for the CloudRecordingService.
// It sets up the API endpoints and applies necessary middleware for request handling.
//
// Parameters:
//   - r: *gin.Engine - The Gin engine instance to register the routes with.
//
// Behavior:
//   - Creates an API group for cloud recording routes.
//   - Applies middleware for NoCache and CORS.
//   - Registers routes for ping, acquireResource, startRecording, stopRecording, getStatus, update subscriber list, and update layout.
//
// Notes:
//   - This function organizes the API routes and ensures that requests are handled with appropriate middleware.
func (s *CloudRecordingService) RegisterRoutes(r *gin.Engine) {
	// set group route
	api := r.Group("/cloud_recording")
	// use middleware headers
	api.Use(s.middleware.NoCache())
	api.Use(s.middleware.CORSMiddleware())
	// routes to functions
	// api.POST("/acquireResource", s.AcquireResource)
	api.POST("/startRecording", s.StartRecording)
	api.POST("/stopRecording", s.StopRecording)
	api.GET("/getStatus", s.GetStatus)
	// set "update" group route
	updateAPI := api.Group("/update")
	updateAPI.POST("/subscriber-list", s.UpdateSubscriptionList)
	updateAPI.POST("/layout", s.UpdateLayout)
}

func (s *CloudRecordingService) AcquireResource(c *gin.Context) {
	var req = c.Request
	var respWriter = c.Writer
	var aquaireReq AcquireResourceRequest
	err := json.NewDecoder(req.Body).Decode(&aquaireReq)
	if err != nil {
		// invalid request
		http.Error(respWriter, err.Error(), http.StatusBadRequest)
		return
	}
	s.HandleAcquireResource(aquaireReq, c.Writer)
}


func (s *CloudRecordingService) StartRecording(c *gin.Context) {
	var req = c.Request
	var respWriter = c.Writer
	var startReq StartRecordingRequest
	err := json.NewDecoder(req.Body).Decode(&startReq)
	if err != nil {
		// invalid request
		http.Error(respWriter, err.Error(), http.StatusBadRequest)
		return
	}
	s.HandleStartRecording(startReq, respWriter)
}

// StopRecording
func (s *CloudRecordingService) StopRecording(c *gin.Context) {
	var req = c.Request
	var respWriter = c.Writer
	var stopReq StartRecordingRequest
	err := json.NewDecoder(req.Body).Decode(&stopReq)
	if err != nil {
		// invalid request
		http.Error(respWriter, err.Error(), http.StatusBadRequest)
		return
	}
	s.HandleStopRecording(stopReq, respWriter)
}

// GetStatus
func (s *CloudRecordingService) GetStatus(c *gin.Context) {
	s.HandleGetStatus(c.Writer, c.Request)
}

// UpdateSubscriptionList
func (s *CloudRecordingService) UpdateSubscriptionList(c *gin.Context) {
	var req = c.Request
	var respWriter = c.Writer
	var updateReq StartRecordingRequest
	err := json.NewDecoder(req.Body).Decode(&updateReq)
	if err != nil {
		// invalid request
		http.Error(respWriter, err.Error(), http.StatusBadRequest)
		return
	}
	s.HandleUpdateSubscriptionList(updateReq, respWriter)
}

// UpdateLayout
func (s *CloudRecordingService) UpdateLayout(c *gin.Context) {
	var req = c.Request
	var respWriter = c.Writer
	var updateReq StartRecordingRequest
	err := json.NewDecoder(req.Body).Decode(&updateReq)
	if err != nil {
		// invalid request
		http.Error(respWriter, err.Error(), http.StatusBadRequest)
		return
	}
	s.HandleUpdateLayout(updateReq, respWriter)
}