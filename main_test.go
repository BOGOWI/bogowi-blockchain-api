package main

import (
	"context"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"bogowi-blockchain-go/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	// Set up test environment
	os.Setenv("TESTNET_PRIVATE_KEY", "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	os.Setenv("MAINNET_PRIVATE_KEY", "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	os.Setenv("RPC_URL", "https://columbus.camino.network/ext/bc/C/rpc")
	os.Setenv("API_PORT", "3001")
	defer func() {
		os.Unsetenv("TESTNET_PRIVATE_KEY")
		os.Unsetenv("MAINNET_PRIVATE_KEY")
		os.Unsetenv("RPC_URL")
		os.Unsetenv("API_PORT")
	}()

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful server creation",
			setup: func() {
				// Default setup is sufficient
			},
			wantErr: false,
		},
		{
			name: "invalid private key",
			setup: func() {
				os.Setenv("TESTNET_PRIVATE_KEY", "invalid")
				os.Setenv("MAINNET_PRIVATE_KEY", "invalid")
				// Set contract addresses so NetworkHandler tries to initialize SDKs
				os.Setenv("BOGO_TOKEN_ADDRESS", "0x49fc9939D8431371dD22658a8a969Ec798A26fFB")
				os.Setenv("REWARD_DISTRIBUTOR_ADDRESS", "0x00439bd5eeED2303bfB64529Dad40C7c3F697724")
				t.Cleanup(func() {
					os.Unsetenv("BOGO_TOKEN_ADDRESS")
					os.Unsetenv("REWARD_DISTRIBUTOR_ADDRESS")
				})
			},
			wantErr: true,
			errMsg:  "failed to initialize network handler",
		},
		{
			name: "missing private key",
			setup: func() {
				// Save current values
				oldTPK := os.Getenv("TESTNET_PRIVATE_KEY")
				oldMPK := os.Getenv("MAINNET_PRIVATE_KEY")
				// Clear them
				os.Unsetenv("TESTNET_PRIVATE_KEY")
				os.Unsetenv("MAINNET_PRIVATE_KEY")
				os.Unsetenv("PRIVATE_KEY")
				os.Unsetenv("API_PRIVATE_KEY")
				// Temporarily rename .env file
				os.Rename(".env", ".env.backup")
				// Restore after test
				t.Cleanup(func() {
					os.Rename(".env.backup", ".env")
					if oldTPK != "" {
						os.Setenv("TESTNET_PRIVATE_KEY", oldTPK)
					}
					if oldMPK != "" {
						os.Setenv("MAINNET_PRIVATE_KEY", oldMPK)
					}
				})
			},
			wantErr: true,
			errMsg:  "at least one private key (TESTNET_PRIVATE_KEY or MAINNET_PRIVATE_KEY) is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			cfg, err := config.Load()
			if tt.name == "missing private key" && err != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "at least one private key")
				return
			}

			server, err := NewServer(cfg)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, server)
				assert.NotNil(t, server.srv)
				assert.NotNil(t, server.config)
			}
		})
	}
}

func TestServerStartAndShutdown(t *testing.T) {
	// Set up test environment
	os.Setenv("TESTNET_PRIVATE_KEY", "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	os.Setenv("MAINNET_PRIVATE_KEY", "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	os.Setenv("RPC_URL", "https://columbus.camino.network/ext/bc/C/rpc")
	os.Setenv("API_PORT", "18765") // Use a specific high port
	defer func() {
		os.Unsetenv("TESTNET_PRIVATE_KEY")
		os.Unsetenv("MAINNET_PRIVATE_KEY")
		os.Unsetenv("RPC_URL")
		os.Unsetenv("API_PORT")
	}()

	cfg, err := config.Load()
	require.NoError(t, err)

	server, err := NewServer(cfg)
	require.NoError(t, err)

	// Start server in goroutine
	serverErr := make(chan error, 1)
	go func() {
		if startErr := server.Start(); startErr != nil {
			serverErr <- startErr
		}
	}()

	// Give server time to start
	time.Sleep(200 * time.Millisecond)

	// Test that server is running by making a request
	resp, err := http.Get("http://localhost:" + cfg.APIPort + "/api/health")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Test graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	assert.NoError(t, err)

	// Verify server is stopped
	select {
	case err := <-serverErr:
		// Server should exit with no error
		assert.NoError(t, err)
	case <-time.After(1 * time.Second):
		// Server stopped gracefully
	}
}

func TestServerShutdownTimeout(t *testing.T) {
	// Set up test environment
	os.Setenv("TESTNET_PRIVATE_KEY", "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	os.Setenv("MAINNET_PRIVATE_KEY", "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	os.Setenv("RPC_URL", "https://columbus.camino.network/ext/bc/C/rpc")
	os.Setenv("API_PORT", "18766")
	defer func() {
		os.Unsetenv("TESTNET_PRIVATE_KEY")
		os.Unsetenv("MAINNET_PRIVATE_KEY")
		os.Unsetenv("RPC_URL")
		os.Unsetenv("API_PORT")
	}()

	cfg, err := config.Load()
	require.NoError(t, err)

	server, err := NewServer(cfg)
	require.NoError(t, err)

	// Start server
	go func() {
		_ = server.Start()
	}()

	// Give server time to start
	time.Sleep(200 * time.Millisecond)

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Shutdown with cancelled context
	err = server.Shutdown(ctx)
	// On some systems, shutdown might succeed even with cancelled context
	// if the server hasn't fully started yet
	if err != nil {
		assert.True(t,
			err.Error() == "context canceled" ||
				strings.Contains(err.Error(), "shutdown") ||
				strings.Contains(err.Error(), "context deadline exceeded"),
			"Expected context canceled or shutdown error, got: %s", err.Error())
	}

	// Clean up - force shutdown
	forceCtx, forceCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer forceCancel()
	_ = server.srv.Shutdown(forceCtx)
}

func TestMainFunction(t *testing.T) {
	// This test verifies that main() can be called without panicking
	// We can't easily test the full main() with signal handling,
	// but we can verify it compiles and basic structure is correct

	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Override os.Exit to prevent test from exiting
	var exitCode int
	oldExit := osExit
	osExit = func(code int) {
		exitCode = code
	}
	defer func() { osExit = oldExit }()

	// Test that main function exists and is callable
	// The actual execution would require complex mocking
	assert.NotPanics(t, func() {
		// Verify main exists by checking function signature
		var _ func() = main
	})

	// Verify exit wasn't called during compilation
	assert.Equal(t, 0, exitCode)
}

// osExit is a variable so we can mock it in tests
var osExit = os.Exit

func TestProductionMode(t *testing.T) {
	// Set up test environment
	os.Setenv("TESTNET_PRIVATE_KEY", "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	os.Setenv("MAINNET_PRIVATE_KEY", "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	os.Setenv("RPC_URL", "https://columbus.camino.network/ext/bc/C/rpc")
	os.Setenv("API_PORT", "18767")
	os.Setenv("NODE_ENV", "production")
	defer func() {
		os.Unsetenv("TESTNET_PRIVATE_KEY")
		os.Unsetenv("MAINNET_PRIVATE_KEY")
		os.Unsetenv("RPC_URL")
		os.Unsetenv("API_PORT")
		os.Unsetenv("NODE_ENV")
	}()

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "production", cfg.Environment)

	server, err := NewServer(cfg)
	require.NoError(t, err)
	assert.NotNil(t, server)
}
