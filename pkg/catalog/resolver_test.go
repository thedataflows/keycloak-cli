package catalog

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

func TestPathParamsResolvesPathPlaceholders(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	tests := []struct {
		name     string
		resource manifest.Resource
		path     string
		want     map[string]string
	}{
		{
			name:     "realm-only path",
			resource: manifest.Resource{Type: "realm", Realm: "demo", Data: map[string]interface{}{"realm": "demo"}},
			path:     "/admin/realms/{realm}",
			want:     map[string]string{"realm": "demo"},
		},
		{
			name:     "user with username fallback",
			resource: manifest.Resource{Type: "user", Realm: "demo", Data: map[string]interface{}{"username": "alice"}},
			path:     "/admin/realms/{realm}/users/{user-id}",
			want:     map[string]string{"realm": "demo", "user-id": "alice"},
		},
		{
			name:     "user with id present",
			resource: manifest.Resource{Type: "user", Realm: "demo", Data: map[string]interface{}{"id": "user-1", "username": "alice"}},
			path:     "/admin/realms/{realm}/users/{user-id}",
			want:     map[string]string{"realm": "demo", "user-id": "user-1"},
		},
		{
			name:     "client with clientId",
			resource: manifest.Resource{Type: "client", Realm: "demo", Data: map[string]interface{}{"clientId": "app"}},
			path:     "/admin/realms/{realm}/clients/{client-uuid}",
			want:     map[string]string{"realm": "demo", "client-uuid": "app"},
		},
		{
			name:     "nested protocolmapper with parent id in data",
			resource: manifest.Resource{Type: "protocolmapper", Realm: "demo", Data: map[string]interface{}{"id": "mapper-1", "clientScopeId": "scope-1"}},
			path:     "/admin/realms/{realm}/client-scopes/{client-scope-id}/protocol-mappers/models/{id}",
			want:     map[string]string{"realm": "demo", "client-scope-id": "scope-1", "id": "mapper-1"},
		},
		{
			name:     "role with role-name path",
			resource: manifest.Resource{Type: "role", Realm: "demo", Data: map[string]interface{}{"name": "developer"}},
			path:     "/admin/realms/{realm}/roles/{role-name}",
			want:     map[string]string{"realm": "demo", "role-name": "developer"},
		},
		{
			name:     "empty path returns realm only",
			resource: manifest.Resource{Type: "group", Realm: "demo", Data: map[string]interface{}{"name": "devs"}},
			path:     "",
			want:     map[string]string{"realm": "demo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := OperationContract{Path: tt.path, Method: "GET"}
			got, err := spec.Resolver().PathParams(tt.resource, op)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPrimaryIdentifierParam(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/admin/realms/{realm}/users/{user-id}", "user-id"},
		{"/admin/realms/{realm}/clients/{client-uuid}", "client-uuid"},
		{"/admin/realms/{realm}/groups/{id}", "id"},
		{"/admin/realms/{realm}", ""},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			assert.Equal(t, tt.want, primaryIdentifierParam(tt.path))
		})
	}
}

func TestExtractPathParams(t *testing.T) {
	tests := []struct {
		path string
		want []string
	}{
		{"/admin/realms/{realm}/users/{user-id}", []string{"realm", "user-id"}},
		{"/admin/realms/{realm}/clients/{client-uuid}/roles/{role-name}", []string{"realm", "client-uuid", "role-name"}},
		{"/admin/realms/{realm}/groups/{id}/children", []string{"realm", "id"}},
		{"", nil},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			assert.Equal(t, tt.want, extractPathParams(tt.path))
		})
	}
}

func TestKebabToCamelCase(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"user-id", "userId"},
		{"client-uuid", "clientUuid"},
		{"client-scope-id", "clientScopeId"},
		{"role-name", "roleName"},
		{"realm", "realm"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.want, kebabToCamelCase(tt.in))
		})
	}
}

func TestPathParamsFallsBackToIdentifier(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	resource := manifest.Resource{
		Type:  "group",
		Realm: "demo",
		Data:  map[string]interface{}{"id": "group-1", "name": "devs"},
	}
	path := "/admin/realms/{realm}/groups/{group-id}"
	scope, err := spec.Resolver().PathParams(resource, OperationContract{Path: path, Method: "GET"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"realm": "demo", "group-id": "group-1"}, scope)
}

func TestParentReferenceFieldsMapsPlaceholderToParentID(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	tests := []struct {
		name   string
		path   string
		parent manifest.Resource
		want   map[string]string
	}{
		{
			name:   "protocolmapper under client scope",
			path:   "/admin/realms/{realm}/client-scopes/{client-scope-id}/protocol-mappers/models",
			parent: manifest.Resource{Type: "clientscope", Realm: "demo", Data: map[string]interface{}{"id": "scope-1", "name": "email"}},
			want:   map[string]string{"clientScopeId": "scope-1"},
		},
		{
			name:   "authz resource under client",
			path:   "/admin/realms/{realm}/clients/{client-uuid}/authz/resource-server/resource",
			parent: manifest.Resource{Type: "client", Realm: "demo", Data: map[string]interface{}{"id": "client-1", "clientId": "app"}},
			want:   map[string]string{"clientUuid": "client-1"},
		},
		{
			name:   "no matching placeholder",
			path:   "/admin/realms/{realm}/users",
			parent: manifest.Resource{Type: "client", Realm: "demo", Data: map[string]interface{}{"id": "client-1"}},
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := spec.Resolver().ParentReferenceFields(tt.path, tt.parent.Type, tt.parent)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPathParamsUsesCamelCaseField(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	resource := manifest.Resource{
		Type:  "protocolmapper",
		Realm: "demo",
		Data:  map[string]interface{}{"id": "mapper-1", "clientScopeId": "scope-1"},
	}
	path := "/admin/realms/{realm}/client-scopes/{client-scope-id}/protocol-mappers/models/{id}"
	scope, err := spec.Resolver().PathParams(resource, OperationContract{Path: path, Method: "GET"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"realm": "demo", "client-scope-id": "scope-1", "id": "mapper-1"}, scope)
}

func TestPathParamsUsesExactFieldName(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	resource := manifest.Resource{
		Type:  "user",
		Realm: "demo",
		Data:  map[string]interface{}{"user-id": "alice", "username": "alice-smith"},
	}
	path := "/admin/realms/{realm}/users/{user-id}"
	scope, err := spec.Resolver().PathParams(resource, OperationContract{Path: path, Method: "GET"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"realm": "demo", "user-id": "alice"}, scope)
}

func TestPathParamsAlwaysAddsPrimaryIdentifier(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	resource := manifest.Resource{
		Type:  "user",
		Realm: "demo",
		Data:  map[string]interface{}{"username": "alice"},
	}
	path := "/admin/realms/{realm}/users/{user-id}"
	scope, err := spec.Resolver().PathParams(resource, OperationContract{Path: path, Method: "GET"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"realm": "demo", "user-id": "alice"}, scope)
}

func TestPathParamsPrimaryIdentifierEmptyWhenNoData(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	resource := manifest.Resource{
		Type:  "user",
		Realm: "demo",
		Data:  map[string]interface{}{},
	}
	path := "/admin/realms/{realm}/users/{user-id}"
	scope, err := spec.Resolver().PathParams(resource, OperationContract{Path: path, Method: "GET"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"realm": "demo", "user-id": ""}, scope)
}

func TestResolveResourceOperationWithParentType(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	resolver := spec.Resolver()

	// protocolmapper under clientscope
	contract, err := resolver.ResolveResourceOperation("protocolmapper", "clientscope", "POST", OperationAny)
	require.NoError(t, err)
	assert.Contains(t, contract.Path, "client-scopes")
	assert.Contains(t, contract.Path, "protocol-mappers")

	// protocolmapper under client
	contract, err = resolver.ResolveResourceOperation("protocolmapper", "client", "POST", OperationAny)
	require.NoError(t, err)
	assert.Contains(t, contract.Path, "clients")
	assert.Contains(t, contract.Path, "protocol-mappers")

	// without ParentType falls back to best endpoint
	contract, err = resolver.ResolveResourceOperation("protocolmapper", "", "POST", OperationAny)
	require.NoError(t, err)
	assert.NotEmpty(t, contract.Path)
}

func TestResolveResourceOperationIdentityProviderMapper(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	contract, err := spec.Resolver().ResolveResourceOperation("identityprovidermapper", "identityprovider", "POST", OperationAny)
	require.NoError(t, err)
	assert.Contains(t, contract.Path, "identity-provider")
	assert.Contains(t, contract.Path, "mappers")
}

func TestResolveResourceOperationUnknownParent(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	// Should fall back to best endpoint when parent type doesn't match any endpoint
	contract, err := spec.Resolver().ResolveResourceOperation("protocolmapper", "unknownparent", "POST", OperationAny)
	require.NoError(t, err)
	assert.NotEmpty(t, contract.Path)
}

func TestResolveResourceOperationClientRole(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	contract, err := spec.Resolver().ResolveResourceOperation("role", "client", "POST", OperationAny)
	require.NoError(t, err)
	assert.Contains(t, contract.Path, "clients")
	assert.Contains(t, contract.Path, "roles")
}

func TestResolveResourceOperationBackwardCompatibility(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	// Should still work for single-location types
	contract, err := spec.Resolver().ResolveResourceOperation("client", "", "POST", OperationAny)
	require.NoError(t, err)
	assert.NotEmpty(t, contract.Path)

	contract, err = spec.Resolver().ResolveResourceOperation("user", "", "POST", OperationAny)
	require.NoError(t, err)
	assert.NotEmpty(t, contract.Path)
}

func TestResolveResourcePathWithParentType(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	resource := manifest.Resource{
		Type:       "protocolmapper",
		Realm:      "demo",
		ParentType: "clientscope",
		Data: map[string]interface{}{
			"name": "email",
		},
	}
	path, _, _, err := spec.Resolver().ResolveResourcePath(resource, "POST", OperationAny)
	require.NoError(t, err)
	assert.Contains(t, path, "client-scopes")
}

func TestResolveResourcePathWithoutParentType(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	resource := manifest.Resource{
		Type:  "client",
		Realm: "demo",
		Data: map[string]interface{}{
			"clientId": "my-app",
		},
	}
	path, _, _, err := spec.Resolver().ResolveResourcePath(resource, "POST", OperationAny)
	require.NoError(t, err)
	assert.Equal(t, "/admin/realms/demo/clients", path)
}

func TestResolveResourcePathParamsWithClientScopeParent(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	resource := manifest.Resource{
		Type:       "protocolmapper",
		Realm:      "demo",
		ParentType: "clientscope",
		Data: map[string]interface{}{
			"name":           "email",
			"protocolMapper": "oidc-usermodel-property-mapper",
			"clientScopeId":  "scope-1",
		},
	}
	path, contract, params, err := spec.Resolver().ResolveResourcePath(resource, "POST", OperationAny)
	require.NoError(t, err)
	assert.Contains(t, path, "client-scopes")
	assert.NotEmpty(t, contract.Path)

	// Verify scope params work with the resolved path
	assert.Equal(t, "scope-1", params["client-scope-id"])
}
