package manifest_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

func TestParseResources(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		wantCount int
		wantType  string
		wantErr   bool
	}{
		{
			name: "json list",
			input: []byte(`[
				{"type":"realm","realm":"demo","data":{"realm":"demo"}},
				{"type":"client","realm":"demo","data":{"clientId":"app"}}
			]`),
			wantCount: 2,
			wantType:  "realm",
		},
		{
			name:      "yaml single",
			input:     []byte("type: user\nrealm: demo\ndata:\n  username: alice\n"),
			wantCount: 1,
			wantType:  "user",
		},
		{
			name:    "invalid resource payload",
			input:   []byte(`[{"type":"realm"}]`),
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(test *testing.T) {
			resources, ok, err := manifest.ParseResources(testCase.input)
			assert.True(test, ok)
			if testCase.wantErr {
				require.Error(test, err)
				return
			}

			require.NoError(test, err)
			require.Len(test, resources, testCase.wantCount)
			assert.Equal(test, testCase.wantType, resources[0].Type)
		})
	}
}

func TestResourceFields(t *testing.T) {
	resource := manifest.Resource{
		Type:  "client",
		Realm: "demo",
		Data: map[string]interface{}{
			"id":       "client-id",
			"clientId": "app",
		},
	}

	assert.Equal(t, "client-id", resource.Identifier())
	assert.Equal(t, "app", resource.Name())
	assert.Equal(t, "app", resource.DisplayName())
	assert.NoError(t, manifest.ValidateResources([]manifest.Resource{resource}))
}

func TestParseRelationshipManifest(t *testing.T) {
	relManifest, ok, err := manifest.ParseRelationshipManifest([]byte(`{
		"relationships": [
			{"path": "/admin/realms/demo/users/user-1/groups/group-1"}
		]
	}`))
	require.NoError(t, err)
	require.True(t, ok)
	require.Len(t, relManifest.Relationships, 1)
	assert.Equal(t, "demo/users/user-1/groups/group-1", relManifest.Relationships[0].Path)
}

func TestNewRelationshipOperation(t *testing.T) {
	relationship, err := manifest.NewRelationshipOperation(
		"/admin/realms/{realm}/users/{user-id}/groups/{groupId}",
		"PUT",
		map[string]string{"realm": "demo", "user-id": "user-1", "groupId": "group-1"},
		nil,
	)
	require.NoError(t, err)
	assert.Equal(t, "demo/users/user-1/groups/group-1", relationship.Path)
	assert.Equal(t, "PUT", relationship.Method)
	assert.False(t, relationship.Delete)
}

func TestLoadPathsSkipsInvalid(t *testing.T) {
	tmpDir := t.TempDir()

	validPath := filepath.Join(tmpDir, "realm.json")
	err := os.WriteFile(validPath, []byte(`[
		{"type":"realm","realm":"demo","data":{"realm":"demo"}}
	]`), 0o644)
	require.NoError(t, err)

	invalidPath := filepath.Join(tmpDir, "invalid.json")
	err = os.WriteFile(invalidPath, []byte(`{"$schema":"http://example.com"}`), 0o644)
	require.NoError(t, err)

	loaded, err := manifest.LoadPaths([]string{tmpDir})
	require.NoError(t, err)
	require.Len(t, loaded.Resources, 1)
	require.Empty(t, loaded.Relationships)
	require.Len(t, loaded.Skipped, 1)
	assert.Equal(t, invalidPath, loaded.Skipped[0].Path)
	assert.Contains(t, loaded.Skipped[0].Reason, "unsupported manifest format")
}

func TestNormalizeRoundTripResources(t *testing.T) {
	resources := []manifest.Resource{
		{
			Type:  "user",
			Realm: "demo",
			Data: map[string]interface{}{
				"id":               "user-1",
				"username":         "alice",
				"createdTimestamp": float64(1234),
				"credentials": []interface{}{
					map[string]interface{}{"id": "cred-2", "type": "otp"},
					map[string]interface{}{"id": "cred-1", "type": "password"},
				},
			},
		},
		{
			Type:  "group",
			Realm: "demo",
			Data: map[string]interface{}{
				"id":   "group-1",
				"name": "devs",
			},
		},
	}

	normalizedResources, normalizedRelationships := manifest.NormalizeRoundTrip(resources, nil)
	require.Len(t, normalizedResources, 2)
	assert.Nil(t, normalizedRelationships)
	assert.Equal(t, "group", normalizedResources[0].Type)
	assert.Equal(t, "user", normalizedResources[1].Type)
	assert.NotContains(t, normalizedResources[1].Data, "id")
	assert.NotContains(t, normalizedResources[1].Data, "createdTimestamp")

	// Credentials are write-only and stripped during normalization
	// because Keycloak never returns them in GET responses.
	assert.NotContains(t, normalizedResources[1].Data, "credentials")
}

func TestNormalizeForApplyPreservesReferencedIDs(t *testing.T) {
	resources := []manifest.Resource{
		{
			Type:  "authenticationflow",
			Realm: "demo",
			Data: map[string]interface{}{
				"id":         "f4509057-d632-42d8-9ce8-4a2f832a8620",
				"alias":      "browser",
				"providerId": "basic-flow",
			},
		},
		{
			Type:  "client",
			Realm: "demo",
			Data: map[string]interface{}{
				"id":       "eb51dd4e-d7bd-40ec-90dc-430df5275c0a",
				"clientId": "my-app",
				"authenticationFlowBindingOverrides": map[string]interface{}{
					"browser": "f4509057-d632-42d8-9ce8-4a2f832a8620",
				},
			},
		},
		{
			Type:  "user",
			Realm: "demo",
			Data: map[string]interface{}{
				"id":       "user-id-1",
				"username": "alice",
			},
		},
	}

	normalizedResources, _ := manifest.NormalizeForApply(resources, nil)
	require.Len(t, normalizedResources, 3)

	// Find the flow and client by iterating; order is sorted by resource key.
	var flow, client, user manifest.Resource
	for _, r := range normalizedResources {
		switch r.Type {
		case "authenticationflow":
			flow = r
		case "client":
			client = r
		case "user":
			user = r
		}
	}

	// The flow id is referenced by the client, so it must be preserved.
	assert.Equal(t, "f4509057-d632-42d8-9ce8-4a2f832a8620", flow.Data["id"])
	// The client id is not referenced by any other resource, so it is stripped.
	assert.NotContains(t, client.Data, "id")
	// The user id is not referenced, so it is stripped.
	assert.NotContains(t, user.Data, "id")
}

func TestNormalizeRoundTripRelationships(t *testing.T) {
	resources := []manifest.Resource{
		{Type: "user", Realm: "demo", Data: map[string]interface{}{"id": "user-1", "username": "alice"}},
		{Type: "group", Realm: "demo", Data: map[string]interface{}{"id": "group-1", "name": "devs"}},
		{Type: "client", Realm: "demo", Data: map[string]interface{}{"id": "client-1", "clientId": "app"}},
		{Type: "role", Realm: "demo", Data: map[string]interface{}{"id": "role-1", "name": "developer"}},
		{Type: "role", Realm: "demo", Data: map[string]interface{}{"id": "role-2", "name": "auditor"}},
		{Type: "organization", Realm: "demo", Data: map[string]interface{}{"id": "org-1", "name": "Platform"}},
		{Type: "identityprovider", Realm: "demo", Data: map[string]interface{}{"id": "idp-1", "alias": "github"}},
	}
	relationships := []manifest.RelationshipOperation{
		{
			Kind:     "organization-member",
			Method:   "post",
			Template: "{realm}/organizations/{org-id}/members",
			Path:     "demo/organizations/org-1/members",
			PathParams: map[string]string{
				"realm":  "demo",
				"org-id": "org-1",
			},
			Data: json.RawMessage(`"user-1"`),
		},
		{
			Kind:     "organization-identity-provider",
			Method:   "POST",
			Template: "{realm}/organizations/{org-id}/identity-providers",
			Path:     "demo/organizations/org-1/identity-providers",
			PathParams: map[string]string{
				"realm":  "demo",
				"org-id": "org-1",
			},
			Data: json.RawMessage(`"idp-1"`),
		},
		{
			Kind:     "user-federated-identity",
			Method:   "POST",
			Template: "{realm}/users/{user-id}/federated-identity/{provider}",
			Path:     "demo/users/user-1/federated-identity/idp-1",
			PathParams: map[string]string{
				"realm":    "demo",
				"user-id":  "user-1",
				"provider": "idp-1",
			},
			Data: json.RawMessage(`{"identityProvider":"idp-1","userId":"external-1","userName":"alice-ext","extra":"ignored"}`),
		},
		{
			Kind:     "user-realm-role-mapping",
			Method:   "POST",
			Template: "{realm}/users/{user-id}/role-mappings/realm",
			Path:     "demo/users/user-1/role-mappings/realm",
			PathParams: map[string]string{
				"realm":   "demo",
				"user-id": "user-1",
			},
			Data: json.RawMessage(`[{"id":"role-1","name":"developer"},{"id":"role-2","name":"auditor"}]`),
		},
		{
			Kind:     "client-scope-client-role-mapping",
			Method:   "POST",
			Template: "{realm}/client-scopes/{client-scope-id}/scope-mappings/clients/{client}",
			Path:     "demo/client-scopes/scope-1/scope-mappings/clients/client-1",
			PathParams: map[string]string{
				"realm":           "demo",
				"client-scope-id": "scope-1",
				"client":          "client-1",
			},
			Data: json.RawMessage(`[{"id":"role-3","name":"app-user","clientRole":true,"containerId":"client-1","description":"ignored"}]`),
		},
		{
			Kind:     "role-composite-mapping",
			Method:   "POST",
			Template: "{realm}/roles-by-id/{role-id}/composites",
			Path:     "demo/roles-by-id/role-1/composites",
			PathParams: map[string]string{
				"realm":   "demo",
				"role-id": "role-1",
			},
			Data: json.RawMessage(`[{"id":"role-2","name":"auditor","description":"ignored"}]`),
		},
		{
			Kind:     "user-group-membership",
			Method:   "PUT",
			Template: "{realm}/users/{user-id}/groups/{groupId}",
			Path:     "demo/users/user-1/groups/group-1",
			PathParams: map[string]string{
				"realm":   "demo",
				"user-id": "user-1",
				"groupId": "group-1",
			},
		},
	}

	_, normalizedRelationships := manifest.NormalizeRoundTrip(resources, relationships)
	require.Len(t, normalizedRelationships, 7)
	byKind := make(map[string]manifest.RelationshipOperation, len(normalizedRelationships))
	for _, relationship := range normalizedRelationships {
		byKind[relationship.Kind] = relationship
	}

	assert.Equal(t, "POST", byKind["organization-member"].Method)
	assert.Equal(t, "demo/organizations/Platform/members", byKind["organization-member"].Path)
	assert.JSONEq(t, `"alice"`, string(byKind["organization-member"].Data))
	assert.JSONEq(t, `"github"`, string(byKind["organization-identity-provider"].Data))
	assert.Equal(t, "demo/users/alice/federated-identity/github", byKind["user-federated-identity"].Path)
	assert.JSONEq(t, `{"identityProvider":"github","userId":"external-1","userName":"alice-ext"}`, string(byKind["user-federated-identity"].Data))
	assert.Equal(t, "demo/users/alice/groups/devs", byKind["user-group-membership"].Path)
	assert.JSONEq(t, `[{"name":"auditor"},{"name":"developer"}]`, string(byKind["user-realm-role-mapping"].Data))
	assert.JSONEq(t, `[{"name":"app-user","clientRole":true,"client":"app"}]`, string(byKind["client-scope-client-role-mapping"].Data))
	assert.Equal(t, "demo/roles-by-id/developer/composites", byKind["role-composite-mapping"].Path)
	assert.JSONEq(t, `[{"name":"auditor"}]`, string(byKind["role-composite-mapping"].Data))
}

func TestNormalizeForApplyPreservesParentType(t *testing.T) {
	resources := []manifest.Resource{{
		Type:       "protocolmapper",
		Realm:      "demo",
		ParentType: "clientscope",
		Data:       map[string]interface{}{"id": "mapper-1", "name": "email"},
	}}
	normalized, _ := manifest.NormalizeForApply(resources, nil)
	require.Len(t, normalized, 1)
	assert.Equal(t, "protocolmapper", normalized[0].Type)
	assert.Equal(t, "clientscope", normalized[0].ParentType)
}

func TestCompareRoundTrip(t *testing.T) {
	expectedResources := []manifest.Resource{{
		Type:  "user",
		Realm: "demo",
		Data: map[string]interface{}{
			"id":       "user-1",
			"username": "alice",
			"enabled":  true,
		},
	}, {
		Type:  "group",
		Realm: "demo",
		Data: map[string]interface{}{
			"id":   "group-1",
			"name": "devs",
		},
	}}
	actualResources := []manifest.Resource{
		{
			Type:  "user",
			Realm: "demo",
			Data: map[string]interface{}{
				"id":       "user-1",
				"username": "alice",
				"enabled":  true,
			},
		},
		{
			Type:  "group",
			Realm: "demo",
			Data: map[string]interface{}{
				"id":   "group-1",
				"name": "devs",
			},
		},
	}
	expectedRelationships := []manifest.RelationshipOperation{{
		Kind:     "user-group-membership",
		Method:   "PUT",
		Template: "{realm}/users/{user-id}/groups/{groupId}",
		Path:     "demo/users/user-1/groups/group-1",
		PathParams: map[string]string{
			"realm":   "demo",
			"user-id": "user-1",
			"groupId": "group-1",
		},
	}}
	actualRelationships := []manifest.RelationshipOperation{{
		Kind:     "user-group-membership",
		Method:   "PUT",
		Template: "{realm}/users/{user-id}/groups/{groupId}",
		Path:     "demo/users/user-1/groups/group-1",
		PathParams: map[string]string{
			"realm":   "demo",
			"user-id": "user-1",
			"groupId": "group-1",
		},
	}}

	report := manifest.CompareRoundTrip(expectedResources, expectedRelationships, actualResources, actualRelationships)
	assert.True(t, report.Match)
	assert.Empty(t, report.UnexpectedResources)
	assert.Empty(t, report.MissingResources)
	assert.Empty(t, report.MismatchedResources)
	assert.Empty(t, report.MissingRelationships)
	assert.Empty(t, report.UnexpectedRelationships)

	mismatchedActual := []manifest.Resource{{
		Type:  "user",
		Realm: "demo",
		Data: map[string]interface{}{
			"id":       "user-1",
			"username": "alice",
			"enabled":  false,
		},
	}}
	report = manifest.CompareRoundTrip(expectedResources, nil, mismatchedActual, nil)
	assert.False(t, report.Match)
	assert.Len(t, report.MismatchedResources, 1)
}

func TestCompareRoundTripIgnoresBuiltInResources(t *testing.T) {
	orig := manifest.IsBuiltInResource
	manifest.IsBuiltInResource = func(r manifest.Resource) bool {
		return r.Type == "client" && r.Data["clientId"] == "account"
	}
	t.Cleanup(func() { manifest.IsBuiltInResource = orig })

	expected := []manifest.Resource{{
		Type:  "client",
		Realm: "demo",
		Data:  map[string]interface{}{"clientId": "my-app"},
	}}
	actual := []manifest.Resource{
		{Type: "client", Realm: "demo", Data: map[string]interface{}{"clientId": "my-app"}},
		{Type: "client", Realm: "demo", Data: map[string]interface{}{"clientId": "account"}},
	}

	report := manifest.CompareRoundTrip(expected, nil, actual, nil)
	assert.True(t, report.Match)
	assert.Empty(t, report.UnexpectedResources)
}

func TestCompareRoundTripMatchReflectsUnexpectedResources(t *testing.T) {
	expected := []manifest.Resource{{
		Type:  "client",
		Realm: "demo",
		Data:  map[string]interface{}{"clientId": "my-app"},
	}}
	actual := []manifest.Resource{
		{Type: "client", Realm: "demo", Data: map[string]interface{}{"clientId": "my-app"}},
		{Type: "client", Realm: "demo", Data: map[string]interface{}{"clientId": "extra"}},
	}

	report := manifest.CompareRoundTrip(expected, nil, actual, nil)
	assert.False(t, report.Match)
	assert.Len(t, report.UnexpectedResources, 1)
}

func TestParseCombinedEnvelope(t *testing.T) {
	tmpDir := t.TempDir()
	envelopePath := filepath.Join(tmpDir, "manifest.json")
	require.NoError(t, os.WriteFile(envelopePath, []byte(`{
		"resources": [
			{"type":"realm","realm":"demo","data":{"realm":"demo"}}
		],
		"relationships": [
			{"path":"demo/users/alice/groups/devs","method":"PUT"}
		]
	}`), 0644))
	loaded, err := manifest.LoadPaths([]string{envelopePath})
	require.NoError(t, err)
	require.Len(t, loaded.Resources, 1)
	require.Len(t, loaded.Relationships, 1)
	assert.Equal(t, "demo/users/alice/groups/devs", loaded.Relationships[0].Path)
}
