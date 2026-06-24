package catalog_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/catalog"
	"github.com/thedataflows/keycloak-cli/pkg/catalog/internal"
)

//go:embed testdata/minimal-spec.json
var minimalSpecJSON []byte

func loadMinimalSpec(t *testing.T) *catalog.Spec {
	t.Helper()
	loaded, err := internal.NewSpecFromBytes(minimalSpecJSON)
	require.NoError(t, err)
	return catalog.WrapSpec(loaded)
}

func TestMinimalSpecLoads(t *testing.T) {
	spec := loadMinimalSpec(t)
	require.NotNil(t, spec)

	schemas, err := spec.GetSchemas()
	require.NoError(t, err)
	assert.Contains(t, schemas, "RealmRepresentation")
	assert.Contains(t, schemas, "ClientDTO")
}

func TestMinimalSpecDiscoversClientDTOAsClient(t *testing.T) {
	spec := loadMinimalSpec(t)

	contracts, err := spec.ResourceContracts()
	require.NoError(t, err)

	clientContract, ok := contracts["client"]
	require.True(t, ok, "expected 'client' resource type to be discovered from ClientDTO")
	assert.Equal(t, "ClientDTO", clientContract.SchemaName)
}

func TestMinimalSpecDiscoversUserOrganizationsRelationship(t *testing.T) {
	spec := loadMinimalSpec(t)

	patterns, err := spec.DiscoverRelationshipPatterns()
	require.NoError(t, err)

	var found *catalog.RelationshipOperationPattern
	for i := range patterns {
		if patterns[i].Path == "/admin/realms/{realm}/users/{user-id}/organizations" {
			found = &patterns[i]
			break
		}
	}
	require.NotNil(t, found, "expected user/organizations relationship to be discovered")
	assert.Equal(t, "user", found.ResourceA)
	assert.Equal(t, "organization", found.ResourceB)
}

func TestMinimalSpecResourceIdentity(t *testing.T) {
	spec := loadMinimalSpec(t)

	identities, err := spec.ResourceIdentities()
	require.NoError(t, err)

	user, ok := identities["user"]
	require.True(t, ok)
	assert.Equal(t, "username", user.PrimaryKey)
	assert.Equal(t, "user-id", user.IDParam)

	client, ok := identities["client"]
	require.True(t, ok)
	assert.Equal(t, "clientId", client.PrimaryKey)
	assert.Equal(t, "client-uuid", client.IDParam)
}

func TestMinimalSpecVolatileFields(t *testing.T) {
	spec := loadMinimalSpec(t)

	volatile, err := spec.VolatileFields("user")
	require.NoError(t, err)
	assert.Contains(t, volatile, "id")
	assert.Contains(t, volatile, "createdTimestamp")

	writeOnly, err := spec.WriteOnlyFields("client")
	require.NoError(t, err)
	assert.Contains(t, writeOnly, "clientSecret")
}
