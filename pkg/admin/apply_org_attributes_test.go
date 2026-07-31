package admin_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/admin"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

// orgAttributesServer records the body of the organization create/update write
// so the test can assert whether attributes reach Keycloak. existing controls
// whether the organization is already present (drives create vs update).
func orgAttributesServer(t *testing.T, capturedBody *map[string]interface{}, existing bool) *httptest.Server {
	t.Helper()
	orgList := []map[string]interface{}{{"id": "org-1", "alias": "acme", "name": "acme"}}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.TrimRight(r.URL.Path, "/") == "/admin/realms/demo/organizations":
			if !existing {
				writeJSON(t, w, []map[string]interface{}{})
				return
			}
			writeJSON(t, w, orgList)
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/organizations/org-1":
			if !existing {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			writeJSON(t, w, orgList[0])
		case r.Method == http.MethodPost && r.URL.Path == "/admin/realms/demo/organizations":
			body, _ := io.ReadAll(r.Body)
			require.NoError(t, json.Unmarshal(body, capturedBody))
			w.Header().Set("Location", r.URL.Path+"/org-1")
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodPut && r.URL.Path == "/admin/realms/demo/organizations/org-1":
			body, _ := io.ReadAll(r.Body)
			require.NoError(t, json.Unmarshal(body, capturedBody))
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
}

func orgResourceWithAttributes() manifest.Resource {
	return manifest.Resource{
		Type:  "organization",
		Realm: "demo",
		Data: map[string]interface{}{
			"name":        "acme",
			"alias":       "acme",
			"enabled":     true,
			"description": "Acme Corp",
			"attributes": map[string]interface{}{
				"department": []interface{}{"ENGINEERING"},
				"costcenter": []interface{}{"OPS-42"},
			},
		},
	}
}

func TestApplyOrganizationCreatePersistsAttributes(t *testing.T) {
	var body map[string]interface{}
	server := orgAttributesServer(t, &body, false)
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{orgResourceWithAttributes()}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "created", report.Results[0].Action)

	require.Contains(t, body, "attributes", "create body must carry attributes")
	attrs, ok := body["attributes"].(map[string]interface{})
	require.True(t, ok, "attributes must be an object, got %T", body["attributes"])
	assert.Equal(t, []interface{}{"ENGINEERING"}, attrs["department"])
	assert.Equal(t, []interface{}{"OPS-42"}, attrs["costcenter"])
}

func TestApplyOrganizationUpdatePersistsAttributes(t *testing.T) {
	var body map[string]interface{}
	server := orgAttributesServer(t, &body, true)
	defer server.Close()

	resource := orgResourceWithAttributes()
	resource.Data["id"] = "org-1"

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{resource}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "updated", report.Results[0].Action)

	require.Contains(t, body, "attributes", "update body must carry attributes")
	attrs, ok := body["attributes"].(map[string]interface{})
	require.True(t, ok, "attributes must be an object, got %T", body["attributes"])
	assert.Equal(t, []interface{}{"ENGINEERING"}, attrs["department"])
}
