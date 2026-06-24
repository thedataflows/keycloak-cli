package admin_test

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
	"github.com/thedataflows/keycloak-cli/pkg/admin"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

func TestNewRejectsMissingBaseURL(t *testing.T) {
	service, err := admin.New(admin.Config{SpecPath: filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json")})
	require.Error(t, err)
	assert.Nil(t, service)
}

func TestNewBuildsClient(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()

	service, err := admin.New(admin.Config{
		BaseURL:  server.URL,
		SpecPath: filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"),
	})
	require.NoError(t, err)
	require.NotNil(t, service)
	require.NotNil(t, service.Spec())

	schemas, err := service.Spec().GetSchemas()
	require.NoError(t, err)
	assert.Contains(t, schemas, "RealmRepresentation")
}

func TestApplyCreatesResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo":
			http.Error(w, "not found", http.StatusNotFound)
		case r.Method == http.MethodPost && r.URL.Path == "/admin/realms":
			w.WriteHeader(http.StatusCreated)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{{
		Type:  "realm",
		Realm: "demo",
		Data:  map[string]interface{}{"realm": "demo"},
	}}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "created", report.Results[0].Action)
	assert.Equal(t, http.StatusCreated, report.Results[0].Status)
	assert.Zero(t, report.Failed)
}

func TestApplyFailureUsesTypedErrorText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/users/alice":
			http.Error(w, "not found", http.StatusNotFound)
		case r.Method == http.MethodPost && r.URL.Path == "/admin/realms/demo/users":
			http.Error(w, "bad payload", http.StatusBadRequest)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{{
		Type:  "user",
		Realm: "demo",
		Data:  map[string]interface{}{"username": "alice"},
	}}, nil, admin.ApplyOptions{ContinueOnError: true})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Equal(t, http.StatusBadRequest, report.Results[0].Status)
	assert.Equal(t, "apply user: validation failure (400): bad payload\n", report.Results[0].Error)
	assert.NotContains(t, report.Results[0].Error, "HTTP 400")
}

func TestApplyRelationshipConflictHandledAsSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/admin/realms/demo/users/user-1/groups/group-1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
		http.Error(w, "already exists", http.StatusConflict)
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), nil, []manifest.RelationshipOperation{{
		Path: "demo/users/user-1/groups/group-1",
	}}, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "unchanged", report.Results[0].Action)
	assert.Equal(t, http.StatusOK, report.Results[0].Status)
	assert.Zero(t, report.Failed)
}

func TestFetchBuildsRealmScopedResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms":
			writeJSON(t, w, []map[string]interface{}{{"realm": "demo"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/users":
			writeJSON(t, w, []map[string]interface{}{{"username": "alice"}})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Fetch(context.Background(), admin.FetchQuery{Resources: "user"})
	require.NoError(t, err)
	require.Len(t, report.Resources, 1)
	assert.Equal(t, "demo", report.Resources[0].Realm)
	assert.Empty(t, report.Failures)
}

func TestApplyRejectsInvalidResourceBeforeNetwork(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	_, err := service.Apply(context.Background(), []manifest.Resource{{
		Type:  "user",
		Realm: "demo",
		Data: map[string]interface{}{
			"username": "alice",
			"enabled":  "yes",
		},
	}}, nil, admin.ApplyOptions{})
	require.Error(t, err)
	assert.Zero(t, requestCount)
}

func TestFetchIncludesRelationships(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms":
			writeJSON(t, w, []map[string]interface{}{{"realm": "demo"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/users":
			writeJSON(t, w, []map[string]interface{}{{"id": "user-1", "username": "alice"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/groups":
			writeJSON(t, w, []map[string]interface{}{{"id": "group-1", "name": "dev"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/roles":
			writeJSON(t, w, []map[string]interface{}{{"id": "role-1", "name": "developer"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/clients":
			writeJSON(t, w, []map[string]interface{}{{"id": "client-1", "clientId": "app"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/client-scopes":
			writeJSON(t, w, []map[string]interface{}{{"id": "scope-1", "name": "scope"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/organizations":
			writeJSON(t, w, []map[string]interface{}{{"id": "org-1", "name": "org"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/default-groups":
			writeJSON(t, w, []map[string]interface{}{{"id": "group-1", "name": "dev"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/default-default-client-scopes":
			writeJSON(t, w, []map[string]interface{}{{"id": "scope-3", "name": "realm-default-scope"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/default-optional-client-scopes":
			writeJSON(t, w, []map[string]interface{}{{"id": "scope-4", "name": "realm-optional-scope"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/users/user-1/groups":
			writeJSON(t, w, []map[string]interface{}{{"id": "group-1", "name": "dev"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/users/user-1/role-mappings/realm":
			writeJSON(t, w, []map[string]interface{}{{"id": "role-1", "name": "developer"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/users/user-1/federated-identity":
			writeJSON(t, w, []map[string]interface{}{{"identityProvider": "github", "userId": "external-1", "userName": "alice"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/groups/group-1/role-mappings/realm":
			writeJSON(t, w, []map[string]interface{}{{"id": "role-1", "name": "developer"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/roles-by-id/role-1/composites":
			writeJSON(t, w, []map[string]interface{}{{"id": "role-2", "name": "employee"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/clients/client-1/default-client-scopes":
			writeJSON(t, w, []map[string]interface{}{{"id": "scope-1", "name": "scope"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/clients/client-1/optional-client-scopes":
			writeJSON(t, w, []map[string]interface{}{{"id": "scope-2", "name": "optional-scope"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/client-scopes/scope-1/scope-mappings/realm":
			writeJSON(t, w, []map[string]interface{}{{"id": "role-1", "name": "developer"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/client-scopes/scope-1/scope-mappings/clients/client-1":
			writeJSON(t, w, []map[string]interface{}{{"id": "role-3", "name": "app-user", "clientRole": true, "containerId": "client-1"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/organizations/org-1/members":
			writeJSON(t, w, []map[string]interface{}{{"id": "user-1", "username": "alice"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/organizations/org-1/identity-providers":
			writeJSON(t, w, []map[string]interface{}{{"alias": "github", "providerId": "oidc"}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Fetch(context.Background(), admin.FetchQuery{Resources: "user", IncludeRelationships: true})
	require.NoError(t, err)
	assert.NotEmpty(t, report.Relationships)
	kinds := relationshipKinds(report.Relationships)
	assert.Subset(t, kinds, []string{
		"user-group-membership",
		"user-realm-role-mapping",
		"user-federated-identity",
		"default-group-membership",
		"realm-default-client-scope",
		"realm-optional-client-scope",
		"group-realm-role-mapping",
		"role-composite-mapping",
		"client-default-scope",
		"client-optional-scope",
		"client-scope-realm-role-mapping",
		"client-scope-client-role-mapping",
		"organization-member",
		"organization-identity-provider",
	})
}

func newServiceForTest(t *testing.T, baseURL string) admin.Service {
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
	require.NoError(t, os.Setenv("KEYCLOAK_ACCESS_TOKEN", validAccessToken()))
	require.NoError(t, os.Setenv("KEYCLOAK_REFRESH_TOKEN", ""))

	service, err := admin.New(admin.Config{
		BaseURL:  baseURL,
		SpecPath: filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"),
		Timeout:  time.Second,
	})
	require.NoError(t, err)
	return service
}

func validAccessToken() string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"exp":%d}`, time.Now().Add(time.Hour).Unix())))
	return strings.Join([]string{header, payload, ""}, ".")
}

func writeJSON(t *testing.T, writer http.ResponseWriter, payload interface{}) {
	t.Helper()
	writer.Header().Set("Content-Type", "application/json")
	require.NoError(t, json.NewEncoder(writer).Encode(payload))
}
func TestApplyUpdateWhenResourceExists(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/groups/group-1":
			writeJSON(t, w, map[string]interface{}{"id": "group-1", "name": "devs"})
		case r.Method == http.MethodPut && r.URL.Path == "/admin/realms/demo/groups/group-1":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	}))
	defer server.Close()
	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{{
		Type:  "group",
		Realm: "demo",
		Data:  map[string]interface{}{"id": "group-1", "name": "devs"},
	}}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "updated", report.Results[0].Action)
	assert.Equal(t, http.StatusNoContent, report.Results[0].Status)
	assert.Zero(t, report.Failed)
}
func TestApplyDeleteNotFoundIsIdempotent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/users/alice":
			http.Error(w, "not found", http.StatusNotFound)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	}))
	defer server.Close()
	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{{
		Type:   "user",
		Realm:  "demo",
		Data:   map[string]interface{}{"username": "alice"},
		Delete: true,
	}}, nil, admin.ApplyOptions{Delete: true})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "not-found", report.Results[0].Action)
	assert.Equal(t, http.StatusOK, report.Results[0].Status)
	assert.Zero(t, report.Failed)
}
func TestApplyPopulatesIDMapFromCreateResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/groups/client-group-id":
			http.Error(w, "not found", http.StatusNotFound)
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/groups":
			writeJSON(t, w, []map[string]interface{}{{"id": "server-group-id", "name": "devs"}})
		case r.Method == http.MethodPost && r.URL.Path == "/admin/realms/demo/groups":
			w.Header().Set("Location", "/admin/realms/demo/groups/server-group-id")
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodPut && r.URL.Path == "/admin/realms/demo/users/alice/groups/server-group-id":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	}))
	defer server.Close()
	service := newServiceForTest(t, server.URL)
	resources := []manifest.Resource{{
		Type:  "group",
		Realm: "demo",
		Data:  map[string]interface{}{"id": "client-group-id", "name": "devs"},
	}}
	relationships := []manifest.RelationshipOperation{{
		Kind:   "user-group-membership",
		Path:   "demo/users/alice/groups/client-group-id",
		Method: "PUT",
	}}
	report, err := service.Apply(context.Background(), resources, relationships, admin.ApplyOptions{})
	require.NoError(t, err)
	assert.Zero(t, report.Failed)
	var relResult *admin.ApplyResult
	for i := range report.Results {
		if report.Results[i].Resource == "relationship" {
			relResult = &report.Results[i]
			break
		}
	}
	require.NotNil(t, relResult)
	assert.Equal(t, "applied", relResult.Action)
}
func TestApplyResolvesRoleClientUuidFromManifest(t *testing.T) {
	sourceClientUUID := "11111111-1111-1111-1111-111111111111"
	targetClientUUID := "22222222-2222-2222-2222-222222222222"
	var receivedPath string
	var receivedBody map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/clients/"+sourceClientUUID:
			http.Error(w, "not found", http.StatusNotFound)
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/clients":
			writeJSON(t, w, []map[string]interface{}{{"id": targetClientUUID, "clientId": "app"}})
		case r.Method == http.MethodPost && r.URL.Path == "/admin/realms/demo/clients":
			w.Header().Set("Location", "/admin/realms/demo/clients/"+targetClientUUID)
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/clients/"+targetClientUUID+"/roles/admin":
			http.Error(w, "not found", http.StatusNotFound)
		case r.Method == http.MethodPost && r.URL.Path == "/admin/realms/demo/clients/"+targetClientUUID+"/roles":
			receivedPath = r.URL.Path
			require.NoError(t, json.NewDecoder(r.Body).Decode(&receivedBody))
			w.WriteHeader(http.StatusCreated)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{
		{
			Type:  "client",
			Realm: "demo",
			Data:  map[string]interface{}{"id": sourceClientUUID, "clientId": "app"},
		},
		{
			Type:       "role",
			Realm:      "demo",
			ParentType: "client",
			Data:       map[string]interface{}{"name": "admin", "clientUuid": "app"},
		},
	}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Zero(t, report.Failed)
	assert.Equal(t, "/admin/realms/demo/clients/"+targetClientUUID+"/roles", receivedPath)
	assert.NotContains(t, receivedBody, "clientUuid")
	assert.NotContains(t, receivedBody, "realm")
}

func TestApplyResolvesRoleClientUuidFromServer(t *testing.T) {
	clientUUID := "22222222-2222-2222-2222-222222222222"
	var receivedPath string
	var receivedBody map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/clients":
			writeJSON(t, w, []map[string]interface{}{{"id": clientUUID, "clientId": "app"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/clients/"+clientUUID+"/roles/admin":
			http.Error(w, "not found", http.StatusNotFound)
		case r.Method == http.MethodPost && r.URL.Path == "/admin/realms/demo/clients/"+clientUUID+"/roles":
			receivedPath = r.URL.Path
			require.NoError(t, json.NewDecoder(r.Body).Decode(&receivedBody))
			w.WriteHeader(http.StatusCreated)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{{
		Type:       "role",
		Realm:      "demo",
		ParentType: "client",
		Data:       map[string]interface{}{"name": "admin", "clientUuid": "app"},
	}}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Zero(t, report.Failed)
	assert.Equal(t, "/admin/realms/demo/clients/"+clientUUID+"/roles", receivedPath)
	assert.NotContains(t, receivedBody, "clientUuid")
}

func TestApplySkipsOrganizationWhenDisabled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/organizations/org-1":
			http.Error(w, "not found", http.StatusNotFound)
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/organizations":
			writeJSON(t, w, []map[string]interface{}{})
		case r.Method == http.MethodPost && r.URL.Path == "/admin/realms/demo/organizations":
			http.Error(w, `{"error":"Organizations not enabled"}`, http.StatusBadRequest)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{{
		Type:  "organization",
		Realm: "demo",
		Data:  map[string]interface{}{"id": "org-1", "name": "Org One"},
	}}, nil, admin.ApplyOptions{ContinueOnError: true})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "skipped", report.Results[0].Action)
	assert.Equal(t, http.StatusOK, report.Results[0].Status)
	assert.Zero(t, report.Failed)
}

func TestApplySingletonSkipsCreateAndGoesToUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/client-policies/policies":
			writeJSON(t, w, map[string]interface{}{"policies": []interface{}{}})
		case r.Method == http.MethodPut && r.URL.Path == "/admin/realms/demo/client-policies/policies":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	}))
	defer server.Close()
	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{{
		Type:  "clientpolicies",
		Realm: "demo",
		Data:  map[string]interface{}{"policies": []interface{}{}},
	}}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "updated", report.Results[0].Action)
	assert.Equal(t, http.StatusNoContent, report.Results[0].Status)
	assert.Zero(t, report.Failed)
}
func TestApplyRemapsInlineIDsInResourceData(t *testing.T) {
	var receivedBody map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/admin/realms/demo/authentication/flows/"):
			http.Error(w, "not found", http.StatusNotFound)
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/authentication/flows":
			writeJSON(t, w, []map[string]interface{}{})
		case r.Method == http.MethodPost && r.URL.Path == "/admin/realms/demo/authentication/flows":
			w.Header().Set("Location", "/admin/realms/demo/authentication/flows/f4509057-d632-42d8-9ce8-4a2f832a8620")
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/admin/realms/demo/clients/"):
			http.Error(w, "not found", http.StatusNotFound)
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/clients":
			writeJSON(t, w, []map[string]interface{}{})
		case r.Method == http.MethodPost && r.URL.Path == "/admin/realms/demo/clients":
			require.NoError(t, json.NewDecoder(r.Body).Decode(&receivedBody))
			w.Header().Set("Location", "/admin/realms/demo/clients/eb51dd4e-d7bd-40ec-90dc-430df5275c0a")
			w.WriteHeader(http.StatusCreated)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	}))
	defer server.Close()
	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{
		{
			Type:  "authenticationflow",
			Realm: "demo",
			Data:  map[string]interface{}{"id": "client-flow-id", "alias": "browser", "providerId": "basic-flow"},
		},
		{
			Type:  "client",
			Realm: "demo",
			Data: map[string]interface{}{
				"clientId": "fcc",
				"name":     "FCC",
				"authenticationFlowBindingOverrides": map[string]interface{}{
					"browser": "client-flow-id",
				},
			},
		},
	}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Zero(t, report.Failed)

	overrides, ok := receivedBody["authenticationFlowBindingOverrides"].(map[string]interface{})
	require.True(t, ok, "expected authenticationFlowBindingOverrides to be a map, got %T", receivedBody["authenticationFlowBindingOverrides"])
	assert.Equal(t, "f4509057-d632-42d8-9ce8-4a2f832a8620", overrides["browser"])
}

func TestApplyInlineIDRemappingIsIdempotent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/admin/realms/demo/authentication/flows/"):
			http.Error(w, "not found", http.StatusNotFound)
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/authentication/flows":
			writeJSON(t, w, []map[string]interface{}{
				{"id": "f4509057-d632-42d8-9ce8-4a2f832a8620", "alias": "browser", "providerId": "basic-flow"},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/admin/realms/demo/authentication/flows":
			w.Header().Set("Location", "/admin/realms/demo/authentication/flows/f4509057-d632-42d8-9ce8-4a2f832a8620")
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/admin/realms/demo/clients/"):
			http.Error(w, "not found", http.StatusNotFound)
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/clients":
			writeJSON(t, w, []map[string]interface{}{
				{"id": "eb51dd4e-d7bd-40ec-90dc-430df5275c0a", "clientId": "fcc"},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/admin/realms/demo/clients":
			w.Header().Set("Location", "/admin/realms/demo/clients/eb51dd4e-d7bd-40ec-90dc-430df5275c0a")
			w.WriteHeader(http.StatusConflict)
		case r.Method == http.MethodPut && r.URL.Path == "/admin/realms/demo/clients/eb51dd4e-d7bd-40ec-90dc-430df5275c0a":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	}))
	defer server.Close()
	service := newServiceForTest(t, server.URL)
	resources := []manifest.Resource{
		{
			Type:  "authenticationflow",
			Realm: "demo",
			Data:  map[string]interface{}{"id": "client-flow-id", "alias": "browser", "providerId": "basic-flow"},
		},
		{
			Type:  "client",
			Realm: "demo",
			Data: map[string]interface{}{
				"clientId": "fcc",
				"name":     "FCC",
				"authenticationFlowBindingOverrides": map[string]interface{}{
					"browser": "client-flow-id",
				},
			},
		},
	}

	first, err := service.Apply(context.Background(), resources, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Zero(t, first.Failed)

	second, err := service.Apply(context.Background(), resources, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Zero(t, second.Failed)
}

func TestApplySingletonDeleteReportsNotSupported(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/client-policies/policies":
			writeJSON(t, w, map[string]interface{}{"policies": []interface{}{}})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	}))
	defer server.Close()
	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{{
		Type:   "clientpolicies",
		Realm:  "demo",
		Data:   map[string]interface{}{"policies": []interface{}{}},
		Delete: true,
	}}, nil, admin.ApplyOptions{Delete: true})
	require.NoError(t, err, "Apply error: %v", err)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "not-supported", report.Results[0].Action)
	assert.Equal(t, http.StatusOK, report.Results[0].Status)
	assert.Zero(t, report.Failed)
}
