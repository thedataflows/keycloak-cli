package catalog

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEveryRegistryReadPathResolvesAgainstSpec is the guardrail that prevents
// ISSUE 0004 from recurring: a relationship kind whose ReadPath placeholder does
// not exist in the loaded spec is classified by nobody, so its edges are
// silently never fetched. Registry.ByPath is an exact lookup after
// normalization, so a placeholder rename in the spec (e.g. {client} -> {client-id})
// disables a kind with no error. Assert that every registered ReadPath still
// resolves to a GET path in the embedded spec.
func TestEveryRegistryReadPathResolvesAgainstSpec(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	specGetPaths := make(map[string]struct{})
	spec.ForEachOperation(func(path, method string, _ *v3.Operation, _ *v3.PathItem) {
		if method == http.MethodGet {
			specGetPaths[normalizeReadPath(path)] = struct{}{}
		}
	})

	for _, kind := range DefaultRegistry().Kinds() {
		normalized := normalizeReadPath(kind.ReadPath)
		if _, ok := specGetPaths[normalized]; !ok {
			t.Errorf("relationship kind %q read path %q does not resolve to any GET path in the spec", kind.Name, kind.ReadPath)
		}
	}
}

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

// TestClientRoleMappingKindsDiscovered pins ISSUE 0004: once the {client-id}
// placeholder matches the spec, the two client-scoped role-mapping kinds must be
// discovered as relationship patterns and must expose the owning client as a
// path parameter so the (parent x client) fetch iteration can render it.
func TestClientRoleMappingKindsDiscovered(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	patterns, err := spec.DiscoverRelationshipPatterns()
	require.NoError(t, err)

	placeholderMap, err := spec.PlaceholderToResourceType()
	require.NoError(t, err)

	byKind := make(map[string]RelationshipOperationPattern)
	for _, p := range patterns {
		byKind[p.Kind] = p
	}

	for _, tc := range []struct {
		kind       string
		parentType string
	}{
		{"user-client-role-mapping", "user"},
		{"group-client-role-mapping", "group"},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			pattern, ok := byKind[tc.kind]
			require.True(t, ok, "kind %q must be discovered from the spec", tc.kind)
			assert.Contains(t, pattern.PathParams, "client-id",
				"client owner must be a path parameter so the fetch can iterate it")

			types, err := pattern.ParentResourceTypes(placeholderMap)
			require.NoError(t, err)
			assert.Contains(t, types, tc.parentType)
			assert.Contains(t, types, "client")
		})
	}
}
