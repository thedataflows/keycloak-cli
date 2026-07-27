package catalog_test

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/catalog"
)

// Realm groups and org-scoped groups are both resource type "group", so the only
// thing keeping them apart is parent-type disambiguation. These tests pin that
// behavior: a regression here silently routes org-scoped writes at the realm
// endpoint, which Keycloak rejects with "Cannot manage organization related group
// via non Organization API."
func TestResolveGroupOperationDisambiguatesByParentType(t *testing.T) {
	spec := loadOrgGroupSpec(t)

	tests := []struct {
		name       string
		parentType string
		method     string
		shape      catalog.OperationShape
		wantPath   string
	}{
		{
			name:       "organization parent resolves the org-scoped collection",
			parentType: "organization",
			method:     http.MethodPost,
			shape:      catalog.OperationCollection,
			wantPath:   "/admin/realms/{realm}/organizations/{org-id}/groups",
		},
		{
			name:       "organization parent resolves the org-scoped collection on GET",
			parentType: "organization",
			method:     http.MethodGet,
			shape:      catalog.OperationCollection,
			wantPath:   "/admin/realms/{realm}/organizations/{org-id}/groups",
		},
		{
			name:       "organization parent resolves the org-scoped single resource on PUT",
			parentType: "organization",
			method:     http.MethodPut,
			shape:      catalog.OperationSingle,
			wantPath:   "/admin/realms/{realm}/organizations/{org-id}/groups/{group-id}",
		},
		{
			name:       "organization parent resolves the org-scoped single resource on DELETE",
			parentType: "organization",
			method:     http.MethodDelete,
			shape:      catalog.OperationSingle,
			wantPath:   "/admin/realms/{realm}/organizations/{org-id}/groups/{group-id}",
		},
		{
			name:       "no parent resolves the realm collection",
			parentType: "",
			method:     http.MethodPost,
			shape:      catalog.OperationCollection,
			wantPath:   "/admin/realms/{realm}/groups",
		},
		{
			name:       "group parent resolves the realm children collection",
			parentType: "group",
			method:     http.MethodPost,
			shape:      catalog.OperationCollection,
			wantPath:   "/admin/realms/{realm}/groups/{group-id}/children",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contract, err := spec.Resolver().ResolveResourceOperation("group", tt.parentType, tt.method, tt.shape)
			require.NoError(t, err)
			assert.Equal(t, tt.wantPath, contract.Path)
		})
	}
}

func TestDownwardGraphExposesOrgScopedGroups(t *testing.T) {
	spec := loadOrgGroupSpec(t)

	graph, err := spec.BuildDownwardGraph()
	require.NoError(t, err)

	assert.Contains(t, childPaths(graph["organization"], "group"),
		"/admin/realms/{realm}/organizations/{org-id}/groups",
		"organization must expose its own groups as a structural child")

	assert.Contains(t, childPaths(graph["group"], "group"),
		"/admin/realms/{realm}/groups/{group-id}/children",
		"realm group children must stay on the realm path")
}

func loadOrgGroupSpec(t *testing.T) *catalog.Spec {
	t.Helper()
	spec, err := catalog.NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)
	return spec
}

func childPaths(children []catalog.DownwardChild, childType string) []string {
	paths := make([]string, 0, len(children))
	for _, child := range children {
		if child.ChildType == childType {
			paths = append(paths, child.Path)
		}
	}
	return paths
}
