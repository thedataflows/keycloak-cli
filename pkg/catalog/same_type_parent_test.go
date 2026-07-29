package catalog_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/catalog"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

// A group nested under a group shares the "group" type with its parent, so the
// only thing distinguishing "parent binding" from "own identifier" on the shared
// {group-id} token is position. These tests pin ISSUE 0005: the parent binding is
// rendered into the create path and stripped from the body, while update/delete
// address the group by its own id even when a stray parent id is in the data.

func TestNestedGroupCreateStripsParentBinding(t *testing.T) {
	resolver := loadOrgGroupSpec(t).Resolver()

	createOp, err := resolver.ResolveResourceOperation("group", "group", http.MethodPost, catalog.OperationCollection)
	require.NoError(t, err)
	assert.Equal(t, "/admin/realms/{realm}/groups/{group-id}/children", createOp.Path)
	assert.Equal(t, []string{"groupId"}, resolver.ParentReferenceFieldNames("group", createOp),
		"the same-type parent binding must be a stripped parent-reference field")

	// No contamination of the realm paths.
	realmCollection, err := resolver.ResolveResourceOperation("group", "", http.MethodPost, catalog.OperationCollection)
	require.NoError(t, err)
	assert.Equal(t, "/admin/realms/{realm}/groups", realmCollection.Path)
	assert.Nil(t, resolver.ParentReferenceFieldNames("group", realmCollection))

	singleOp, err := resolver.ResolveResourceOperation("group", "", http.MethodPut, catalog.OperationSingle)
	require.NoError(t, err)
	assert.Equal(t, "/admin/realms/{realm}/groups/{group-id}", singleOp.Path)
	assert.Nil(t, resolver.ParentReferenceFieldNames("group", singleOp),
		"the single-group path addresses the group by its own id; nothing to strip")
}

// TestNestedGroupCreateRendersParentPath proves the parent binding is still
// load-bearing on create: the parent's id renders {group-id} in the children path.
func TestNestedGroupCreateRendersParentPath(t *testing.T) {
	resolver := loadOrgGroupSpec(t).Resolver()

	child := manifest.Resource{
		Type: "group", Realm: "demo", ParentType: "group",
		Data: map[string]interface{}{"name": "platform", "groupId": "parent-1"},
	}
	path, _, params, err := resolver.ResolveResourcePath(child, http.MethodPost, catalog.OperationCollection)
	require.NoError(t, err)
	assert.Equal(t, "parent-1", params["group-id"])
	assert.Equal(t, "/admin/realms/demo/groups/parent-1/children", path)
}

// TestNestedGroupUpdateDeleteAddressOwnID guards the destructive Gap 3 trap: even
// with the parent id present in Data, update and delete must resolve /groups/{id}
// to the child's own id, never the parent's.
func TestNestedGroupUpdateDeleteAddressOwnID(t *testing.T) {
	resolver := loadOrgGroupSpec(t).Resolver()

	child := manifest.Resource{
		Type: "group", Realm: "demo", ParentType: "group",
		Data: map[string]interface{}{"id": "child-1", "name": "platform", "groupId": "parent-1"},
	}
	for _, method := range []string{http.MethodPut, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			path, _, params, err := resolver.ResolveResourcePath(child, method, catalog.OperationSingle)
			require.NoError(t, err)
			assert.Equal(t, "child-1", params["group-id"], "must address the group by its own id")
			assert.Equal(t, "/admin/realms/demo/groups/child-1", path)
			assert.NotContains(t, path, "parent-1")
		})
	}
}

// TestNestedGroupReadEdgeUnchanged pins that the structural read edge is intact —
// the fix must not register the children path as a relationship kind.
func TestNestedGroupReadEdgeUnchanged(t *testing.T) {
	graph, err := loadOrgGroupSpec(t).BuildDownwardGraph()
	require.NoError(t, err)
	assert.Contains(t, childPaths(graph["group"], "group"),
		"/admin/realms/{realm}/groups/{group-id}/children",
		"realm group children must remain a structural child, not a relationship kind")
}
