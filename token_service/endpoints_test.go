package token_service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UrlCodePair struct {
	url  string
	code int
	body []byte
}

func TestGetTokenValidAndInvalid(t *testing.T) {

	tests := []UrlCodePair{
		{"/getToken", http.StatusOK, []byte(`{"tokenType": "rtc", "channel": "channel123", "role": "publisher", "uid": "user123", "expire": 3600}`)},
		{"/getToken", http.StatusOK, []byte(`{"tokenType": "rtm", "uid": "user456", "expire": 1800}`)},
		{"/getToken", http.StatusOK, []byte(`{"tokenType": "rtm", "uid": "user456", "expire": 1800, "channel": "test"}`)},
		{"/getToken", http.StatusOK, []byte(`{"tokenType": "chat", "uid": "user789", "expire": 900}`)},
		{"/getToken", http.StatusBadRequest, []byte(`{"channel": "channel456", "role": "subscriber", "expire": 1800}`)},
		{"/getToken", http.StatusBadRequest, []byte(`{"tokenType": "invalid_type", "channel": "channel789", "role": "publisher", "uid": "user123", "expire": 3600}`)},
		{"/getToken", http.StatusBadRequest, []byte(`{"tokenType": "rtc", "role": "publisher", "uid": "user123", "expire": 3600}`)},
		{"/getToken", http.StatusBadRequest, []byte(`{"tokenType": "rtm", "expire": 1800}`)},
		{"/getToken", http.StatusBadRequest, []byte(``)},
		{"/getToken", http.StatusOK, []byte(`{"tokenType": "chat"}`)},
		{"/getToken", http.StatusOK, []byte(`{"tokenType": "chat", "uid": "user123"}`)},
	}
	for _, httpTest := range tests {
		testApi, err := http.NewRequest(http.MethodPost, httpTest.url, bytes.NewBuffer(httpTest.body))
		log.Println(bytes.NewBuffer(httpTest.body))
		if err != nil {
			t.Fatal(err)
		}
		resp := httptest.NewRecorder()
		// Call the endpoint
		testService.Server.Handler.ServeHTTP(resp, testApi)
		log.Println(httpTest.code)
		log.Println(resp.Code)
		assert.Equal(t, httpTest.code, resp.Code, resp.Body)
	}
}

func TestPingPong(t *testing.T) {
	testApi, err := http.NewRequest(http.MethodGet, "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp := httptest.NewRecorder()
	// Call the endpoint
	testService.Server.Handler.ServeHTTP(resp, testApi)

	// Check the response body
	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.Code)
	}

	assert.Equal(t, "pong", response["message"])
}
