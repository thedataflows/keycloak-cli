package catalog

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

func TestLoadFieldOverrides(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "field-overrides.yaml")
	data := []byte(`
overrides:
  - type: user
    addVolatile:
      - customTimestamp
    removeVolatile:
      - origin
    addWriteOnly:
      - totpSecret
`)
	require.NoError(t, os.WriteFile(path, data, 0o644))

	overrides, err := LoadFieldOverrides(path)
	require.NoError(t, err)
	require.Len(t, overrides, 1)
	assert.Equal(t, "user", overrides[0].Type)
	assert.Equal(t, []string{"customTimestamp"}, overrides[0].AddVolatile)
	assert.Equal(t, []string{"origin"}, overrides[0].RemoveVolatile)
	assert.Equal(t, []string{"totpSecret"}, overrides[0].AddWriteOnly)
}

func TestApplyFieldOverrides(t *testing.T) {
	origVolatile := cloneStringSliceMap(defaultVolatileFields)
	origWriteOnly := cloneStringSliceMap(defaultWriteOnlyFields)
	t.Cleanup(func() {
		defaultVolatileFields = origVolatile
		defaultWriteOnlyFields = origWriteOnly
	})

	err := ApplyFieldOverrides([]FieldOverride{
		{Type: "user", AddVolatile: []string{"customField"}, RemoveVolatile: []string{"origin"}, AddWriteOnly: []string{"totp"}},
	})
	require.NoError(t, err)

	volatile := defaultVolatileFields["user"]
	assert.Contains(t, volatile, "customField")
	assert.NotContains(t, volatile, "origin")

	writeOnly := defaultWriteOnlyFields["user"]
	assert.Contains(t, writeOnly, "totp")
}

func TestApplyFieldOverridesRequiresType(t *testing.T) {
	err := ApplyFieldOverrides([]FieldOverride{{AddVolatile: []string{"x"}}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type")
}

func TestInstallDefaultFieldOverridesMissingFileIsNoOp(t *testing.T) {
	require.NoError(t, InstallDefaultFieldOverrides(filepath.Join(t.TempDir(), "spec.json")))
}

func TestInstallDefaultFieldOverridesWiresManifest(t *testing.T) {
	origVolatile := cloneStringSliceMap(defaultVolatileFields)
	origWriteOnly := cloneStringSliceMap(defaultWriteOnlyFields)
	t.Cleanup(func() {
		defaultVolatileFields = origVolatile
		defaultWriteOnlyFields = origWriteOnly
	})

	dir := t.TempDir()
	path := filepath.Join(dir, "field-overrides.yaml")
	data := []byte(`
overrides:
  - type: user
    addVolatile:
      - wiredField
    addWriteOnly:
      - wiredSecret
`)
	require.NoError(t, os.WriteFile(path, data, 0o644))

	require.NoError(t, InstallDefaultFieldOverrides(filepath.Join(dir, "26.6.2.spec.json")))

	stripped := manifest.StripVolatileFields(manifest.Resource{Type: "user", Data: map[string]interface{}{
		"username":    "alice",
		"wiredField":  "x",
		"wiredSecret": "y",
	}})
	assert.NotContains(t, stripped.Data, "wiredField")
	assert.NotContains(t, stripped.Data, "wiredSecret")
}

func TestInstallDefaultFieldOverridesWiresDefaultsWhenFileMissing(t *testing.T) {
	require.NoError(t, InstallDefaultFieldOverrides(filepath.Join(t.TempDir(), "26.6.2.spec.json")))

	stripped := manifest.StripVolatileFields(manifest.Resource{Type: "identityprovider", Data: map[string]interface{}{
		"alias": "idp-1",
		"types": []string{"USER_AUTHENTICATION"},
		"config": map[string]interface{}{
			"clientId":     "idp-client",
			"clientSecret": "should-be-stripped",
		},
	}})

	assert.NotContains(t, stripped.Data, "types")
	config, ok := stripped.Data["config"].(map[string]interface{})
	require.True(t, ok)
	assert.NotContains(t, config, "clientSecret")
}

func cloneStringSliceMap(src map[string][]string) map[string][]string {
	clone := make(map[string][]string, len(src))
	for k, v := range src {
		clone[k] = append([]string(nil), v...)
	}
	return clone
}
