# Design: `pkg/kcapi` — spec-driven Keycloak client library, generic CLI, MCP server

Status: **shipped.** The design below was implemented on `feat/kcapi-library` (plan:
`docs/superpowers/plans/2026-09-22-kcapi-library.md`); the body remains as originally written for
historical context. Deviations from the design as actually built:

- The vendored `keycloak-oapi/26.6.2.spec.json` defines **no operationIds**, so resource+verb is
  the primary resolution mode everywhere (`Call{Resource, Verb}`, `invoke --resource --verb`,
  `graph` types). `Op` mode still works for specs that carry operationIds.
- `Neighbors` traversal instantiates each edge's **collection prefix** — the spec path ending
  directly before the child's placeholder (e.g. `/admin/realms/{realm}/users/{user-id}/groups`) —
  rather than the edge's own path, so the walk lists collections instead of having to guess child
  ids the parent cannot know.
- The legacy curated relationship registry was **retained in full** inside `kcapi`
  (`relationship_registry.go`: kinds, path matching, param types, delete-operation building), not
  reduced to display-name-only overrides; it still powers the manifest relationship flows pending
  a successor design.
- `Config.Credentials` **validates the intended grant shape only** (password pair vs client
  secret, mutual exclusivity); the actual token values are resolved from the environment
  (`KEYCLOAK_ACCESS_TOKEN`/`KEYCLOAK_REFRESH_TOKEN`) exactly as the CLI's `.env` flow does.
  Programmatic token sourcing goes through `Config.Auth` (`auth.Service`;
  kcapi only calls its `AccessToken` method).
- The MCP server was **not started**; it is deferred to its own future plan.

## Context

Today the repo is a CLI whose intelligence lives in `pkg/admin` (runtime client, apply/fetch flows)
and `pkg/catalog` (contracts, identity resolution, relationship patterns). A full Keycloak OpenAPI
spec ships at `keycloak-oapi/26.6.2.spec.json` and is updated over time. The ponytail audit found
dead code (`pkg/models/models.gen.go`, `cmd/gen-rel-schema` + `generated/relationships.json`
schema output) and duplicated logic inside `pkg/catalog` and `pkg/admin`.

Goal: a first-class Go library that can drive **any** Keycloak instance matching the given spec —
every operation, not just manifest-supported types — plus the relationships between Keycloak
objects, a CLI exposing both, and (designed, then built) an MCP server so LLMs can work with
Keycloak instances through a lightweight tool surface.

## Goals

- `pkg/kcapi`: spec-loaded generic client. All operations discovered from the spec, invocable by
  `operationId` or resource+verb, validated against the live catalog.
- Relationships as data: edges inferred structurally from path templates; public graph queries
  (`Resolve`, `Neighbors`).
- Tolerates spec updates without code changes; supports atomic reload.
- CLI keeps current behavior (migrated onto `kcapi`), gains `invoke` and `graph` commands.
- MCP server design included; implementation is the final phase.
- Phase 0 applies the audit fixes as part of the migration.

## Non-goals

- No generated typed models (the audit deletes `pkg/models/models.gen.go`; it stays deleted).
- No manifest apply exposed over MCP in v1.
- No offline mock server product; tests use `httptest` fakes only.

## Architecture (end state)

```
pkg/kcapi        library: spec-driven client, generic invoke, catalog, graph
pkg/manifest     declarative fetch/upload/compare, rebuilt on kcapi (RealmCache lives here)
pkg/auth         unchanged (token sources); kcapi consumes it
internal/cli     commands migrate to kcapi + manifest; new `invoke` and `graph` commands
cmd/keycloak-cli unchanged entry point
```

`pkg/admin` dissolves: runtime client, identity resolution, and relationship fetching move into
`kcapi`; manifest flows move to `pkg/manifest`. `pkg/realmgen` and `pkg/output` are unaffected.

Deleted in phase 0: `pkg/models/`, `cmd/gen-rel-schema/`, `generated/relationships.json`
(no external consumer was flagged; the schema it produced is replaced by catalog validation).

## kcapi public API

```go
client, err := kcapi.New(kcapi.Config{
    BaseURL: "https://kc.example.com",
    Spec:    kcapi.SpecSource{Path: "keycloak-oapi/26.6.2.spec.json"}, // Path, URL, or Raw bytes
    Auth:    kcapi.Credentials{ClientID: "...", ClientSecret: "..."}, // client-credentials or password grant
})

// Discovery — reflects the loaded spec.
ops, err := client.Operations(kcapi.OpFilter{Resource: "users", Method: kcapi.Get})

// Generic invoke.
out, err := client.Invoke(ctx, kcapi.Call{
    Op:     "getUser",              // operationId; or kcapi.ByVerb("users", kcapi.Get)
    Realm:  "demo",
    Params: kcapi.P{"id": "…"},     // path + query params, validated against the spec
    Body:   bodyJSON,               // optional; any JSON-serializable value
})

// Relationships & graph.
node,  err := client.Resolve(ctx, kcapi.Ref{Type: kcapi.User, Name: "alice", Realm: "demo"})
edges, err := client.Neighbors(ctx, node, kcapi.EdgeFilter{Child: kcapi.Role})
```

### Catalog

Moved from `pkg/catalog`, generalized. Every path+method becomes an `Operation{id, resource, verb,
pathTemplate, params, requestSchema, responseSchema, tags}`. `Operations()` enumerates all of them;
nothing is restricted to manifest-supported types. Resource names derive from path segments
(`{realm}/users/{id}` → resource `users`), with a small override table for odd cases.

### Relationships

Inferred structurally: any path template containing a parent-collection placeholder plus child
identifiers yields an edge `{ParentType, ChildType, PathPattern, Cardinality, CascadeOwner}`.
The existing curated registry (`catalog/relationship_registry.go`) drops to optional display-name
overrides only. Edges power `Neighbors`, which instantiates child patterns from a resolved parent
resource (reusing the proven pattern-instantiation and search-or-list identity resolution code,
now public).

### Spec updates

`client.Reload(ctx)` re-fetches the spec and swaps the catalog atomically (single pointer swap;
in-flight calls finish against the old catalog). Calls are validated against the live catalog, so
an updated spec takes effect without process restart. Additive spec changes require no code.

### Errors

`*kcapi.Error{Kind, Op, Status, Body}` with sentinels `ErrNotFound`, `ErrConflict`,
`ErrValidation`, `ErrAuth`. Kinds map from the existing contract error classification in
`pkg/catalog`/`pkg/admin/errors.go`. `errors.Is` works against both sentinels and wrapped errors.

### Responses

`Invoke` returns `json.RawMessage`; callers decode into their own types. A pagination helper
wraps Keycloak's `first`/`max` list parameters.

## CLI

- Existing commands (`fetch`, `upload`/`apply`, `compare`, `generate`, `admin-token`) keep exact
  behavior; their implementations migrate to `kcapi` + `manifest`.
- New: `keycloak-cli invoke <operation-id> --realm R [--param k=v ...] [--body @file]`.
- New: `keycloak-cli graph resolve|neighbors <type> <name> --realm R [--edge-type T]`.
- Both print JSON via the existing `pkg/output`.

## MCP server design

- Go, official `modelcontextprotocol/go-sdk`, stdio transport first.
- Tools: `kc_operations` (list/filter — LLM discovery), `kc_invoke` (generic call),
  `kc_resolve`, `kc_neighbors`. Manifest apply deliberately not exposed in v1.
- Safety: server-level read-only mode (HTTP method allowlist), tool annotations
  (`readOnlyHint`, `destructiveHint`), and destructive `kc_invoke` calls require an explicit
  `confirm: true` argument or they are rejected.
- Configuration via env/flags: base URL, client credentials, spec path, read-only flag.
- The server is a thin adapter: every tool maps onto one `kcapi` call; no business logic.

## Testing strategy

- Table-driven unit tests with `httptest` fake Keycloak per `kcapi` concern: invoke param
  validation, body handling, pattern instantiation, edge inference, reload atomicity, error
  mapping.
- Existing `pkg/admin` tests migrate with their code into `kcapi`/`manifest`.
- Integration tests against a live instance stay env-gated, skipped by default.
- TDD throughout per repo conventions (acceptance tests first for each plan task).

## Risks

- Spec drift: guarded by catalog-driven validation (unknown ops fail loudly) and reload.
- Relationship inference noise: curated overrides table bounds display naming; edge direction is
  structural, so wrong inferences are visible in `kc_operations`/tests and fixable in one table.
- Migration regression risk: behavior-preservation is enforced by migrating existing tests
  unchanged where possible before refactoring their subjects.
