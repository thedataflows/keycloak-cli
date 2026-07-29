package catalog

import (
	"net/http"
	"path/filepath"
	"testing"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuiltinOrgGroupKinds pins ISSUE 0008: the organization-group membership and
// containment kinds ship in the built-in registry (no consumer override), with
// the exact override definitions, and their read/write/delete paths validate
// against the embedded spec.
func TestBuiltinOrgGroupKinds(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	// method -> set of normalized spec paths
	specPaths := map[string]map[string]struct{}{}
	spec.ForEachOperation(func(path, method string, _ *v3.Operation, _ *v3.PathItem) {
		if specPaths[method] == nil {
			specPaths[method] = map[string]struct{}{}
		}
		specPaths[method][normalizeReadPath(path)] = struct{}{}
	})
	resolves := func(method, template string) bool {
		_, ok := specPaths[method][normalizeReadPath(template)]
		return ok
	}

	reg := DefaultRegistry()

	member, ok := reg.ByName("organization-group-member")
	require.True(t, ok, "organization-group-member must be a built-in kind")
	assert.Equal(t, "group", member.ResourceA)
	assert.Equal(t, "user", member.ResourceB)
	assert.Equal(t, "id", member.ItemParamName)
	assert.Empty(t, member.PayloadField)
	assert.Equal(t, map[string]string{"org-id": "organization", "group-id": "group", "userId": "user"}, member.ParamTypes)
	assert.True(t, resolves(http.MethodGet, member.ReadPath), "member read path must be a spec GET")
	assert.True(t, resolves(http.MethodPut, member.WriteTemplate), "member write path must be a spec PUT")
	assert.True(t, resolves(http.MethodDelete, member.DeleteTemplate), "member delete path must be a spec DELETE")

	child, ok := reg.ByName("organization-group-child")
	require.True(t, ok, "organization-group-child must be a built-in kind")
	assert.Equal(t, "group", child.ResourceA)
	assert.Equal(t, "group", child.ResourceB)
	assert.Empty(t, child.ItemParamName, "child hands back the whole item")
	assert.Empty(t, child.PayloadField, "child hands back the whole item")
	assert.Equal(t, map[string]string{"org-id": "organization", "group-id": "group"}, child.ParamTypes)
	assert.True(t, resolves(http.MethodGet, child.ReadPath), "child read path must be a spec GET")
	assert.True(t, resolves(http.MethodPost, child.WriteTemplate), "child write path must be a spec POST")

	// Both resolve by path too.
	_, ok = reg.ByPath("/organizations/{org-id}/groups/{group-id}/members")
	assert.True(t, ok)
	_, ok = reg.ByPath("/organizations/{org-id}/groups/{group-id}/children")
	assert.True(t, ok)
}
