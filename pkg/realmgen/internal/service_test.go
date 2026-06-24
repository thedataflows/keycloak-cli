package internal

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

type fakeGenerator struct {
	bundle generatedBundle
	err    error
	path   string
	seen   Options
}

func (f *fakeGenerator) GenerateBundle(specPath string, options Options) (generatedBundle, error) {
	f.path = specPath
	f.seen = options
	return f.bundle, f.err
}

func TestServiceGenerateUsesGeneratorBoundary(t *testing.T) {
	fake := &fakeGenerator{bundle: generatedBundle{
		Resources:     []manifest.Resource{{Type: "realm", Data: map[string]interface{}{"realm": "demo"}}},
		Relationships: []manifest.RelationshipOperation{{Path: "demo/users/a/groups/b"}},
	}}
	service := &Service{generator: fake}

	result, err := service.Generate(filepath.Join("..", "..", "..", "keycloak-oapi", "26.6.2.spec.json"), Options{Realm: "demo", WithUsers: 2})
	require.NoError(t, err)
	assert.Equal(t, "demo", result.Summary.Realm)
	assert.Equal(t, 2, result.Summary.ResourceCounts["user"])
	require.Len(t, result.Resources, 1)
	require.Len(t, result.Relationships, 1)
	assert.Equal(t, "demo", fake.seen.Realm)
	assert.Equal(t, 2, fake.seen.WithUsers)
	assert.Contains(t, fake.path, "26.6.2.spec.json")
}
