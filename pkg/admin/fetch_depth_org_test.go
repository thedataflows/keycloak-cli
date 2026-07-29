package admin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/admin"
)

// TestFetchDepthDescendsOrgGroupHierarchy pins that the depth traversal reaches
// an organization's nested groups by routing org-scoped descent through the
// scoped child-collection selection (FetchChildren's mechanism). The fake serves
// ONLY the org children paths, so the test passes only if the traversal uses
// them — the realm children path /groups/{group-id}/children (which Keycloak
// rejects for org groups) 404s here and would leave the grandchild unreachable.
func TestFetchDepthDescendsOrgGroupHierarchy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/organizations":
			writeJSON(t, w, []map[string]interface{}{{"id": "org-1", "name": "demo-org", "alias": "demo-org"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/organizations/org-1/groups":
			writeJSON(t, w, []map[string]interface{}{{"id": "g1", "name": "orggroup"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/organizations/org-1/groups/g1/children":
			writeJSON(t, w, []map[string]interface{}{{"id": "c1", "name": "orgchild"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/organizations/org-1/groups/c1/children":
			writeJSON(t, w, []map[string]interface{}{{"id": "gr1", "name": "orggrandchild"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/organizations/org-1/groups/gr1/children":
			writeJSON(t, w, []map[string]interface{}{})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Fetch(context.Background(), admin.FetchQuery{
		Realm: "demo", Resources: "organization", Depth: 3,
	})
	require.NoError(t, err)

	byName := map[string]struct {
		parentType string
		orgID      interface{}
	}{}
	for _, res := range report.Resources {
		if res.Type == "group" {
			byName[res.Data["name"].(string)] = struct {
				parentType string
				orgID      interface{}
			}{res.ParentType, res.Data["orgId"]}
		}
	}

	for _, name := range []string{"orggroup", "orgchild", "orggrandchild"} {
		got, ok := byName[name]
		require.True(t, ok, "depth traversal must reach %q", name)
		assert.Equal(t, "organization", got.parentType, "%q must keep the org scope marker for further descent", name)
	}
	assert.Equal(t, "org-1", byName["orggrandchild"].orgID, "orgId must propagate down the hierarchy")
}
