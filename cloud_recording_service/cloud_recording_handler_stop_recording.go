package cloud_recording_service

import (
	"encoding/json"
	"fmt"
	"time"
)

// HandleStopRecording processes the request to stop a cloud recording session in Agora's cloud recording service.
// It validates the request parameters, constructs the request URL, and use makeRequest to send data to Agora Endpoint
//
// Parameters:
//   - stopReq: StopRecordingRequest object containing details for the stop request.
//   - resourceId: The unique identifier of the resource (i.e., channel) being recorded.
//   - recordingId: The unique identifier of the recording session.
//   - modeType: The recording mode indicating the type of recording, such as individual or mixed.
//
// Returns:
//   - A byte array of the JSON response from the cloud recording service.
//   - An error if the operation fails at any stage, including request validation, API request sending, or response parsing.
//
// Notes:
//   - It is critical to ensure that all identifiers and request parameters are valid and not nil.
//   - This function uses s.baseURL, s.appID, s.customerID, and s.customerCertificate to construct the API request.
func (s *CloudRecordingService) HandleStopRecording(stopReq StopRecordingRequest, resourceId string, recordingId string, modeType string) (json.RawMessage, error) {
	// build stop recording endpoint
	url := fmt.Sprintf("%s/%s/cloud_recording/resourceid/%s/sid/%s/mode/%s/stop", s.baseURL, s.appID, resourceId, recordingId, modeType)

	// send request to stop recording endpoint
	body, err := s.makeRequest("POST", url, stopReq)
	if err != nil {
		return []byte{}, err
	}

	// Parse the response body to ensure it conforms to the expected structure
	var response StopRecordingResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return []byte{}, fmt.Errorf("error parsing response body into StopRecordingResponse: %v", err)
	}

	// Validate the FileList conforms to the expected structure based on FileListMode.
	_, err = response.ServerResponse.UnmarshalFileList()
	if err != nil {
		return nil, fmt.Errorf("error parsing ServerResponse: %v", err)
	}

	// Add timestamp to Agora response
	timestampBody, err := s.AddTimestamp(response)
	if err != nil {
		return nil, fmt.Errorf("error encoding timestamped response: %v", err)
	}

	return timestampBody, nil
}

// UnmarshalFileList interprets the file list from the server response based on the specified FileListMode.
// It supports different data representations as per the mode set in the response (e.g., string or JSON).
//
// Returns:
//   - An interface holding the parsed file list which can be a slice of FileDetail or FileListEntry.
//   - An error if parsing fails or the mode is not supported.
//
// Notes:
//   - This function assumes that FileListMode and FileList are not nil before proceeding with parsing.

func (sr *ServerResponse) UnmarshalFileList() (interface{}, error) {
	if sr.FileListMode == nil || sr.FileList == nil {
		// No FileListMode set or FileList is nil
		return nil, fmt.Errorf("FileListMode or FileList are empty")
	}
	switch *sr.FileListMode {
	case "string":
		var fileList []FileDetail
		if err := json.Unmarshal(*sr.FileList, &fileList); err != nil {
			return nil, fmt.Errorf("error parsing FileList into []FileDetail: %v", err)
		}
		return fileList, nil
	case "json":
		var fileList []FileListEntry
		if err := json.Unmarshal(*sr.FileList, &fileList); err != nil {
			return nil, fmt.Errorf("error parsing FileList into []FileListEntry: %v", err)
		}
		return fileList, nil
	default:
		return nil, fmt.Errorf("unknown FileListMode: %s", *sr.FileListMode)
	}
}

func (s *CloudRecordingService) AddTimestamp(response StopRecordingResponse) (json.RawMessage, error) {
	// Set the current timestamp
	now := time.Now().UTC().Format(time.RFC3339)
	response.Timestamp = &now

	// Marshal the response back to JSON
	timestampedBody, err := json.Marshal(response)
	if err != nil {
		return []byte{}, fmt.Errorf("error marshaling final response: %v", err)
	}
	return timestampedBody, nil
}
