package admin

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

func TestSanitizeResourceDataStripsNilAndEmptyMaps(t *testing.T) {
	tests := []struct {
		name string
		data map[string]interface{}
		want map[string]interface{}
	}{
		{
			name: "nil map returns nil",
			data: nil,
			want: nil,
		},
		{
			name: "top-level null removed",
			data: map[string]interface{}{"enabled": true, "displayName": nil},
			want: map[string]interface{}{"enabled": true},
		},
		{
			name: "nested null removed",
			data: map[string]interface{}{
				"clientId": "my-app",
				"authenticationFlowBindingOverrides": map[string]interface{}{
					"browser": nil,
				},
			},
			want: map[string]interface{}{"clientId": "my-app"},
		},
		{
			name: "empty map removed",
			data: map[string]interface{}{
				"clientId":   "my-app",
				"attributes": map[string]interface{}{},
			},
			want: map[string]interface{}{"clientId": "my-app"},
		},
		{
			name: "arrays and primitives preserved",
			data: map[string]interface{}{
				"clientId":     "my-app",
				"redirectUris": []interface{}{},
				"protocolMappers": []interface{}{
					map[string]interface{}{
						"name":   "mapper",
						"config": nil,
					},
				},
			},
			want: map[string]interface{}{
				"clientId":     "my-app",
				"redirectUris": []interface{}{},
				"protocolMappers": []interface{}{
					map[string]interface{}{"name": "mapper"},
				},
			},
		},
		{
			name: "nil array items removed",
			data: map[string]interface{}{
				"items": []interface{}{"a", nil, "b"},
			},
			want: map[string]interface{}{
				"items": []interface{}{"a", "b"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, sanitizeResourceData("client", tt.data))
		})
	}
}

func TestSanitizeResourceDataStripsClientTimeoutAttributes(t *testing.T) {
	data := map[string]interface{}{
		"clientId": "demo-sa",
		"attributes": map[string]interface{}{
			"client.session.idle.timeout": "86400",
			"client.session.max.lifespan": "2592000",
			"other.attr":                  "keep",
		},
	}
	want := map[string]interface{}{
		"clientId": "demo-sa",
		"attributes": map[string]interface{}{
			"other.attr": "keep",
		},
	}
	assert.Equal(t, want, sanitizeResourceData("client", data))
	// Other resource types should leave attributes untouched.
	assert.Equal(t, data, sanitizeResourceData("user", data))
}

func TestOperationContextIgnoresParentDeadline(t *testing.T) {
	type ctxKey string
	key := ctxKey("test-key")

	parent, parentCancel := context.WithTimeout(context.WithValue(context.Background(), key, "preserved"), time.Nanosecond)
	defer parentCancel()
	// Ensure the parent is already expired.
	time.Sleep(5 * time.Millisecond)
	assert.Error(t, parent.Err())

	s := &service{timeout: 30 * time.Second}
	child, cancel := s.operationContext(parent)
	defer cancel()

	assert.NoError(t, child.Err(), "child context should not inherit parent deadline")
	assert.Equal(t, "preserved", child.Value(key))

	deadline, ok := child.Deadline()
	assert.True(t, ok, "child should have a deadline")
	assert.True(t, deadline.After(time.Now().Add(25*time.Second)), "child deadline should be ~30s in the future")
}

func TestSanitizeResourceDataStripsRealmForNonRealmResources(t *testing.T) {
	assert.Equal(t, map[string]interface{}{"alias": "browser"}, sanitizeResourceData("authenticationflow", map[string]interface{}{"alias": "browser", "realm": "demo"}))
	assert.Equal(t, map[string]interface{}{"alias": "browser", "realm": "demo"}, sanitizeResourceData("realm", map[string]interface{}{"alias": "browser", "realm": "demo"}))
}

func TestPriorityMapWithInlineReferencesOrdersReferencerAfterReferenced(t *testing.T) {
	resources := []manifest.Resource{
		{
			Type:  "client",
			Realm: "demo",
			Data: map[string]interface{}{
				"clientId": "fcc",
				"authenticationFlowBindingOverrides": map[string]interface{}{
					"browser": "f4509057-d632-42d8-9ce8-4a2f832a8620",
				},
			},
		},
		{
			Type:  "authenticationflow",
			Realm: "demo",
			Data: map[string]interface{}{
				"id":    "f4509057-d632-42d8-9ce8-4a2f832a8620",
				"alias": "X509 Browser",
			},
		},
	}

	specGraph := map[string][]string{
		"realm":              {},
		"client":             {"realm"},
		"authenticationflow": {"realm"},
	}

	priorityMap, err := priorityMapWithInlineReferences(resources, specGraph)
	require.NoError(t, err)
	assert.Less(t, priorityMap["authenticationflow"], priorityMap["client"], "authenticationflow must be applied before client")
}

func TestPriorityMapWithInlineReferencesDetectsCycle(t *testing.T) {
	resources := []manifest.Resource{
		{
			Type: "client",
			Data: map[string]interface{}{
				"id":       "eb51dd4e-d7bd-40ec-90dc-430df5275c0a",
				"clientId": "a",
				"authenticationFlowBindingOverrides": map[string]interface{}{
					"browser": "f4509057-d632-42d8-9ce8-4a2f832a8620",
				},
			},
		},
		{
			Type: "authenticationflow",
			Data: map[string]interface{}{
				"id":    "f4509057-d632-42d8-9ce8-4a2f832a8620",
				"alias": "browser",
				"config": map[string]interface{}{
					"clientId": "eb51dd4e-d7bd-40ec-90dc-430df5275c0a",
				},
			},
		},
	}

	specGraph := map[string][]string{
		"client":             {},
		"authenticationflow": {},
	}

	_, err := priorityMapWithInlineReferences(resources, specGraph)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cycle")
}
