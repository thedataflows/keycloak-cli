package cmd

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/auth"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

func TestCompareCmdRun(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/users":
			writeCompareJSON(t, w, []map[string]interface{}{{"id": "user-1", "username": "alice", "enabled": true}})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	}))
	defer server.Close()

	setValidAdminTokens(t)
	manifestPath := filepath.Join(t.TempDir(), "user.json")
	require.NoError(t, os.WriteFile(manifestPath, []byte(`[{"type":"user","realm":"demo","data":{"username":"alice","enabled":true}}]`), 0o644))

	cmd := &CompareCmd{InputFiles: []string{manifestPath}, Realm: "demo", Format: "json"}
	cli := &CLI{Globals: Globals{KeycloakBaseURL: server.URL, SpecPath: "../keycloak-oapi/26.6.2.spec.json", Timeout: time.Second}}

	err := cmd.Run(nil, cli)
	assert.NoError(t, err)
}

func TestCompareCmdRunReturnsMismatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/users":
			writeCompareJSON(t, w, []map[string]interface{}{{"id": "user-1", "username": "alice", "enabled": false}})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	}))
	defer server.Close()

	setValidAdminTokens(t)
	outputPath := filepath.Join(t.TempDir(), "report.json")
	manifestPath := filepath.Join(t.TempDir(), "user.json")
	require.NoError(t, os.WriteFile(manifestPath, []byte(`[{"type":"user","realm":"demo","data":{"username":"alice","enabled":true}}]`), 0o644))

	cmd := &CompareCmd{InputFiles: []string{manifestPath}, Realm: "demo", Format: "json", Output: outputPath, Force: true}
	cli := &CLI{Globals: Globals{KeycloakBaseURL: server.URL, SpecPath: "../keycloak-oapi/26.6.2.spec.json", Timeout: time.Second}}

	err := cmd.Run(nil, cli)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "comparison failed")

	report, readErr := os.ReadFile(outputPath)
	require.NoError(t, readErr)
	assert.Contains(t, string(report), "mismatchedResources")
	assert.Contains(t, string(report), "\"match\": false")
}

func TestCompareCmdResolveRealm(t *testing.T) {
	cmd := &CompareCmd{}
	realm, err := cmd.resolveRealm(manifest.LoadResult{Resources: []manifest.Resource{{Type: "realm", Realm: "demo", Data: map[string]interface{}{"realm": "demo"}}}})
	require.NoError(t, err)
	assert.Equal(t, "demo", realm)

	_, err = cmd.resolveRealm(manifest.LoadResult{Resources: []manifest.Resource{
		{Type: "realm", Realm: "demo", Data: map[string]interface{}{"realm": "demo"}},
		{Type: "realm", Realm: "other", Data: map[string]interface{}{"realm": "other"}},
	}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "multiple realms")
}

func TestCompareFetchResources(t *testing.T) {
	resources := compareFetchResources(manifest.LoadResult{Resources: []manifest.Resource{
		{Type: "user", Realm: "demo", Data: map[string]interface{}{"username": "alice"}},
		{Type: "group", Realm: "demo", Data: map[string]interface{}{"name": "devs"}},
	}})
	assert.Equal(t, "group,user", resources)
	loaded := manifest.LoadResult{Resources: []manifest.Resource{{Type: "user", Realm: "demo", Data: map[string]interface{}{"username": "alice"}}}}
	loaded.Relationships = []manifest.RelationshipOperation{{Kind: "user-group-membership", Path: "demo/users/a/groups/b"}}
	result := compareFetchResources(loaded)
	assert.Contains(t, result, "user")
	assert.Contains(t, result, "group")
}

func setValidAdminTokens(t *testing.T) {
	t.Helper()
	previousAccessToken, hadAccessToken := os.LookupEnv(auth.AccessTokenEnvVar)
	previousRefreshToken, hadRefreshToken := os.LookupEnv(auth.RefreshTokenEnvVar)
	t.Cleanup(func() {
		if hadAccessToken {
			require.NoError(t, os.Setenv(auth.AccessTokenEnvVar, previousAccessToken))
		} else {
			require.NoError(t, os.Unsetenv(auth.AccessTokenEnvVar))
		}
		if hadRefreshToken {
			require.NoError(t, os.Setenv(auth.RefreshTokenEnvVar, previousRefreshToken))
		} else {
			require.NoError(t, os.Unsetenv(auth.RefreshTokenEnvVar))
		}
	})
	require.NoError(t, os.Setenv(auth.AccessTokenEnvVar, validCompareAccessToken()))
	require.NoError(t, os.Setenv(auth.RefreshTokenEnvVar, ""))
}

func validCompareAccessToken() string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"exp":4102444800}`))
	return strings.Join([]string{header, payload, ""}, ".")
}

func writeCompareJSON(t *testing.T, writer http.ResponseWriter, payload interface{}) {
	t.Helper()
	writer.Header().Set("Content-Type", "application/json")
	require.NoError(t, json.NewEncoder(writer).Encode(payload))
}
