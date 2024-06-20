package main

import (
	"github.com/AgoraIO-Community/agora-backend-service/token_service"
)

func main() {
	s := token_service.NewTokenService()
	// Stop is called on another thread, but waits for an interrupt
	go s.Stop()
	s.Start()
}
