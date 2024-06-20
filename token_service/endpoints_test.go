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

func TestRtcValidAndInvalid(t *testing.T) {

	tests := []UrlCodePair{
		{"/rtc/fsda/publisher/uid/0/?expiry=600", http.StatusOK, nil},
		{"/rtc/fsda/publisher/uid//?expiry=600", http.StatusOK, nil},
		{"/rtc/fsda/publisher/uid/test/?expiry=600", http.StatusBadRequest, nil},
		{"/rtc/fsda/publisher/uid/0/?expiry=failing", http.StatusBadRequest, nil},
	}
	for _, httpTest := range tests {
		testApi, err := http.NewRequest(http.MethodGet, httpTest.url, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp := httptest.NewRecorder()
		// Call the endpoint
		testService.Server.Handler.ServeHTTP(resp, testApi)
		assert.Equal(t, httpTest.code, resp.Code, resp.Body)
	}
}

func TestRtmValidAndInvalid(t *testing.T) {

	// Create a new gin engine for testing
	reqValid, err := http.NewRequest(http.MethodGet, "/rtm/username/?expiry=600", nil)
	if err != nil {
		t.Fatal(err)
	}
	reqInvalid, err := http.NewRequest(http.MethodGet, "/rtm/0/?expiry=600", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a response recorder to inspect the response
	resp := httptest.NewRecorder()

	// Call the endpoint
	testService.Server.Handler.ServeHTTP(resp, reqValid)

	assert.Equal(t, http.StatusOK, resp.Code)

	resp = httptest.NewRecorder()

	// Call the endpoint
	testService.Server.Handler.ServeHTTP(resp, reqInvalid)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

type UrlCodePair struct {
	url  string
	code int
	body []byte
}

func TestChatValidAndInvalid(t *testing.T) {

	tests := []UrlCodePair{
		{"/chat/app/", http.StatusOK, nil},
		{"/chat/account/username/", http.StatusOK, nil},
		{"/chat/account/", http.StatusNotFound, nil},
		{"/chat/invalid/", http.StatusNotFound, nil},
		{"/chat/account/username/?expiry=600", http.StatusOK, nil},
		{"/chat/account/username/?expiry=fail", http.StatusBadRequest, nil},
	}
	for _, httpTest := range tests {
		testApi, err := http.NewRequest(http.MethodGet, httpTest.url, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp := httptest.NewRecorder()
		// Call the endpoint
		testService.Server.Handler.ServeHTTP(resp, testApi)
		assert.Equal(t, httpTest.code, resp.Code, resp.Body)
	}
}

func TestRteValidAndInvalid(t *testing.T) {

	tests := []UrlCodePair{
		{"/rte/channelName/publisher/uid/0/rtmid/?expiry=600", http.StatusOK, nil},
		{"/rte/channelName/publisher/uid/2345/?expiry=600", http.StatusOK, nil},
		{"/rte/channelName/subscriber/uid/0/rtmid/?expiry=600", http.StatusOK, nil},
		{"/rte/channelName/subscriber/uid/2345/?expiry=600", http.StatusOK, nil},
		{"/rte/channelName/publisher/uid/0/?expiry=600", http.StatusBadRequest, nil},
		{"/rte/channelName/publisher/uid/2345/?expiry=failing", http.StatusBadRequest, nil},
	}
	for _, httpTest := range tests {
		testApi, err := http.NewRequest(http.MethodGet, httpTest.url, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp := httptest.NewRecorder()
		// Call the endpoint
		testService.Server.Handler.ServeHTTP(resp, testApi)
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

func TestTokentypes(t *testing.T) {

	// Create a new gin engine for testing
	reqValidUid, err := http.NewRequest(http.MethodGet, "/rtc/fsda/publisher/uid/0/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a new gin engine for testing
	reqValidUserAcc, err := http.NewRequest(http.MethodGet, "/rtc/fsda/publisher/userAccount/0/", nil)
	if err != nil {
		t.Fatal(err)
	}
	reqInvalid, err := http.NewRequest(http.MethodGet, "/rtc/fsda/publisher/nonsense/0/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a response recorder to inspect the response
	resp := httptest.NewRecorder()

	// Call the endpoint
	testService.Server.Handler.ServeHTTP(resp, reqValidUid)

	assert.Equal(t, http.StatusOK, resp.Code)

	resp = httptest.NewRecorder()

	// Call the endpoint
	testService.Server.Handler.ServeHTTP(resp, reqValidUserAcc)

	assert.Equal(t, http.StatusOK, resp.Code)

	resp = httptest.NewRecorder()

	// Call the endpoint
	testService.Server.Handler.ServeHTTP(resp, reqInvalid)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
