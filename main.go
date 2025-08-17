package main

import (
	"context"
	"fmt"
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
// @host web3.bogowi.com
// @BasePath /api
// Server represents the application server
type Server struct {
	srv    *http.Server
	sdk    *sdk.BOGOWISDK
	config *config.Config
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config) (*Server, error) {
	// Initialize NetworkHandler (it creates SDKs internally)
	networkHandler, err := api.NewNetworkHandler(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize network handler: %w", err)
	}

	// Get default SDK based on environment
	var defaultSDK api.SDKInterface
	if cfg.Environment == "development" {
		defaultSDK, _ = networkHandler.GetSDK("testnet")
	} else {
		defaultSDK, _ = networkHandler.GetSDK("mainnet")
	}

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize API server with unified router
	routerConfig := &api.RouterConfig{
		SDK:            defaultSDK,
		NetworkHandler: networkHandler,
		AppConfig:      cfg,
		// Storage will use default (in-memory)
	}
	router := api.CreateRouter(routerConfig)

	// Create HTTP server
	srv := &http.Server{
		Addr:              ":" + cfg.APIPort,
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Server{
		srv:    srv,
		sdk:    nil, // We're using NetworkHandler now
		config: cfg,
	}, nil
}

// Start starts the server
func (s *Server) Start() error {
	log.Printf("🚀 BOGOWI API Server starting on port %s", s.config.APIPort)
	log.Printf("📚 Swagger documentation available at http://localhost:%s/docs", s.config.APIPort)
	log.Printf("🌍 Environment: %s", s.config.Environment)

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("🛑 Server shutting down...")
	err := s.srv.Shutdown(ctx)
	if err == nil {
		log.Println("✅ Server exited")
	}
	return err
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create server
	server, err := NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("⚠️ Server forced to shutdown: %v", err)
		os.Exit(1)
	}
}
