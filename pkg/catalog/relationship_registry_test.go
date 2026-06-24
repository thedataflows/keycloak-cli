package catalog

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistryLookupByName(t *testing.T) {
	r := NewRegistry()
	r.Register(RelationshipKind{Name: "test-kind", ReadPath: "/test", WriteTemplate: "{realm}/test"})

	kind, ok := r.ByName("test-kind")
	require.True(t, ok)
	assert.Equal(t, "test-kind", kind.Name)

	_, ok = r.ByName("missing")
	assert.False(t, ok)
}

func TestRegistryLookupByPath(t *testing.T) {
	r := NewRegistry()
	r.Register(RelationshipKind{Name: "user-group", ReadPath: "/users/{user-id}/groups"})

	tests := []string{
		"/admin/realms/{realm}/users/{user-id}/groups",
		"users/{user-id}/groups",
		"{realm}/users/{user-id}/groups",
	}
	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			kind, ok := r.ByPath(path)
			require.True(t, ok, "path %s should match", path)
			assert.Equal(t, "user-group", kind.Name)
		})
	}

	_, ok := r.ByPath("/admin/realms/{realm}/users/{user-id}/roles")
	assert.False(t, ok)
}

func TestRegistryLookupByWriteTemplate(t *testing.T) {
	r := NewRegistry()
	r.Register(RelationshipKind{Name: "role-composite", WriteTemplate: "{realm}/roles-by-id/{role-id}/composites"})

	kind, ok := r.ByWriteTemplate("{realm}/roles-by-id/{role-id}/composites")
	require.True(t, ok)
	assert.Equal(t, "role-composite", kind.Name)
}

func TestDefaultRegistryContainsKnownKinds(t *testing.T) {
	r := DefaultRegistry()

	for _, name := range []string{
		"user-group-membership",
		"user-realm-role-mapping",
		"role-composite-mapping",
		"client-scope-realm-role-mapping",
		"client-scope-client-role-mapping",
		"organization-member",
	} {
		t.Run(name, func(t *testing.T) {
			kind, ok := r.ByName(name)
			require.True(t, ok, "expected %q in default registry", name)
			assert.NotEmpty(t, kind.ReadPath)
			assert.NotEmpty(t, kind.WriteTemplate)
			assert.NotEmpty(t, kind.ParamTypes)
		})
	}
}

func TestDefaultRegistryParamTypes(t *testing.T) {
	r := DefaultRegistry()

	paramTypes := r.ParamTypes("client-default-scope")
	require.NotNil(t, paramTypes)
	assert.Equal(t, "client", paramTypes["client-uuid"])
	assert.Equal(t, "clientscope", paramTypes["clientScopeId"])
}

func TestApplyRelationshipOverrides(t *testing.T) {
	r := NewRegistry()
	err := ApplyRelationshipOverrides(r, []RelationshipOverride{
		{
			Name:          "custom-relationship",
			ResourceA:     "user",
			ResourceB:     "group",
			ReadPath:      "/users/{user-id}/custom-groups",
			WriteTemplate: "{realm}/users/{user-id}/custom-groups/{groupId}",
			WriteMethod:   "PUT",
			ItemParamName: "id",
			ParamTypes:    map[string]string{"user-id": "user", "groupId": "group"},
		},
	})
	require.NoError(t, err)

	kind, ok := r.ByName("custom-relationship")
	require.True(t, ok)
	assert.Equal(t, "user", kind.ResourceA)
	assert.Equal(t, "group", kind.ResourceB)
	assert.Equal(t, "PUT", kind.WriteMethod)

	_, ok = r.ByPath("/admin/realms/{realm}/users/{user-id}/custom-groups")
	assert.True(t, ok)
}

func TestApplyRelationshipOverridesDisable(t *testing.T) {
	r := NewRegistry()
	r.Register(RelationshipKind{Name: "remove-me", ReadPath: "/to-remove", WriteTemplate: "{realm}/to-remove"})

	err := ApplyRelationshipOverrides(r, []RelationshipOverride{
		{Name: "remove-me", Disabled: true},
	})
	require.NoError(t, err)

	_, ok := r.ByName("remove-me")
	assert.False(t, ok)
}

func TestLoadRelationshipOverrides(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "overrides.yaml")
	data := []byte(`
overrides:
  - name: custom-relationship
    resourceA: user
    resourceB: group
    readPath: /users/{user-id}/custom-groups
    writeTemplate: "{realm}/users/{user-id}/custom-groups/{groupId}"
    writeMethod: PUT
    paramTypes:
      user-id: user
      groupId: group
`)
	require.NoError(t, os.WriteFile(path, data, 0o644))

	overrides, err := LoadRelationshipOverrides(path)
	require.NoError(t, err)
	require.Len(t, overrides, 1)
	assert.Equal(t, "custom-relationship", overrides[0].Name)
	assert.Equal(t, "user", overrides[0].ParamTypes["user-id"])
}

func TestApplyRelationshipOverridesRequiresName(t *testing.T) {
	r := NewRegistry()
	err := ApplyRelationshipOverrides(r, []RelationshipOverride{{ReadPath: "/x", WriteTemplate: "{realm}/x"}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name")
}

func TestApplyRelationshipOverridesRequiresReadPath(t *testing.T) {
	r := NewRegistry()
	err := ApplyRelationshipOverrides(r, []RelationshipOverride{{Name: "x", WriteTemplate: "{realm}/x"}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "readPath")
}

func TestApplyRelationshipOverridesRequiresWriteTemplate(t *testing.T) {
	r := NewRegistry()
	err := ApplyRelationshipOverrides(r, []RelationshipOverride{{Name: "x", ReadPath: "/x"}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "writeTemplate")
}
