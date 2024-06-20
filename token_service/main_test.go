package token_service

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var testService *TokenService

// Create a service which can be fetched from any of the tests
func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	// os.Setenv("APP_ID", "18c2dfa12345678987654321fb84931d")
	// os.Setenv("APP_CERTIFICATE", "12345712345678765432117b84ad9ef9")
	// os.Setenv("SERVER_PORT", "8080")
	// os.Setenv("PORT", "8080")
	testService = NewTokenService()
	os.Exit(m.Run())
}
