// Package mcpserver exposes the kcapi library to LLM agents as an MCP server
// speaking JSON-RPC over stdio. It is a thin adapter: every tool is one kcapi
// call, and the only safety rule lives in kc_invoke — state-changing verbs
// require an explicit confirm:true argument.
package mcpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
		filter := kcapi.OpFilter{Resource: in.Resource, Tag: in.Tag, Search: in.Search}
		if in.Method != "" {
			filter.Method = kcapi.Verb(in.Method)
		}
		ops, err := client.Operations(filter)
		if err != nil {
			return nil, nil, err
		}
		if ops == nil {
			ops = []kcapi.Operation{}
		}
		return jsonResult(ops), nil, nil
	})

	mcp.AddTool(srv, &mcp.Tool{
		Name: "kc_invoke",
		Description: "Invoke one Keycloak admin API operation. Select it with resource+verb " +
			"(e.g. resource=users, verb=GET) or an explicit op (operationId; the bundled spec has none). " +
			"params carries path placeholders (realm, id, ...) and query parameters as key/value strings; " +
			"body is a JSON request body. " +
			"State-changing verbs (" + destructiveVerbs + ") require confirm=true; " +
			"without it the call is rejected before anything is sent. Returns the raw JSON response.",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: boolPtr(true)},
	}, func(ctx context.Context, req *mcp.CallToolRequest, in invokeIn) (*mcp.CallToolResult, any, error) {
		call, err := invokeCall(in)
		if err != nil {
			return nil, nil, err
		}
		verb, known := callVerb(client, call)
		if err := confirmGate(in.Confirm, verb, known); err != nil {
			return nil, nil, err
		}
		raw, err := client.Invoke(ctx, call)
		if err != nil {
			return nil, nil, err
		}
		return jsonResult(raw), nil, nil
	})

	mcp.AddTool(srv, &mcp.Tool{
		Name: "kc_resolve",
		Description: "Resolve one Keycloak resource by name or id: type uses the plural path " +
			"vocabulary (users, clients, roles, realms, organizations, ...), name is the human " +
			"identity (username, clientId, role name, realm name). Returns the resource object.",
		Annotations: readOnly(),
	}, func(ctx context.Context, req *mcp.CallToolRequest, in refIn) (*mcp.CallToolResult, any, error) {
		node, err := client.Resolve(ctx, kcapi.Ref{Type: in.Type, Name: in.Name, ID: in.ID, Realm: in.Realm})
		if err != nil {
			return nil, nil, toErrorResult(err)
		}
		return jsonResult(node.Resource.Data), nil, nil
	})

	mcp.AddTool(srv, &mcp.Tool{
		Name: "kc_neighbors",
		Description: "List the child collections of one resolved Keycloak resource: resolves the " +
			"parent (type+name or id), then walks the relationship edges whose collection prefix " +
			"matches, fetching each child collection. Optional child/parent narrow the edge set. " +
			"Use it to discover what hangs off a resource (e.g. an organization's groups).",
		Annotations: readOnly(),
	}, func(ctx context.Context, req *mcp.CallToolRequest, in neighborsIn) (*mcp.CallToolResult, any, error) {
		node, err := client.Resolve(ctx, kcapi.Ref{Type: in.Type, Name: in.Name, ID: in.ID, Realm: in.Realm})
		if err != nil {
			return nil, nil, toErrorResult(err)
		}
		children, edges, err := client.Neighbors(ctx, node, kcapi.EdgeFilter{Child: in.Child, Parent: in.Parent})
		if err != nil {
			return nil, nil, toErrorResult(err)
		}
		return jsonResult(neighborsOut(children, edges)), nil, nil
	})

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "kc_reload",
		Description: "Reload the OpenAPI spec from disk. Call it if the spec file changed since the server started.",
		Annotations: &mcp.ToolAnnotations{IdempotentHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, in reloadIn) (*mcp.CallToolResult, any, error) {
		if err := client.Reload(ctx); err != nil {
			return nil, nil, toErrorResult(err)
		}
		return jsonResult(map[string]string{"status": "reloaded"}), nil, nil
	})

	return srv
}

// invokeCall maps kc_invoke's typed input onto a kcapi.Call, mirroring the
// CLI's resolution modes: op alone, or resource+verb together. Body stays nil
// when empty: a typed-nil json.RawMessage inside interface{} would make kcapi
// treat the call as carrying a body.
func invokeCall(in invokeIn) (kcapi.Call, error) {
	call := kcapi.Call{Realm: in.Realm, Params: kcapi.P(in.Params)}
	if in.Body != "" {
		call.Body = json.RawMessage(in.Body)
	}
	switch {
	case in.Op != "" && in.Resource != "":
		return kcapi.Call{}, fmt.Errorf("use either op or resource with verb, not both")
	case in.Resource != "":
		if in.Verb == "" {
			return kcapi.Call{}, fmt.Errorf("resource %q requires verb", in.Resource)
		}
		call.Resource = in.Resource
		call.Verb = kcapi.Verb(in.Verb)
	case in.Op != "":
		call.Op = in.Op
	default:
		return kcapi.Call{}, fmt.Errorf("provide op or resource (with verb) to select an operation")
	}
	return call, nil
}

// readOnlyVerbs pass the confirm gate freely; every other verb — including
// unknown ones — changes state and needs confirm:true.
var readOnlyVerbs = map[string]bool{"GET": true, "HEAD": true}

// confirmGate rejects state-changing invokes that lack confirm:true, before
// anything is sent. An unknown verb (op mode whose operation is not in the
// spec) is treated as state-changing: the safe default.
func confirmGate(confirmed bool, verb string, known bool) error {
	if confirmed || known && readOnlyVerbs[verb] {
		return nil
	}
	return fmt.Errorf(
		"refused: %s changes Keycloak state; retry the same call with \"confirm\": true to send it",
		verbLabel(verb, known))
}

func verbLabel(verb string, known bool) string {
	if known {
		return verb
	}
	return "this operation (unknown verb)"
}

// callVerb reports the call's HTTP verb: the resource+verb mode's own verb,
// or the verb of the spec operation named by op mode. known is false when the
// verb cannot be determined (unknown operationId).
func callVerb(client *kcapi.Client, call kcapi.Call) (verb string, known bool) {
	if call.Resource != "" {
		return strings.ToUpper(strings.TrimSpace(string(call.Verb))), true
	}
	ops, err := client.Operations(kcapi.OpFilter{Search: call.Op})
	if err != nil {
		return "", false
	}
	for _, op := range ops {
		if op.ID == call.Op {
			return strings.ToUpper(strings.TrimSpace(string(op.Verb))), true
		}
	}
	return "", false
}

// neighborsOut shapes a Neighbors walk as the children's resource objects
// plus the edges used, with empty slices (never null) for stable JSON.
func neighborsOut(children []kcapi.Node, edges []kcapi.Edge) map[string]any {
	kids := make([]map[string]any, 0, len(children))
	for _, child := range children {
		kids = append(kids, child.Resource.Data)
	}
	if edges == nil {
		edges = []kcapi.Edge{}
	}
	return map[string]any{"children": kids, "edges": edges}
}

// toErrorResult prefixes kcapi's typed errors with their taxonomy kind, so
// the calling agent sees e.g. "[not_found] ..." and can react.
func toErrorResult(err error) error {
	var kerr *kcapi.Error
	if errors.As(err, &kerr) {
		return fmt.Errorf("[%s] %w", kerr.Kind, err)
	}
	return err
}

// Run serves the MCP server over stdio until the client disconnects or ctx
// is cancelled. stdout carries only protocol frames; logging stays on stderr.
func Run(ctx context.Context, client *kcapi.Client) error {
	return New(client).Run(ctx, &mcp.StdioTransport{})
}

func readOnly() *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{ReadOnlyHint: true}
}

// errorResultf builds a tool error result the calling agent can read and
// react to.
func errorResultf(format string, args ...any) *mcp.CallToolResult {
	var res mcp.CallToolResult
	res.SetError(fmt.Errorf(format, args...))
	return &res
}

// jsonResult marshals v as the tool result's JSON text content.
func jsonResult(v any) *mcp.CallToolResult {
	data, err := json.Marshal(v)
	if err != nil {
		return errorResultf("marshal result: %v", err)
	}
	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: string(data)}}}
}

func boolPtr(b bool) *bool { return &b }
