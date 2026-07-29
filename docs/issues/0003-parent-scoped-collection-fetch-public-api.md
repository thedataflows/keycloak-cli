---
type: Issue
title: Parent-scoped collection fetch on the public admin.Service
description: Reading one child collection of one parent (e.g. a client's roles) is only reachable through Depth traversal, which fans out to every child type of every seed plus a realm-wide reference sweep. RuntimeClient.FetchResourcesWithParent already does the targeted fetch but is unexported.
tags: [issue, enhancement, admin, fetch, roles, clients]
timestamp: 2026-07-28T00:00:00Z
---

# ISSUE 0003: Parent-scoped collection fetch on the public admin.Service

- **Type**: feature
- **Status**: done
- **Priority**: high
- **Labels**: [admin, fetch, enhancement]
- **Assignee**: none
- **Related**: [ISSUE 0001](0001-full-representation-collection-fetch.md) (representation flag must be honoured by the new path), [ISSUE 0002](0002-org-scoped-groups.md) (established the parent-scoping machinery this builds on), downstream consumer: iga-dash `syncengine/docs/issues/0022-client-roles-never-synced-to-target-realm.md`
- **Related code**: [`pkg/admin/admin.go`](../../pkg/admin/admin.go), [`pkg/admin/fetch.go`](../../pkg/admin/fetch.go), [`pkg/admin/internal/client.go`](../../pkg/admin/internal/client.go)
- **Closing commits**: `db1ed41` (FetchChildren + tests, verified live). Released in `v1.3.0`.

## Summary

A library consumer cannot ask for *one* child collection of *one* parent — for example the roles of a single client (`GET /admin/realms/{realm}/clients/{client-uuid}/roles`). The public `admin.Service` exposes only `Spec`, `Fetch` and `Apply`, and `Fetch` resolves a resource type to its realm-level collection. The only route is `FetchQuery{Depth: 1}`, which fans out to every child type of every seed resource *and* runs a realm-wide reference-resolution sweep. The targeted fetch already exists — `RuntimeClient.FetchResourcesWithParent` — but lives in `pkg/admin/internal` and is therefore unreachable from outside the module.

## Details

### Already working — verified, no work needed

Verified 2026-07-28 against a live Keycloak **26.6.3** instance with the `26.6.2` spec, using `fetch client demo-a-client-1 --realm sync-source --depth 1`:

| Claim under test | Result |
|---|---|
| `Depth: 1` on `client` reaches client-scoped roles | Yes — `demo-a-clientrole-1` is returned |
| Fetched child is disambiguated from a realm role | Yes — `parentType: "client"`, `clientRole: true` |
| Parent reference field is injected so the child can be re-applied standalone | Yes — `clientUuid: "0bd3b65d-…"` |
| Absent optional children are benign, not errors | Yes — `Optional resources not found: resource (1), scope (1)` (404 → `FetchFailure.NotFound`) |

Raw fetched resource:

```json
{
  "type": "role",
  "realm": "sync-source",
  "parentType": "client",
  "data": {
    "name": "demo-a-clientrole-1",
    "description": "Seeded demo client role",
    "clientRole": true,
    "clientUuid": "0bd3b65d-40c6-46ed-8ffe-81610c031ce1",
    "composite": false,
    "id": "e6f79b96-9640-4e88-bd5d-4e6488452eac"
  }
}
```

So the parent-scoping machinery ISSUE 0002 pinned for org-scoped groups covers client-scoped roles too. **This issue is not "the library cannot read client roles".** It is: the only way to read them is disproportionately expensive and cannot be scoped to the child type the caller wants.

### Gap 1 (P1) — no targeted parent-scoped fetch in the public API

`internal.RuntimeClient.FetchResourcesWithParent(ctx, resource, params...)` (`pkg/admin/internal/client.go:160`) does exactly the right thing: it resolves the operation with `resource.ParentType`, renders the nested path from the parent's path params, and issues one GET. It is already in production use by apply (`pkg/admin/apply.go:692` and `:765`). It is simply not exposed:

- `admin.Service` (`pkg/admin/admin.go:25-29`) declares only `Spec`, `Fetch`, `Apply`.
- `RuntimeClient.FetchResources` — the one reachable through `Fetch` — hardcodes the parent as empty: `ResolveResourceOperation(resourceType, "", http.MethodGet, catalog.OperationCollection)` (`pkg/admin/internal/client.go:94`). Passing `client-uuid` in the `scope` map cannot change the resolved path, so the nested endpoint is unreachable that way.

### Gap 2 (P1) — Depth traversal is the only alternative, and it is expensive

Observed request pattern for `--depth 1` on clients (`--log-level trace`, URLs deduplicated):

```
/clients                      (seed list)
/authentication/flows         reference-resolution sweep
/client-scopes                reference-resolution sweep
/groups                       reference-resolution sweep
/roles                        reference-resolution sweep
/identity-provider/instances  reference-resolution sweep
/components                   reference-resolution sweep
/clients                      reference-resolution sweep (second time)
/users                        reference-resolution sweep
/organizations                reference-resolution sweep
```

Plus, per seed client, one GET for **each** downward child type — observed live as four (`roles`, `protocol-mappers/models`, `authz/…/resource`, `authz/…/scope`, the last two 404 on a client without authorization services).

For the 9-client realm used above, obtaining one client's roles costs roughly `1 + 9×4 + 9 ≈ 46` GETs, including a full `/users` listing. The downstream consumer (iga-dash syncengine) polls both a source and a target realm every 60 s, so this is not a one-off import cost.

Two further consequences of routing a targeted read through `Depth`:

- **Blast radius.** The caller cannot restrict the traversal to `role`, so an unrelated non-404 failure on an authz or protocol-mapper endpoint lands in `FetchReport.Failures` for a caller that only asked for roles. Consumers that (correctly) treat non-NotFound failures as fatal — see the syncengine `checkFetchFailures` policy — then fail a role pull because of an unrelated endpoint.
- **Cost is invisible.** `FetchPathCollection` (`pkg/admin/internal/client.go:230`) issues the nested GETs without a single log line, unlike `FetchResources`/`FetchResourcesWithParent`, which log `Fetching resources` at debug. The nested requests cannot be observed at `--log-level trace`; the numbers above had to be reconstructed from the code. Worth fixing while in here.

### Proposed shape

A thin public wrapper over the existing internal method, on `admin.Service`:

```go
// FetchChildren returns one child collection of one parent resource, e.g. the
// roles of a client. parent must carry enough identity to render the nested
// path (typically Realm + Data["id"]).
FetchChildren(ctx context.Context, parent manifest.Resource, childType string, query ChildFetchQuery) (FetchReport, error)
```

Requirements on the returned resources, to match what `Depth` already produces:

- `Type = childType`, `Realm = parent.Realm`, `ParentType = parent.Type`.
- Parent reference fields injected via `Resolver().ParentReferenceFields`, so each child can be applied standalone (the `clientUuid` above).
- `FullRepresentation` honoured (ISSUE 0001 parity — `briefRepresentation=false`), since client roles do carry `attributes`.
- 404 classified as `FetchFailure{NotFound: true}` rather than a hard error, consistent with `fetchDepthLevels`.
- Exactly one HTTP GET, no reference-resolution sweep, no sibling child types.

An alternative shape is extending `FetchQuery` with a parent selector (its `Parent` field is currently only consumed by the `authenticationexecution` branch, `pkg/admin/fetch.go:130-141`). A dedicated method is preferred: `Fetch` is realm-sweep-shaped (it iterates realms and resource lists and owns the depth/relationship logic), whereas this is a single addressed collection. The implementer may still choose the `FetchQuery` route if it keeps the public surface smaller — record the reasoning either way.

### Downstream requirement (why this is P1)

iga-dash ISSUE 0022: client-scoped roles never reach the target realm. syncengine needs, per pull:

```go
// source: enumerate one client's roles
report, err := svc.FetchChildren(ctx,
    manifest.Resource{Type: "client", Realm: realm, Data: map[string]any{"id": clientUUID}},
    "role", admin.ChildFetchQuery{FullRepresentation: true})
```

The write direction is expected to need **no** library change — `manifest.Resource{Type: "role", ParentType: "client", Data: {"name": …, "clientUuid": …}}` should resolve to `POST /clients/{client-uuid}/roles`, with `clientUuid` stripped from the body by `ParentReferenceFieldNames`, exactly as ISSUE 0002 pinned for org-scoped groups. That expectation is derived from the resolver code, **not** yet verified live for roles — pin it with a test as part of this issue so the consumer can rely on it.

## Acceptance Criteria

- [x] `admin.Service` exposes a parent-scoped collection fetch (`FetchChildren`) that issues exactly one GET for one `(parent, childType)` pair — no sibling child types, no reference-resolution sweep (`TestFetchChildrenReturnsOneClientRoleCollection`)
- [x] Returned children carry `Type`, `Realm`, `ParentType` and the injected parent reference field(s), matching what `Depth: 1` produces (asserts the `clientUuid` shape); verified live on `sync-source`
- [x] `FullRepresentation` is honoured on the new path (ISSUE 0001 parity), verified against a child collection that has `attributes`
- [x] A 404 on the child collection is reported as `FetchFailure{NotFound: true}`, not a hard error (`TestFetchChildrenClassifies404AsNotFound`)
- [x] Resolver behaviour for roles is pinned by tests (`client_role_resolver_test.go`): `("role", "client", POST, collection)` → `/clients/{client-uuid}/roles`, `("role", "", …)` → `/roles`, single-resource `PUT`/`DELETE` with `ParentType: "client"` → `/clients/{client-uuid}/roles/{role-name}`. **Recorded divergence:** a bare `("role","client",GET,collection)` resolves to `/clients/{client-uuid}/roles/{role-name}/composites` (operationPriority), so `FetchChildren` resolves the path via `BuildDownwardGraph` (as `Depth` does), not the resolver
- [x] An apply test proves a `role` with `ParentType: "client"` creates via `POST /clients/{client-uuid}/roles` with `clientUuid` resolved from the parent and stripped from the body (spec-validated), and updates/deletes via the single-resource endpoint (`apply_client_role_test.go`)
- [x] Realm-role and existing `Depth` behaviour unchanged — no regression (a realm role still resolves to `/roles`)
- [x] `FetchPathCollection` logs its request at debug like the other fetch paths, so nested/child GET volume is observable (`Fetching path collection`)
- [x] `go vet ./...` clean, `go test ./...` green (except a pre-existing, unrelated `pkg/output` TOML failure)
- [x] Tagged release (`v1.3.0`) so iga-dash can bump `syncengine/go.mod`. Note for the consumer: iga-dash vendors three trees (`igadash/vendor`, `syncengine/vendor`, and the root `go work vendor`) — all three must be regenerated together — released in `v1.3.0`

## Out of Scope

- Client roles' **composite** membership (`/clients/{client-uuid}/roles/{role-name}/composites`) — a relationship kind, already registered, unchanged here.
- Children addressed by **two** parent path params (org-scoped subgroups, org-group membership) — still the deferred `ParentType`-generalization work from ISSUE 0002 Gap 2/3.
- Any change to realm-level `Fetch` semantics, the depth traversal, or the reference-resolution sweep itself. This issue adds a targeted alternative; it does not retune `Depth`.
- The consumer-side sync model (how syncengine names, correlates and provisions client roles cross-realm) — that is iga-dash ISSUE 0022.

## Notes

- Reproduction of everything above:

  ```bash
  KEYCLOAK_ACCESS_TOKEN=$(master admin token) \
    keycloak-cli fetch client demo-a-client-1 --realm sync-source --depth 1 \
      -u http://<keycloak> --log-level trace -f json -o depth1.json
  jq '.resources[] | select(.type=="role")' depth1.json
  ```

  Note the CLI's default `--exclude-fields=containerId` hides `containerId` in
  displayed output; it is present in the underlying data.

- `FetchResourcesWithParent` is already exercised in production by the apply
  conflict-resolution path, so exposing it is low-risk: the behaviour under test
  is not new code, only newly reachable.
- The 9-collection reference sweep listed in Gap 2 is emitted once per `Fetch`
  call regardless of `Filter`; filtering to a single client narrows the *seed*
  set but not the sweep. Verified with and without the positional filter.

### Implementation notes (resolved)

- Added `Service.FetchChildren(ctx, parent, childType, ChildFetchQuery)` — a thin
  wrapper reusing the existing `fetchNestedResourceCollection`, so it issues
  exactly one GET, injects the parent reference field(s), tags
  `Type`/`Realm`/`ParentType`, honours `FullRepresentation`
  (`briefRepresentation=false`), and reports a 404 as `FetchFailure{NotFound}`.
- **Resolved the child path via `BuildDownwardGraph`, not `ResolveResourceOperation`.**
  Recording actual resolver behaviour (as the issue asked) surfaced a divergence:
  a bare `("role","client",GET,collection)` resolves to
  `/clients/{client-uuid}/roles/{role-name}/composites` because `operationPriority`
  prefers composites among the GET-collection candidates. The structural graph
  correctly maps `client → role` to the plain `/clients/{client-uuid}/roles`, and
  it is the same source `Depth` uses — so `FetchChildren` returns exactly what
  `Depth: 1` produces. Pinned in `TestClientScopedRoleGetCollectionDivergence`.
- Added a debug log to `FetchPathCollection` (`Fetching path collection`) so
  nested/child GET volume is observable, as Gap 2 requested.
- Tests: `fetch_children_test.go` (one GET, shape, FullRepresentation, 404),
  `client_role_resolver_test.go` (resolver pins + divergence record),
  `apply_client_role_test.go` (create via `POST /clients/{client-uuid}/roles`
  with `clientUuid` stripped and spec-validated; update/delete via the
  single-resource endpoint).
- **Verified live** against Keycloak 26.6.3, realm `sync-source`: the read
  `FetchChildren` wraps returns `demo-a-clientrole-1` with `parentType: client`
  and the injected `clientUuid`, via a single observable
  `Fetching path collection` GET to `/clients/{client-uuid}/roles`. FetchChildren
  itself is a library API (no CLI surface) and is covered by the unit tests above.
