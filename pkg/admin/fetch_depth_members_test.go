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

// TestFetchDepthSurfacesOrgGroupMembers pins the CLI exposure for ISSUE 0007: the
// depth traversal, when descending an org group, also reads that group's members
// (a GET-only scoped collection) and emits them as member resources keyed to the
// group. So `fetch organization --depth N` surfaces members at every level.
func TestFetchDepthSurfacesOrgGroupMembers(t *testing.T) {
	const org = "org-1"
	// The SAME user is a member at every level (matches the live fixture:
	// demo-a-user-1 on the org group, its child, and its grandchild). Each
	// membership is a distinct edge and must not collapse under resource dedup.
	members := map[string][]map[string]interface{}{
		"g1":  {{"id": "u1", "username": "alice"}},
		"c1":  {{"id": "u1", "username": "alice"}},
		"gr1": {{"id": "u1", "username": "alice"}},
	}
	children := map[string][]map[string]interface{}{
		"g1":  {{"id": "c1", "name": "orgchild"}},
		"c1":  {{"id": "gr1", "name": "orggrandchild"}},
		"gr1": {},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		base := "/admin/realms/demo/organizations/" + org
		switch {
		case r.URL.Path == "/admin/realms/demo/organizations":
			writeJSON(t, w, []map[string]interface{}{{"id": org, "name": "demo-org", "alias": "demo-org"}})
		case r.URL.Path == base+"/groups":
			writeJSON(t, w, []map[string]interface{}{{"id": "g1", "name": "orggroup"}})
		default:
			for gid, kids := range children {
				if r.URL.Path == base+"/groups/"+gid+"/children" {
					writeJSON(t, w, kids)
					return
				}
			}
			for gid, mem := range members {
				if r.URL.Path == base+"/groups/"+gid+"/members" {
					writeJSON(t, w, mem)
					return
				}
			}
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Fetch(context.Background(), admin.FetchQuery{
		Realm: "demo", Resources: "organization", Depth: 4,
	})
	require.NoError(t, err)

	// The same user's membership at each level is a distinct edge: collect the
	// set of groups it is keyed to. All three must survive.
	groups := map[string]bool{}
	for _, res := range report.Resources {
		if res.Type == "member" {
			assert.Equal(t, "alice", res.Data["username"])
			assert.Equal(t, "organization", res.ParentType)
			assert.Equal(t, org, res.Data["orgId"])
			groups[res.Data["groupId"].(string)] = true
		}
	}
	assert.Equal(t, map[string]bool{"g1": true, "c1": true, "gr1": true}, groups,
		"the same user's membership must surface at every level, keyed to each group")
}
