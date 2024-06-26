package cloud_recording_service

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/AgoraIO-Community/agora-backend-service/token_service"
	"github.com/gin-gonic/gin"
)

// CloudRecordingService represents the cloud recording service.
// It holds the necessary configurations and dependencies for managing cloud recordings.
type CloudRecordingService struct {
	appID         string                      // The Agora app ID
	baseURL       string                      // The base URL for the Agora cloud recording API
	basicAuth     string                      // Middleware for handling requests
	tokenService  *token_service.TokenService // Token service for generating tokens
	storageConfig StorageConfig
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
//   - Initializes and returns a CloudRecordingService struct with the given configurations.
//
// Notes:
//   - Logs a fatal error and exits if any required environment variables are missing.
func NewCloudRecordingService(appID string, baseURL string, basicAuth string, tokenService *token_service.TokenService, storageConfig StorageConfig) *CloudRecordingService {

	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Return a new instance of the service
	return &CloudRecordingService{
		appID:         appID,
		baseURL:       baseURL,
		basicAuth:     basicAuth,
		tokenService:  tokenService,
		storageConfig: storageConfig,
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
	// group route
	api := r.Group("/cloud_recording")
	// routes
	api.POST("/startRecording", s.StartRecording)
	api.POST("/stopRecording", s.StopRecording)
	api.GET("/getStatus", s.GetStatus)
	// "update" group route
	updateAPI := api.Group("/update")
	updateAPI.POST("/subscriber-list", s.UpdateSubscriptionList)
	updateAPI.POST("/layout", s.UpdateLayout)
}

func (s *CloudRecordingService) StartRecording(c *gin.Context) {

	var clientReq ClientStartRecordingRequest
	if err := c.ShouldBindJSON(&clientReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate recording mode
	modes := []string{"individual", "mix", "web"}
	if !Contains(modes, clientReq.RecordingMode) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid recording mode."})
		return
	}

	// Generate a unique UID for this recording session
	uid := generateUID()

	// Acquire Resource
	acquireReq := AcquireResourceRequest{
		Cname:         clientReq.ChannelName,
		Uid:           uid,
		ClientRequest: make(map[string]interface{}), // Initialize as an empty map
	}
	resourceID, err := s.HandleAcquireResourceReq(acquireReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acquire resource: " + err.Error()})
		return
	}

	// Generate token for recording using token_service
	tokenRequest := token_service.TokenRequest{
		TokenType: "rtc",
		Channel:   clientReq.ChannelName,
		Uid:       uid,
	}
	token, err := s.tokenService.GenRtcToken(tokenRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Build the full StartRecordingRequest
	startReq := StartRecordingRequest{
		Cname: clientReq.ChannelName,
		Uid:   uid,
		ClientRequest: ClientRequest{
			Scene:               1,  // Assuming default scene, adjust as needed
			ResourceExpiredHour: 24, // Assuming 24 hours, adjust as needed
			StartParameter: StartParameter{
				Token:           token,
				StorageConfig:   s.storageConfig,
				RecordingConfig: clientReq.RecordingConfig,
			},
		},
	}

	// Step 3: Start Recording
	recordingID, err := s.HandleStartRecordingReq(startReq, resourceID, clientReq.RecordingMode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start recording: " + err.Error()})
		return
	}

	// Step 4: Return Resource ID and Recording ID
	c.JSON(http.StatusOK, gin.H{
		"UID":         uid,
		"resourceId":  resourceID,
		"recordingId": recordingID,
	})
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
