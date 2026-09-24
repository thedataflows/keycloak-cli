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
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func TestServiceClientCredentialsToken(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse map[string]interface{}
		serverStatus   int
		wantErr        bool
		errContains    string
	}{
		{
			name: "successful client_credentials exchange",
			serverResponse: map[string]interface{}{
				"access_token": "cc-access-token",
				"token_type":   "Bearer",
				"expires_in":   300,
			},
			serverStatus: http.StatusOK,
		},
		{
			name:           "missing access token",
			serverResponse: map[string]interface{}{"token_type": "Bearer"},
			serverStatus:   http.StatusOK,
			wantErr:        true,
			errContains:    "no access token returned",
		},
		{
			name:           "non-ok status",
			serverResponse: map[string]interface{}{"error": "invalid_client"},
			serverStatus:   http.StatusUnauthorized,
			wantErr:        true,
			errContains:    "401",
		},
	}

	service := auth.New()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
				assert.Equal(t, "client_credentials", r.FormValue("grant_type"))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.serverStatus)
				require.NoError(t, json.NewEncoder(w).Encode(tc.serverResponse))
			}))
			defer server.Close()

			token, err := service.ClientCredentialsToken(context.Background(), server.URL, "master", "my-client", "secret")
			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errContains)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "cc-access-token", token.AccessToken)
		})
	}
}

// setEnv sets key to value for the duration of the test and restores the
// previous state afterwards. An empty value unsets the variable.
func setEnv(t *testing.T, key, value string) {
	t.Helper()
	previous, had := os.LookupEnv(key)
	if value == "" {
		require.NoError(t, os.Unsetenv(key))
	} else {
		require.NoError(t, os.Setenv(key, value))
	}
	t.Cleanup(func() {
		if had {
			require.NoError(t, os.Setenv(key, previous))
			return
		}
		require.NoError(t, os.Unsetenv(key))
	})
}

// resetTokenEnv clears every credential and token variable the auth service
// reads or writes, so tests are independent of the developer's environment.
func resetTokenEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"KEYCLOAK_USERNAME", "KC_BOOTSTRAP_ADMIN_USERNAME",
		"KEYCLOAK_PASSWORD", "KC_BOOTSTRAP_ADMIN_PASSWORD",
		"KEYCLOAK_REALM",
		auth.AccessTokenEnvVar, auth.RefreshTokenEnvVar,
	} {
		setEnv(t, key, "")
	}
}

// chdirTemp runs the test in a fresh directory so SetEnvToken's default
// ".env" writes land there instead of in the repository.
func chdirTemp(t *testing.T) {
	t.Helper()
	t.Chdir(t.TempDir())
}

// makeTestJWT builds an unsigned JWT carrying only an exp claim, which is all
// TokenValid inspects (it parses without verifying the signature).
func makeTestJWT(t *testing.T, exp time.Time) string {
	t.Helper()
	claims := jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(exp)}
	raw, err := jwt.NewWithClaims(jwt.SigningMethodNone, claims).SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	return raw
}

// tokenResponse builds a token-endpoint response body; an empty refreshToken
// omits the claim entirely.
func tokenResponse(accessToken, refreshToken string) map[string]interface{} {
	response := map[string]interface{}{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   3600,
	}
	if refreshToken != "" {
		response["refresh_token"] = refreshToken
	}
	return response
}

// writeJSONToken writes response as a JSON token-endpoint reply.
func writeJSONToken(t *testing.T, w http.ResponseWriter, response map[string]interface{}, status int) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	require.NoError(t, json.NewEncoder(w).Encode(response))
}

func TestServiceAccessTokenReturnsValidTokenUnchanged(t *testing.T) {
	chdirTemp(t)

	valid := makeTestJWT(t, time.Now().Add(time.Hour))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("token endpoint must not be called for a valid access token, got %s %s", r.Method, r.URL.Path)
	}))
	defer server.Close()

	token, err := auth.New().AccessToken(context.Background(), server.URL, valid, "unused-refresh")
	require.NoError(t, err)
	assert.Equal(t, valid, token)
}

func TestServiceAccessTokenRefreshesExpiredToken(t *testing.T) {
	chdirTemp(t)
	resetTokenEnv(t)

	expired := makeTestJWT(t, time.Now().Add(-time.Hour))
	fresh := makeTestJWT(t, time.Now().Add(time.Hour))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "refresh_token", r.FormValue("grant_type"))
		writeJSONToken(t, w, tokenResponse(fresh, "rotated-refresh"), http.StatusOK)
	}))
	defer server.Close()

	token, err := auth.New().AccessToken(context.Background(), server.URL, expired, "stale-refresh")
	require.NoError(t, err)
	assert.Equal(t, fresh, token)
	assert.Equal(t, fresh, os.Getenv(auth.AccessTokenEnvVar))
	assert.Equal(t, "rotated-refresh", os.Getenv(auth.RefreshTokenEnvVar))
}

func TestServiceAccessTokenAutoFetchesWhenNoTokens(t *testing.T) {
	chdirTemp(t)
	resetTokenEnv(t)
	setEnv(t, "KEYCLOAK_USERNAME", "env-user")
	setEnv(t, "KEYCLOAK_PASSWORD", "env-pass")

	fresh := makeTestJWT(t, time.Now().Add(time.Hour))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "password", r.FormValue("grant_type"))
		assert.Equal(t, "env-user", r.FormValue("username"))
		assert.Equal(t, "env-pass", r.FormValue("password"))
		assert.Contains(t, r.URL.Path, "/realms/master/")
		writeJSONToken(t, w, tokenResponse(fresh, "new-refresh"), http.StatusOK)
	}))
	defer server.Close()

	token, err := auth.New().AccessToken(context.Background(), server.URL, "", "")
	require.NoError(t, err)
	assert.Equal(t, fresh, token)
	assert.Equal(t, fresh, os.Getenv(auth.AccessTokenEnvVar))
	assert.Equal(t, "new-refresh", os.Getenv(auth.RefreshTokenEnvVar))

	envFile, err := os.ReadFile(".env")
	require.NoError(t, err)
	assert.Contains(t, string(envFile), auth.AccessTokenEnvVar+"="+fresh)
	assert.Contains(t, string(envFile), auth.RefreshTokenEnvVar+"=new-refresh")
}

func TestServiceAccessTokenFallsBackToPasswordGrantWhenRefreshFails(t *testing.T) {
	chdirTemp(t)
	resetTokenEnv(t)
	setEnv(t, "KEYCLOAK_USERNAME", "env-user")
	setEnv(t, "KEYCLOAK_PASSWORD", "env-pass")

	expired := makeTestJWT(t, time.Now().Add(-time.Hour))
	fresh := makeTestJWT(t, time.Now().Add(time.Hour))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.FormValue("grant_type") {
		case "refresh_token":
			writeJSONToken(t, w, map[string]interface{}{"error": "invalid_grant"}, http.StatusBadRequest)
		case "password":
			writeJSONToken(t, w, tokenResponse(fresh, "new-refresh"), http.StatusOK)
		default:
			t.Errorf("unexpected grant type %q", r.FormValue("grant_type"))
		}
	}))
	defer server.Close()

	token, err := auth.New().AccessToken(context.Background(), server.URL, expired, "exhausted-refresh")
	require.NoError(t, err)
	assert.Equal(t, fresh, token)
	assert.Equal(t, "new-refresh", os.Getenv(auth.RefreshTokenEnvVar))
}

func TestServiceAccessTokenFailsLoudWhenAutoFetchRejected(t *testing.T) {
	chdirTemp(t)
	resetTokenEnv(t)
	setEnv(t, "KEYCLOAK_USERNAME", "env-user")
	setEnv(t, "KEYCLOAK_PASSWORD", "wrong-pass")

	expired := makeTestJWT(t, time.Now().Add(-time.Hour))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSONToken(t, w, map[string]interface{}{"error": "invalid_grant"}, http.StatusUnauthorized)
	}))
	defer server.Close()

	_, err := auth.New().AccessToken(context.Background(), server.URL, expired, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "auto-fetch admin token")
	assert.Contains(t, err.Error(), "admin-token")
	assert.Contains(t, err.Error(), "password credentials token")
}

func TestServiceAccessTokenUsesBootstrapAdminEnvFallback(t *testing.T) {
	chdirTemp(t)
	resetTokenEnv(t)
	setEnv(t, "KC_BOOTSTRAP_ADMIN_USERNAME", "boot-user")
	setEnv(t, "KC_BOOTSTRAP_ADMIN_PASSWORD", "boot-pass")
	setEnv(t, "KEYCLOAK_REALM", "boot-realm")

	fresh := makeTestJWT(t, time.Now().Add(time.Hour))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "password", r.FormValue("grant_type"))
		assert.Equal(t, "boot-user", r.FormValue("username"))
		assert.Equal(t, "boot-pass", r.FormValue("password"))
		assert.Contains(t, r.URL.Path, "/realms/boot-realm/")
		writeJSONToken(t, w, tokenResponse(fresh, "new-refresh"), http.StatusOK)
	}))
	defer server.Close()

	token, err := auth.New().AccessToken(context.Background(), server.URL, "", "")
	require.NoError(t, err)
	assert.Equal(t, fresh, token)
}
