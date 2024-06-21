package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AgoraIO-Community/agora-backend-service/cloud_recording_service"
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

	// Initialize Gin router
	r := gin.Default()

	// Create instances of your services
	tokenService := token_service.NewTokenService()
	cloudRecordingService := cloud_recording_service.NewCloudRecordingService(tokenService)

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