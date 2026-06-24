package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/auth"
)

func TestServicePasswordToken(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse map[string]interface{}
		serverStatus   int
		wantErr        bool
		errContains    string
	}{
		{
			name: "successful token exchange",
			serverResponse: map[string]interface{}{
				"access_token":  "access-token",
				"token_type":    "Bearer",
				"expires_in":    3600,
				"refresh_token": "refresh-token",
			},
			serverStatus: http.StatusOK,
		},
		{
			name: "missing refresh token",
			serverResponse: map[string]interface{}{
				"access_token": "access-token",
				"token_type":   "Bearer",
				"expires_in":   3600,
			},
			serverStatus: http.StatusOK,
			wantErr:      true,
			errContains:  "no refresh token returned",
		},
	}

	service := auth.New()
	for _, testCase := range tests {
		t.Run(testCase.name, func(test *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(test, http.MethodPost, r.Method)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(testCase.serverStatus)
				require.NoError(test, json.NewEncoder(w).Encode(testCase.serverResponse))
			}))
			defer server.Close()

			response, err := service.PasswordToken(context.Background(), server.URL, "master", "admin", "admin")
			if testCase.wantErr {
				require.Error(test, err)
				assert.Contains(test, err.Error(), testCase.errContains)
				return
			}

			require.NoError(test, err)
			assert.Equal(test, "access-token", response.AccessToken)
			assert.Equal(test, "refresh-token", response.RefreshToken)
		})
	}
}

func TestServiceSetEnvToken(t *testing.T) {
	service := auth.New()
	envFile := filepath.Join(t.TempDir(), ".env")

	err := service.SetEnvToken("TEST_AUTH_TOKEN", "value", envFile)
	require.NoError(t, err)

	content, err := os.ReadFile(envFile)
	require.NoError(t, err)
	assert.True(t, strings.Contains(string(content), "TEST_AUTH_TOKEN=value"))
	assert.Equal(t, "value", os.Getenv("TEST_AUTH_TOKEN"))
}
