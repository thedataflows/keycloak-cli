package mcpserver_test

import (
	"context"
	"net"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/mcpserver"
)

// Scenario 12: the same server also speaks streamable HTTP — a remote MCP
// client completes the initialize handshake over the handler's HTTP surface,
// lists the five tools and calls kc_resolve end-to-end against the fake
// Keycloak. RunHTTP serves on an already-bound listener so tests pin the port.
func TestServerServesToolsOverStreamableHTTP(t *testing.T) {
	client, _ := newTestClient(t)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- mcpserver.RunHTTP(ctx, client, ln) }()

	mcpClient := mcp.NewClient(&mcp.Implementation{Name: "acceptance-test", Version: "0"}, nil)
	session, err := mcpClient.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: "http://" + ln.Addr().String()}, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = session.Close() })

	result, err := session.ListTools(t.Context(), nil)
	require.NoError(t, err)
	assert.Len(t, result.Tools, 5, "streamable HTTP session must see the five tools")

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
