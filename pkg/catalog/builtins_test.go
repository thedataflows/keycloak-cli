package catalog

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

func TestInstallDefaultBuiltInResourcesWiresMatcher(t *testing.T) {
	t.Cleanup(func() {
		manifest.IsBuiltInResource = func(manifest.Resource) bool { return false }
	})

	require.NoError(t, InstallDefaultBuiltInResources(filepath.Join(t.TempDir(), "26.6.2.spec.json")))

	assert.True(t, manifest.IsBuiltInResource(manifest.Resource{
		Type: "client",
		Data: map[string]interface{}{"clientId": "account"},
	}))
	assert.True(t, manifest.IsBuiltInResource(manifest.Resource{
		Type: "role",
		Data: map[string]interface{}{"name": "default-roles-demo123"},
	}))
	assert.False(t, manifest.IsBuiltInResource(manifest.Resource{
		Type: "client",
		Data: map[string]interface{}{"clientId": "my-app"},
	}))
}

func TestBuiltInResourceNameFallbacks(t *testing.T) {
	assert.Equal(t, "my-app", builtInResourceName(manifest.Resource{Data: map[string]interface{}{"clientId": "my-app"}}))
	assert.Equal(t, "devs", builtInResourceName(manifest.Resource{Data: map[string]interface{}{"name": "devs"}}))
	assert.Equal(t, "github", builtInResourceName(manifest.Resource{Data: map[string]interface{}{"alias": "github"}}))
	assert.Equal(t, "alice", builtInResourceName(manifest.Resource{Data: map[string]interface{}{"username": "alice"}}))
	assert.Equal(t, "demo", builtInResourceName(manifest.Resource{Data: map[string]interface{}{"realm": "demo"}}))
}

func TestMatchSimplePattern(t *testing.T) {
	assert.True(t, matchSimplePattern("default-roles-demo123", "default-roles-*"))
	assert.False(t, matchSimplePattern("other-role", "default-roles-*"))
	assert.True(t, matchSimplePattern("my-suffix", "*-suffix"))
	assert.True(t, matchSimplePattern("exact", "exact"))
	assert.False(t, matchSimplePattern("not-exact", "exact"))
}
