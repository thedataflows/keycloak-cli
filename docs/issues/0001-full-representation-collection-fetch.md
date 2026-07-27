---
type: Issue
title: Full-representation collection fetch
description: Add a FetchQuery option that sends briefRepresentation=false so collection fetches return the attributes map for organizations, groups, and users.
tags: [issue, enhancement, admin, fetch]
timestamp: 2026-07-27T00:00:00Z
---

# ISSUE 0001: Full-representation collection fetch

- **Type**: feature
- **Status**: in-progress
- **Priority**: high
- **Labels**: [admin, fetch, enhancement]
- **Assignee**: none
- **Related**: none
- **Related code**: [`pkg/admin/fetch.go`](../../pkg/admin/fetch.go), [`pkg/admin/internal/client.go`](../../pkg/admin/internal/client.go), [`cmd/fetch.go`](../../cmd/fetch.go)
- **Closing commits**: none

## Summary

`admin.Service.Fetch` cannot ask Keycloak for the full representation of a
resource collection, so `attributes` never come back for organizations,
groups, or users. Add a `FetchQuery` option that emits
`briefRepresentation=false` on collection fetches.

## Details

Keycloak's list/search endpoints return a *brief* representation by default
that omits the `attributes` map (and, for groups, subgroup/role detail). The
only way to get attributes from a list is the `briefRepresentation=false`
query parameter, which this library never sends and offers no way to request.

Downstream consumer iga-dash / syncengine had to bypass the library with a raw
`net/http` GET to hydrate organization attributes
(`syncengine/connectors/keycloak/admin_adapter.go`,
`hydrateOrgAttributes`), and the planned group-sync feature
(`syncengine/docs/plans/2026-07-27-keycloak-group-sync.md`, Task 3a) would
have to do the same. This issue removes the need for that workaround.

### Evidence (verified live against Keycloak 26.6.3)

| Endpoint | `description` | `domains` | `attributes` |
|---|---|---|---|
| `GET /admin/realms/{realm}/organizations` (default) | present | present | **omitted** |
| `GET /admin/realms/{realm}/organizations?briefRepresentation=false` | present | present | **present** |
| `GET /admin/realms/{realm}/groups` (default) | present | n/a | **omitted** |
| `GET /admin/realms/{realm}/groups?briefRepresentation=false` | present | n/a | **present** |

### Root cause

1. `FetchQuery` ([`pkg/admin/fetch.go:16`](../../pkg/admin/fetch.go)) has no
   field to request a full representation. Its fields are `Realm`,
   `Resources`, `Search`, `Max`, `Parent`, `IncludeRelationships`, `Depth`,
   `Filter`.
2. `buildQueryParams` ([`pkg/admin/fetch.go:237`](../../pkg/admin/fetch.go))
   only ever emits `search` and `max`, and short-circuits to `nil` when both
   are unset.
3. Every collection fetch therefore uses Keycloak's brief default.

The plumbing already exists end to end; only the value is missing:

```
Service.Fetch
  -> buildQueryParams(query)                                   // <-- add briefRepresentation here
  -> fetchResourceCollection(ctx, resource, scope, realm, params...)
  -> specClient.FetchResources(ctx, resource, scope, params...)
  -> mergeQueryParams + url.Values encode                      // pkg/admin/internal/client.go:105
```

A new param produced by `buildQueryParams` reaches the HTTP request with no
other layer changes.

### Corrections to the original upstream request

The request was written against a vendored copy and asserts two things that
are **not** true of this repository:

- There is no `ExactMatch` field on `FetchQuery` and no `exact` query
  parameter anywhere in the codebase. Do not "preserve" it — either leave it
  out or add it as separate work.
- `buildQueryParams` returns early (`nil`) when `Search` and `Max` are both
  unset, so the new flag must be part of that guard or the guard must be
  dropped.

### Validation layers are already safe

- Request validation (`catalog.validateOperationInput`,
  [`pkg/catalog/contracts.go:428`](../../pkg/catalog/contracts.go)) iterates
  only the *spec-declared* parameters, so an undeclared extra query param is
  never rejected. Keycloak likewise ignores unrecognized query parameters.
- The 26.6.2 spec declares `briefRepresentation` on the GET collection
  operations for `organizations`, `groups`, and `users`, so it is a legal
  param there.
- `OrganizationRepresentation`, `GroupRepresentation`, and
  `UserRepresentation` all declare `attributes`, so
  `ValidateOperationResponse` accepts the richer payload.

### Proposed change (backward compatible)

Add a `FullRepresentation bool` to `FetchQuery`, documented as "requests
Keycloak's complete resource representation on collection fetches by sending
`briefRepresentation=false`; default false preserves brief behavior", and emit
it from `buildQueryParams`:

```go
if query.FullRepresentation {
    params["briefRepresentation"] = "false"
}
```

A tri-state `BriefRepresentation *bool` was considered and rejected: the only
need is "give me everything", so a single bool is simpler and unambiguous.

Spec-gating the parameter (only emitting it for operations that declare it) is
optional. The operation contract is available in `RuntimeClient.FetchResources`
but not in `buildQueryParams`; since Keycloak ignores unknown params and
request validation does not reject them, emit unconditionally and mark it with
a `ponytail:` comment naming the gate as the upgrade path.

Wire a `--full-representation` flag through `FetchCmd`
([`cmd/fetch.go:44`](../../cmd/fetch.go)) so the capability is reachable from
the CLI, not only the library API.

## Acceptance Criteria

- [x] `FetchQuery.FullRepresentation` exists and is documented with a doc comment.
- [x] `buildQueryParams` emits `briefRepresentation=false` when the field is set, including when `Search` and `Max` are both unset.
- [x] Table-driven unit test over `buildQueryParams` covers the flag alone and combined with `Search` and `Max`, and asserts the key is absent when the flag is unset. (`pkg/admin/fetch_internal_test.go`)
- [x] A transport test asserts the outgoing request URL contains `briefRepresentation=false`, confirming the param survives `FetchResources` → URL. (`pkg/admin/fetch_test.go`)
- [x] `Fetch(FetchQuery{Resources: "organization", FullRepresentation: true})` returns organizations whose `Data["attributes"]` is populated when the org has attributes.
- [x] `Fetch(FetchQuery{Resources: "group", FullRepresentation: true})` returns groups whose `Data["attributes"]` is populated.
- [x] With the flag unset, request behavior is unchanged — no new query param, existing fetch tests green.
- [x] `cmd fetch` exposes `--full-representation`.
- [x] `go vet ./...` clean.
- [x] Structural child collections receive the flag: `fetch organization --depth 1 --full-representation` returns org-scoped groups with `attributes` populated.
- [x] `search` and `max` do not leak into nested collection requests (asserted at transport level).

## Scope extension: nested (depth-expanded) collections

The first implementation only reached top-level collections, because
`fetchNestedResourceCollection` did not accept or forward query params. That left
org-scoped groups — reachable *only* via
`/organizations/{org-id}/groups` — without attributes, blocking the primary
downstream consumer. Scope was therefore extended to cover structural children:

- `buildNestedQueryParams` emits the representation flag for child collections.
- `fetchDepthLevels` and `fetchNestedResourceCollection` thread it through to
  `FetchPathCollection`, which already accepted variadic params.
- **`search`, `max` and `exact` are deliberately NOT forwarded.** They scope the
  *requested* resources; forwarding them would silently filter children too
  (a `--search acme` would drop every child whose name does not match).

## Out of Scope

- Relationship inclusion (`IncludeRelationships`) — already handled by an existing `FetchQuery` field.
- The org-scoped subgroup hierarchy and org-group membership — see [ISSUE 0002](0002-org-scoped-groups.md).
- Single-resource (non-collection) GET paths; brief representation only affects list/search endpoints.
- Adding an `exact` / `ExactMatch` query option.
- Pagination (`first`) and the `q` attribute-query parameter.

## Notes

### Live verification (Keycloak 26.6.2, throwaway realm, since deleted)

Verified end to end through the built CLI against a running instance, using a
temporary realm with an attributed group, subgroup, organization and user:

| Fetch | attributes with `--full-representation` | attributes without |
|---|---|---|
| `group` | `{costCenter, owner}` | absent |
| `group --depth 2` (subgroup `backend`) | `{squad}` | absent |
| `organization` | `{tier, externalId}` | absent |
| `user` (default realm profile) | absent | absent |
| `user` (`unmanagedAttributePolicy: ENABLED`) | `{department}` | `{department}` |

- **Correction to the original request: `user` does *not* benefit identically.**
  On 26.6.2 the users list returns `attributes` regardless of
  `briefRepresentation`; what actually gates them is the realm's declarative
  user profile `unmanagedAttributePolicy`, which defaults to disabled and drops
  unmanaged attributes on *every* endpoint, including `GET /users/{id}`. The flag
  is harmless for users but is not the fix there — a downstream consumer needing
  user attributes must enable that policy on the realm.
- **Depth needs no extra work for groups.** With `briefRepresentation=false`
  Keycloak inlines `subGroups` complete with their attributes, so a single
  `/groups` request covers the whole hierarchy; `--depth 2` produced the
  subgroup's attributes without any second parameterized request.
- Immediate downstream need is `organization` and `group`.
- The `apply` conflict-resolution fetches ([`pkg/admin/apply.go`](../../pkg/admin/apply.go)) call `FetchResources` directly with their own params and are unaffected; consider whether they should also request full representations once this lands.
- Requires a version bump and tag for the downstream consumer to drop its raw-HTTP workaround.
