// Package mcpserver_test holds the acceptance tests for the MCP server: an
// in-memory MCP client session against a real kcapi client whose HTTP traffic
// lands on a fake Keycloak. The "user" in these scenarios is the calling
// agent; every test asserts what it sees through the MCP protocol.
package mcpserver_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/internal/testutil"
	"github.com/thedataflows/keycloak-cli/pkg/kcapi"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
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
	mux.HandleFunc("GET /admin/realms/{realm}", func(w http.ResponseWriter, r *http.Request) {
		// Echo the requested realm: a real Keycloak's realm representation
		// names the realm it serves, and the neighbor walk anchors on it.
		fake.record(r)
		name := r.PathValue("realm")
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"realm":%q,"id":%q,"displayName":%q,"enabled":true}`, name, name, "Realm "+name)
	})
	mux.HandleFunc("GET /admin/realms/{realm}/organizations", fake.jsonReply(`[{"id":"org-1","alias":"acme-org","name":"Acme Org"}]`))
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
		Spec:    kcapi.SpecSource{Path: testutil.KeycloakSpecPath(t)},
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

// newManifestService builds a manifest service pointing at the same fake
// Keycloak as the kcapi client, so the manifest tools and the raw tools see
// one backend. The static token provider satisfies auth.Service, keeping the
// fake auth-free.
func newManifestService(t *testing.T, fake *fakeKeycloak) manifest.Service {
	t.Helper()
	svc, err := manifest.NewService(manifest.Config{
		BaseURL:  fake.server.URL,
		SpecPath: testutil.KeycloakSpecPath(t),
		Timeout:  5 * time.Second,
		Auth:     staticTokenProvider("test-token"),
	})
	require.NoError(t, err)
	return svc
}

// newTestServer builds the MCP server with the manifest service alongside the
// kcapi client — the production wiring, against the fake Keycloak.
func newTestServer(t *testing.T, client *kcapi.Client, fake *fakeKeycloak) *mcp.Server {
	t.Helper()
	srv, err := mcpserver.New(client, func() (manifest.Service, error) {
		return newManifestService(t, fake), nil
	})
	require.NoError(t, err)
	return srv
}

// Scenario 1: the server exposes the kcapi library surface — the five raw
// tools plus the paging walk, the edge vocabulary, and both manifest
// directions.
func TestServerExposesAllLibraryTools(t *testing.T) {
	client, fake := newTestClient(t)
	session := newSession(t, newTestServer(t, client, fake))

	result, err := session.ListTools(t.Context(), nil)
	require.NoError(t, err)

	names := make([]string, 0, len(result.Tools))
	for _, tool := range result.Tools {
		names = append(names, tool.Name)
	}
	assert.ElementsMatch(t,
		[]string{
			"kc_invoke", "kc_neighbors", "kc_operations", "kc_reload", "kc_resolve",
			"kc_list", "kc_edges", "kc_fetch", "kc_apply",
		},
		names)
}

// Scenario 2: kc_operations lists spec operations narrowed by resource and
// method filters, as JSON the agent can parse.
func TestKcOperationsFiltersByResourceAndMethod(t *testing.T) {
	client, fake := newTestClient(t)
	session := newSession(t, newTestServer(t, client, fake))

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
	session := newSession(t, newTestServer(t, client, fake))

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
	session := newSession(t, newTestServer(t, client, fake))

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
	session := newSession(t, newTestServer(t, client, fake))

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

// Scenario 6: kc_resolve resolves a resource by name through the real kcapi
// flow and returns the resource object.
func TestKcResolveByName(t *testing.T) {
	client, fake := newTestClient(t)
	session := newSession(t, newTestServer(t, client, fake))

	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "kc_resolve", Arguments: map[string]any{
		"type":  "users",
		"name":  "alice",
		"realm": "master",
	}})
	require.NoError(t, err)
	require.False(t, res.IsError, toolText(t, res))

	var user map[string]any
	require.NoError(t, json.Unmarshal([]byte(toolText(t, res)), &user))
	assert.Equal(t, "alice", user["username"])
	assert.Contains(t, fake.requests(), "GET /admin/realms/master/users")
}

// Scenario 7: kc_neighbors walks an organization's child collection and
// returns the fetched children plus the edges used.
func TestKcNeighborsWalksChildCollection(t *testing.T) {
	client, fake := newTestClient(t)
	session := newSession(t, newTestServer(t, client, fake))

	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "kc_neighbors", Arguments: map[string]any{
		"type":  "organizations",
		"name":  "acme-org",
		"realm": "master",
		"child": "groups",
	}})
	require.NoError(t, err)
	require.False(t, res.IsError, toolText(t, res))

	var out struct {
		Children []map[string]any `json:"children"`
		Edges    []map[string]any `json:"edges"`
	}
	require.NoError(t, json.Unmarshal([]byte(toolText(t, res)), &out))
	require.Len(t, out.Children, 1)
	assert.Equal(t, "eng", out.Children[0]["name"])
	require.NotEmpty(t, out.Edges)
	assert.Equal(t, "groups", out.Edges[0]["Child"])
	assert.Contains(t, fake.requests(), "GET /admin/realms/master/organizations/acme-org/groups")
}

// Scenario 8: kc_reload completes and keeps the tools working afterwards.
func TestKcReloadCompletes(t *testing.T) {
	client, fake := newTestClient(t)
	session := newSession(t, newTestServer(t, client, fake))

	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "kc_reload"})
	require.NoError(t, err)
	require.False(t, res.IsError, toolText(t, res))

	after, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "kc_operations", Arguments: map[string]any{"resource": "users", "method": "GET"}})
	require.NoError(t, err)
	assert.False(t, after.IsError, toolText(t, after))
}

// Scenario 9: kcapi failures surface as tool error results carrying kind and
// status, so the agent can self-correct.
func TestKcapiErrorsBecomeToolErrors(t *testing.T) {
	client, fake := newTestClient(t)
	session := newSession(t, newTestServer(t, client, fake))

	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "kc_resolve", Arguments: map[string]any{
		"type": "nosuchtype",
		"name": "x",
	}})
	require.NoError(t, err, "domain errors are tool errors, not transport failures")
	require.True(t, res.IsError)
	text := toolText(t, res)
	assert.Contains(t, text, "unknown resource type")
	assert.Contains(t, text, "nosuchtype")
}

// Scenario 10: kc_resolve resolves a realm by its own name. The realm is the
// one resource whose identity field IS the {realm} path parameter, so the
// lookup is the single-resource GET /admin/realms/{name} — one exact request,
// no collection search.
func TestKcResolveRealmByName(t *testing.T) {
	client, fake := newTestClient(t)
	session := newSession(t, newTestServer(t, client, fake))

	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "kc_resolve", Arguments: map[string]any{
		"type":  "realms",
		"name":  "acme",
		"realm": "acme",
	}})
	require.NoError(t, err)
	require.False(t, res.IsError, toolText(t, res))

	var realm map[string]any
	require.NoError(t, json.Unmarshal([]byte(toolText(t, res)), &realm))
	assert.Equal(t, "acme", realm["realm"])
	assert.Contains(t, fake.requests(), "GET /admin/realms/acme")
}

// Scenario 11: kc_neighbors from a realm node lists the child collections
// hanging off the realm root, anchored by the realm's own identifier.
func TestKcNeighborsFromRealmNode(t *testing.T) {
	client, fake := newTestClient(t)
	session := newSession(t, newTestServer(t, client, fake))

	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "kc_neighbors", Arguments: map[string]any{
		"type":  "realms",
		"name":  "master",
		"realm": "master",
		"child": "users",
	}})
	require.NoError(t, err)
	require.False(t, res.IsError, toolText(t, res))

	var out struct {
		Children []map[string]any `json:"children"`
		Edges    []map[string]any `json:"edges"`
	}
	require.NoError(t, json.Unmarshal([]byte(toolText(t, res)), &out))
	require.Len(t, out.Children, 1)
	assert.Equal(t, "alice", out.Children[0]["username"])
	require.NotEmpty(t, out.Edges)
	assert.Equal(t, "users", out.Edges[0]["Child"])
	assert.Contains(t, fake.requests(), "GET /admin/realms/master/users")
}

// Scenario 14: kc_list walks a GET collection and returns every item as a
// JSON array — the agent does not hand-page kc_invoke.
func TestKcListWalksCollection(t *testing.T) {
	client, fake := newTestClient(t)
	session := newSession(t, newTestServer(t, client, fake))

	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "kc_list", Arguments: map[string]any{
		"resource": "users",
		"verb":     "GET",
		"realm":    "master",
	}})
	require.NoError(t, err)
	require.False(t, res.IsError, toolText(t, res))

	var users []map[string]any
	require.NoError(t, json.Unmarshal([]byte(toolText(t, res)), &users), "result must be a JSON array")
	require.Len(t, users, 1)
	assert.Equal(t, "alice", users[0]["username"])
	assert.Contains(t, fake.requests(), "GET /admin/realms/master/users")
}

// Scenario 15: kc_list is read-only by construction — it refuses any verb a
// GET collection walk cannot serve, before anything is sent.
func TestKcListRefusesNonGet(t *testing.T) {
	client, fake := newTestClient(t)
	session := newSession(t, newTestServer(t, client, fake))

	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "kc_list", Arguments: map[string]any{
		"resource": "users",
		"verb":     "POST",
		"realm":    "master",
	}})
	require.NoError(t, err, "refusal is a tool error, not a transport failure")
	require.True(t, res.IsError, "non-GET must be refused, got: %v", res.Content)
	assert.Contains(t, toolText(t, res), "GET")
	assert.Contains(t, toolText(t, res), "kc_invoke")
	for _, req := range fake.requests() {
		assert.NotEqual(t, "POST /admin/realms/master/users", req, "nothing may reach Keycloak")
	}
}

// Scenario 16: kc_edges exposes the relationship vocabulary so an agent can
// plan walks without guessing path shapes.
func TestKcEdgesListsVocabulary(t *testing.T) {
	client, fake := newTestClient(t)
	session := newSession(t, newTestServer(t, client, fake))

	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "kc_edges"})
	require.NoError(t, err)
	require.False(t, res.IsError, toolText(t, res))

	var edges []map[string]any
	require.NoError(t, json.Unmarshal([]byte(toolText(t, res)), &edges), "result must be a JSON array")
	require.NotEmpty(t, edges, "vendored spec implies parent-child edges")

	kinds := make(map[string]bool, len(edges))
	for _, e := range edges {
		kinds[e["Parent"].(string)+">"+e["Child"].(string)] = true
	}
	assert.Contains(t, kinds, "realms>users", "realm-rooted user collection must be an edge")
}

// Scenario 17: kc_fetch returns the manifest export — resources plus readable
// failures — that kc_apply can consume.
func TestKcFetchReturnsManifestReport(t *testing.T) {
	client, fake := newTestClient(t)
	session := newSession(t, newTestServer(t, client, fake))

	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "kc_fetch", Arguments: map[string]any{
		"realm":     "master",
		"resources": "realm,user",
	}})
	require.NoError(t, err)
	require.False(t, res.IsError, toolText(t, res))

	var out struct {
		Resources     []map[string]any `json:"resources"`
		Relationships []map[string]any `json:"relationships"`
		Failures      []map[string]any `json:"failures"`
	}
	require.NoError(t, json.Unmarshal([]byte(toolText(t, res)), &out))
	types := make(map[string]bool, len(out.Resources))
	for _, r := range out.Resources {
		types[r["type"].(string)] = true
	}
	assert.True(t, types["realm"] && types["user"], "fetch must return realm and user, got %v", types)
	assert.Empty(t, out.Failures, "fake serves both collections")
	assert.Contains(t, fake.requests(), "GET /admin/realms/master/users")
}

// Scenario 18: kc_apply is gated like kc_invoke — without confirm:true
// nothing is sent.
func TestKcApplyRequiresConfirm(t *testing.T) {
	client, fake := newTestClient(t)
	session := newSession(t, newTestServer(t, client, fake))

	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "kc_apply", Arguments: map[string]any{
		"resources": []map[string]any{
			{"type": "user", "realm": "master", "data": map[string]any{"username": "bob"}},
		},
	}})
	require.NoError(t, err, "refusal is a tool error, not a transport failure")
	require.True(t, res.IsError, "apply without confirm must be a tool error, got: %v", res.Content)
	assert.Contains(t, toolText(t, res), "confirm")
	for _, req := range fake.requests() {
		assert.NotEqual(t, "POST /admin/realms/master/users", req, "nothing may reach Keycloak without confirm")
	}
}

// Scenario 19: kc_apply with confirm:true applies the manifest and returns
// the per-resource report.
func TestKcApplyAppliesWithConfirm(t *testing.T) {
	client, fake := newTestClient(t)
	session := newSession(t, newTestServer(t, client, fake))

	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "kc_apply", Arguments: map[string]any{
		"resources": []map[string]any{
			{"type": "user", "realm": "master", "data": map[string]any{"username": "bob"}},
		},
		"confirm": true,
	}})
	require.NoError(t, err)
	require.False(t, res.IsError, toolText(t, res))

	var report struct {
		Results []struct {
			Resource string `json:"resource"`
			Action   string `json:"action"`
			Status   int    `json:"status"`
		} `json:"results"`
		Failed int `json:"failed"`
	}
	require.NoError(t, json.Unmarshal([]byte(toolText(t, res)), &report))
	require.Len(t, report.Results, 1)
	assert.Equal(t, "user", report.Results[0].Resource)
	assert.Equal(t, "created", report.Results[0].Action)
	assert.Zero(t, report.Failed)
	assert.Contains(t, fake.requests(), "POST /admin/realms/master/users")
}

// Scenario 20: kc_reload rebuilds the manifest service too — otherwise
// kc_fetch/kc_apply would keep serving the stale spec after a reload.
func TestKcReloadRebuildsManifestService(t *testing.T) {
	client, fake := newTestClient(t)
	builds := 0
	srv, err := mcpserver.New(client, func() (manifest.Service, error) {
		builds++
		return newManifestService(t, fake), nil
	})
	require.NoError(t, err)
	session := newSession(t, srv)

	require.Equal(t, 1, builds, "manifest service builds once at startup")

	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{Name: "kc_reload"})
	require.NoError(t, err)
	require.False(t, res.IsError, toolText(t, res))

	assert.Equal(t, 2, builds, "kc_reload must rebuild the manifest service")
}

// Scenario 21: the startup guide tells a fresh operator what the server
// exposes, what the safety rule is, and how to wire it into the common
// harnesses — without reading the docs.
func TestGuideNamesToolsAndSetup(t *testing.T) {
	stdio := mcpserver.Guide("stdio", "")
	for _, tool := range []string{"kc_operations", "kc_invoke", "kc_list", "kc_resolve", "kc_neighbors", "kc_edges", "kc_fetch", "kc_apply", "kc_reload"} {
		assert.Contains(t, stdio, tool, "guide must name every tool")
	}
	assert.Contains(t, stdio, `"confirm": true`, "guide must state the safety rule")
	assert.Contains(t, stdio, "claude mcp add", "guide must show the Claude Code wiring")
	assert.Contains(t, stdio, `"mcpServers"`, "guide must show the generic mcp.json wiring")

	httpGuide := mcpserver.Guide("http", "127.0.0.1:8081")
	assert.Contains(t, httpGuide, "http://127.0.0.1:8081", "guide must name the HTTP endpoint")
	assert.Contains(t, httpGuide, "loopback", "guide must carry the no-auth caveat")
}
