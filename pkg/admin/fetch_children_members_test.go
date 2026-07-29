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

// TestFetchChildrenReadsOrgGroupMembers pins ISSUE 0007: FetchChildren reads the
// members of an org group (a GET-only scoped collection) with one addressed GET,
// keyed to that group. Members are tagged member/organization, carry the org
// scope (orgId), and carry the immediate parent group id (groupId) — which is
// safe here (a member is never recursed as a group parent, so the 0006 collision
// does not apply).
func TestFetchChildrenReadsOrgGroupMembers(t *testing.T) {
	const org = "org-1"
	var gets []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gets = append(gets, r.URL.Path)
		if r.URL.Path == "/admin/realms/demo/organizations/"+org+"/groups/g1/members" {
			writeJSON(t, w, []map[string]interface{}{{"id": "user-1", "username": "alice"}})
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	orgGroup := manifest.Resource{Type: "group", Realm: "demo", ParentType: "organization",
		Data: map[string]interface{}{"id": "g1", "orgId": org, "name": "orggroup"}}

	report, err := service.FetchChildren(context.Background(), orgGroup, "member", admin.ChildFetchQuery{})
	require.NoError(t, err)
	require.Empty(t, report.Failures)

	assert.Equal(t, []string{"/admin/realms/demo/organizations/" + org + "/groups/g1/members"}, gets,
		"exactly one addressed GET to the group's members")
	require.Len(t, report.Resources, 1)
	m := report.Resources[0]
	assert.Equal(t, "member", m.Type)
	assert.Equal(t, "organization", m.ParentType)
	assert.Equal(t, "alice", m.Data["username"])
	assert.Equal(t, org, m.Data["orgId"], "org scope propagated")
	assert.Equal(t, "g1", m.Data["groupId"], "member keyed to its group")
}
