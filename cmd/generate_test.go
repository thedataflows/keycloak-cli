package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

func TestGenerateCmd_Run(t1 *testing.T) {
	// Create temp output directory
	tmpDir := t1.TempDir()

	tests := []struct {
		name        string
		cmd         *GenerateCmd
		expectError bool
		validate    func(t *testing.T, outputDir string)
	}{
		{
			name: "generate complete realm",
			cmd: &GenerateCmd{
				Output:      tmpDir,
				Format:      "json",
				Realm:       "test-realm",
				WithUsers:   5,
				WithClients: 2,
				WithRoles:   3,
				WithGroups:  1,
				Overwrite:   true,
			},
			expectError: false,
			validate: func(t *testing.T, outputDir string) {
				// Check realm file exists
				realmPath := filepath.Join(outputDir, "realm.json")
				require.FileExists(t, realmPath)

				// Parse and validate realm - now it's an array of GenericResource
				data, err := os.ReadFile(realmPath)
				require.NoError(t, err)

				var resources []manifest.Resource
				err = json.Unmarshal(data, &resources)
				require.NoError(t, err)

				// Relationships manifest should be generated when dependent resources exist
				relPath := filepath.Join(outputDir, "relationships.json")
				require.FileExists(t, relPath)

				relData, err := os.ReadFile(relPath)
				require.NoError(t, err)

				var relManifest manifest.RelationshipManifest
				err = json.Unmarshal(relData, &relManifest)
				require.NoError(t, err)
				assert.NotEmpty(t, relManifest.Relationships, "Should generate relationship operations")

				// Find realm resource
				var realmResource *manifest.Resource
				for i := range resources {
					if resources[i].Type == "realm" {
						realmResource = &resources[i]
						break
					}
				}
				require.NotNil(t, realmResource, "realm resource should exist")

				// Validate realm structure
				assert.Equal(t, "test-realm", realmResource.Data["realm"])
				assert.Equal(t, true, realmResource.Data["enabled"])

				// Check resources exist - count each type
				userCount := 0
				clientCount := 0
				for _, res := range resources {
					switch res.Type {
					case "user":
						userCount++
					case "client":
						clientCount++
					}
				}
				assert.Equal(t, 5, userCount, "Should have 5 users")
				assert.Equal(t, 2, clientCount, "Should have 2 clients")
			},
		},
		{
			cmd: &GenerateCmd{
				Output:    filepath.Join(tmpDir, "minimal") + string(os.PathSeparator),
				Format:    "json",
				Realm:     "minimal-realm",
				Overwrite: true,
			},
			expectError: false,
			validate: func(t *testing.T, outputDir string) {
				realmPath := filepath.Join(outputDir, "realm.json")
				require.FileExists(t, realmPath)

				data, err := os.ReadFile(realmPath)
				require.NoError(t, err)

				var resources []manifest.Resource
				err = json.Unmarshal(data, &resources)
				require.NoError(t, err)

				// Find realm resource
				var realmResource *manifest.Resource
				for i := range resources {
					if resources[i].Type == "realm" {
						realmResource = &resources[i]
						break
					}
				}
				require.NotNil(t, realmResource, "realm resource should exist")

				assert.Equal(t, "minimal-realm", realmResource.Data["realm"])
				assert.Equal(t, true, realmResource.Data["enabled"])
				assert.Len(t, resources, 1, "minimal generation should stay realm-only")
			},
		},
		{
			cmd: &GenerateCmd{
				Output:    tmpDir,
				Format:    "json",
				Realm:     "test-realm",
				Overwrite: false,
			},
			expectError: true, // Should fail if file exists
		},
	}

	for _, tc := range tests {
		t1.Run(tc.name, func(t *testing.T) {
			// Set spec path to test spec
			cli := &CLI{
				Globals: Globals{
					SpecPath: "../keycloak-oapi/26.6.2.spec.json",
				},
			}

			err := tc.cmd.Run(nil, cli)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tc.validate != nil {
					tc.validate(t, tc.cmd.Output)
				}
			}
		})
	}
}

func TestGenerateCmd_ValidateOptions(t1 *testing.T) {
	tests := []struct {
		name        string
		cmd         *GenerateCmd
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid options",
			cmd: &GenerateCmd{
				Output:    "test",
				Realm:     "test-realm",
				WithUsers: 5,
			},
			expectError: false,
		},
		{
			name: "empty realm name",
			cmd: &GenerateCmd{
				Output: "test",
				Realm:  "",
			},
			expectError: true,
			errorMsg:    "realm name is required",
		},
		{
			name: "negative counts",
			cmd: &GenerateCmd{
				Output:    "test",
				Realm:     "test",
				WithUsers: -1,
			},
			expectError: true,
			errorMsg:    "resource counts must be non-negative",
		},
	}

	for _, tc := range tests {
		t1.Run(tc.name, func(t *testing.T) {
			err := tc.cmd.Validate()

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func BenchmarkGenerateCmd_Run(b *testing.B) {
	tmpDir := b.TempDir()

	cmd := &GenerateCmd{
		Output:      tmpDir,
		Format:      "json",
		Realm:       "bench-realm",
		WithUsers:   100,
		WithClients: 10,
		WithRoles:   20,
		WithGroups:  5,
		Overwrite:   true,
	}

	cli := &CLI{
		Globals: Globals{
			SpecPath: "../keycloak-oapi/26.6.2.spec.json",
		},
	}

	for b.Loop() {
		_ = cmd.Run(nil, cli)
	}
}
