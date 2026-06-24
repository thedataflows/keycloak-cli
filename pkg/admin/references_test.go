package admin

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

func TestResolveReferencesIsIdempotent(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/authentication/flows":
			writeJSONForReferences(t, w, []map[string]interface{}{
				{"id": "f4509057-d632-42d8-9ce8-4a2f832a8620", "alias": "browser", "providerId": "basic-flow"},
			})
		case r.Method == http.MethodGet:
			writeJSONForReferences(t, w, []map[string]interface{}{})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	service := newServiceForReferencesTest(t, server.URL)

	client := manifest.Resource{
		Type:  "client",
		Realm: "demo",
		Data: map[string]interface{}{
			"id":       "client-1",
			"clientId": "fcc",
			"authenticationFlowBindingOverrides": map[string]interface{}{
				"browser": "f4509057-d632-42d8-9ce8-4a2f832a8620",
			},
		},
	}

	first, failures := service.resolveReferences(context.Background(), []string{"demo"}, []manifest.Resource{client})
	require.Empty(t, failures)
	require.Len(t, first, 1)
	assert.Equal(t, "authenticationflow", first[0].Type)

	combined := append([]manifest.Resource{client}, first...)
	second, failures := service.resolveReferences(context.Background(), []string{"demo"}, combined)
	require.Empty(t, failures)
	assert.Empty(t, second, "second resolution should not return already-fetched resources")
}

func TestResolveReferencesDeduplicatesResources(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/authentication/flows":
			writeJSONForReferences(t, w, []map[string]interface{}{
				{"id": "f4509057-d632-42d8-9ce8-4a2f832a8620", "alias": "browser", "providerId": "basic-flow"},
			})
		case r.Method == http.MethodGet:
			writeJSONForReferences(t, w, []map[string]interface{}{})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	service := newServiceForReferencesTest(t, server.URL)

	resources := []manifest.Resource{
		{
			Type:  "client",
			Realm: "demo",
			Data: map[string]interface{}{
				"id":       "client-1",
				"clientId": "fcc",
				"authenticationFlowBindingOverrides": map[string]interface{}{
					"browser": "f4509057-d632-42d8-9ce8-4a2f832a8620",
				},
			},
		},
		{
			Type:  "client",
			Realm: "demo",
			Data: map[string]interface{}{
				"id":       "client-2",
				"clientId": "other",
				"authenticationFlowBindingOverrides": map[string]interface{}{
					"browser": "f4509057-d632-42d8-9ce8-4a2f832a8620",
				},
			},
		},
	}

	resolved, failures := service.resolveReferences(context.Background(), []string{"demo"}, resources)
	require.Empty(t, failures)
	require.Len(t, resolved, 1)
	assert.Equal(t, "authenticationflow", resolved[0].Type)
}

func newServiceForReferencesTest(t *testing.T, baseURL string) *service {
	t.Helper()
	previousAccessToken, hadAccessToken := os.LookupEnv("KEYCLOAK_ACCESS_TOKEN")
	previousRefreshToken, hadRefreshToken := os.LookupEnv("KEYCLOAK_REFRESH_TOKEN")
	t.Cleanup(func() {
		if hadAccessToken {
			require.NoError(t, os.Setenv("KEYCLOAK_ACCESS_TOKEN", previousAccessToken))
		} else {
			require.NoError(t, os.Unsetenv("KEYCLOAK_ACCESS_TOKEN"))
		}
		if hadRefreshToken {
			require.NoError(t, os.Setenv("KEYCLOAK_REFRESH_TOKEN", previousRefreshToken))
		} else {
			require.NoError(t, os.Unsetenv("KEYCLOAK_REFRESH_TOKEN"))
		}
	})
	require.NoError(t, os.Setenv("KEYCLOAK_ACCESS_TOKEN", validAccessTokenForReferences()))
	require.NoError(t, os.Setenv("KEYCLOAK_REFRESH_TOKEN", ""))

	svc, err := New(Config{
		BaseURL:  baseURL,
		SpecPath: filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"),
		Timeout:  time.Second,
	})
	require.NoError(t, err)
	typed, ok := svc.(*service)
	require.True(t, ok, "expected *service")
	return typed
}

func validAccessTokenForReferences() string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"exp":%d}`, time.Now().Add(time.Hour).Unix())))
	return strings.Join([]string{header, payload, ""}, ".")
}

func writeJSONForReferences(t *testing.T, w http.ResponseWriter, payload interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	require.NoError(t, json.NewEncoder(w).Encode(payload))
}
