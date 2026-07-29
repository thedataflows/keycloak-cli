package admin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/admin"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

// TestFetchChildrenWalksOrgGroupHierarchy walks org -> child -> grandchild via
// repeated FetchChildren, feeding each returned child straight back in as the
// next parent. Recursion must work out of the box: correct path per level,
// ParentType=organization preserved, orgId propagated, and no groupId that would
// collide and re-fetch the parent's children (ISSUE 0006).
func TestFetchChildrenWalksOrgGroupHierarchy(t *testing.T) {
	const org = "org-1"
	var gets []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gets = append(gets, r.URL.Path)
		switch r.URL.Path {
		case "/admin/realms/demo/organizations/" + org + "/groups/parent-1/children":
			writeJSON(t, w, []map[string]interface{}{{"id": "child-1", "name": "child"}})
		case "/admin/realms/demo/organizations/" + org + "/groups/child-1/children":
			writeJSON(t, w, []map[string]interface{}{{"id": "grand-1", "name": "grandchild"}})
		case "/admin/realms/demo/organizations/" + org + "/groups/grand-1/children":
			writeJSON(t, w, []map[string]interface{}{})
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	// A top-level org group as ISSUE 0002 produces it: ParentType=organization + orgId.
	parent := manifest.Resource{Type: "group", Realm: "demo", ParentType: "organization",
		Data: map[string]interface{}{"id": "parent-1", "orgId": org, "name": "parent"}}

	names := []string{}
	cur := parent
	for {
		rep, err := service.FetchChildren(context.Background(), cur, "group", admin.ChildFetchQuery{})
		require.NoError(t, err)
		require.Empty(t, rep.Failures)
		if len(rep.Resources) == 0 {
			break
		}
		require.Len(t, rep.Resources, 1)
		child := rep.Resources[0]
		assert.Equal(t, "group", child.Type)
		assert.Equal(t, "organization", child.ParentType, "org scope marker must propagate for recursion")
		assert.Equal(t, org, child.Data["orgId"], "orgId must be propagated down the subtree")
		_, hasGroupID := child.Data["groupId"]
		assert.False(t, hasGroupID, "must not inject a colliding parent groupId")
		names = append(names, child.Data["name"].(string))
		cur = child
	}

	assert.Equal(t, []string{"child", "grandchild"}, names, "recursion must reach the grandchild")
	assert.Equal(t, []string{
		"/admin/realms/demo/organizations/" + org + "/groups/parent-1/children",
		"/admin/realms/demo/organizations/" + org + "/groups/child-1/children",
		"/admin/realms/demo/organizations/" + org + "/groups/grand-1/children",
	}, gets, "each level addresses the group by its own id")
}
