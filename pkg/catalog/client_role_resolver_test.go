package catalog_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/catalog"
)

// Client-scoped roles are resource type "role" told apart from realm roles by
// ParentType. These tests pin the endpoints ISSUE 0003's FetchChildren and the
// apply write direction depend on, and record the one place the resolver's bare
// GET-collection choice diverges from the structural graph.
func TestClientScopedRoleResolution(t *testing.T) {
	resolver := loadOrgGroupSpec(t).Resolver()

	tests := []struct {
		name       string
		parentType string
		method     string
		shape      catalog.OperationShape
		wantPath   string
		wantStrip  []string
	}{
		{"client parent create", "client", http.MethodPost, catalog.OperationCollection,
			"/admin/realms/{realm}/clients/{client-uuid}/roles", []string{"clientUuid"}},
		{"client parent update", "client", http.MethodPut, catalog.OperationSingle,
			"/admin/realms/{realm}/clients/{client-uuid}/roles/{role-name}", []string{"clientUuid"}},
		{"client parent delete", "client", http.MethodDelete, catalog.OperationSingle,
			"/admin/realms/{realm}/clients/{client-uuid}/roles/{role-name}", []string{"clientUuid"}},
		{"no parent create stays realm", "", http.MethodPost, catalog.OperationCollection,
			"/admin/realms/{realm}/roles", nil},
		{"no parent update stays realm", "", http.MethodPut, catalog.OperationSingle,
			"/admin/realms/{realm}/roles/{role-name}", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, err := resolver.ResolveResourceOperation("role", tt.parentType, tt.method, tt.shape)
			require.NoError(t, err)
			assert.Equal(t, tt.wantPath, op.Path)
			assert.Equal(t, tt.wantStrip, resolver.ParentReferenceFieldNames("role", op))
		})
	}
}

// TestClientScopedRoleGetCollectionDivergence records the reason FetchChildren
// resolves its path via BuildDownwardGraph rather than ResolveResourceOperation:
// a bare GET collection for role-under-client resolves to the composites path,
// because role-under-client matches more than one GET collection and
// operationPriority prefers composites. The structural graph, by contrast,
// correctly maps client -> role to the plain /roles collection.
func TestClientScopedRoleGetCollectionDivergence(t *testing.T) {
	spec := loadOrgGroupSpec(t)

	op, err := spec.Resolver().ResolveResourceOperation("role", "client", http.MethodGet, catalog.OperationCollection)
	require.NoError(t, err)
	assert.Equal(t, "/admin/realms/{realm}/clients/{client-uuid}/roles/{role-name}/composites", op.Path,
		"documented: bare GET collection prefers composites — do not use it for a parent-scoped role fetch")

	graph, err := spec.BuildDownwardGraph()
	require.NoError(t, err)
	assert.Contains(t, childPaths(graph["client"], "role"),
		"/admin/realms/{realm}/clients/{client-uuid}/roles",
		"the structural graph maps client -> role to the plain roles collection")
}
