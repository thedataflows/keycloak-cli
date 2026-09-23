package kcapi

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstallDefaultBuiltInResourcesWiresMatcher(t *testing.T) {
	t.Cleanup(func() {
		IsBuiltInResource = func(Resource) bool { return false }
	})

	require.NoError(t, InstallDefaultBuiltInResources(filepath.Join(t.TempDir(), "missing.spec.json")))

	assert.True(t, IsBuiltInResource(Resource{
		Type: "client",
		Data: map[string]interface{}{"clientId": "account"},
	}))
	assert.True(t, IsBuiltInResource(Resource{
		Type: "role",
		Data: map[string]interface{}{"name": "default-roles-demo123"},
	}))
	assert.False(t, IsBuiltInResource(Resource{
		Type: "client",
		Data: map[string]interface{}{"clientId": "my-app"},
	}))
}

func TestBuiltInResourceNameFallbacks(t *testing.T) {
	assert.Equal(t, "my-app", builtInResourceName(Resource{Data: map[string]interface{}{"clientId": "my-app"}}))
	assert.Equal(t, "devs", builtInResourceName(Resource{Data: map[string]interface{}{"name": "devs"}}))
	assert.Equal(t, "github", builtInResourceName(Resource{Data: map[string]interface{}{"alias": "github"}}))
	assert.Equal(t, "alice", builtInResourceName(Resource{Data: map[string]interface{}{"username": "alice"}}))
	assert.Equal(t, "demo", builtInResourceName(Resource{Data: map[string]interface{}{"realm": "demo"}}))
}

func TestMatchSimplePattern(t *testing.T) {
	assert.True(t, matchSimplePattern("default-roles-demo123", "default-roles-*"))
	assert.False(t, matchSimplePattern("other-role", "default-roles-*"))
	assert.True(t, matchSimplePattern("my-suffix", "*-suffix"))
	assert.True(t, matchSimplePattern("exact", "exact"))
	assert.False(t, matchSimplePattern("not-exact", "exact"))
}
