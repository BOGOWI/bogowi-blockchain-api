package api

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Note: GetNetworkFromContext is defined in router_v2.go
func TestGetNetworkFromContext(t *testing.T) {
	tests := []struct {
		name         string
		setupContext func(*gin.Context)
		want         string
	}{
		{
			name: "testnet from context",
			setupContext: func(c *gin.Context) {
				c.Set("network", "testnet")
			},
			want: "testnet",
		},
		{
			name: "mainnet from context",
			setupContext: func(c *gin.Context) {
				c.Set("network", "mainnet")
			},
			want: "mainnet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			
			tt.setupContext(c)
			
			if got := GetNetworkFromContext(c); got != tt.want {
				t.Errorf("GetNetworkFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorResponse(t *testing.T) {
	err := ErrorResponse{
		Error: "test error",
	}
	
	if err.Error != "test error" {
		t.Errorf("ErrorResponse.Error = %v, want %v", err.Error, "test error")
	}
}

func TestSuccessResponse(t *testing.T) {
	resp := SuccessResponse{
		Message: "Operation successful",
		Data: map[string]interface{}{
			"id":   123,
			"name": "test",
		},
	}
	
	if resp.Message != "Operation successful" {
		t.Errorf("SuccessResponse.Message = %v, want %v", resp.Message, "Operation successful")
	}
	
	if data, ok := resp.Data.(map[string]interface{}); !ok {
		t.Errorf("SuccessResponse.Data type assertion failed")
	} else {
		if data["id"] != 123 {
			t.Errorf("SuccessResponse.Data[id] = %v, want %v", data["id"], 123)
		}
		if data["name"] != "test" {
			t.Errorf("SuccessResponse.Data[name] = %v, want %v", data["name"], "test")
		}
	}
}

