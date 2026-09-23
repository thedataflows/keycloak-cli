# Spec: MCP server (`keycloak-cli mcp`)

Status: approved design 2025, implemented against kcapi as built on `main`.

## Goal

Expose the kcapi library to LLM agents as an MCP (Model Context Protocol)
server speaking JSON-RPC over stdio or streamable HTTP. The server is a thin
adapter: every tool is one kcapi call; no business logic lives here.

## Entry point and configuration

- Subcommand `keycloak-cli mcp` (kong, package `cmd`), server in
  `pkg/mcpserver`.
- Transport: **stdio** (default) or **streamable HTTP** via
  `--transport=http`; `--http-addr` sets the listen address (default
  `127.0.0.1:8081`). The HTTP surface has no authentication — the loopback
  default plus the SDK's DNS-rebinding protection are the safety net; do not
  expose it without adding auth. stdio's stdout carries only protocol frames;
  all logging goes to stderr.
- Config reuses the existing globals: `--base-url`/env, `--spec-path`
  (required, no default), `--timeout`. Credentials resolve from the
  environment at request time (kcapi's as-built flow); the server never takes
  a secret on its command line.
- One `kcapi.Client` is constructed at startup and shared by all tools; the
  manifest tools run on a `manifest.Service` built from the same flags and
  rebuilt by `kc_reload`. The server is built once and shared.
- The HTTP transport runs the MCP **2026-07-28 stateless core**: no
  `initialize` handshake and no sessions on the wire — every POST declares its
  protocol version (header + per-request `_meta`) and any request can be
  served by any instance, so a plain round-robin load balancer works. GET and
  DELETE return `405` with `Allow: POST`. Old-protocol clients that still
  handshake are served through the SDK's compatibility path. The tool catalog
  is static, so list results advertise the cacheable ttl hint (5 min,
  scope `public`).
- On startup the command prints an operator quick guide to stderr: the tool
  list, the safety rule, and wiring snippets for Claude Code, generic
  `mcp.json` harnesses, and the HTTP transport (`mcpserver.Guide`).

## Tools

| Tool | Arguments | Maps to | Annotation hints |
|---|---|---|---|
| `kc_operations` | `resource?`, `method?`, `tag?`, `search?` | `Operations(OpFilter)` | read-only |
| `kc_invoke` | `resource` XOR `op`, `verb`, `realm?`, `params?` (k=v map), `body?` (JSON), `confirm?` (bool) | `Invoke(Call)` | read-only=false, destructive |
| `kc_list` | `resource` XOR `op`, `verb` (GET only), `realm?`, `params?` | `ListAll(Call)` | read-only |
| `kc_resolve` | `type`, `name`, `realm?` (or `id`) | `Resolve(Ref)` | read-only |
| `kc_neighbors` | `type`, `name`, `realm?` (or `id`), `child?`, `parent?` | `Resolve(Ref)` + `Neighbors(node, EdgeFilter)` | read-only |
| `kc_edges` | — | `Edges()` | read-only |
| `kc_fetch` | `realm?`, `resources?`, `search?`, `max?`, `parent?`, `includeRelationships?`, `depth?`, `filter?`, `exactMatch?`, `fullRepresentation?` | `manifest.Service.Fetch(FetchQuery)` | read-only |
| `kc_apply` | `resources?` (kcapi.Resource shapes), `relationships?`, `options?` (`dryRun`, `delete`, `continueOnError`, `reconcile`), `confirm?` (bool, required) | `manifest.Service.Apply(...)` | read-only=false, destructive |
| `kc_reload` | — | `Reload(ctx)` + manifest service rebuild | safe, not read-only of Keycloak |

`kc_neighbors` walks the edge's **collection prefix**: for a parent it lists
child collections (e.g. an organization yields its `groups` collection), it
does not instantiate single-child item paths. Tool descriptions say so.

## Safety model (as chosen)

Two rules:

- **`kc_list` is read-only by construction**: it refuses any operation whose
  verb is not a known GET, before anything is sent.
- **State-changing work requires an explicit `confirm: true` tool argument.**
  For `kc_invoke`, GET/HEAD pass freely; POST/PUT/PATCH/DELETE without
  `confirm: true` are rejected before any request is sent, with an error
  telling the caller to retry with `confirm`. `kc_apply` requires
  `confirm: true` on every call, dry-run included (locating existing resources
  already sends requests). No server-level read-only mode, no flags.

Concurrent manifest tool calls are serialized inside the server
(`manifestDeps`): `manifest.Service` caches identity maps lazily without
locking. The kcapi tools run fully parallel.

## Error mapping

`*kcapi.Error` becomes an MCP tool error result (`isError: true`) carrying the
kind, operation label, HTTP status, and body — never a protocol-level error —
so the calling agent can read "unknown resource type" or a 404 and self-correct
on its next turn. Missing/invalid tool arguments are rejected by the SDK's
schema validation with the same tool-error shape.

## Dependencies

Official SDK `github.com/modelcontextprotocol/go-sdk` v1.8.0, vendored.

## Acceptance scenarios (executable: `pkg/mcpserver`)

1. Server exposes the library surface: kc_invoke, kc_neighbors, kc_operations,
   kc_reload, kc_resolve, kc_list, kc_edges, kc_fetch, kc_apply.
2. `kc_operations` lists spec operations; `resource`/`method` filters narrow.
3. `kc_invoke` GET reaches the fake Keycloak and returns its JSON body.
4. `kc_invoke` POST without `confirm` is rejected client-side; the fake
   Keycloak records no request.
5. `kc_invoke` POST with `confirm: true` reaches the fake Keycloak.
6. `kc_resolve` resolves a realm by name through the real kcapi flow.
7. `kc_neighbors` walks a parent's child collection and returns edges.
8. `kc_reload` completes without error and rebuilds the manifest service.
9. kcapi failures surface as tool error results carrying kind and status.
10. `kc_resolve` resolves a realm by its own name onto
    `GET /admin/realms/{name}`; `kc_neighbors` from the realm node lists the
    realm-rooted child collections.
11. Stdio smoke: the built binary completes an MCP initialize handshake over
    stdin/stdout, lists the nine tools, and prints the quick guide to stderr
    while stdout stays protocol-clean.
12. Streamable HTTP speaks the 2026-07-28 stateless core: a raw POST of
    `tools/list` with the protocol-version header + `_meta` is answered
    without any initialize handshake and without a session id; GET and DELETE
    are refused with `405` + `Allow: POST`; an old-protocol client that still
    initializes is served unchanged. `tools/list` advertises
    `ttlMs: 300000` / `cacheScope: "public"`.
13. HTTP smoke: the built binary with `mcp --transport=http` opens the
    listener and serves the nine tools to a streamable-HTTP MCP client;
    ctx cancellation shuts the server down cleanly despite open client
    sessions.
14. `kc_list` walks a GET collection into one JSON array and refuses non-GET.
15. `kc_edges` lists the parent→child vocabulary (realms>users present).
16. `kc_fetch` returns `{resources, relationships, failures}` for a realm.
17. `kc_apply` without `confirm` sends nothing; with `confirm: true` it
    applies and returns the per-resource report.
18. `mcpserver.Guide` names every tool, the safety rule, the Claude Code and
    `mcp.json` wiring, and the HTTP endpoint with its loopback caveat.

## Out of scope

Resources/prompts/sampling, HTTP transport authentication, per-tool authn,
server-side read-only mode, operationId mode beyond passing `op` through (the
vendored spec has none).
