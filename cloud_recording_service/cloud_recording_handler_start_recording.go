package cloud_recording_service

import (
	"encoding/json"
	"fmt"
)

// HandleStartRecordingReq processes the request to start a cloud recording session in Agora's cloud recording service.
// It validates the request parameters, constructs the request URL, and use makeRequest to send data to Agora Endpoint
//
// Parameters:
//   - startReq: StartRecordingRequest - The request payload for starting a recording.
//   - resourceId: string - The resource ID acquired previously.
//   - modeType: string - The recording mode type.
//
// Returns:
//   - string: The recording ID (sid) acquired from the Agora cloud recording API.
//   - error: An error object if any issues occurred during the request process.
func (s *CloudRecordingService) HandleStartRecordingReq(startReq StartRecordingRequest, resourceId string, modeType string) (string, error) {

	// build start recording endpoint
	url := fmt.Sprintf("%s/%s/cloud_recording/resourceid/%s/mode/%s/start", s.baseURL, s.appID, resourceId, modeType)

	// send request to start endpoint
	body, err := s.makeRequest("POST", url, startReq)
	if err != nil {
		return "", err
	}

	// Parse the response body to extract the necessary information
	var response StartRecordingResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	return response.Sid, nil
}
