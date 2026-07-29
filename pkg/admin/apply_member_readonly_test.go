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

// TestApplySkipsReadOnlyMember pins that a member resource (surfaced read-only by
// a depth fetch of org groups; its endpoint is GET-only) is skipped on apply
// rather than failing — there is no write path for it (ISSUE 0007). This keeps a
// fetch -> upload round-trip of org hierarchies from erroring.
func TestApplySkipsReadOnlyMember(t *testing.T) {
	var writes int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writes++
		}
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{{
		Type:       "member",
		Realm:      "demo",
		ParentType: "organization",
		Data:       map[string]interface{}{"id": "u1", "username": "alice", "orgId": "org-1", "groupId": "g1"},
	}}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "skipped", report.Results[0].Action)
	assert.Zero(t, report.Failed)
	assert.Zero(t, writes, "a read-only member must trigger no write requests")
}
