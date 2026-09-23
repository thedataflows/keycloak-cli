package testutil

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fatalTB records Fatalf instead of ending the test binary, so the helper's
// failure contract is observable from a test.
type fatalTB struct {
	testing.TB
	fatalf string
}

func (f *fatalTB) Helper() {}

func (f *fatalTB) Fatalf(format string, args ...any) {
	f.fatalf = format
	panic("fatal")
}

func TestKeycloakSpecPathResolvesFromEnv(t *testing.T) {
	t.Setenv("KEYCLOAK_VERSION", "9.9.9-test")

	got := KeycloakSpecPath(t)

	assert.True(t, filepath.IsAbs(got), "want an absolute path, got %q", got)
	assert.Equal(
		t,
		filepath.Join("keycloak-oapi", "9.9.9-test.spec.json"),
		filepath.Join(filepath.Base(filepath.Dir(got)), filepath.Base(got)),
		"the version must come from KEYCLOAK_VERSION",
	)
}

func TestKeycloakSpecPathFailsWhenEnvUnset(t *testing.T) {
	t.Setenv("KEYCLOAK_VERSION", "")

	tb := &fatalTB{TB: t}
	func() {
		defer func() { _ = recover() }()
		KeycloakSpecPath(tb)
	}()

	assert.Contains(t, tb.fatalf, "KEYCLOAK_VERSION", "the failure must name the variable")
}

func TestKeycloakSpecPathPointsAtVendoredSpec(t *testing.T) {
	require.FileExists(t, KeycloakSpecPath(t))
}
