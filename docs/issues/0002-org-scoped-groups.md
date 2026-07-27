---
type: Issue
title: Org-scoped groups
description: Organization-nested groups as a parent-scoped resource — the flat set is fetched, applied and regression-tested; the subgroup hierarchy and org-group membership are deferred.
tags: [issue, enhancement, admin, catalog, groups, organizations]
timestamp: 2026-07-27T00:00:00Z
---

# ISSUE 0002: Org-scoped groups

- **Type**: feature
- **Status**: in-progress
- **Priority**: medium
- **Labels**: [admin, catalog, groups, organizations, enhancement]
- **Assignee**: none
- **Related**: [ISSUE 0001](0001-full-representation-collection-fetch.md) (its nested-collection extension closed Gap 1 below)
- **Related code**: [`pkg/admin/fetch.go`](../../pkg/admin/fetch.go), [`pkg/catalog/dependencies.go`](../../pkg/catalog/dependencies.go), [`pkg/catalog/relationship_registry.go`](../../pkg/catalog/relationship_registry.go)
- **Closing commits**: none

## Summary

In Keycloak 26.x an organization owns its own group hierarchy under
`/admin/realms/{realm}/organizations/{org-id}/groups/...`, a distinct location
from realm groups that shares the `GroupRepresentation` schema. Fetch and apply
of the **flat** org-scoped group set now work, including `attributes`, and are
covered by regression tests. The **subgroup** hierarchy and org-group
**membership** are deferred — both need a two-parent-path-param model that does
not exist yet.

## Details

This issue was raised as "the library cannot fetch or apply org-scoped groups at
all". **That premise is wrong** — the investigation below was run against the
26.6.2 spec and a live 26.6.3 instance, and most of the requested P1 already
works. What follows records what was verified so the remaining work is scoped
honestly.

### Already working — verified, no work needed

| Claim under test | Result |
|---|---|
| `ResolveResourceOperation("group", "organization", POST/GET, collection)` | → `/organizations/{org-id}/groups` |
| …`PUT`/`DELETE`, single | → `/organizations/{org-id}/groups/{group-id}` |
| No contamination: parent `""` / `"group"` | → `/groups` / `/groups/{group-id}/children` |
| `BuildDownwardGraph()["organization"]` | `group → /organizations/{org-id}/groups`, coexisting with `group → /groups/{group-id}/children` |
| Fetch traversal of an org's groups | `fetch organization --depth 1` returns them; scoping to one org works via the positional filter |
| Type-collision disambiguation | fetched org groups carry `ParentType: "organization"` and an injected `orgId` |
| Apply CRUD | create `201`, re-apply `204` (idempotent update), attribute change propagated — resolved to the org-scoped endpoint |

So the parent-scoping machinery (`identity.go` placeholder map, resolver
`filterByParentType`, `BuildDownwardGraph`, `fetchNestedResourceCollection`
`ParentType` tagging, `resolveParentReferences`) already covers the flat case
end to end. Only the items below are outstanding.

### ~~Gap 1 (P1) — attributes are missing on org-scoped group fetch~~ — RESOLVED

Closed by the nested-collection extension in
[ISSUE 0001](0001-full-representation-collection-fetch.md), as predicted:
`fetchNestedResourceCollection` now forwards the representation flag to
`FetchPathCollection`. Verified live —
`fetch organization acme --depth 1 --full-representation` returns:

```text
organization acme              attributes= {'tier': ['gold']}
group        acme-engineering  parentType= organization  attributes= {'squad': ['core']}
```

With P1 fetch and apply both working, **the flat org→group sync is unblocked**.
The subgroup hierarchy (Gap 2) and membership (Gap 3) are deferred — see the
decision below.

### Gap 2 (P2) — org-scoped subgroups are unreachable

Confirmed real, and harder than the issue assumed:

- The only downward edge for `group → group` is the **realm** path
  `/groups/{group-id}/children`. Keycloak **rejects** it for an org-scoped
  group: `400 {"errorMessage":"Cannot manage organization related group via non
  Organization API."}`
- `fetch organization --depth 3` does **not** return an org-scoped subgroup that
  demonstrably exists.
- The correct path `/organizations/{org-id}/groups/{group-id}/children` needs
  **two** parent path params (`org-id` *and* `group-id`), which exceeds the
  current single-parent `ParentType` assumption. This is the core work item.
- `briefRepresentation=false` does **not** help here: unlike realm groups, the
  org-scoped collection does not inline `subGroups` (`populateHierarchy=true`
  also does not; `subGroupsCount=true` does report `subGroupCount: 1`, so the
  children exist and are detectable).
- The org children **listing** does not declare `briefRepresentation` at all and
  never returns `attributes`; only the single `GET
  /organizations/{org-id}/groups/{group-id}` does. So fetching org-scoped
  subgroups *with* attributes requires a per-child single GET, not a listing.

### Gap 3 (P3) — org-group membership is not a relationship

`pkg/catalog/relationship_registry.go` has **zero** entries for
`/organizations/{org-id}/groups` (only `organization-member` and
`organization-identity-provider`, lines 349–350). The membership edge
`PUT/DELETE /organizations/{org-id}/groups/{group-id}/members/{userId}` again
needs two fixed path params, so it requires the same registry/param-model
extension as Gap 2. Hardest item; defer unless Gap 2 lands.

### Decision: P2 and P3 are deferred

Both Gap 2 and Gap 3 require the same prerequisite — a child path addressed by
**two** parent path params (`org-id` *and* `group-id`). `ParentType` is a single
string on `manifest.Resource`, and `PathParams`/`ParentReferenceFields` resolve
placeholders against one parent resource, so neither gap can be built without
first generalizing that model. That is a catalog design change, not an increment
of this issue.

They are therefore **explicitly scoped out here** and should be raised as a
follow-up issue once (and if) a consumer needs org-scoped group *hierarchies* or
membership. The flat org→group set — which is what iga-dash/syncengine asked for
— is fully served without them. Everything needed to resume is recorded in the
Gap 2 and Gap 3 sections above, including the Keycloak-side constraints that
make P2 more than a path change (no `subGroups` inlining, no
`briefRepresentation` on the children listing, attributes dropped on child
create).

## Acceptance Criteria

- [x] Org-scoped groups fetched via `--full-representation` include their `attributes` (closed by 0001's nested-collection extension).
- [x] Resolver and downward-graph behavior is pinned by tests so the currently-working disambiguation cannot regress: `("group", "organization", POST)` → org path, `("group", "")` → realm path, `("group", "group")` → realm children path. (`pkg/catalog/org_scoped_groups_test.go`)
- [x] An apply test proves a `group` with `ParentType: "organization"` creates via `POST /organizations/{org-id}/groups` with the org id resolved from the parent, and updates/deletes via the `{group-id}` endpoint. (`pkg/admin/apply_org_group_test.go`)
- [x] Realm-group behavior unchanged — no regression from any type-collision handling (`TestApplyRealmGroupStillUsesRealmEndpoint`).
- [x] Org-scoped subgroups (Gap 2) either traverse correctly via `/organizations/{org-id}/groups/{group-id}/children`, or are explicitly scoped out with the two-path-param limitation documented — **scoped out**, see the decision above.
- [x] Org-group membership (Gap 3) either registered as a relationship kind, or explicitly scoped out with a documented reason — **scoped out**, see the decision above.
- [x] `go vet ./...` clean.
- [ ] Library version bumped and tagged for release.

## Out of Scope

- Realm-group fetch/apply behavior, which already works and must not change.
- The top-level `briefRepresentation` plumbing itself — that is [ISSUE 0001](0001-full-representation-collection-fetch.md).
- Organization members and org identity-providers, already registered relationship kinds.

## Notes

- **Org groups are invisible to the realm-wide list.** `GET /groups` returns
  `[]` while an org-scoped group exists, and fetching its `parentId` directly is
  refused with `"Cannot manage organization related group via non Organization
  API."` The nested org path is the *only* way to enumerate them, so plain
  `fetch group` will never show them — worth stating in user docs.
- **Keycloak quirk:** `POST /organizations/{org-id}/groups/{group-id}/children`
  silently drops `attributes` from the payload (created child comes back with
  `attributes: {}`); a follow-up `PUT` on the child persists them. Any P2 apply
  implementation needs the create-then-update dance or the attributes are lost.
- The `identity.go` placeholder map already has `"org-id" -> "organization"` and
  `"group-id"/"groupId" -> "group"`, as the original request stated.
