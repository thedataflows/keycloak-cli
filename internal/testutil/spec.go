// Package testutil holds shared test helpers.
package testutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// KeycloakSpecPath resolves the vendored Keycloak admin OpenAPI spec from
// KEYCLOAK_VERSION: keycloak-oapi/${KEYCLOAK_VERSION}.spec.json, anchored at
// the repository root so it works from any package's test working directory.
// An unset or empty KEYCLOAK_VERSION fails the test with the fix in the
// message rather than falling back to a pinned version.
func KeycloakSpecPath(t testing.TB) string {
	t.Helper()
	version := os.Getenv("KEYCLOAK_VERSION")
	if version == "" {
		t.Fatalf("KEYCLOAK_VERSION is not set: export it so tests can resolve keycloak-oapi/${KEYCLOAK_VERSION}.spec.json")
	}
	// thisFile = <repo>/internal/testutil/spec.go → three Dir calls reach <repo>.
	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot := filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
	return filepath.Join(repoRoot, "keycloak-oapi", version+".spec.json")
}
