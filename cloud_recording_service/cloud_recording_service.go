package cloud_recording_service

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
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

	var clientStartReq ClientStartRecordingRequest
	if err := c.ShouldBindJSON(&clientStartReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sceneMode := 0
	sceneModes := map[string]int{"realtime": 0, "web": 1, "postponed": 2}
	if clientStartReq.SceneMode != nil {
		if mode, ok := sceneModes[*clientStartReq.SceneMode]; ok {
			sceneMode = mode
		}
	}

	// Default RecordingMode to "composite" if nil
	recordingMode := "mix"
	if clientStartReq.RecordingMode != nil {
		recordingMode = *clientStartReq.RecordingMode
	}

	// Validate recording mode against a set list
	validRecordingModes := []string{"individual", "mix", "web"}
	if !Contains(validRecordingModes, recordingMode) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recording mode."})
		return
	}

	// Generate a unique UID for this recording session
	uid := generateUID()

	// Generate token for recording using token_service
	tokenRequest := token_service.TokenRequest{
		TokenType: "rtc",
		Channel:   clientStartReq.ChannelName,
		Uid:       uid,
	}
	token, err := s.tokenService.GenRtcToken(tokenRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add dynamic directory structure ChannelName/YYYYMMDD/HHMMSS
	currentTimeUTC := time.Now().UTC()
	dateStr := currentTimeUTC.Format("20060102")
	hrsMinSecStr := currentTimeUTC.Format("150405")
	s.storageConfig.FileNamePrefix = &[]string{strings.ReplaceAll(clientStartReq.ChannelName, "-", ""), dateStr, hrsMinSecStr}

	// Check if RecordingConfig is nil, if so, create a default one
	if clientStartReq.RecordingConfig == nil {
		// recording config defaults
		channelType := 0
		streamTypes := 2
		videoStreamType := 0
		maxIdleTime := 120
		subscribeUidGroup := 0
		streamMode := "standard" // Default to "standard", ensure to modify based on your actual logic
		subscribeAudioUids := []string{"#allstream#"}
		subscribeVideoUids := []string{"#allstream#"}
		// create default recording config
		clientStartReq.RecordingConfig = &RecordingConfig{
			ChannelType:        channelType,
			StreamTypes:        &streamTypes,
			VideoStreamType:    &videoStreamType,
			StreamMode:         &streamMode,
			MaxIdleTime:        &maxIdleTime,
			SubscribeAudioUids: &subscribeAudioUids,
			SubscribeVideoUids: &subscribeVideoUids,
			SubscribeUidGroup:  &subscribeUidGroup,
		}
	}

	// Assemble recording client request
	recClientReq := AquireClientRequest{
		Scene:               sceneMode,
		ResourceExpiredHour: 24, // Assuming 24 hours, adjust as needed
		StartParameter: ClientRequest{
			Token:           token,
			StorageConfig:   s.storageConfig,
			RecordingConfig: *clientStartReq.RecordingConfig,
		},
		ExcludeResourceIds: clientStartReq.ExcludeResourceIds,
	}

	// Acquire Resource
	acquireReq := AcquireResourceRequest{
		Cname:         clientStartReq.ChannelName,
		Uid:           uid,
		ClientRequest: &recClientReq, // Initialize as an empty map
	}
	resourceID, err := s.HandleAcquireResourceReq(acquireReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acquire resource: " + err.Error()})
		return
	}

	log.Println("resourceID:", (resourceID))

	// Build the full StartRecordingRequest
	startReq := StartRecordingRequest{
		Cname: clientStartReq.ChannelName,
		Uid:   uid,
		ClientRequest: ClientRequest{
			Token:           token,
			StorageConfig:   s.storageConfig,
			RecordingConfig: *clientStartReq.RecordingConfig,
		},
	}

	// configJSON, err := json.MarshalIndent(startReq, "", "  ")
	// if err != nil {
	// 	log.Fatalf("Error marshalling default config: %v", err)
	// }
	// log.Println("startReq:")
	// log.Println(string(configJSON))

	// Start Recording
	recordingID, err := s.HandleStartRecordingReq(startReq, resourceID, recordingMode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start recording: " + err.Error()})
		return
	}

	// Return Resource ID and Recording ID
	c.JSON(http.StatusOK, gin.H{
		"UID":         uid,
		"resourceId":  resourceID,
		"recordingId": recordingID,
		"timestamp":   time.Now().UTC(),
	})
}

// StopRecording
func (s *CloudRecordingService) StopRecording(c *gin.Context) {
	var req = c.Request
	var respWriter = c.Writer
	var clientStopReq ClientStopRecordingRequest
	err := json.NewDecoder(req.Body).Decode(&clientStopReq)
	if err != nil {
		// invalid request
		http.Error(respWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// build the stop request from user request
	stopReq := StopRecordingRequest{
		Cname: clientStopReq.Cname,
		Uid:   clientStopReq.Uid,
		ClientRequest: StopClientRequest{
			AsyncStop: clientStopReq.AsyncStop,
		},
	}
	recordingMode := "mix"
	if clientStopReq.RecordingMode != nil {
		recordingMode = *clientStopReq.RecordingMode
	}
	// Send Stop Recording Request to Agora
	response, err := s.HandleStopRecording(stopReq, clientStopReq.ResourceId, clientStopReq.RecordingId, recordingMode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the wrapped Agora response
	c.Data(http.StatusOK, "application/json", response)
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
