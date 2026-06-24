package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/catalog"
)

func TestBuildRelationshipVariantRequiresData(t *testing.T) {
	template := catalog.RelationshipTemplate{
		Template:     "users/{user-id}/groups",
		Method:       http.MethodPost,
		Summary:      "assign user to group",
		Body:         map[string]interface{}{"type": "array"},
		RequiresBody: true,
	}

	variant, err := buildRelationshipVariant(template)
	require.NoError(t, err)

	required, ok := variant["required"].([]string)
	require.True(t, ok)
	assert.ElementsMatch(t, []string{"path", "data"}, required)

	props, ok := variant["properties"].(map[string]interface{})
	require.True(t, ok)

	pathSchema, ok := props["path"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "^users/([^/]+)/groups$", pathSchema["pattern"])

	deleteSchema, ok := props["delete"].(map[string]interface{})
	require.True(t, ok)
	deleteEnum, ok := deleteSchema["enum"].([]interface{})
	require.True(t, ok)
	assert.ElementsMatch(t, []interface{}{false}, deleteEnum)

	dataSchema, ok := props["data"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "array", dataSchema["type"])

	assert.Equal(t, "assign user to group", variant["description"])
}

func TestBuildRelationshipVariantOptionalData(t *testing.T) {
	template := catalog.RelationshipTemplate{
		Template: "users/{user-id}/groups/{group-id}",
		Method:   http.MethodDelete,
	}

	variant, err := buildRelationshipVariant(template)
	require.NoError(t, err)

	required, ok := variant["required"].([]string)
	require.True(t, ok)
	assert.ElementsMatch(t, []string{"path", "delete"}, required)

	props, ok := variant["properties"].(map[string]interface{})
	require.True(t, ok)

	_, hasData := props["data"]
	assert.False(t, hasData)

	pathSchema, ok := props["path"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "^users/([^/]+)/groups/([^/]+)$", pathSchema["pattern"])

	deleteSchema, ok := props["delete"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, true, deleteSchema["const"])
}

func TestBuildRelationshipJSONSchemaAggregatesErrors(t *testing.T) {
	valid := catalog.RelationshipTemplate{
		Template:     "users/{user-id}/groups",
		Method:       http.MethodPost,
		Body:         map[string]interface{}{"type": "array"},
		RequiresBody: true,
	}

	invalid := catalog.RelationshipTemplate{Method: http.MethodDelete}

	schema, err := buildRelationshipJSONSchema([]catalog.RelationshipTemplate{valid, invalid})
	require.Error(t, err)

	props, ok := schema["properties"].(map[string]interface{})
	require.True(t, ok)

	relationships, ok := props["relationships"].(map[string]interface{})
	require.True(t, ok)

	items, ok := relationships["items"].(map[string]interface{})
	require.True(t, ok)

	variants, ok := items["oneOf"].([]interface{})
	require.True(t, ok)
	assert.Len(t, variants, 1)
}

func TestWriteSchema(t *testing.T) {
	schema := map[string]interface{}{
		"type": "object",
	}

	targetDir := filepath.Join("..", "..", "tmp", "rel-schema-test")
	require.NoError(t, os.MkdirAll(targetDir, 0o755))
	t.Cleanup(func() {
		_ = os.RemoveAll(targetDir)
	})

	destination := filepath.Join(targetDir, "out.json")
	require.NoError(t, writeSchema(destination, schema))

	data, err := os.ReadFile(destination)
	require.NoError(t, err)

	var decoded map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, "object", decoded["type"])
}
