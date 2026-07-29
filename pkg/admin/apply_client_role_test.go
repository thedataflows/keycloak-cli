package admin_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/admin"
	"github.com/thedataflows/keycloak-cli/pkg/catalog"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

// clientRoleServer serves the parent client, an existence probe, and role writes.
// existing controls whether the client role already exists.
func clientRoleServer(t *testing.T, recorded *[]recordedRequest, postBody *map[string]interface{}, existing bool) *httptest.Server {
	t.Helper()
	const uuid = demoClientUUID
	role := map[string]interface{}{"id": "role-1", "name": "r1", "clientRole": true, "containerId": uuid}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*recorded = append(*recorded, recordedRequest{method: r.Method, path: r.URL.Path})
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/clients/"+uuid:
			writeJSON(t, w, map[string]interface{}{"id": uuid, "clientId": "demo-a-client-1"})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/clients/"+uuid+"/roles/r1":
			if !existing {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			writeJSON(t, w, role)
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/clients/"+uuid+"/roles"):
			if !existing {
				writeJSON(t, w, []map[string]interface{}{})
				return
			}
			writeJSON(t, w, []map[string]interface{}{role})
		case r.Method == http.MethodPost && r.URL.Path == "/admin/realms/demo/clients/"+uuid+"/roles":
			if postBody != nil {
				body, _ := io.ReadAll(r.Body)
				require.NoError(t, json.Unmarshal(body, postBody))
			}
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodPut, r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
}

func clientRoleResource(delete bool) manifest.Resource {
	return manifest.Resource{
		Type:       "role",
		Realm:      "demo",
		ParentType: "client",
		Delete:     delete,
		Data:       map[string]interface{}{"name": "r1", "clientUuid": demoClientUUID},
	}
}

// TestApplyClientRoleCreatesViaClientEndpoint proves ISSUE 0003's write
// direction: a role with ParentType client creates via POST
// /clients/{client-uuid}/roles, with clientUuid resolved from the parent and
// stripped from the body (validated against the embedded spec).
func TestApplyClientRoleCreatesViaClientEndpoint(t *testing.T) {
	var recorded []recordedRequest
	var body map[string]interface{}
	server := clientRoleServer(t, &recorded, &body, false)
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{clientRoleResource(false)}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Zero(t, report.Failed)
	assert.Equal(t, "created", report.Results[0].Action)

	assert.Contains(t, writePaths(recorded), "POST /admin/realms/demo/clients/"+demoClientUUID+"/roles")
	assert.NotContains(t, writePaths(recorded), "POST /admin/realms/demo/roles")

	require.NotNil(t, body)
	assert.NotContains(t, body, "clientUuid", "parent binding must be stripped from the body")
	assert.Equal(t, "r1", body["name"])

	spec, err := catalog.NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)
	require.NoError(t, spec.ValidateOperationRequest(
		"/admin/realms/{realm}/clients/{client-uuid}/roles", http.MethodPost,
		catalog.RequestValidation{
			PathParams: map[string]string{"realm": "demo", "client-uuid": demoClientUUID},
			Body:       body,
		}))
}

func TestApplyClientRoleUpdatesViaClientEndpoint(t *testing.T) {
	var recorded []recordedRequest
	server := clientRoleServer(t, &recorded, nil, true)
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{clientRoleResource(false)}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	assert.Zero(t, report.Failed)

	assert.Contains(t, writePaths(recorded), "PUT /admin/realms/demo/clients/"+demoClientUUID+"/roles/r1")
	for _, p := range writePaths(recorded) {
		assert.NotContains(t, p, "/admin/realms/demo/roles/", "client role must never touch the realm role endpoint")
	}
}

func TestApplyClientRoleDeletesViaClientEndpoint(t *testing.T) {
	var recorded []recordedRequest
	server := clientRoleServer(t, &recorded, nil, true)
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{clientRoleResource(true)}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	assert.Zero(t, report.Failed)

	assert.Contains(t, writePaths(recorded), "DELETE /admin/realms/demo/clients/"+demoClientUUID+"/roles/r1")
}
