package catalog_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ScopedChildCollection selects a child collection path by matching the parent's
// placeholder chain, so it can reach the two-parent org children path that the
// deduped downward graph collapses. These cases pin all three parent shapes plus
// the arbitrary-depth org marker.
func TestScopedChildCollection(t *testing.T) {
	spec := loadOrgGroupSpec(t)

	tests := []struct {
		name             string
		parentType       string
		parentParentType string
		childType        string
		wantPath         string
		wantInherited    []string
	}{
		{"client role", "client", "", "role",
			"/admin/realms/{realm}/clients/{client-uuid}/roles", nil},
		{"realm group top", "group", "", "group",
			"/admin/realms/{realm}/groups/{group-id}/children", nil},
		{"realm subgroup", "group", "group", "group",
			"/admin/realms/{realm}/groups/{group-id}/children", nil},
		{"org group", "group", "organization", "group",
			"/admin/realms/{realm}/organizations/{org-id}/groups/{group-id}/children", []string{"orgId"}},
		// GET-only read collection (org group members has no POST): the selector
		// must still resolve it so the read channel can reach it (ISSUE 0007).
		{"org group members", "group", "organization", "member",
			"/admin/realms/{realm}/organizations/{org-id}/groups/{group-id}/members", []string{"orgId"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, inherited, ok := spec.ScopedChildCollection(tt.parentType, tt.parentParentType, tt.childType)
			require.True(t, ok)
			assert.Equal(t, tt.wantPath, path)
			assert.Equal(t, tt.wantInherited, inherited)
		})
	}

	_, _, ok := spec.ScopedChildCollection("client", "organization", "role")
	assert.False(t, ok, "no org-scoped client role collection exists")
}
