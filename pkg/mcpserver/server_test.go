// Package mcpserver_test holds the acceptance tests for the MCP server: an
// in-memory MCP client session against a real kcapi client whose HTTP traffic
// lands on a fake Keycloak. The "user" in these scenarios is the calling
// agent; every test asserts what it sees through the MCP protocol.
package mcpserver_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/kcapi"
	"github.com/thedataflows/keycloak-cli/pkg/mcpserver"
)

// newFakeKeycloak starts an httptest server standing in for Keycloak's admin
// API. It serves the minimal endpoints the vendored spec's resolve and
// neighbor walks hit, and records every request for write-assertions.
func newFakeKeycloak(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /admin/realms", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, `[{"realm":"master","id":"master","displayName":"Master"},{"realm":"acme","id":"acme"}]`)
	})
	mux.HandleFunc("GET /admin/realms/{realm}/organizations", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, `[{"id":"org-1","name":"acme-org"}]`)
	})
	mux.HandleFunc("GET /admin/realms/{realm}/organizations/{orgid}/groups", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, `[{"id":"g1","name":"eng"}]`)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func writeJSON(w http.ResponseWriter, body string) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(body))
}

// newTestClient builds a real kcapi client against the fake Keycloak, loading
// the vendored production spec — the acceptance layer stays end-to-end.
func newTestClient(t *testing.T) (*kcapi.Client, *httptest.Server) {
	t.Helper()
	fake := newFakeKeycloak(t)
	client, err := kcapi.New(kcapi.Config{
		BaseURL: fake.URL,
		Spec:    kcapi.SpecSource{Path: "../../keycloak-oapi/26.6.2.spec.json"},
		Timeout: 5 * time.Second,
	})
	require.NoError(t, err)
	return client, fake
}

// newSession connects an in-memory MCP client to srv, completing the MCP
// initialize handshake.
func newSession(t *testing.T, srv *mcp.Server) *mcp.ClientSession {
	t.Helper()
	ctx := context.Background()
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	go func() { _ = srv.Run(ctx, serverTransport) }() // blocks until the session closes
	client := mcp.NewClient(&mcp.Implementation{Name: "acceptance-test", Version: "0"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })
	return session
}

// Scenario 1: the server exposes exactly the five approved tools.
func TestServerExposesFiveTools(t *testing.T) {
	client, _ := newTestClient(t)
	session := newSession(t, mcpserver.New(client))

	result, err := session.ListTools(t.Context(), nil)
	require.NoError(t, err)

	names := make([]string, 0, len(result.Tools))
	for _, tool := range result.Tools {
		names = append(names, tool.Name)
	}
	assert.ElementsMatch(t,
		[]string{"kc_invoke", "kc_neighbors", "kc_operations", "kc_reload", "kc_resolve"},
		names)
}
