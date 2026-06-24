package realmgen_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/realmgen"
)

func TestServiceGenerate(t *testing.T) {
	service := realmgen.New()

	result, err := service.Generate(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"), realmgen.Options{
		Realm:       "svc-realm",
		WithUsers:   2,
		WithClients: 1,
		WithGroups:  1,
		WithRoles:   1,
	})
	require.NoError(t, err)
	require.NotEmpty(t, result.Resources)
	assert.Equal(t, "svc-realm", result.Summary.Realm)
	assert.Equal(t, 2, result.Summary.ResourceCounts["user"])
	assert.NotZero(t, result.Summary.SchemaCount)
	assert.NotEmpty(t, result.Relationships)

	var userCount int
	var relationshipKinds int
	for _, resource := range result.Resources {
		if resource.Type == "user" {
			userCount++
		}
		assert.NotEqual(t, "component", resource.Type)
	}
	for _, relationship := range result.Relationships {
		if relationship.Kind != "" {
			relationshipKinds++
		}
	}
	assert.Equal(t, 2, userCount)
	assert.NotZero(t, relationshipKinds)
}

func TestServiceGenerateRejectsInvalidOptions(t *testing.T) {
	service := realmgen.New()

	_, err := service.Generate(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"), realmgen.Options{
		Realm:     "",
		WithUsers: 1,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "realm name is required")
}

func TestValidateOptions(t *testing.T) {
	tests := []struct {
		name        string
		options     realmgen.Options
		expectError bool
		errorText   string
	}{
		{
			name: "valid options",
			options: realmgen.Options{
				Realm:     "demo",
				WithUsers: 1,
			},
		},
		{
			name: "missing realm",
			options: realmgen.Options{
				WithUsers: 1,
			},
			expectError: true,
			errorText:   "realm name is required",
		},
		{
			name: "negative counts",
			options: realmgen.Options{
				Realm:     "demo",
				WithUsers: -1,
			},
			expectError: true,
			errorText:   "resource counts must be non-negative",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(test *testing.T) {
			err := realmgen.ValidateOptions(testCase.options)
			if testCase.expectError {
				require.Error(test, err)
				assert.Contains(test, err.Error(), testCase.errorText)
				return
			}

			require.NoError(test, err)
		})
	}
}
