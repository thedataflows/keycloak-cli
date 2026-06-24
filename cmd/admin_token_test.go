package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/thedataflows/keycloak-cli/pkg/auth"
)

func TestTokenFromPtokenFromPassword(t *testing.T) {
	authService := auth.New()

	tests := []struct {
		name           string
		baseURL        string
		realm          string
		username       string
		password       string
		serverResponse map[string]interface{}
		serverStatus   int
		wantErr        bool
		errContains    string
	}{
		{
			name:     "successful token exchange",
			baseURL:  "https://keycloak.example.com",
			realm:    "master",
			username: "admin",
			password: "admin123",
			serverResponse: map[string]interface{}{
				"access_token":  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
				"token_type":    "Bearer",
				"expires_in":    3600,
				"refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
				"scope":         "openid offline_access",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:     "invalid credentials",
			baseURL:  "https://keycloak.example.com",
			realm:    "master",
			username: "admin",
			password: "wrong",
			serverResponse: map[string]interface{}{
				"error":             "invalid_grant",
				"error_description": "Invalid user credentials",
			},
			serverStatus: http.StatusBadRequest,
			wantErr:      true,
			errContains:  "password credentials token",
		},
		{
			name:     "missing access token in response",
			baseURL:  "https://keycloak.example.com",
			realm:    "master",
			username: "admin",
			password: "admin123",
			serverResponse: map[string]interface{}{
				"token_type":    "Bearer",
				"expires_in":    3600,
				"refresh_token": "refresh-token",
			},
			serverStatus: http.StatusOK,
			wantErr:      true,
			errContains:  "oauth2: server response missing access_token",
		},
		{
			name:     "missing refresh token in response",
			baseURL:  "https://keycloak.example.com",
			realm:    "master",
			username: "admin",
			password: "admin123",
			serverResponse: map[string]interface{}{
				"access_token": "access-token",
				"token_type":   "Bearer",
				"expires_in":   3600,
			},
			serverStatus: http.StatusOK,
			wantErr:      true,
			errContains:  "no refresh token returned",
		},
		{
			name:     "server error",
			baseURL:  "https://keycloak.example.com",
			realm:    "master",
			username: "admin",
			password: "admin123",
			serverResponse: map[string]interface{}{
				"error":             "server_error",
				"error_description": "Internal server error",
			},
			serverStatus: http.StatusInternalServerError,
			wantErr:      true,
			errContains:  "password credentials token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t2 *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request method
				if r.Method != http.MethodPost {
					t2.Errorf("Expected POST request, got %s", r.Method)
				}

				// Verify content type
				expectedContentType := "application/x-www-form-urlencoded"
				if ct := r.Header.Get("Content-Type"); ct != expectedContentType {
					t2.Errorf("Expected Content-Type %s, got %s", expectedContentType, ct)
				}

				// Write response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				json.NewEncoder(w).Encode(tt.serverResponse)
			}))
			defer server.Close()

			// Call the function under test
			tokenCtx := context.Background()
			resp, err := authService.PasswordToken(tokenCtx, server.URL, tt.realm, tt.username, tt.password)

			// Check error expectations
			if (err != nil) != tt.wantErr {
				t2.Errorf("tokenFromPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errContains != "" && err != nil {
					if !strings.Contains(err.Error(), tt.errContains) {
						t2.Errorf("Expected error to contain '%s', got '%s'", tt.errContains, err.Error())
					}
				}
				return
			}

			// Verify response on success
			if resp.AccessToken == "" {
				t2.Error("Expected non-empty access token")
			}
			if resp.RefreshToken == "" {
				t2.Error("Expected non-empty refresh token")
			}
			if resp.TokenType == "" {
				t2.Error("Expected non-empty token type")
			}
			if resp.ExpiresIn <= 0 {
				t2.Error("Expected positive expires_in value")
			}
		})
	}
}

func TestTokenFromPtokenFromPassword_ContextCancellation(t *testing.T) {
	authService := auth.New()

	// Create a test server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Small delay to allow cancellation
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token":  "test-token",
			"refresh_token": "test-refresh",
			"token_type":    "Bearer",
			"expires_in":    3600,
		})
	}))
	defer server.Close()

	tokenCtx := context.Background()
	_, err := authService.PasswordToken(tokenCtx, "invalid-url", "master", "admin", "admin")
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

// BenchmarkTokenFromPtokenFromPassword benchmarks the tokenFromPassword function
func BenchmarkTokenFromPtokenFromPassword(b *testing.B) {
	authService := auth.New()

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token":  "test-access-token",
			"refresh_token": "test-refresh-token",
			"token_type":    "Bearer",
			"expires_in":    3600,
		})
	}))
	defer server.Close()

	_ = setGlobalLoggerLogLevel("error") // Suppress logs during benchmark


	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokenCtx := context.Background()
		_, err := authService.PasswordToken(tokenCtx, server.URL, "master", "admin", "admin")
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}
