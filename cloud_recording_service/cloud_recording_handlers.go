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

type AcquireResourceRequest struct {
	Cname string `json:"cname"`
	Uid   string `json:"uid"`
}

type StartRecordingRequest struct {
	Cname         string          `json:"cname"`
	Uid           string          `json:"uid"`
	ClientRequest ClientRequest   `json:"clientRequest"`
}

type ClientRequest struct {
	Token           string           `json:"token,omitempty"`
	StorageConfig   StorageConfig    `json:"storageConfig,omitempty"`
	RecordingConfig RecordingConfig  `json:"recordingConfig,omitempty"`
}

type StorageConfig struct {
	SecretKey  string `json:"secretKey"`
	Vendor     int    `json:"vendor"`
	Region     int    `json:"region"`
	Bucket     string `json:"bucket"`
	AccessKey  string `json:"accessKey"`
}

type RecordingConfig struct {
	ChannelType int `json:"channelType"`
}

// AcquireResource handles the acquire resource request.
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