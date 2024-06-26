package cloud_recording_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

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
func (s *CloudRecordingService) makeRequest(method, url string, body interface{}) ([]byte, error) {
	// Marshal the request body into JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.basicAuth)

	// Send the request with a 10-second timeout
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Check for a successful status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}
