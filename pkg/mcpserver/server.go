// Package mcpserver exposes the kcapi library to LLM agents as an MCP server
// speaking JSON-RPC over stdio. It is a thin adapter: every tool is one kcapi
// call, and the only safety rule lives in kc_invoke — state-changing verbs
// require an explicit confirm:true argument.
package mcpserver

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/thedataflows/keycloak-cli/pkg/kcapi"
)

const (
	serverName = "keycloak-cli-mcp"
	// destructiveVerbs are the HTTP methods kc_invoke refuses without
	// confirm:true. GET and HEAD pass freely.
	destructiveVerbs = "POST/PUT/PATCH/DELETE"
)

// operationsIn are kc_operations' optional ANDed filters.
type operationsIn struct {
	Resource string `json:"resource,omitempty"`
	Method   string `json:"method,omitempty"`
	Tag      string `json:"tag,omitempty"`
	Search   string `json:"search,omitempty"`
}

// invokeIn is kc_invoke's input. Exactly one resolution mode may be set:
// op alone, or resource together with verb (the bundled spec has no
// operationIds, so resource+verb is the primary mode).
type invokeIn struct {
	Op       string            `json:"op,omitempty"`
	Resource string            `json:"resource,omitempty"`
	Verb     string            `json:"verb,omitempty"`
	Realm    string            `json:"realm,omitempty"`
	Params   map[string]string `json:"params,omitempty"`
	Body     string            `json:"body,omitempty"` // raw JSON text
	// Confirm gates state-changing verbs: POST/PUT/PATCH/DELETE are rejected
	// without it, before any request is sent.
	Confirm bool `json:"confirm,omitempty"`
}

// refIn names one resource for kc_resolve and kc_neighbors.
type refIn struct {
	Type  string `json:"type"`
	Name  string `json:"name,omitempty"`
	ID    string `json:"id,omitempty"`
	Realm string `json:"realm,omitempty"`
}

// neighborsIn extends refIn with the optional EdgeFilter fields.
type neighborsIn struct {
	Type   string `json:"type"`
	Name   string `json:"name,omitempty"`
	ID     string `json:"id,omitempty"`
	Realm  string `json:"realm,omitempty"`
	Child  string `json:"child,omitempty"`
	Parent string `json:"parent,omitempty"`
}

// reloadIn is kc_reload's empty input.
type reloadIn struct{}

// New builds the MCP server for one kcapi client. All tools share the
// client's spec snapshot; agents pick up spec changes through kc_reload.
func New(client *kcapi.Client) *mcp.Server {
	srv := mcp.NewServer(&mcp.Implementation{Name: serverName, Version: "dev"}, nil)

	mcp.AddTool(srv, &mcp.Tool{
		Name: "kc_operations",
		Description: "List the Keycloak admin API operations of the loaded OpenAPI spec. " +
			"Optional filters narrow the list (ANDed): resource (path segment, e.g. users), " +
			"method (GET/POST/PUT/DELETE/PATCH), tag, search (substring over path/summary).",
		Annotations: readOnly(),
	}, func(ctx context.Context, req *mcp.CallToolRequest, in operationsIn) (*mcp.CallToolResult, any, error) {
		return unimplemented(req), nil, nil
	})

	mcp.AddTool(srv, &mcp.Tool{
		Name: "kc_invoke",
		Description: "Invoke one Keycloak admin API operation. Select it with resource+verb " +
			"(e.g. resource=users, verb=GET) or an explicit op (operationId; the bundled spec has none). " +
			"params carries path placeholders (realm, id, ...) and query parameters as key/value strings; " +
			"body is a JSON request body. " +
			"State-changing verbs (" + destructiveVerbs + ") require confirm=true; " +
			"without it the call is rejected before anything is sent. Returns the raw JSON response.",
		Annotations: destructive(),
	}, func(ctx context.Context, req *mcp.CallToolRequest, in invokeIn) (*mcp.CallToolResult, any, error) {
		return unimplemented(req), nil, nil
	})

	mcp.AddTool(srv, &mcp.Tool{
		Name: "kc_resolve",
		Description: "Resolve one Keycloak resource by name or id: type uses the plural path " +
			"vocabulary (users, clients, roles, realms, organizations, ...), name is the human " +
			"identity (username, clientId, role name, realm name). Returns the resource object.",
		Annotations: readOnly(),
	}, func(ctx context.Context, req *mcp.CallToolRequest, in refIn) (*mcp.CallToolResult, any, error) {
		return unimplemented(req), nil, nil
	})

	mcp.AddTool(srv, &mcp.Tool{
		Name: "kc_neighbors",
		Description: "List the child collections of one resolved Keycloak resource: resolves the " +
			"parent (type+name or id), then walks the relationship edges whose collection prefix " +
			"matches, fetching each child collection. Optional child/parent narrow the edge set. " +
			"Use it to discover what hangs off a resource (e.g. an organization's groups).",
		Annotations: readOnly(),
	}, func(ctx context.Context, req *mcp.CallToolRequest, in neighborsIn) (*mcp.CallToolResult, any, error) {
		return unimplemented(req), nil, nil
	})

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "kc_reload",
		Description: "Reload the OpenAPI spec from disk. Call it if the spec file changed since the server started.",
		Annotations: &mcp.ToolAnnotations{IdempotentHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, in reloadIn) (*mcp.CallToolResult, any, error) {
		return unimplemented(req), nil, nil
	})

	return srv
}

func readOnly() *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{ReadOnlyHint: true}
}

func destructive() *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{DestructiveHint: boolPtr(true)}
}

// unimplemented is the cycle-1 stub every acceptance test replaces.
func unimplemented(req *mcp.CallToolRequest) *mcp.CallToolResult {
	return errorResultf("%s is not implemented yet", req.Params.Name)
}

// errorResultf builds a tool error result the calling agent can read and
// react to.
func errorResultf(format string, args ...any) *mcp.CallToolResult {
	var res mcp.CallToolResult
	res.SetError(fmt.Errorf(format, args...))
	return &res
}

func boolPtr(b bool) *bool { return &b }
