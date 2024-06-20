package token_service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// TokenService represents the main application TokenService.
type TokenService struct {
	// Server is the HTTP server for the application.
	Server *http.Server

	// Sigint is a channel to handle OS signals, such as Ctrl+C.
	Sigint chan os.Signal

	// appID is the identifier for the application.
	appID string

	// appCertificate is the certificate used by the application.
	appCertificate string

	// allowOrigin specifies the allowed origin for Cross-Origin Resource Sharing (CORS).
	allowOrigin string
}

// Stop TokenService safely, closing additional connections if needed.
func (s *TokenService) Stop() {
	// Will continue once an interrupt has occurred
	signal.Notify(s.Sigint, os.Interrupt)
	<-s.Sigint

	// cancel would be useful if we had to close third party connection first
	// Like connections to a db or cache
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	cancel()
	err := s.Server.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}
}

// Start runs the TokenService by listening to the specified port
func (s *TokenService) Start() {
	log.Println("Listening to port " + s.Server.Addr)
	if err := s.Server.ListenAndServe(); err != nil {
		panic(err)
	}
}

// NewTokenService returns a TokenService pointer with all configurations set
func NewTokenService() *TokenService {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	appIDEnv, appIDExists := os.LookupEnv("APP_ID")
	appCertEnv, appCertExists := os.LookupEnv("APP_CERTIFICATE")
	serverPort, serverPortExists := os.LookupEnv("SERVER_PORT")
	corsAllowOrigin, _ := os.LookupEnv("CORS_ALLOW_ORIGIN")

	if !appIDExists || !appCertExists || len(appIDEnv) == 0 || len(appCertEnv) == 0 {
		log.Fatal("FATAL ERROR: ENV not properly configured, check .env file or APP_ID and APP_CERTIFICATE")
	}
	if !serverPortExists || len(serverPort) == 0 {
		// Check $PORT, this is used by Railway.
		port, portExists := os.LookupEnv("PORT")
		if portExists && len(port) > 0 {
			serverPort = port
		} else {
			serverPort = "8080"
		}
	}

	s := &TokenService{
		Sigint: make(chan os.Signal, 1),
		Server: &http.Server{
			Addr: fmt.Sprintf(":%s", serverPort),
		},
		appID:          appIDEnv,
		appCertificate: appCertEnv,
		allowOrigin:    corsAllowOrigin,
	}

	api := gin.Default()

	api.Use(s.nocache())
	api.Use(s.CORSMiddleware())
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	api.POST("/getToken", s.getToken)
	s.Server.Handler = api
	return s
}
