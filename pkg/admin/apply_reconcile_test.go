package admin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

func TestReconcileRelationshipSetsNoChurnForPathIdentified(t *testing.T) {
	desired := []manifest.RelationshipOperation{
		{Kind: "user-group-membership", Path: "demo/users/alice/groups/devs", PathParams: map[string]string{"realm": "demo", "user-id": "alice", "groupId": "devs"}, Method: "PUT"},
	}
	actual := []manifest.RelationshipOperation{
		{Kind: "user-group-membership", Path: "demo/users/alice/groups/devs", PathParams: map[string]string{"realm": "demo", "user-id": "alice", "groupId": "devs"}, Method: "PUT", Data: []byte(`{"id":"devs","name":"devs"}`)},
	}

	toAdd, toRemove := reconcileRelationshipSets(desired, actual)
	assert.Empty(t, toAdd)
	assert.Empty(t, toRemove)
}

func TestBuildRelationshipDeleteOperationBulkPayload(t *testing.T) {
	actual := manifest.RelationshipOperation{
		Kind:       "user-realm-role-mapping",
		PathParams: map[string]string{"realm": "demo", "user-id": "alice"},
		Data:       []byte(`[{"name":"admin"}]`),
	}
	op, ok := buildRelationshipDeleteOperation(actual)
	assert.True(t, ok)
	assert.Equal(t, "DELETE", op.Method)
	assert.True(t, op.Delete)
	assert.JSONEq(t, `[{"name":"admin"}]`, string(op.Data))
}

func TestBuildRelationshipDeleteOrganizationMember(t *testing.T) {
	actual := manifest.RelationshipOperation{
		Kind:       "organization-member",
		PathParams: map[string]string{"realm": "demo", "org-id": "org-1"},
		Data:       []byte(`"user-1"`),
	}
	op, ok := buildRelationshipDeleteOperation(actual)
	assert.True(t, ok)
	assert.Equal(t, "DELETE", op.Method)
	assert.Equal(t, "demo/organizations/org-1/members/user-1", op.Path)
}

func TestBuildRelationshipDeleteOrganizationIdentityProvider(t *testing.T) {
	actual := manifest.RelationshipOperation{
		Kind:       "organization-identity-provider",
		PathParams: map[string]string{"realm": "demo", "org-id": "org-1"},
		Data:       []byte(`"github"`),
	}
	op, ok := buildRelationshipDeleteOperation(actual)
	assert.True(t, ok)
	assert.Equal(t, "DELETE", op.Method)
	assert.Equal(t, "demo/organizations/org-1/identity-providers/github", op.Path)
}

func TestReconcileRelationshipSetsNoChanges(t *testing.T) {
	desired := []manifest.RelationshipOperation{
		{Kind: "user-group-membership", Path: "demo/users/alice/groups/devs", PathParams: map[string]string{"realm": "demo", "user-id": "alice", "groupId": "devs"}, Method: "PUT"},
	}
	actual := []manifest.RelationshipOperation{
		{Kind: "user-group-membership", Path: "demo/users/alice/groups/devs", PathParams: map[string]string{"realm": "demo", "user-id": "alice", "groupId": "devs"}, Method: "PUT"},
	}

	toAdd, toRemove := reconcileRelationshipSets(desired, actual)
	assert.Empty(t, toAdd)
	assert.Empty(t, toRemove)
}

func TestReconcileRelationshipSetsAddsMissing(t *testing.T) {
	desired := []manifest.RelationshipOperation{
		{Kind: "user-group-membership", Path: "demo/users/alice/groups/devs", PathParams: map[string]string{"realm": "demo", "user-id": "alice", "groupId": "devs"}, Method: "PUT"},
	}

	toAdd, toRemove := reconcileRelationshipSets(desired, nil)
	assert.Len(t, toAdd, 1)
	assert.Empty(t, toRemove)
}

func TestReconcileRelationshipSetsRemovesUnexpected(t *testing.T) {
	actual := []manifest.RelationshipOperation{
		{Kind: "user-group-membership", Path: "demo/users/alice/groups/old", PathParams: map[string]string{"realm": "demo", "user-id": "alice", "groupId": "old"}, Method: "PUT"},
	}

	toAdd, toRemove := reconcileRelationshipSets(nil, actual)
	assert.Empty(t, toAdd)
	assert.Len(t, toRemove, 1)
	assert.Equal(t, "DELETE", toRemove[0].Method)
	assert.True(t, toRemove[0].Delete)
}

func TestReconcileRelationshipSetsBulkPayload(t *testing.T) {
	desired := []manifest.RelationshipOperation{
		{Kind: "user-realm-role-mapping", Path: "demo/users/alice/role-mappings/realm", PathParams: map[string]string{"realm": "demo", "user-id": "alice"}, Method: "POST", Data: []byte(`[{"name":"admin"}]`)},
	}
	actual := []manifest.RelationshipOperation{
		{Kind: "user-realm-role-mapping", Path: "demo/users/alice/role-mappings/realm", PathParams: map[string]string{"realm": "demo", "user-id": "alice"}, Method: "POST", Data: []byte(`[{"name":"user"}]`)},
	}

	toAdd, toRemove := reconcileRelationshipSets(desired, actual)
	assert.Len(t, toAdd, 1)
	assert.Len(t, toRemove, 1)
}

func TestBuildRelationshipDeleteOperation(t *testing.T) {
	actual := manifest.RelationshipOperation{
		Kind:       "user-group-membership",
		PathParams: map[string]string{"realm": "demo", "user-id": "alice", "groupId": "devs"},
		Template:   "{realm}/users/{user-id}/groups/{groupId}",
	}
	op, ok := buildRelationshipDeleteOperation(actual)
	assert.True(t, ok)
	assert.Equal(t, "DELETE", op.Method)
	assert.True(t, op.Delete)
	assert.Equal(t, "demo/users/alice/groups/devs", op.Path)
}

func TestRelationshipKeyCanonicalizesJSON(t *testing.T) {
	a := manifest.RelationshipOperation{Kind: "k", Path: "p", Data: []byte(`{"b":2,"a":1}`)}
	b := manifest.RelationshipOperation{Kind: "k", Path: "p", Data: []byte(`{"a":1,"b":2}`)}
	assert.Equal(t, relationshipKey(a, false), relationshipKey(b, false))
}

func TestRelationshipRealms(t *testing.T) {
	rels := []manifest.RelationshipOperation{
		{PathParams: map[string]string{"realm": "demo"}},
		{PathParams: map[string]string{"realm": "prod"}},
		{PathParams: map[string]string{"realm": "demo"}},
	}
	assert.Equal(t, []string{"demo", "prod"}, relationshipRealms(rels))
}
