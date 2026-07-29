---
type: Issue
title: Organization-scoped group hierarchies are unreachable below one level
description: The organization-group-child relationship override only ever fires with a top-level org group as its parent, because the traversal builds its parent index from the flat /organizations/{org-id}/groups collection. So a grandchild (a child of an org child) never surfaces regardless of Depth, and no targeted read reaches it — FetchChildren resolves via the structural downward graph, which excludes the org children path. Organization group trees deeper than one level cannot be fetched.
tags: [issue, enhancement, admin, catalog, groups, organizations]
timestamp: 2026-07-29T00:00:00Z
---

# ISSUE 0006: Organization-scoped group hierarchies are unreachable below one level

- **Type**: feature
- **Status**: open
- **Priority**: medium
- **Labels**: [admin, catalog, groups, organizations, enhancement]
- **Assignee**: none
- **Related**: [ISSUE 0002](0002-org-scoped-groups.md) (deferred the org subgroup hierarchy — this is that gap, now with a concrete consumer and a repro), [ISSUE 0003](0003-parent-scoped-collection-fetch-public-api.md) (the parent-scoped-fetch shape a fix likely extends), [ISSUE 0005](0005-same-type-parent-binding.md) (the realm half of deep group nesting, which is reachable and fixed), downstream consumer: iga-dash `syncengine/docs/plans/2026-07-29-nested-group-hierarchy.md` (Task 4 spike) and `syncengine/docs/issues/0027-realm-group-subgroups-are-never-synced.md`
- **Related code**: [`pkg/admin/fetch.go`](../../pkg/admin/fetch.go), [`pkg/catalog/dependencies.go`](../../pkg/catalog/dependencies.go), [`pkg/catalog/relationship_registry.go`](../../pkg/catalog/relationship_registry.go)
- **Closing commits**: none

## Summary

A Keycloak 26 organization owns a group hierarchy under `/organizations/{org-id}/groups`, and those groups can nest arbitrarily (`…/groups/{group-id}/children`). The library reads exactly **one** level of that hierarchy: the children of a top-level org group. A grandchild — a child of an org child — is never returned, at any `Depth`, and there is no targeted API that reaches it. Organization group trees deeper than one level are therefore unfetchable, so a consumer cannot mirror them.

## Details

Verified 2026-07-29 against a live Keycloak **26.6.3** with the `26.6.2` spec, using a consumer that has registered the `organization-group-child` relationship override on `/organizations/{org-id}/groups/{group-id}/children` (the pattern [ISSUE 0002](0002-org-scoped-groups.md) established; the org children path is not in the built-in registry, so the consumer supplies it).

### The observation

Source realm tree:

```
org demo-a-org-1
  demo-a-orggroup-1
    demo-a-orgchild-1            (a child — one level down)
      demo-a-orggrandchild-1     (a grandchild — two levels down)
  demo-a-orggroup-2
    demo-a-orgchild-1
      demo-a-orggrandchild-1
```

Fetching organizations with the override installed, at `Depth: 2`, `3`, and `4`, and counting `organization-group-child` relationship operations in the report:

| Depth | org-group-child ops returned | which |
|---|---|---|
| 2 | 2 | `demo-a-orgchild-1` under each of the two org groups |
| 3 | 2 | same two children |
| 4 | 2 | same two children |

The grandchild `demo-a-orggrandchild-1` **never appears**, at any depth. Raising `Depth` does not descend past the first level of org-group children.

### Why — the parent index is flat

The relationship pass fires the `organization-group-child` override once per **parent** in its parent index, and that index is built from the **flat** `/organizations/{org-id}/groups` collection — which returns only the organizations' *top-level* groups (verified: the list carries `subGroups: []` and `subGroupCount: null` for every entry, so it advertises no hierarchy). The override therefore only ever binds `{group-id}` to a top-level org group and reads *its* children. An org child (e.g. `demo-a-orgchild-1`) never enters the parent index, so the override never fires with the child as parent, so the grandchild is never read — independent of `Depth`, because deepening the traversal does not add org children to that index.

### Why there is no escape hatch

`admin.Service` exposes `Spec`, `Fetch`, `FetchChildren`, `Apply`. None reaches an org child's children:

- **`Fetch`** is realm-sweep-shaped; the org children arrive only through the override during the org traversal, indexed as above. There is no way to inject a specific org child as a parent.
- **`FetchChildren`** (ISSUE 0003) resolves the child collection path from `BuildDownwardGraph`, which **excludes registered relationship read paths** (`pkg/catalog/dependencies.go`) — the org children path is exactly such a path. And the graph's only `group → group` edge is the **realm** children path `/groups/{group-id}/children`, which Keycloak rejects for an org-related group (`400 "Cannot manage organization related group via non Organization API."`). So `FetchChildren("group" under "group")` cannot target the org path.

### Contrast with the realm half (ISSUE 0005, reachable)

Realm group nesting to arbitrary depth *is* reachable, and the consumer implemented it: realm groups sit in `BuildDownwardGraph` (`group → /groups/{group-id}/children`), so `FetchChildren` targets a specific group's children and the consumer recurses per parent, terminating on `subGroupCount == 0`. The organization half has neither property — its children path is an override the graph excludes, and its groups advertise no `subGroupCount` — so the same recursion is not expressible against the current API.

## Proposed shape

Either (the implementer picks one and records the reasoning):

1. **A targeted parent-scoped relationship read.** Extend `FetchChildren` (or add a sibling) so it can read a *relationship/override* child collection for one `(parent, kind)` pair — not only structural-graph paths — issuing one addressed GET on `/organizations/{org-id}/groups/{group-id}/children` for a caller-supplied org group id. The consumer then reads each discovered org group's children and recurses, exactly as it does for realm groups via `FetchChildren`. This mirrors ISSUE 0003's spirit and keeps depth on the caller's side.
2. **Recursive parent indexing in the traversal.** When `Depth > N`, seed the `organization-group-child` parent index with the org children discovered at depth `N-1`, so the override fires again with each child as parent. This keeps depth server-driven but is a deeper change to the traversal's index-building.

Whichever lands, note the terminator: org groups expose **no `subGroupCount`** and an **empty inline `subGroups`** in the collection, so a consumer's recursion must terminate on an *empty children read*, not on a count.

## Acceptance Criteria

- [ ] A grandchild under an org child (`/organizations/{org-id}/groups/{parent}/{child}/{grandchild}`) is reachable through a supported `admin.Service` API, returned with `ParentType`/parent-reference identifying its immediate parent — verified against the fixture above (two grandchildren, one per top-level org group)
- [ ] Arbitrary depth works, not just two levels (a great-grandchild is reachable by the same mechanism), bounded only by the caller
- [ ] The org children **create** path (`POST /organizations/{org-id}/groups/{group-id}/children` with `{group-id}` = an org *child*) is confirmed to nest under a child, not only a top-level org group — pin it with a test (the read is the blocker; confirm the write is not a second one)
- [ ] Realm group nesting (ISSUE 0005) and one-level org groups (ISSUE 0002) are unchanged — no regression
- [ ] `go vet ./...` clean, `go test ./...` green
- [ ] Tagged release so iga-dash can bump `syncengine/go.mod`; note the three vendor trees must be regenerated together (`igadash/vendor` is currently absent — `syncengine/vendor` + the root `go work vendor`)

## Out of Scope

- Realm-group deep nesting — reachable today and implemented downstream (ISSUE 0005).
- Organization-group **membership** at depth, org invitations, org identity-provider links.
- The consumer-side sync model (how iga-dash keys and provisions nested org groups cross-realm) — that is iga-dash's plan; this issue only unblocks the read.

## Notes

- **Reproduction.** With the `organization-group-child` override registered, seed a grandchild under an org child, then:
  ```go
  for _, depth := range []int{2, 3, 4} {
      rep, _ := svc.Fetch(ctx, admin.FetchQuery{
          Realm: "sync-source", Resources: "organization", Depth: depth, FullRepresentation: true,
      })
      // count rep.Relationships where Kind == "organization-group-child"
  }
  ```
  Each depth returns the same two first-level children; the grandchild is absent.
- The org children collection advertises no hierarchy signal (`subGroupCount: null`, `subGroups: []`), unlike realm groups — worth surfacing whatever `subGroupCount`/`subGroupsCount` Keycloak *does* expose on the org path so a consumer can avoid a GET on a known-leaf, but the blocker is reachability, not the terminator.
