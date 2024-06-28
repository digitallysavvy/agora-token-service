package cloud_recording_service

import (
	"encoding/json"
	"fmt"
)

// HandleAcquireResourceReq handles the acquire resource request.
// It constructs the URL, marshals the request, sends it to the Agora cloud recording API, and processes the response.
//
// Parameters:
//   - acquireReq: AcquireResourceRequest - The request payload for acquiring a resource.
//
// Returns:
//   - string: The resource ID acquired from the Agora cloud recording API.
//   - error: An error object if any issues occurred during the request process.
//
// Behavior:
//   - Marshals the acquireReq struct into JSON.
//   - Constructs the URL for the Agora cloud recording API request.
//   - Sends the Acquire request to the Agora cloud recording API using makeRequest.
//   - Reads and processes the response body to extract the resource ID if the request is successful.
//
// Notes:
//   - This function assumes the presence of s.baseURL, s.appID, s.customerID, and s.customerCertificate for constructing the API request.
func (s *CloudRecordingService) HandleAcquireResourceReq(acquireReq AcquireResourceRequest) (string, error) {
	// Construct the URL
	url := fmt.Sprintf("%s/%s/cloud_recording/acquire", s.baseURL, s.appID)

	body, err := s.makeRequest("POST", url, acquireReq)
	if err != nil {
		return "", err
	}

	// Parse the response body to extract the resource ID
	var response struct {
		ResourceId string `json:"resourceId"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	return response.ResourceId, nil
}
