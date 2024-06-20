package token_service

// Add these imports if they are not already present
import (
	"os"
	"testing"
)

func CreateTestTokenService(t *testing.T) *TokenService {
	appIdEnv, appIDExists := os.LookupEnv("APP_ID")
	appCertEnv, appCertExists := os.LookupEnv("APP_CERTIFICATE")
	if !appIDExists || !appCertExists {
		t.Errorf("check appId or appCertificate")
	}
	return &TokenService{
		appID:          appIdEnv,
		appCertificate: appCertEnv,
	}
}

// TestGenRtcToken tests the GenRtcToken function.
func TestGenRtcToken(t *testing.T) {
	tokenService := CreateTestTokenService(t)

	// Test valid RTC token generation
	tokenReq := TokenRequest{
		TokenType:         "rtc",
		Channel:           "my_channel",
		Uid:               "user123",
		RtcRole:           "publisher",
		ExpirationSeconds: 3600,
	}

	token, err := tokenService.GenRtcToken(tokenReq)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// Add assertions for the expected token value
	// For example: assert token is not empty, or has a specific format, etc.
	if token == "" {
		t.Error("Expected a non-empty RTC token")
	}

	// Test missing expiration, subscriber, int Uid (optional param)
	validWithoutExpiration := TokenRequest{
		TokenType: "rtc",
		Channel:   "my_channel",
		Uid:       "123",
		RtcRole:   "subscriber",
	}

	_, err = tokenService.GenRtcToken(validWithoutExpiration)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test invalid RTC token generation (missing channel name)
	invalidTokenReq := TokenRequest{
		TokenType:         "rtc",
		RtcRole:           "publisher",
		ExpirationSeconds: 3600,
	}

	_, err = tokenService.GenRtcToken(invalidTokenReq)
	if err == nil {
		t.Error("Expected error, but got nil")
	}

	// Test invalid RTC token generation (missing channel name)
	invalidTokenReq2 := TokenRequest{
		TokenType: "rtc",
		Channel:   "my_channel",
		RtcRole:   "subscriber",
	}

	_, err = tokenService.GenRtcToken(invalidTokenReq2)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
}

// TestGenRtmToken tests the genRtmToken function.
func TestGenRtmToken(t *testing.T) {
	tokenService := CreateTestTokenService(t)

	// Test valid RTM token generation
	tokenReq := TokenRequest{
		TokenType:         "rtm",
		Uid:               "user123",
		ExpirationSeconds: 3600,
	}

	token, err := tokenService.GenRtmToken(tokenReq)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// Add assertions for the expected token value
	// For example: assert token is not empty, or has a specific format, etc.
	if token == "" {
		t.Error("Expected a non-empty RTM token")
	}

	// Test valid RTM token generation with channel, but without expiration
	tokenChannelReq := TokenRequest{
		TokenType: "rtm",
		Uid:       "user123",
		Channel:   "test_channel",
	}

	token, err = tokenService.GenRtmToken(tokenChannelReq)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// Add assertions for the expected token value
	// For example: assert token is not empty, or has a specific format, etc.
	if token == "" {
		t.Error("Expected a non-empty RTM token")
	}

	// Test invalid RTM token generation (missing user ID)
	invalidTokenReq := TokenRequest{
		TokenType:         "rtm",
		ExpirationSeconds: 3600,
	}

	_, err = tokenService.GenRtmToken(invalidTokenReq)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	// Add assertions for the expected error message, for example:
	// assert error contains "missing user ID or account"
}

// TestGenChatToken tests the genChatToken function.
func TestGenChatToken(t *testing.T) {
	service := CreateTestTokenService(t)

	// Test valid chat token generation (chat app token)
	tokenReqApp := TokenRequest{
		TokenType:         "chat",
		ExpirationSeconds: 3600,
	}

	tokenApp, err := service.GenChatToken(tokenReqApp)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// Add assertions for the expected tokenApp value
	// For example: assert tokenApp is not empty, or has a specific format, etc.
	if tokenApp == "" {
		t.Error("Expected a non-empty chat app token")
	}

	// Test valid chat token generation (chat user token)
	tokenReqUser := TokenRequest{
		TokenType:         "chat",
		Uid:               "user123",
		ExpirationSeconds: 3600,
	}

	tokenUser, err := service.GenChatToken(tokenReqUser)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// Add assertions for the expected tokenUser value
	// For example: assert tokenUser is not empty, or has a specific format, etc.
	if tokenUser == "" {
		t.Error("Expected a non-empty chat user token")
	}

	// Test chat token generation without expiration
	invalidTokenReq := TokenRequest{
		TokenType: "chat",
		Uid:       "user123",
	}

	_, err = service.GenChatToken(invalidTokenReq)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
