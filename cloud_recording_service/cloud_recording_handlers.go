package cloud_recording_service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AgoraIO-Community/agora-backend-service/token_service"
	"github.com/gin-gonic/gin"
)

// AcquireResourceRequest represents the JSON payload structure for acquiring a cloud recording resource.
// It contains the channel name and UID necessary for resource acquisition.
type AcquireResourceRequest struct {
	Cname string `json:"cname"` // The channel name for the cloud recording
	Uid   string `json:"uid"`   // The UID for the cloud recording session
}

// StartRecordingRequest represents the JSON payload structure for starting a cloud recording.
// It includes the channel name, UID, and the client request configuration.
type StartRecordingRequest struct {
	Cname         string        `json:"cname"`         // The channel name for the cloud recording
	Uid           string        `json:"uid"`           // The UID for the cloud recording session
	ClientRequest ClientRequest `json:"clientRequest"` // The client request configuration for the cloud recording
}

// ClientRequest represents the client request configuration for starting or updating a cloud recording.
// It includes the token, storage configuration, and recording configuration.
type ClientRequest struct {
	Token           string          `json:"token,omitempty"`          // The token for the cloud recording session
	StorageConfig   StorageConfig   `json:"storageConfig,omitempty"`  // The storage configuration for the cloud recording
	RecordingConfig RecordingConfig `json:"recordingConfig,omitempty"`// The recording configuration for the cloud recording
}

// StorageConfig represents the storage configuration for cloud recording.
// It includes the secret key, vendor, region, bucket, and access key for storage.
type StorageConfig struct {
	SecretKey  string `json:"secretKey"` // The secret key for storage authentication
	Vendor     int    `json:"vendor"`    // The storage vendor identifier
	Region     int    `json:"region"`    // The storage region identifier
	Bucket     string `json:"bucket"`    // The storage bucket name
	AccessKey  string `json:"accessKey"` // The access key for storage authentication
}

// RecordingConfig represents the recording configuration for cloud recording.
// It includes the channel type for the recording session.
type RecordingConfig struct {
	ChannelType int `json:"channelType"` // The channel type for the cloud recording
}

// AcquireResource handles the acquire resource request.
// It validates the request, constructs the URL, and sends the request to the Agora cloud recording API.
//
// Parameters:
//   - c: *gin.Context - The Gin context representing the HTTP request and response.
//
// Behavior:
//   - Parses the request body into an AcquireResourceRequest struct.
//   - Validates the request fields.
//   - Constructs the URL and authentication header for the API request.
//   - Sends the request to the Agora cloud recording API and returns the response.
//
// Notes:
//   - This function assumes the presence of s.baseURL, s.appID, s.customerID, and s.customerCertificate for constructing the API request.
func (s *CloudRecordingService) AcquireResource(c *gin.Context) {
	var req AcquireResourceRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	url := fmt.Sprintf("%s/%s/cloud_recording/acquire", s.baseURL, s.appID)
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.customerID, s.customerCertificate)))

	body, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal request"})
		return
	}

	resp, err := makeRequest("POST", url, auth, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", resp)
}

// StartRecording handles the start recording request.
// It validates the request, generates a token, constructs the URL, and sends the request to the Agora cloud recording API.
//
// Parameters:
//   - c: *gin.Context - The Gin context representing the HTTP request and response.
//
// Behavior:
//   - Parses the request body into a StartRecordingRequest struct.
//   - Validates the request fields.
//   - Generates a token using the token_service.
//   - Constructs the URL and authentication header for the API request.
//   - Sends the request to the Agora cloud recording API and returns the response.
//
// Notes:
//   - This function assumes the presence of s.baseURL, s.appID, s.customerID, s.customerCertificate, and s.tokenService for constructing the API request and generating the token.
func (s *CloudRecordingService) StartRecording(c *gin.Context) {
	var req StartRecordingRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resourceID := c.Param("resourceid")
	mode := c.Param("mode")

	// Generate token for recording using token_service
	tokenRequest := token_service.TokenRequest{
		TokenType: "rtc",
		Channel:   req.Cname,
		Uid:       req.Uid,
	}
	token, err := s.tokenService.GenRtcToken(tokenRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.ClientRequest.Token = token

	url := fmt.Sprintf("%s/%s/cloud_recording/resourceid/%s/mode/%s/start", s.baseURL, s.appID, resourceID, mode)
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.customerID, s.customerCertificate)))

	body, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal request"})
		return
	}

	resp, err := makeRequest("POST", url, auth, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", resp)
}

// StopRecording handles the stop recording request.
// It validates the request, constructs the URL, and sends the request to the Agora cloud recording API.
//
// Parameters:
//   - c: *gin.Context - The Gin context representing the HTTP request and response.
//
// Behavior:
//   - Parses the request body into a StartRecordingRequest struct.
//   - Validates the request fields.
//   - Constructs the URL and authentication header for the API request.
//   - Sends the request to the Agora cloud recording API and returns the response.
//
// Notes:
//   - This function assumes the presence of s.baseURL, s.appID, s.customerID, and s.customerCertificate for constructing the API request.
func (s *CloudRecordingService) StopRecording(c *gin.Context) {
	var req StartRecordingRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resourceID := c.Param("resourceid")
	sid := c.Param("sid")
	mode := c.Param("mode")

	url := fmt.Sprintf("%s/%s/cloud_recording/resourceid/%s/sid/%s/mode/%s/stop", s.baseURL, s.appID, resourceID, sid, mode)
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.customerID, s.customerCertificate)))

	body, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal request"})
		return
	}

	resp, err := makeRequest("POST", url, auth, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", resp)
}

// GetStatus handles the get status request.
// It constructs the URL and sends the request to the Agora cloud recording API.
//
// Parameters:
//   - c: *gin.Context - The Gin context representing the HTTP request and response.
//
// Behavior:
//   - Retrieves the resource ID, SID, and mode from the URL parameters.
//   - Constructs the URL and authentication header for the API request.
//   - Sends the request to the Agora cloud recording API and returns the response.
//
// Notes:
//   - This function assumes the presence of s.baseURL, s.appID, s.customerID, and s.customerCertificate for constructing the API request.
func (s *CloudRecordingService) GetStatus(c *gin.Context) {
	resourceID := c.Param("resourceid")
	sid := c.Param("sid")
	mode := c.Param("mode")

	url := fmt.Sprintf("%s/%s/cloud_recording/resourceid/%s/sid/%s/mode/%s/query", s.baseURL, s.appID, resourceID, sid, mode)
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.customerID, s.customerCertificate)))

	resp, err := makeRequest("GET", url, auth, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", resp)
}

// UpdateSubscriptionList handles the update subscription list request.
// It validates the request, constructs the URL, and sends the request to the Agora cloud recording API.
//
// Parameters:
//   - c: *gin.Context - The Gin context representing the HTTP request and response.
//
// Behavior:
//   - Parses the request body into a StartRecordingRequest struct.
//   - Validates the request fields.
//   - Constructs the URL and authentication header for the API request.
//   - Sends the request to the Agora cloud recording API and returns the response.
//
// Notes:
//   - This function assumes the presence of s.baseURL, s.appID, s.customerID, and s.customerCertificate for constructing the API request.
func (s *CloudRecordingService) UpdateSubscriptionList(c *gin.Context) {
	var req StartRecordingRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resourceID := c.Param("resourceid")
	sid := c.Param("sid")

	url := fmt.Sprintf("%s/%s/cloud_recording/resourceid/%s/sid/%s/update", s.baseURL, s.appID, resourceID, sid)
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.customerID, s.customerCertificate)))

	body, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal request"})
		return
	}

	resp, err := makeRequest("POST", url, auth, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", resp)
}

// UpdateLayout handles the update video layout request.
// It validates the request, constructs the URL, and sends the request to the Agora cloud recording API.
//
// Parameters:
//   - c: *gin.Context - The Gin context representing the HTTP request and response.
//
// Behavior:
//   - Parses the request body into a StartRecordingRequest struct.
//   - Validates the request fields.
//   - Constructs the URL and authentication header for the API request.
//   - Sends the request to the Agora cloud recording API and returns the response.
//
// Notes:
//   - This function assumes the presence of s.baseURL, s.appID, s.customerID, and s.customerCertificate for constructing the API request.
func (s *CloudRecordingService) UpdateLayout(c *gin.Context) {
	var req StartRecordingRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resourceID := c.Param("resourceid")
	sid := c.Param("sid")

	url := fmt.Sprintf("%s/%s/cloud_recording/resourceid/%s/sid/%s/updateLayout", s.baseURL, s.appID, resourceID, sid)
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.customerID, s.customerCertificate)))

	body, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal request"})
		return
	}

	resp, err := makeRequest("POST", url, auth, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", resp)
}

// makeRequest is a helper function to make HTTP requests with basic authentication.
//
// Parameters:
//   - method: string - The HTTP method to use for the request (e.g., "GET", "POST").
//   - url: string - The URL to send the request to.
//   - auth: string - The base64-encoded authorization header value.
//   - body: []byte - The request body to send (can be nil for GET requests).
//
// Returns:
//   - []byte: The response body from the server.
//   - error: An error if there are any issues during the request.
//
// Behavior:
//   - Creates a new HTTP request with the specified method, URL, and body.
//   - Sets the Authorization and Content-Type headers.
//   - Sends the request using an HTTP client.
//   - Reads and returns the response body, or an error if the request fails.
func makeRequest(method, url, auth string, body []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("error response from server: %s", respBody)
	}

	return respBody, nil
}