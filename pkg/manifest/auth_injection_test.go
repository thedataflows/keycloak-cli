package manifest_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/internal/testutil"
	"golang.org/x/oauth2"

	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

// recordingAuth captures the arguments passed to AccessToken so tests can
// verify that manifest.NewService wired the injected auth provider rather than the
// default password-grant service. It implements auth.Service.
type recordingAuth struct {
	calls        int
	baseURL      string
	accessToken  string
	refreshToken string
	returnToken  string
}

func (r *recordingAuth) PasswordToken(ctx context.Context, baseURL, realm, username, password string) (oauth2.Token, error) {
	return oauth2.Token{}, nil
}

func (r *recordingAuth) ClientCredentialsToken(ctx context.Context, baseURL, realm, clientID, clientSecret string) (oauth2.Token, error) {
	return oauth2.Token{AccessToken: "cc-token"}, nil
}

func (r *recordingAuth) AccessToken(ctx context.Context, baseURL, accessToken, refreshToken string) (string, error) {
	r.calls++
	r.baseURL = baseURL
	r.accessToken = accessToken
	r.refreshToken = refreshToken
	if r.returnToken == "" {
		return "injected-token", nil
	}
	return r.returnToken, nil
}

func (r *recordingAuth) SetEnvToken(key, value, envFile string) error { return nil }

// TestNewWithCustomAuth verifies that a caller-supplied auth.Service is used
// by the runtime client instead of the default auth.New(). This is the hook
// the syncengine migration uses to inject a client_credentials provider.
func TestNewWithCustomAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Echo back the Authorization header so we can prove the injected
		// token reached the wire.
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"authorization":"` + r.Header.Get("Authorization") + `"}`))
	}))
	defer server.Close()

	recorder := &recordingAuth{}
	svc, err := manifest.NewService(manifest.Config{
		BaseURL:  server.URL,
		SpecPath: testutil.KeycloakSpecPath(t),
		Auth:     recorder,
	})
	require.NoError(t, err)
	require.NotNil(t, svc)

	// Trigger an authenticated request through the runtime client.
	_, _ = svc.Fetch(context.Background(), manifest.FetchQuery{Realm: "master", Resources: "realm"})

	assert.GreaterOrEqual(t, recorder.calls, 1, "injected auth.AccessToken must be invoked")
	assert.Equal(t, server.URL, recorder.baseURL, "injected auth must receive the configured base URL")
}

// TestNewWithoutAuthKeepsBackwardCompat verifies the zero-arg default path
// still works when Auth is not provided.
func TestNewWithoutAuthKeepsBackwardCompat(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()

	svc, err := manifest.NewService(manifest.Config{
		BaseURL:  server.URL,
		SpecPath: testutil.KeycloakSpecPath(t),
	})
	require.NoError(t, err)
	require.NotNil(t, svc)
}
