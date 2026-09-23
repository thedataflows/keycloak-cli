// Package mcpserver_test holds the acceptance tests for the MCP server: an
// in-memory MCP client session against a real kcapi client whose HTTP traffic
// lands on a fake Keycloak. The "user" in these scenarios is the calling
// agent; every test asserts what it sees through the MCP protocol.
package mcpserver_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/kcapi"
	"github.com/thedataflows/keycloak-cli/pkg/mcpserver"
	"golang.org/x/oauth2"
)

// fakeKeycloak is the stand-in Keycloak: it serves the endpoints the vendored
// spec's resolve and neighbor walks hit and records every request.
type fakeKeycloak struct {
	server *httptest.Server
	mu     sync.Mutex
	reqs   []string
}

func (f *fakeKeycloak) requests() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.reqs...)
}

func (f *fakeKeycloak) record(r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.reqs = append(f.reqs, r.Method+" "+r.URL.Path)
}

func (f *fakeKeycloak) jsonReply(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f.record(r)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}
}

func newFakeKeycloak(t *testing.T) *fakeKeycloak {
	t.Helper()
	fake := &fakeKeycloak{}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /admin/realms", fake.jsonReply(`[{"realm":"master","id":"master","displayName":"Master"},{"realm":"acme","id":"acme"}]`))
	mux.HandleFunc("GET /admin/realms/{realm}/organizations", fake.jsonReply(`[{"id":"org-1","name":"acme-org"}]`))
	mux.HandleFunc("GET /admin/realms/{realm}/organizations/{orgid}/groups", fake.jsonReply(`[{"id":"g1","name":"eng"}]`))
	mux.HandleFunc("GET /admin/realms/{realm}/users", fake.jsonReply(`[{"id":"u1","username":"alice"}]`))
	mux.HandleFunc("POST /admin/realms/{realm}/users", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		fake.record(r)
		t.Logf("fake keycloak POST users body: %s", body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"u2"}`))
	})
	fake.server = httptest.NewServer(mux)
	t.Cleanup(fake.server.Close)
	return fake
}

// newTestClient builds a real kcapi client against the fake Keycloak, loading
// the vendored production spec — the acceptance layer stays end-to-end.
func newTestClient(t *testing.T) (*kcapi.Client, *fakeKeycloak) {
	t.Helper()
	fake := newFakeKeycloak(t)
	client, err := kcapi.New(kcapi.Config{
		BaseURL: fake.server.URL,
		Spec:    kcapi.SpecSource{Path: "../../keycloak-oapi/26.6.2.spec.json"},
		Timeout: 5 * time.Second,
		Auth:    staticTokenProvider("test-token"),
	})
	require.NoError(t, err)
	return client, fake
}

// staticTokenProvider is a fixed-token kcapi.TokenProvider: the acceptance
// tests exercise tool behavior, not token exchange.
type staticTokenProvider string

func (p staticTokenProvider) AccessToken(_ context.Context, _, _, _ string) (string, error) {
	return string(p), nil
}

func (p staticTokenProvider) PasswordToken(context.Context, string, string, string, string) (oauth2.Token, error) {
	return oauth2.Token{}, nil
}

func (p staticTokenProvider) ClientCredentialsToken(context.Context, string, string, string, string) (oauth2.Token, error) {
	return oauth2.Token{}, nil
}

func (p staticTokenProvider) SetEnvToken(string, string, string) error { return nil }

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

// Scenario 2: kc_operations lists spec operations narrowed by resource and
// method filters, as JSON the agent can parse.
func TestKcOperationsFiltersByResourceAndMethod(t *testing.T) {
	client, _ := newTestClient(t)
	session := newSession(t, mcpserver.New(client))

	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "kc_operations", Arguments: map[string]any{"resource": "users", "method": "GET"}})
	require.NoError(t, err)
	require.False(t, res.IsError, toolText(t, res))

	var ops []map[string]any
	require.NoError(t, json.Unmarshal([]byte(toolText(t, res)), &ops), "result must be a JSON array")
	require.NotEmpty(t, ops, "vendored spec has users operations")
	for _, op := range ops {
		assert.Equal(t, "users", op["Resource"])
		assert.Equal(t, "GET", op["Verb"])
	}
}

// Scenario 3: kc_invoke GET needs no confirmation and reaches Keycloak.
func TestKcInvokeGetWithoutConfirm(t *testing.T) {
	client, fake := newTestClient(t)
	session := newSession(t, mcpserver.New(client))

	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "kc_invoke", Arguments: map[string]any{
		"resource": "users",
		"verb":     "GET",
		"realm":    "master",
	}})
	require.NoError(t, err)
	require.False(t, res.IsError, toolText(t, res))

	var users []map[string]any
	require.NoError(t, json.Unmarshal([]byte(toolText(t, res)), &users), "result must be the JSON response body")
	require.Len(t, users, 1)
	assert.Equal(t, "alice", users[0]["username"])
	assert.Contains(t, fake.requests(), "GET /admin/realms/master/users")
}

// Scenario 4: kc_invoke POST without confirm is rejected before anything is
// sent.
func TestKcInvokeWriteWithoutConfirmIsRejected(t *testing.T) {
	client, fake := newTestClient(t)
	session := newSession(t, mcpserver.New(client))

	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "kc_invoke", Arguments: map[string]any{
		"resource": "users",
		"verb":     "POST",
		"realm":    "master",
		"body":     `{"username":"bob"}`,
	}})
	require.NoError(t, err, "rejection is a tool error, not a transport failure")
	require.True(t, res.IsError, "POST without confirm must be a tool error, got: %v", res.Content)
	assert.Contains(t, toolText(t, res), "confirm")
	for _, req := range fake.requests() {
		assert.NotEqual(t, "POST /admin/realms/master/users", req, "nothing may reach Keycloak without confirm")
	}
}

// Scenario 5: kc_invoke POST with confirm:true performs the write.
func TestKcInvokeWriteWithConfirmPasses(t *testing.T) {
	client, fake := newTestClient(t)
	session := newSession(t, mcpserver.New(client))

	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "kc_invoke", Arguments: map[string]any{
		"resource": "users",
		"verb":     "POST",
		"realm":    "master",
		"body":     `{"username":"bob"}`,
		"confirm":  true,
	}})
	require.NoError(t, err)
	require.False(t, res.IsError, toolText(t, res))
	assert.Contains(t, fake.requests(), "POST /admin/realms/master/users")
	assert.Contains(t, toolText(t, res), `"id":"u2"`)
}

// toolText returns the first text content of a tool result.
func toolText(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	require.NotEmpty(t, res.Content, "tool result has no content")
	text, ok := res.Content[0].(*mcp.TextContent)
	require.True(t, ok, "first content is %T, want TextContent", res.Content[0])
	return text.Text
}
