package output_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/admin"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
	"github.com/thedataflows/keycloak-cli/pkg/output"
)

func TestDestination(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "out.json")

	writer, shouldClose, err := output.Destination(outputPath, false)
	require.NoError(t, err)
	assert.True(t, shouldClose)
	require.NotNil(t, writer)
	assert.NoError(t, writer.Close())
}

func TestWriteJSON(t *testing.T) {
	var buffer bytes.Buffer

	err := output.WriteJSON(&buffer, map[string]string{"name": "demo"})
	require.NoError(t, err)
	assert.Contains(t, buffer.String(), "\"name\": \"demo\"")
}

func TestSanitizeResources(t *testing.T) {
	resources := []manifest.Resource{{
		Type:       "protocolmapper",
		Realm:      "demo",
		ParentType: "client",
		Delete:     true,
		Data: map[string]interface{}{
			"username":    "alice",
			"accessToken": "secret",
			"items": []interface{}{
				map[string]interface{}{"name": "one", "accessToken": "hidden"},
			},
		},
	}}

	sanitized := output.SanitizeResources(resources, []string{"accessToken"})
	require.Len(t, sanitized, 1)
	assert.Equal(t, "protocolmapper", sanitized[0].Type)
	assert.Equal(t, "client", sanitized[0].ParentType)
	assert.True(t, sanitized[0].Delete)
	assert.Equal(t, "alice", sanitized[0].Data["username"])
	_, ok := sanitized[0].Data["accessToken"]
	assert.False(t, ok)

	items, ok := sanitized[0].Data["items"].([]interface{})
	require.True(t, ok)
	first, ok := items[0].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "hidden", first["accessToken"])
}

func TestWriteResourceTable(t *testing.T) {
	resources := []manifest.Resource{
		{Type: "realm", Data: map[string]interface{}{"realm": "demo"}},
		{Type: "user", Realm: "demo", Data: map[string]interface{}{"username": "alice"}},
	}

	var compact bytes.Buffer
	err := output.WriteResourceTable(&compact, resources, false)
	require.NoError(t, err)
	assert.Contains(t, compact.String(), "NAME")
	assert.Contains(t, compact.String(), "demo/alice")

	var detailed bytes.Buffer
	err = output.WriteResourceTable(&detailed, resources, true)
	require.NoError(t, err)
	assert.Contains(t, detailed.String(), "DETAILS")
	assert.Contains(t, detailed.String(), "{\"username\":\"alice\"}")
}

func TestWriteApplyResults(t *testing.T) {
	results := []admin.ApplyResult{{
		Resource: "user",
		Realm:    "demo",
		Name:     "alice",
		Action:   "created",
		Status:   201,
	}}

	var table bytes.Buffer
	err := output.WriteApplyResults(&table, results, "table")
	require.NoError(t, err)
	assert.Contains(t, table.String(), "RESOURCE")
	assert.Contains(t, table.String(), "alice")

	var jsonOut bytes.Buffer
	err = output.WriteApplyResults(&jsonOut, results, "json")
	require.NoError(t, err)
	assert.Contains(t, jsonOut.String(), "\"resource\": \"user\"")

	var yamlOut bytes.Buffer
	err = output.WriteApplyResults(&yamlOut, results, "yaml")
	require.NoError(t, err)
	assert.Contains(t, yamlOut.String(), "resource: user")

	var tomlOut bytes.Buffer
	err = output.WriteApplyResults(&tomlOut, results, "toml")
	require.NoError(t, err)
	assert.Contains(t, tomlOut.String(), "Resource = 'user'")

	err = output.WriteApplyResults(&bytes.Buffer{}, results, "xml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported format")
}
func TestWriteResourcesToDir(t *testing.T) {
	t.Run("creates one file per resource", func(t *testing.T) {
		dir := t.TempDir()
		resources := []manifest.Resource{
			{Type: "user", Realm: "demo", Data: map[string]interface{}{"username": "alice", "id": "uid-1"}},
			{Type: "client", Realm: "demo", Data: map[string]interface{}{"clientId": "myapp", "id": "cid-1"}},
		}
		err := output.WriteResourcesToDir(dir, resources, nil, "json", false)
		require.NoError(t, err)
		files, err := os.ReadDir(dir)
		require.NoError(t, err)
		require.Len(t, files, 2)
		names := make([]string, len(files))
		for i, f := range files {
			names[i] = f.Name()
		}
		sort.Strings(names)
		assert.Equal(t, []string{"demo__client__cid-1.json", "demo__user__uid-1.json"}, names)
		// verify content contains wrapped resource envelope
		data, err := os.ReadFile(filepath.Join(dir, "demo__user__uid-1.json"))
		require.NoError(t, err)
		var payload map[string]interface{}
		require.NoError(t, json.Unmarshal(data, &payload))
		resList, ok := payload["resources"].([]interface{})
		require.True(t, ok)
		require.Len(t, resList, 1)
		resMap, ok := resList[0].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "user", resMap["type"])
		assert.Equal(t, "alice", resMap["data"].(map[string]interface{})["username"])
	})
	t.Run("writes relationships file when present", func(t *testing.T) {
		dir := t.TempDir()
		resources := []manifest.Resource{
			{Type: "user", Realm: "demo", Data: map[string]interface{}{"username": "alice"}},
		}
		relationships := []manifest.RelationshipOperation{{
			Kind:   "user",
			Method: "POST",
			Path:   "/admin/realms/demo/users/uid/groups",
		}}
		err := output.WriteResourcesToDir(dir, resources, relationships, "yaml", false)
		require.NoError(t, err)
		files, err := os.ReadDir(dir)
		require.NoError(t, err)
		require.Len(t, files, 2)
		_, err = os.ReadFile(filepath.Join(dir, "relationships.yaml"))
		require.NoError(t, err)
	})
	t.Run("handles collisions with suffix", func(t *testing.T) {
		dir := t.TempDir()
		resources := []manifest.Resource{
			{Type: "user", Realm: "demo", Data: map[string]interface{}{"username": "alice", "id": "x"}},
			{Type: "user", Realm: "demo", Data: map[string]interface{}{"username": "alice", "id": "x"}},
		}
		err := output.WriteResourcesToDir(dir, resources, nil, "json", false)
		require.NoError(t, err)
		files, err := os.ReadDir(dir)
		require.NoError(t, err)
		require.Len(t, files, 2)
		names := make([]string, len(files))
		for i, f := range files {
			names[i] = f.Name()
		}
		sort.Strings(names)
		assert.Equal(t, []string{"demo__user__x-1.json", "demo__user__x.json"}, names)
	})
	t.Run("force controls overwrite", func(t *testing.T) {
		dir := t.TempDir()
		resources := []manifest.Resource{
			{Type: "user", Realm: "demo", Data: map[string]interface{}{"username": "alice"}},
		}
		err := output.WriteResourcesToDir(dir, resources, nil, "json", false)
		require.NoError(t, err)
		// second write without force should error
		err = output.WriteResourcesToDir(dir, resources, nil, "json", false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
		// with force should succeed
		err = output.WriteResourcesToDir(dir, resources, nil, "json", true)
		require.NoError(t, err)
	})
}
