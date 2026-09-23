package kcapi

import (
	"os"
	"testing"

	"github.com/thedataflows/keycloak-cli/internal/testutil"
)

func testClient(t *testing.T) *Client {
	t.Helper()
	spec, err := os.ReadFile(testutil.KeycloakSpecPath(t))
	if err != nil {
		t.Skipf("repo spec not available: %v", err)
	}
	c, err := New(Config{BaseURL: "http://test", Spec: SpecSource{Raw: spec}})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c
}

func filterOps(ops []Operation, keep func(Operation) bool) []Operation {
	var out []Operation
	for _, o := range ops {
		if keep(o) {
			out = append(out, o)
		}
	}
	return out
}
