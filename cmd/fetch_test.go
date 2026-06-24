package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

func TestFilterResources(t *testing.T) {
	resources := []manifest.Resource{
		{Type: "user", Realm: "demo", Data: map[string]interface{}{"username": "alice", "id": "uid-1"}},
		{Type: "user", Realm: "demo", Data: map[string]interface{}{"username": "bob", "id": "uid-2"}},
		{Type: "client", Realm: "demo", Data: map[string]interface{}{"clientId": "myapp", "id": "cid-1"}},
	}

	tests := []struct {
		name     string
		filter   string
		expected []string
	}{
		{
			name:     "match by name case-insensitive",
			filter:   "ALICE",
			expected: []string{"alice"},
		},
		{
			name:     "match by identifier",
			filter:   "uid-2",
			expected: []string{"bob"},
		},
		{
			name:     "match by clientId",
			filter:   "myapp",
			expected: []string{"myapp"},
		},
		{
			name:     "no match",
			filter:   "charlie",
			expected: []string{},
		},
		{
			name:     "empty filter returns empty",
			filter:   "  ",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterResources(resources, tt.filter)
			names := make([]string, len(got))
			for i, r := range got {
				names[i] = r.Name()
			}
			assert.Equal(t, tt.expected, names)
		})
	}
}

func TestIsDirTarget(t *testing.T) {
	tmp := t.TempDir()
	filePath := filepath.Join(tmp, "file.txt")
	assert.NoError(t, os.WriteFile(filePath, []byte("x"), 0o644))

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"empty string", "", false},
		{"plain file", filePath, false},
		{"existing dir", tmp, true},
		{"trailing slash", tmp + "/", true},
		{"trailing backslash on non-windows", "out\\", false},
		{"nonexistent path", filepath.Join(tmp, "nope"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isDirTarget(tt.path))
		})
	}
}
