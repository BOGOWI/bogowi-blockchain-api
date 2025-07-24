package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "bogowi-blockchain-go/docs" // Import generated docs
	"bogowi-blockchain-go/internal/api"
	"bogowi-blockchain-go/internal/config"
	"bogowi-blockchain-go/internal/sdk"

	"github.com/gin-gonic/gin"
)

// @title BOGOWI Blockchain API
// @version 1.0
// @description API for BOGOWI blockchain operations including tokens, NFTs, and DAO functionality
// @host localhost:3001
// @BasePath /api
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize SDK
	bogoSDK, err := sdk.NewBOGOWISDK(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize SDK: %v", err)
	}

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize API server
	router := api.NewRouter(bogoSDK, cfg)

	// Create HTTP server
	srv := &http.Server{
		Addr:              ":" + cfg.APIPort,
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 BOGOWI API Server starting on port %s", cfg.APIPort)
		log.Printf("📚 Swagger documentation available at http://localhost:%s/docs", cfg.APIPort)
		log.Printf("🌍 Environment: %s", cfg.Environment)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Server shutting down...")

	// The context is used to inform the server it has 5 seconds to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("⚠️ Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	log.Println("✅ Server exited")
}
