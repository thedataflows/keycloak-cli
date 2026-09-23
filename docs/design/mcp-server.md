# Spec: MCP server (`keycloak-cli mcp`)

Status: approved design 2025, implemented against kcapi as built on `main`.

## Goal

Expose the kcapi library to LLM agents as an MCP (Model Context Protocol)
server speaking JSON-RPC over stdio. The server is a thin adapter: every tool
is one kcapi call; no business logic lives here.

## Entry point and configuration

- Subcommand `keycloak-cli mcp` (kong, package `cmd`), server in
  `pkg/mcpserver`.
- Transport: **stdio**. stdout carries only protocol frames; all logging goes
  to stderr.
- Config reuses the existing globals: `--base-url`/env, `--spec` (default
  `keycloak-oapi/26.6.2.spec.json`), `--timeout`. Credentials resolve from the
  environment at request time (kcapi's as-built flow); the server never takes
  a secret on its command line.
- One `kcapi.Client` is constructed at startup and shared by all tools.

## Tools

| Tool | Arguments | Maps to | Annotation hints |
|---|---|---|---|
| `kc_operations` | `resource?`, `method?`, `tag?`, `search?` | `Operations(OpFilter)` | read-only |
| `kc_invoke` | `resource` XOR `op`, `verb`, `realm?`, `params?` (k=v map), `body?` (JSON), `confirm?` (bool) | `Invoke(Call)` | read-only=false, destructive |
| `kc_resolve` | `type`, `name`, `realm?` (or `id`) | `Resolve(Ref)` | read-only |
| `kc_neighbors` | `type`, `name`, `realm?` (or `id`), `child?`, `parent?` | `Resolve(Ref)` + `Neighbors(node, EdgeFilter)` | read-only |
| `kc_reload` | — | `Reload(ctx)` | safe, not read-only of Keycloak |

`kc_neighbors` walks the edge's **collection prefix**: for a parent it lists
child collections (e.g. an organization yields its `groups` collection), it
does not instantiate single-child item paths. Tool descriptions say so.

## Safety model (as chosen)

One rule: **state-changing verbs require an explicit `confirm: true` tool
argument.** GET/HEAD pass freely; POST/PUT/PATCH/DELETE without
`confirm: true` are rejected before any request is sent, with an error telling
the caller to retry with `confirm`. No server-level read-only mode, no flags.

## Error mapping

`*kcapi.Error` becomes an MCP tool error result (`isError: true`) carrying the
kind, operation label, HTTP status, and body — never a protocol-level error —
so the calling agent can read "unknown resource type" or a 404 and self-correct
on its next turn. Missing/invalid tool arguments are rejected by the SDK's
schema validation with the same tool-error shape.

## Dependencies

Official SDK `github.com/modelcontextprotocol/go-sdk` v1.8.0, vendored.

## Acceptance scenarios (executable: `pkg/mcpserver`)

1. Server exposes exactly five tools: kc_invoke, kc_neighbors, kc_operations,
   kc_reload, kc_resolve.
2. `kc_operations` lists spec operations; `resource`/`method` filters narrow.
3. `kc_invoke` GET reaches the fake Keycloak and returns its JSON body.
4. `kc_invoke` POST without `confirm` is rejected client-side; the fake
   Keycloak records no request.
5. `kc_invoke` POST with `confirm: true` reaches the fake Keycloak.
6. `kc_resolve` resolves a realm by name through the real kcapi flow.
7. `kc_neighbors` walks a parent's child collection and returns edges.
8. `kc_reload` completes without error.
9. kcapi failures surface as tool error results carrying kind and status.
10. Stdio smoke: the built binary completes an MCP initialize handshake over
    stdin/stdout and lists the five tools.

## Out of scope

Resources/prompts/sampling, HTTP transport, per-tool authn, server-side
read-only mode, operationId mode beyond passing `op` through (the vendored
spec has none).
