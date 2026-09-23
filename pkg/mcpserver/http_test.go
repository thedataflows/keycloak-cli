package mcpserver_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/kcapi"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
	"github.com/thedataflows/keycloak-cli/pkg/mcpserver"
)

// startHTTP runs RunHTTP on a loopback listener pinned by the test and returns
// its base URL. RunHTTP exits cleanly on test cleanup via the context.
func startHTTP(t *testing.T, client *kcapi.Client, fake *fakeKeycloak) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- mcpserver.RunHTTP(ctx, client, func() (manifest.Service, error) {
			return newManifestService(t, fake), nil
		}, ln)
	}()
	// Cleanups run LIFO: cancel first, then collect RunHTTP, then the listener.
	t.Cleanup(func() { _ = <-done })
	t.Cleanup(cancel)
	t.Cleanup(func() { _ = ln.Close() })
	return "http://" + ln.Addr().String()
}

// postRPC sends one raw JSON-RPC request over the streamable HTTP surface the
// way a 2026-07-28 client does: no initialize handshake, and every request
// declares its protocol version twice — the MCP-Protocol-Version header and
// the per-request _meta field, which must match.
func postRPC(t *testing.T, base string, method string, id int) (*http.Response, string) {
	t.Helper()
	const proto = "2026-07-28"
	body := fmt.Sprintf(
		`{"jsonrpc":"2.0","id":%d,"method":%q,"params":{"_meta":{"io.modelcontextprotocol/protocolVersion":%q,"io.modelcontextprotocol/clientCapabilities":{}}}}`,
		id, method, proto)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, base, strings.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	req.Header.Set("MCP-Protocol-Version", proto)
	req.Header.Set("Mcp-Method", method)
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	return res, strings.TrimSpace(string(data))
}

// rpcResult is the shape the tests read out of a JSON-RPC response.
type rpcResult struct {
	Result struct {
		Tools      []json.RawMessage `json:"tools"`
		TTLMs      int               `json:"ttlMs"`
		CacheScope string            `json:"cacheScope"`
	} `json:"result"`
}

// decodeRPC extracts the JSON-RPC response: the body is either a bare JSON
// object or an SSE stream whose data lines carry it.
func decodeRPC(t *testing.T, body string) rpcResult {
	t.Helper()
	if strings.HasPrefix(body, "{") {
		var out rpcResult
		require.NoError(t, json.Unmarshal([]byte(body), &out))
		return out
	}
	for _, line := range strings.Split(body, "\n") {
		data, ok := strings.CutPrefix(line, "data: ")
		if !ok {
			continue
		}
		var out rpcResult
		if err := json.Unmarshal([]byte(data), &out); err == nil {
			return out
		}
	}
	require.Fail(t, "no JSON-RPC response in body", "body: %s", body)
	return rpcResult{}
}

// Scenario 22: the HTTP surface speaks the 2026-07-28 stateless core — a
// client POSTs any request with its protocol version, no initialize handshake,
// no session id, no per-client state.
func TestStreamableHTTPAnswersWithoutInitialize(t *testing.T) {
	client, fake := newTestClient(t)
	base := startHTTP(t, client, fake)

	res, body := postRPC(t, base, "tools/list", 1)
	if !assert.Equal(t, http.StatusOK, res.StatusCode, body) {
		return
	}
	assert.Empty(t, res.Header.Get("Mcp-Session-Id"), "stateless servers never issue a session id")

	out := decodeRPC(t, body)
	assert.Len(t, out.Result.Tools, 9, "tools must be listed without any handshake")
}

// Scenario 23: in stateless mode the streamable transport serves POST only —
// GET and DELETE are refused with 405 and an Allow header, per RFC 9110.
func TestStreamableHTTPStatelessServesPostOnly(t *testing.T) {
	client, fake := newTestClient(t)
	base := startHTTP(t, client, fake)

	for _, method := range []string{http.MethodGet, http.MethodDelete} {
		req, err := http.NewRequestWithContext(t.Context(), method, base, nil)
		require.NoError(t, err)
		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		res.Body.Close()
		assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode, "%s must be refused", method)
		assert.Equal(t, "POST", res.Header.Get("Allow"), "%s refusal must carry Allow", method)
	}
}

// Scenario 24: the 2026-07-28 cacheable-list mechanism — the server advertises
// a ttl for tools/list so clients may cache it, and scopes it public (the
// catalog is static; kc_reload does not change the tool set).
func TestStreamableHTTPToolsListIsCacheable(t *testing.T) {
	client, fake := newTestClient(t)
	base := startHTTP(t, client, fake)

	res, body := postRPC(t, base, "tools/list", 1)
	if !assert.Equal(t, http.StatusOK, res.StatusCode, body) {
		return
	}

	out := decodeRPC(t, body)
	assert.Equal(t, 300000, out.Result.TTLMs, "tools/list must advertise a 5-minute ttl")
	assert.Equal(t, "public", out.Result.CacheScope)
}

// Scenario 12: the same server also speaks streamable HTTP — a remote MCP
// client completes the initialize handshake over the handler's HTTP surface,
// lists the library tools and calls kc_resolve end-to-end against the fake
// Keycloak. RunHTTP serves on an already-bound listener so tests pin the port.
func TestServerServesToolsOverStreamableHTTP(t *testing.T) {
	client, fake := newTestClient(t)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- mcpserver.RunHTTP(ctx, client, func() (manifest.Service, error) {
			return newManifestService(t, fake), nil
		}, ln)
	}()

	mcpClient := mcp.NewClient(&mcp.Implementation{Name: "acceptance-test", Version: "0"}, nil)
	session, err := mcpClient.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: "http://" + ln.Addr().String()}, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })

	result, err := session.ListTools(t.Context(), nil)
	require.NoError(t, err)
	assert.Len(t, result.Tools, 9, "streamable HTTP session must see the library tools")

	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{
		Name:      "kc_resolve",
		Arguments: map[string]any{"type": "realms", "name": "master"},
	})
	require.NoError(t, err)
	assert.False(t, res.IsError)
	assert.Contains(t, toolText(t, res), `"realm":"master"`,
		"kc_resolve over HTTP must return the resolved realm")

	cancel()
	assert.NoError(t, <-done, "RunHTTP must exit cleanly on ctx cancellation")
}
