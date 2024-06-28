package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/AgoraIO-Community/agora-backend-service/cloud_recording_service"
	"github.com/AgoraIO-Community/agora-backend-service/middleware"
	"github.com/AgoraIO-Community/agora-backend-service/token_service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	appIDEnv, appIDExists := os.LookupEnv("APP_ID")
	appCertEnv, appCertExists := os.LookupEnv("APP_CERTIFICATE")
	customerIDEnv, customerIDExists := os.LookupEnv("CUSTOMER_ID")
	customerSecretEnv, customerSecretExists := os.LookupEnv("CUSTOMER_SECRET")
	corsAllowOrigin, _ := os.LookupEnv("CORS_ALLOW_ORIGIN")
	baseURLEnv, baseURLExists := os.LookupEnv("AGORA_BASE_URL")
	cloudRecordingURLEnv, cloudRecordingURLExists := os.LookupEnv("AGORA_CLOUD_RECORDING_URL")
	storageVendorEnv, vendorExists := os.LookupEnv("STORAGE_VENDOR")
	storageRegionEnv, regionExists := os.LookupEnv("STORAGE_REGION")
	storageBucketEnv, bucketExists := os.LookupEnv("STORAGE_BUCKET")
	storageAccessKeyEnv, accessKeyExists := os.LookupEnv("STORAGE_BUCKET_ACCESS_KEY")
	storageSecretKeyEnv, secretKeyExists := os.LookupEnv("STORAGE_BUCKET_SECRET_KEY")

	if !appIDExists || !appCertExists || !customerIDExists || !customerSecretExists || !baseURLExists || !cloudRecordingURLExists ||
		!secretKeyExists || !vendorExists || !regionExists || !bucketExists || !accessKeyExists {
		log.Fatal("FATAL ERROR: ENV not properly configured, check .env file for all required variables")
	}

	storageVenderInt, storageVenderErr := strconv.Atoi(storageVendorEnv)
	storageRegionInt, storageRegionErr := strconv.Atoi(storageRegionEnv)

	if storageVenderErr != nil || storageRegionErr != nil {
		log.Fatal("FATAL ERROR: Invalid STORAGE_VENDOR / STORAGE_REGION not properly configured")
	}

	// Set Storage Config
	storageConfig := cloud_recording_service.StorageConfig{
		Vendor:    storageVenderInt,
		Region:    storageRegionInt,
		Bucket:    storageBucketEnv,
		AccessKey: storageAccessKeyEnv,
		SecretKey: storageSecretKeyEnv,
	}

	// Initialize Gin router
	r := gin.Default()

	// add headers
	var middleware = middleware.NewMiddleware(corsAllowOrigin)
	r.Use(middleware.NoCache())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.TimestampMiddleware())

	// Create instances of your services
	tokenService := token_service.NewTokenService(appIDEnv, appCertEnv, corsAllowOrigin)
	cloudRecordingService := cloud_recording_service.NewCloudRecordingService(appIDEnv, baseURLEnv+cloudRecordingURLEnv, getBasicAuth(customerIDEnv, customerSecretEnv), tokenService, storageConfig)

	// Register routes for each service
	tokenService.RegisterRoutes(r)
	cloudRecordingService.RegisterRoutes(r)
	r.GET("/ping", Ping)

	// Get the server port from environment variables or use a default
	serverPort, exists := os.LookupEnv("SERVER_PORT")
	if !exists {
		serverPort = "8080"
	}

	// Configure and start the HTTP server
	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: r,
	}

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Create a buffered channel to receive OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

// Ping is a simple handler for the /ping route.
// It responds with a "pong" message to indicate that the service is running.
//
// Parameters:
//   - c: *gin.Context - The Gin context representing the HTTP request and response.
//
// Behavior:
//   - Sends a JSON response with a "pong" message.
//
// Notes:
//   - This function is useful for health checks and ensuring that the service is up and running.
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func getBasicAuth(customerID string, customerSecret string) string {
	auth := fmt.Sprintf("%s:%s", customerID, customerSecret)
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
