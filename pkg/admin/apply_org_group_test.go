package admin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/admin"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

// Org-scoped groups share resource type "group" with realm groups and are told
// apart only by ParentType. These tests pin the endpoint each variant writes to:
// Keycloak rejects realm-path writes against an org-scoped group with
// "Cannot manage organization related group via non Organization API."

type recordedRequest struct {
	method string
	path   string
}

// orgGroupServer serves the parent organization plus an existence probe, records
// every request, and returns the supplied status for writes. existing controls
// whether the group already exists.
func orgGroupServer(t *testing.T, recorded *[]recordedRequest, existing bool) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*recorded = append(*recorded, recordedRequest{method: r.Method, path: r.URL.Path})
		orgGroup := map[string]interface{}{"id": "group-1", "name": "acme-engineering", "path": "/acme-engineering"}
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/organizations":
			writeJSON(t, w, []map[string]interface{}{{"id": "org-1", "alias": "acme", "name": "acme"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/organizations/org-1":
			writeJSON(t, w, map[string]interface{}{"id": "org-1", "alias": "acme", "name": "acme"})
		// Existence probe: the resolver looks the group up by path before writing.
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/groups/group-by-path/"),
			r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/groups/group-1"):
			if !existing {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			writeJSON(t, w, orgGroup)
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/groups"):
			if !existing {
				writeJSON(t, w, []map[string]interface{}{})
				return
			}
			writeJSON(t, w, []map[string]interface{}{orgGroup})
		case r.Method == http.MethodPost:
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodPut, r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
}

func writePaths(recorded []recordedRequest) []string {
	paths := make([]string, 0, len(recorded))
	for _, req := range recorded {
		switch req.method {
		case http.MethodPost, http.MethodPut, http.MethodDelete:
			paths = append(paths, req.method+" "+req.path)
		}
	}
	return paths
}

func TestApplyOrgScopedGroupCreatesViaOrganizationEndpoint(t *testing.T) {
	var recorded []recordedRequest
	server := orgGroupServer(t, &recorded, false)
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{{
		Type:       "group",
		Realm:      "demo",
		ParentType: "organization",
		Data:       map[string]interface{}{"name": "acme-engineering", "orgId": "org-1"},
	}}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Zero(t, report.Failed)
	assert.Equal(t, "created", report.Results[0].Action)

	assert.Contains(t, writePaths(recorded), "POST /admin/realms/demo/organizations/org-1/groups")
	assert.NotContains(t, writePaths(recorded), "POST /admin/realms/demo/groups")
}

func TestApplyOrgScopedGroupUpdatesViaOrganizationEndpoint(t *testing.T) {
	var recorded []recordedRequest
	server := orgGroupServer(t, &recorded, true)
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{{
		Type:       "group",
		Realm:      "demo",
		ParentType: "organization",
		Data:       map[string]interface{}{"id": "group-1", "name": "acme-engineering", "path": "/acme-engineering", "orgId": "org-1"},
	}}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Zero(t, report.Failed)
	assert.Equal(t, "updated", report.Results[0].Action)

	assert.Contains(t, writePaths(recorded), "PUT /admin/realms/demo/organizations/org-1/groups/group-1")
	assert.NotContains(t, writePaths(recorded), "PUT /admin/realms/demo/groups/group-1")
}

func TestApplyOrgScopedGroupDeletesViaOrganizationEndpoint(t *testing.T) {
	var recorded []recordedRequest
	server := orgGroupServer(t, &recorded, true)
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{{
		Type:       "group",
		Realm:      "demo",
		ParentType: "organization",
		Delete:     true,
		Data:       map[string]interface{}{"id": "group-1", "name": "acme-engineering", "path": "/acme-engineering", "orgId": "org-1"},
	}}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Zero(t, report.Failed)

	assert.Contains(t, writePaths(recorded), "DELETE /admin/realms/demo/organizations/org-1/groups/group-1")
	assert.NotContains(t, writePaths(recorded), "DELETE /admin/realms/demo/groups/group-1")
}

func TestApplyRealmGroupStillUsesRealmEndpoint(t *testing.T) {
	var recorded []recordedRequest
	server := orgGroupServer(t, &recorded, false)
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{{
		Type:  "group",
		Realm: "demo",
		Data:  map[string]interface{}{"name": "realm-engineering"},
	}}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Zero(t, report.Failed)

	assert.Contains(t, writePaths(recorded), "POST /admin/realms/demo/groups")
	for _, path := range writePaths(recorded) {
		assert.NotContains(t, path, "/organizations/", "realm groups must never touch the org-scoped endpoint")
	}
}
