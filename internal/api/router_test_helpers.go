package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"bogowi-blockchain-go/internal/config"
	"bogowi-blockchain-go/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

// TestRouterDependencies creates minimal dependencies for testing
func TestRouterDependencies(t *testing.T) *RouterDependencies {
	cfg := &config.Config{
		Environment:       "test",
		APIPort:           "8080",
		FirebaseProjectID: "test-project",
	}

	// Create mock SDK
	mockSDK := NewSimpleMockSDK()

	return &RouterDependencies{
		NetworkHandler: nil, // Can be nil for basic tests
		DefaultSDK:     mockSDK,
		Config:         cfg,
		Storage:        storage.NewInMemoryRewardsStorage(),
		RateLimiter:    nil, // Skip rate limiting in tests
		AuthMiddleware: nil, // Skip auth in tests
		CORSConfig:     nil, // Skip CORS in tests
		TrustedProxies: nil,
		NFTHandlerFunc: nil, // Skip NFT routes in tests
	}
}

// TestRequest helps make test requests
type TestRequest struct {
	Method  string
	Path    string
	Body    interface{}
	Headers map[string]string
	Query   map[string]string
}

// TestResponse captures response data
type TestResponse struct {
	Code   int
	Body   []byte
	Result interface{}
}

// MakeRequest performs a test request
func MakeRequest(router *gin.Engine, req TestRequest) (*TestResponse, error) {
	// Prepare body
	var bodyReader io.Reader
	if req.Body != nil {
		jsonBytes, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest(req.Method, req.Path, bodyReader)
	if err != nil {
		return nil, err
	}

	// Add headers
	if req.Body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// Add query params
	if len(req.Query) > 0 {
		q := httpReq.URL.Query()
		for k, v := range req.Query {
			q.Add(k, v)
		}
		httpReq.URL.RawQuery = q.Encode()
	}

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httpReq)

	return &TestResponse{
		Code: w.Code,
		Body: w.Body.Bytes(),
	}, nil
}

// ParseJSON parses the response body into the result
func (tr *TestResponse) ParseJSON(result interface{}) error {
	return json.Unmarshal(tr.Body, result)
}

// AssertStatus checks the response status code
func (tr *TestResponse) AssertStatus(t *testing.T, expected int) {
	require.Equal(t, expected, tr.Code, "unexpected status code: %s", string(tr.Body))
}

// AssertJSON checks that response contains expected JSON
func (tr *TestResponse) AssertJSON(t *testing.T, expected interface{}) {
	expectedJSON, err := json.Marshal(expected)
	require.NoError(t, err)

	var actualData, expectedData interface{}
	require.NoError(t, json.Unmarshal(tr.Body, &actualData))
	require.NoError(t, json.Unmarshal(expectedJSON, &expectedData))

	require.Equal(t, expectedData, actualData)
}

// TestRouterOptions configures test router creation
type TestRouterOptions struct {
	SkipMiddleware     bool
	WithAuth           bool
	WithRateLimiter    bool
	WithNFTRoutes      bool
	WithNetworkSupport bool
	MockSDK            *SimpleMockSDK
	Storage            storage.RewardsStorage
}

// CreateTestRouter creates a router configured for testing
func CreateTestRouter(t *testing.T, opts TestRouterOptions) *gin.Engine {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	deps := TestRouterDependencies(t)

	// Apply options
	if opts.MockSDK != nil {
		deps.DefaultSDK = opts.MockSDK
	}
	if opts.Storage != nil {
		deps.Storage = opts.Storage
	}
	if opts.WithRateLimiter {
		deps.RateLimiter = rate.NewLimiter(rate.Every(1), 10)
	}
	if opts.WithAuth {
		// Create mock auth middleware - for now, skip auth in tests
		// deps.AuthMiddleware = CreateMockAuthMiddleware()
	}
	if opts.WithNFTRoutes {
		deps.NFTHandlerFunc = func(h *Handler) *NFTHandler {
			return &NFTHandler{
				Handler: h,
			}
		}
	}
	if opts.WithNetworkSupport {
		// Create mock network handler
		deps.NetworkHandler = CreateMockNetworkHandler()
	}

	builder := NewRouterBuilder(deps)
	if opts.SkipMiddleware {
		builder.SkipMiddleware()
	}

	return builder.Build()
}

// CreateMockAuthMiddleware creates a mock auth middleware for testing
func CreateMockAuthMiddleware() *MockAuthMiddleware {
	return &MockAuthMiddleware{
		allowAll: true,
	}
}

// MockAuthMiddleware for testing
type MockAuthMiddleware struct {
	allowAll bool
	userID   string
}

// Authenticate implements the auth interface
func (m *MockAuthMiddleware) Authenticate(c *gin.Context) {
	if m.allowAll {
		if m.userID != "" {
			c.Set("userID", m.userID)
		}
		c.Next()
	} else {
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
	}
}

// SetUserID sets the user ID for authenticated requests
func (m *MockAuthMiddleware) SetUserID(userID string) {
	m.userID = userID
}

// SetAllowAll controls whether all requests are allowed
func (m *MockAuthMiddleware) SetAllowAll(allow bool) {
	m.allowAll = allow
}

// CreateMockNetworkHandler creates a mock network handler
func CreateMockNetworkHandler() *NetworkHandler {
	// This would be a full mock implementation
	// For now, returning nil as placeholder
	return nil
}

// AssertRouteExists checks if a route exists
func AssertRouteExists(t *testing.T, router *gin.Engine, method, path string) {
	routes := router.Routes()
	found := false
	for _, route := range routes {
		if route.Method == method && route.Path == path {
			found = true
			break
		}
	}
	require.True(t, found, "route %s %s not found", method, path)
}

// AssertRouteNotExists checks if a route does not exist
func AssertRouteNotExists(t *testing.T, router *gin.Engine, method, path string) {
	routes := router.Routes()
	found := false
	for _, route := range routes {
		if route.Method == method && route.Path == path {
			found = true
			break
		}
	}
	require.False(t, found, "route %s %s should not exist", method, path)
}
