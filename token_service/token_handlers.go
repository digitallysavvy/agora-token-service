/** token_handlers.go */

package token_service

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/AgoraIO-Community/go-tokenbuilder/chatTokenBuilder"
	rtctokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	rtmtokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtmtokenbuilder"
)

// HandleGetToken handles the HTTP request to generate a token based on the provided tokenType.
// It checks the tokenType from the query parameters and calls the appropriate token generation method.
// The generated token is sent as a JSON response to the client.
//
// Parameters:
//   - w: http.ResponseWriter - The HTTP response writer to send the response to the client.
//   - r: *http.Request - The HTTP request received from the client.
//
// Behavior:
//  1. Retrieves the tokenType from the query parameters. Error if invalid entry or not provided.
//  2. Uses a switch statement to handle different tokenType cases:
//     - "rtc": Calls the GenRtcToken method to generate the RTC token and sends it as a JSON response.
//     - "rtm": Calls the GenRtmToken method to generate the RTM token and sends it as a JSON response.
//     - "chat": Calls the GenChatToken method to generate the chat token and sends it as a JSON response.
//     - Default: Returns an error response indicating an unsupported token type.
//
// Notes:
//   - The actual token generation methods (GenRtcToken, GenRtmToken, and GenChatToken) are part of the TokenService struct.
//   - The generated token is sent as a JSON response with appropriate HTTP status codes.
//
// Example usage:
//
//	router.POST("/getNew", TokenService.GetToken)
func (s *TokenService) HandleGetToken(tokenReq TokenRequest, w http.ResponseWriter) {

	var token string
	var tokenErr error

	switch tokenReq.TokenType {
	case "rtc":
		token, tokenErr = s.GenRtcToken(tokenReq)
	case "rtm":
		token, tokenErr = s.GenRtmToken(tokenReq)
	case "chat":
		token, tokenErr = s.GenChatToken(tokenReq)
	default:
		http.Error(w, "Unsupported tokenType", http.StatusBadRequest)
		return
	}
	if tokenErr != nil {
		http.Error(w, tokenErr.Error(), http.StatusBadRequest)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{Token: token}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GenRtcToken generates an RTC token based on the provided TokenRequest and returns it.
//
// Parameters:
//   - tokenRequest: TokenRequest - The TokenRequest struct containing the required information for RTC token generation.
//
// Returns:
//   - string: The generated RTC token.
//   - error: An error if there are any issues during token generation or validation.
//
// Behavior:
//  1. Validates the required fields in the TokenRequest (channel and UID).
//  2. Sets a default expiration time of 3600 seconds (1 hour) if not provided in the request.
//  3. Determines the user's role (publisher or subscriber) based on the "Role" field in the request.
//  4. Generates the RTC token using the rtctokenbuilder2 package.
//
// Notes:
//   - The rtctokenbuilder2 package is used for generating RTC tokens.
//   - The "Role" field can be "publisher" or "subscriber"; other values are considered invalid.
//
// Example usage:
//
//	tokenReq := TokenRequest{
//	    TokenType:  "rtc",
//	    Channel:    "my_channel",
//	    Uid:        "user123",
//	    RtcRole:    "publisher",
//	    ExpirationSeconds: 3600,
//	}
//	token, err := TokenService.GenRtcToken(tokenReq)
func (s *TokenService) GenRtcToken(tokenRequest TokenRequest) (string, error) {
	if tokenRequest.Channel == "" {
		return "", errors.New("invalid: missing channel name")
	}
	if tokenRequest.Uid == "" {
		return "", errors.New("invalid: missing user ID or account")
	}

	var userRole rtctokenbuilder2.Role
	if tokenRequest.RtcRole == "publisher" {
		userRole = rtctokenbuilder2.RolePublisher
	} else {
		userRole = rtctokenbuilder2.RoleSubscriber
	}

	if tokenRequest.ExpirationSeconds == 0 {
		tokenRequest.ExpirationSeconds = 3600
	}

	uid64, parseErr := strconv.ParseUint(tokenRequest.Uid, 10, 64)
	if parseErr != nil {
		return rtctokenbuilder2.BuildTokenWithAccount(
			s.appID, s.appCertificate, tokenRequest.Channel,
			tokenRequest.Uid, userRole, uint32(tokenRequest.ExpirationSeconds),
		)
	}

	return rtctokenbuilder2.BuildTokenWithUid(
		s.appID, s.appCertificate, tokenRequest.Channel,
		uint32(uid64), userRole, uint32(tokenRequest.ExpirationSeconds),
	)
}

// GenRtmToken generates an RTM (Real-Time Messaging) token based on the provided TokenRequest and returns it.
//
// Parameters:
//   - tokenRequest: TokenRequest - The TokenRequest struct containing the required information for RTM token generation.
//
// Returns:
//   - string: The generated RTM token.
//   - error: An error if there are any issues during token generation or validation.
//
// Behavior:
//  1. Validates the required field in the TokenRequest (UID).
//  2. Sets a default expiration time of 3600 seconds (1 hour) if not provided in the request.
//  3. Generates the RTM token using the rtmtokenbuilder2 package.
//
// Notes:
//   - The rtmtokenbuilder2 package is used for generating RTM tokens.
//   - The "UID" field in TokenRequest is mandatory for RTM token generation.
//
// Example usage:
//
//	tokenReq := TokenRequest{
//	    TokenType:  "rtm",
//	    Uid:        "user123",
//	    ExpirationSeconds: 3600,
//	}
//	token, err := TokenService.GenRtmToken(tokenReq)
func (s *TokenService) GenRtmToken(tokenRequest TokenRequest) (string, error) {
	if tokenRequest.Uid == "" {
		return "", errors.New("invalid: missing user ID or account")
	}
	if tokenRequest.ExpirationSeconds == 0 {
		tokenRequest.ExpirationSeconds = 3600
	}

	return rtmtokenbuilder2.BuildToken(
		s.appID, s.appCertificate,
		tokenRequest.Uid,
		uint32(tokenRequest.ExpirationSeconds),
		tokenRequest.Channel,
	)
}

// GenChatToken generates a chat token based on the provided TokenRequest and returns it.
//
// Parameters:
//   - tokenRequest: TokenRequest - The TokenRequest struct containing the required information for chat token generation.
//
// Returns:
//   - string: The generated chat token.
//   - error: An error if there are any issues during token generation or validation.
//
// Behavior:
//  1. Sets a default expiration time of 3600 seconds (1 hour) if not provided in the request.
//  2. Determines whether to generate a chat app token or a chat user token based on the "UID" field in the request.
//  3. Generates the chat token using the chatTokenBuilder package.
//
// Notes:
//   - The chatTokenBuilder package is used for generating chat tokens.
//   - If the "UID" field is empty, a chat app token is generated; otherwise, a chat user token is generated.
//
// Example usage:
//
//	// Generate a chat app token
//	tokenReq := TokenRequest{
//	    TokenType:  "chat",
//	    ExpirationSeconds: 3600,
//	}
//	token, err := TokenService.GenChatToken(tokenReq)
//
//	// Generate a chat user token
//	tokenReq := TokenRequest{
//	    TokenType:  "chat",
//	    Uid:        "user123",
//	    ExpirationSeconds: 3600,
//	}
//	token, err := TokenService.GenChatToken(tokenReq)
func (s *TokenService) GenChatToken(tokenRequest TokenRequest) (string, error) {
	if tokenRequest.ExpirationSeconds == 0 {
		tokenRequest.ExpirationSeconds = 3600
	}

	var chatToken string
	var tokenErr error

	if tokenRequest.Uid == "" {
		chatToken, tokenErr = chatTokenBuilder.BuildChatAppToken(
			s.appID, s.appCertificate, uint32(tokenRequest.ExpirationSeconds),
		)
	} else {
		chatToken, tokenErr = chatTokenBuilder.BuildChatUserToken(
			s.appID, s.appCertificate,
			tokenRequest.Uid,
			uint32(tokenRequest.ExpirationSeconds),
		)

	}

	return chatToken, tokenErr
}