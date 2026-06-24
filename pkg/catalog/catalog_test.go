package catalog_test

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/catalog"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

func TestNewSpecLoadsDocument(t *testing.T) {
	spec, err := catalog.NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)
	require.NotNil(t, spec)

	schemas, err := spec.GetSchemas()
	require.NoError(t, err)
	require.Contains(t, schemas, "RealmRepresentation")
}

func TestValidateRelationshipOperations(t *testing.T) {
	spec, err := catalog.NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	rels := []manifest.RelationshipOperation{{Path: "test-realm/users/test-user/groups/test-group"}}
	require.NoError(t, catalog.ValidateRelationshipOperations(spec, rels))
	require.Equal(t, "{realm}/users/{user-id}/groups/{groupId}", rels[0].Template)
	require.Equal(t, "PUT", rels[0].Method)
}

func TestCollectRelationshipTemplates(t *testing.T) {
	spec, err := catalog.NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	templates, err := catalog.CollectRelationshipTemplates(spec)
	require.NoError(t, err)
	require.NotEmpty(t, templates)

	found := false
	for _, tmpl := range templates {
		if tmpl.Template == "{realm}/users/{user-id}/groups/{groupId}" && tmpl.Method == http.MethodPut {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestValidateRelationshipOperationsDelete(t *testing.T) {
	spec, err := catalog.NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	rels := []manifest.RelationshipOperation{{
		Path:   "test-realm/users/test-user/groups/test-group",
		Delete: true,
	}}

	require.NoError(t, catalog.ValidateRelationshipOperations(spec, rels))
	assert.Equal(t, http.MethodDelete, rels[0].Method)
	assert.True(t, rels[0].Delete)
	assert.Equal(t, map[string]string{
		"realm":   "test-realm",
		"user-id": "test-user",
		"groupId": "test-group",
	}, rels[0].PathParams)
}

func TestValidateRelationshipOperationsAssignsRealmScopeKinds(t *testing.T) {
	spec, err := catalog.NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	rels := []manifest.RelationshipOperation{{Path: "test-realm/default-default-client-scopes/scope-1"}, {Path: "test-realm/default-optional-client-scopes/scope-2"}}

	require.NoError(t, catalog.ValidateRelationshipOperations(spec, rels))
	assert.Equal(t, http.MethodPut, rels[0].Method)
	assert.Equal(t, "realm-default-client-scope", rels[0].Kind)
	assert.Equal(t, http.MethodPut, rels[1].Method)
	assert.Equal(t, "realm-optional-client-scope", rels[1].Kind)
}

func TestOperationContractIncludesQueryParams(t *testing.T) {
	spec, err := catalog.NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	contract, err := spec.OperationContract("/admin/realms/{realm}/users", http.MethodGet)
	require.NoError(t, err)
	assert.Equal(t, http.MethodGet, contract.Method)
	assert.NotEmpty(t, contract.Parameters)

	foundQuery := false
	for _, parameter := range contract.Parameters {
		if parameter.In == "query" && parameter.Name == "search" {
			foundQuery = true
			break
		}
	}
	assert.True(t, foundQuery)
}

func TestValidateResourceRejectsWrongTypes(t *testing.T) {
	spec, err := catalog.NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	err = spec.ValidateResource(manifest.Resource{
		Type:  "user",
		Realm: "demo",
		Data: map[string]interface{}{
			"username":      "alice",
			"enabled":       "yes",
			"emailVerified": true,
		},
	}, http.MethodPost)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "enabled")
}

func TestValidateOperationRequestRejectsWrongRelationshipBody(t *testing.T) {
	spec, err := catalog.NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	err = spec.ValidateOperationRequest("/admin/realms/{realm}/users/{user-id}/role-mappings/realm", http.MethodPost, catalog.RequestValidation{
		PathParams: map[string]string{"realm": "demo", "user-id": "user-1"},
		Body:       map[string]interface{}{"id": "role-1"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "request body")
}

func TestResourceContractsPreferRoleNameEndpoint(t *testing.T) {
	spec, err := catalog.NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	contracts, err := spec.ResourceContracts()
	require.NoError(t, err)

	roleContract, ok := contracts["role"]
	require.True(t, ok)
	putOperation, ok := roleContract.Operations[http.MethodPut]
	require.True(t, ok)
	assert.Contains(t, putOperation.Path, "/roles/{role-name}")
	assert.NotContains(t, putOperation.Path, "/roles-by-id/")

	deleteOperation, ok := roleContract.Operations[http.MethodDelete]
	require.True(t, ok)
	assert.Contains(t, deleteOperation.Path, "/roles/{role-name}")
	assert.NotContains(t, deleteOperation.Path, "/roles-by-id/")
}
