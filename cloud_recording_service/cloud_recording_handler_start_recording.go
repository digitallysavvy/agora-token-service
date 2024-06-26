package cloud_recording_service

import (
	"encoding/json"
	"fmt"
)

// HandleStartRecordingReq handles the start recording request.
// It uses the makeRequest function to send the request to the Agora cloud recording API.
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
	url := fmt.Sprintf("%s/%s/cloud_recording/resourceid/%s/mode/%s/start", s.baseURL, s.appID, resourceId, modeType)
	body, err := s.makeRequest("POST", url, startReq)
	if err != nil {
		return "", err
	}

	// Parse the response body to extract the necessary information
	var response struct {
		Cname      string `json:"cname"`
		Uid        string `json:"uid"`
		ResourceId string `json:"resourceId"`
		Sid        string `json:"sid"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	return response.Sid, nil
}
