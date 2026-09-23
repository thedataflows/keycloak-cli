// Package mcpserver exposes the kcapi library to LLM agents as an MCP server
// speaking JSON-RPC over stdio or streamable HTTP. It is a thin adapter: every
// tool is one library call, and the safety rules are kc_list (GET-only walks),
// kc_invoke and kc_apply (state-changing work requires confirm:true).
package mcpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/thedataflows/keycloak-cli/pkg/kcapi"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
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

// listIn is kc_list's input: the kc_invoke selection shape (resource+verb, or
// op) minus everything a read-only collection walk cannot use.
type listIn struct {
	Op       string            `json:"op,omitempty"`
	Resource string            `json:"resource,omitempty"`
	Verb     string            `json:"verb,omitempty"`
	Realm    string            `json:"realm,omitempty"`
	Params   map[string]string `json:"params,omitempty"`
}

// fetchIn is kc_fetch's input, mirroring manifest.FetchQuery.
type fetchIn struct {
	Realm                string `json:"realm,omitempty"`
	Resources            string `json:"resources,omitempty"`
	Search               string `json:"search,omitempty"`
	Max                  int    `json:"max,omitempty"`
	Parent               string `json:"parent,omitempty"`
	IncludeRelationships bool   `json:"includeRelationships,omitempty"`
	Depth                int    `json:"depth,omitempty"`
	Filter               string `json:"filter,omitempty"`
	ExactMatch           bool   `json:"exactMatch,omitempty"`
	FullRepresentation   bool   `json:"fullRepresentation,omitempty"`
}

// fetchOut shapes a FetchReport for the tool result: failures carry their
// error text (the error value itself does not marshal), slices stay
// non-null.
type fetchOut struct {
	Resources     []manifest.Resource              `json:"resources"`
	Relationships []manifest.RelationshipOperation `json:"relationships"`
	Failures      []fetchFailureOut                `json:"failures"`
}

type fetchFailureOut struct {
	Resource string `json:"resource"`
	Detail   string `json:"detail,omitempty"`
	NotFound bool   `json:"notFound,omitempty"`
	Error    string `json:"error,omitempty"`
}

// applyIn is kc_apply's input: the manifest shapes kc_fetch returns, plus
// options and the confirm gate.
type applyIn struct {
	Resources     []manifest.Resource              `json:"resources,omitempty"`
	Relationships []manifest.RelationshipOperation `json:"relationships,omitempty"`
	Options       applyOptionsIn                   `json:"options,omitempty"`
	// Confirm gates every apply, dry-run included: locating existing resources
	// already sends requests, so the agent states intent once, up front.
	Confirm bool `json:"confirm,omitempty"`
}

type applyOptionsIn struct {
	DryRun          bool `json:"dryRun,omitempty"`
	Delete          bool `json:"delete,omitempty"`
	ContinueOnError bool `json:"continueOnError,omitempty"`
	Reconcile       bool `json:"reconcile,omitempty"`
}

// manifestDeps owns the manifest service backing kc_fetch and kc_apply: built
// once at startup and rebuilt by kc_reload. mu serializes service access —
// manifest.Service caches identity maps lazily without locking, so concurrent
// tool calls on one service would race. The kcapi tools stay fully parallel.
type manifestDeps struct {
	mu    sync.Mutex
	build func() (manifest.Service, error)
	svc   manifest.Service
}

func (d *manifestDeps) reload() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	svc, err := d.build()
	if err != nil {
		return fmt.Errorf("rebuild manifest service: %w", err)
	}
	d.svc = svc
	return nil
}

func (d *manifestDeps) fetch(ctx context.Context, q manifest.FetchQuery) (manifest.FetchReport, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.svc.Fetch(ctx, q)
}

func (d *manifestDeps) apply(ctx context.Context, resources []manifest.Resource, relationships []manifest.RelationshipOperation, opts manifest.ApplyOptions) (manifest.ApplyReport, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.svc.Apply(ctx, resources, relationships, opts)
}

// New builds the MCP server exposing the kcapi library: the raw client tools
// plus the manifest tools, which run on a manifest.Service built by
// newManifest (and rebuilt by kc_reload). All tools share the client's spec
// snapshot; agents pick up spec changes through kc_reload.
func New(client *kcapi.Client, newManifest func() (manifest.Service, error)) (*mcp.Server, error) {
	svc, err := newManifest()
	if err != nil {
		return nil, fmt.Errorf("build manifest service: %w", err)
	}
	deps := &manifestDeps{build: newManifest, svc: svc}

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
			"Use it to discover what hangs off a resource (e.g. a realm's users, an organization's groups). " +
			"A realm node is the graph root: resolve type \"realms\" by realm name, then walk from there.",
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
		Name: "kc_list",
		Description: "Fetch ALL pages of one GET collection operation in a single call: Keycloak list " +
			"endpoints page small (20-100 items); kc_list follows first/max until the collection ends and " +
			"returns every item as a JSON array. Select the operation like kc_invoke (resource+verb or op; " +
			"verb must be GET — kc_list is read-only by construction). Use this instead of hand-paging kc_invoke.",
		Annotations: readOnly(),
	}, func(ctx context.Context, req *mcp.CallToolRequest, in listIn) (*mcp.CallToolResult, any, error) {
		call, err := invokeCall(invokeIn{Op: in.Op, Resource: in.Resource, Verb: in.Verb, Realm: in.Realm, Params: in.Params})
		if err != nil {
			return nil, nil, err
		}
		verb, known := callVerb(client, call)
		if !known || verb != "GET" {
			return nil, nil, fmt.Errorf(
				"refused: kc_list only walks GET collections (%s); use kc_invoke for other verbs",
				verbLabel(verb, known))
		}
		items, err := client.ListAll(ctx, call)
		if err != nil {
			return nil, nil, toErrorResult(err)
		}
		if items == nil {
			items = []json.RawMessage{}
		}
		return jsonResult(items), nil, nil
	})

	mcp.AddTool(srv, &mcp.Tool{
		Name: "kc_edges",
		Description: "List every parent→child relationship edge implied by the spec: parent and child in the " +
			"plural path vocabulary (realms, users, clients, roles, ...), the path template, cardinality, " +
			"cascadeOwner, and a display name. Use it to plan kc_resolve/kc_neighbors walks or kc_fetch resource lists.",
		Annotations: readOnly(),
	}, func(ctx context.Context, req *mcp.CallToolRequest, in struct{}) (*mcp.CallToolResult, any, error) {
		edges := client.Edges()
		if edges == nil {
			edges = []kcapi.Edge{}
		}
		return jsonResult(edges), nil, nil
	})

	mcp.AddTool(srv, &mcp.Tool{
		Name: "kc_fetch",
		Description: "Fetch a realm-scoped export of Keycloak resources: returns {resources, relationships, " +
			"failures} where resources carry {type, realm, data} ready to feed back into kc_apply. " +
			"resources is a comma list of singular types (realm, user, client, group, role, ...). " +
			"Optional: search/filter/exactMatch narrow the set, max caps items per collection, parent scopes " +
			"to one parent, depth>1 or includeRelationships add cross-resource relationships, " +
			"fullRepresentation requests briefRepresentation=false. Read-only.",
		Annotations: readOnly(),
	}, func(ctx context.Context, req *mcp.CallToolRequest, in fetchIn) (*mcp.CallToolResult, any, error) {
		report, err := deps.fetch(ctx, manifest.FetchQuery{
			Realm:                in.Realm,
			Resources:            in.Resources,
			Search:               in.Search,
			Max:                  in.Max,
			Parent:               in.Parent,
			IncludeRelationships: in.IncludeRelationships,
			Depth:                in.Depth,
			Filter:               in.Filter,
			ExactMatch:           in.ExactMatch,
			FullRepresentation:   in.FullRepresentation,
		})
		if err != nil {
			return nil, nil, toErrorResult(err)
		}
		out := fetchOut{
			Resources:     report.Resources,
			Relationships: report.Relationships,
		}
		if out.Resources == nil {
			out.Resources = []manifest.Resource{}
		}
		if out.Relationships == nil {
			out.Relationships = []manifest.RelationshipOperation{}
		}
		out.Failures = make([]fetchFailureOut, 0, len(report.Failures))
		for _, f := range report.Failures {
			msg := ""
			if f.Err != nil {
				msg = f.Err.Error()
			}
			out.Failures = append(out.Failures, fetchFailureOut{Resource: f.Resource, Detail: f.Detail, NotFound: f.NotFound, Error: msg})
		}
		return jsonResult(out), nil, nil
	})

	mcp.AddTool(srv, &mcp.Tool{
		Name: "kc_apply",
		Description: "Create/update/delete Keycloak resources from a manifest. resources is an array of " +
			"{type, realm, data} objects (the shape kc_fetch returns; data holds the Keycloak representation, " +
			"e.g. {\"type\":\"user\",\"realm\":\"demo\",\"data\":{\"username\":\"bob\"}}); mark a resource " +
			"\"delete\":true to remove it, and relationships is an array of {kind, path, data}. " +
			"Options: dryRun (validate and resolve without writing), delete (force-delete the listed " +
			"resources even when not marked), continueOnError, reconcile. " +
			"REQUIRES confirm:true — without it nothing is sent. Returns a per-resource report " +
			"(action created/updated/unchanged/deleted/skipped, status, error, createdId).",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: boolPtr(true)},
	}, func(ctx context.Context, req *mcp.CallToolRequest, in applyIn) (*mcp.CallToolResult, any, error) {
		if !in.Confirm {
			return nil, nil, fmt.Errorf(
				"refused: kc_apply changes Keycloak state; retry the same call with \"confirm\": true to send it")
		}
		report, err := deps.apply(ctx, in.Resources, in.Relationships, manifest.ApplyOptions{
			DryRun:          in.Options.DryRun,
			Delete:          in.Options.Delete,
			ContinueOnError: in.Options.ContinueOnError,
			Reconcile:       in.Options.Reconcile,
		})
		if err != nil {
			return nil, nil, toErrorResult(err)
		}
		if report.Results == nil {
			report.Results = []manifest.ApplyResult{}
		}
		return jsonResult(report), nil, nil
	})

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "kc_reload",
		Description: "Reload the OpenAPI spec from disk — both the raw client and the manifest tools (kc_fetch/kc_apply) pick up the new spec. Call it if the spec file changed since the server started.",
		Annotations: &mcp.ToolAnnotations{IdempotentHint: true},
	}, func(ctx context.Context, req *mcp.CallToolRequest, in reloadIn) (*mcp.CallToolResult, any, error) {
		if err := client.Reload(ctx); err != nil {
			return nil, nil, toErrorResult(err)
		}
		if err := deps.reload(); err != nil {
			// The raw client already serves the new spec; only the manifest
			// tools kept the old one. Say so — the next kc_reload retries just
			// the rebuild.
			return nil, nil, fmt.Errorf("spec reloaded, but %w (manifest tools still serve the old spec)", err)
		}
		return jsonResult(map[string]string{"status": "reloaded"}), nil, nil
	})

	return srv, nil
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
func Run(ctx context.Context, client *kcapi.Client, newManifest func() (manifest.Service, error)) error {
	srv, err := New(client, newManifest)
	if err != nil {
		return err
	}
	return srv.Run(ctx, &mcp.StdioTransport{})
}

// RunHTTP serves the same server over the streamable HTTP transport on an
// already-bound listener until ctx is cancelled or the listener fails. The
// server is built once and shared across HTTP sessions — the tools' shared
// state (the kcapi client, the manifest service) is built for the process,
// not per connection. Sessions carry the SDK's default localhost DNS-rebinding
// protection; authentication is out of scope — bind to loopback.
func RunHTTP(ctx context.Context, client *kcapi.Client, newManifest func() (manifest.Service, error), ln net.Listener) error {
	srv, err := New(client, newManifest)
	if err != nil {
		return err
	}
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return srv }, nil)
	// BaseContext ties request lifetimes — the client's SSE streams included —
	// to ctx, so cancellation ends streams at once and Shutdown never waits
	// behind an open session.
	httpSrv := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}
	serveErr := make(chan error, 1)
	go func() { serveErr <- httpSrv.Serve(ln) }()

	select {
	case err := <-serveErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpSrv.Shutdown(shutdownCtx)
	}
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
