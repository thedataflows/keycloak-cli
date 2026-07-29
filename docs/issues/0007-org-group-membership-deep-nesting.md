---
type: Issue
title: Organization-group membership is unreadable below the top level
description: The organization-group-member read returns members only for an organization's TOP-LEVEL groups, because the traversal that fires it builds its parent index from the flat /organizations/{org-id}/groups collection. A member of a nested org group (a child or grandchild) is never read, so it cannot be mirrored. This is the membership sibling of ISSUE 0006, which fixed only the containment read (FetchChildren).
tags: [issue, enhancement, admin, catalog, groups, organizations, membership]
timestamp: 2026-07-29T00:00:00Z
---

# ISSUE 0007: Organization-group membership is unreadable below the top level

- **Type**: feature
- **Status**: open
- **Priority**: medium
- **Labels**: [admin, catalog, groups, organizations, membership, enhancement]
- **Assignee**: none
- **Related**: [ISSUE 0006](0006-org-scoped-group-deep-nesting.md) (the containment sibling — fixed org group *reads* to arbitrary depth via FetchChildren; this is the same gap for *membership*), [ISSUE 0002](0002-org-scoped-groups.md), [ISSUE 0005](0005-same-type-parent-binding.md), downstream consumer: iga-dash `syncengine/docs/issues/0026-org-group-membership-and-child-groups.md` (G1) and the nested-membership fixtures in its seeder
- **Related code**: [`pkg/admin/fetch.go`](../../pkg/admin/fetch.go), [`pkg/catalog/dependencies.go`](../../pkg/catalog/dependencies.go), [`pkg/catalog/relationship_registry.go`](../../pkg/catalog/relationship_registry.go)
- **Closing commits**: none

## Summary

An organization's groups can nest, and a user can be a member of a group at any level (`/organizations/{org-id}/groups/{group-id}/members`). The library reads members only for an organization's **top-level** groups: the traversal that fires the org-group membership read binds `{group-id}` from the flat `/organizations/{org-id}/groups` collection, which lists top-level org groups only. A member of a nested org group — a child or a grandchild — is never returned, so a consumer cannot mirror it. This is the membership sibling of [ISSUE 0006](0006-org-scoped-group-deep-nesting.md), which fixed the *containment* read (`FetchChildren`) but not the *membership* read.

## Details

Verified 2026-07-29 against a live Keycloak **26.6.3** with the `26.6.2` spec (keycloak-cli **v1.4.0**), using a consumer that registers the org-group membership and child kinds and reads org groups recursively via the v1.4.0 `FetchChildren`.

Source tree (an org group, its child, its grandchild — a member seeded at every level):

```
org demo-a-org-1
  demo-a-orggroup-1                         member: demo-a-user-1
    demo-a-orgchild-1                        member: demo-a-user-1   <- nested
      demo-a-orggrandchild-1                 member: demo-a-user-1   <- nested
```

Result after sync: the top-level org group's member crosses; the child's and grandchild's members do **not**. Containment (the child and grandchild groups themselves) *does* cross — that is ISSUE 0006's fix — so the target has the nested groups but they are memberless, while the source has the members.

### Why — same flat parent index as ISSUE 0006

The org-group membership read fires per **parent** in a parent index built from the flat `/organizations/{org-id}/groups` collection (top-level org groups only). A nested org group is never a parent in that index, so its members endpoint is never called — independent of `Depth`. ISSUE 0006 addressed this for the *children* collection by giving `FetchChildren` a scoped selector for the org children path; the *members* collection did not get the same treatment.

### Why FetchChildren does not already cover it

`FetchChildren(parent, childType, …)` resolves a child collection via `ScopedChildCollection`, keyed on structural containment. The members endpoint returns **users**, not groups, and is a relationship/override path, not a structural child — so there is no `ScopedChildCollection("group","organization","user")` selecting `/organizations/{org-id}/groups/{group-id}/members`. A consumer that has each nested org group (from ISSUE 0006's recursion) still has no supported call to read *its* members.

### Contrast with realm groups (works)

Realm-group membership is read from the **user** side (`/users/{user-id}/groups`), which returns a user's realm groups at any depth in one call — so a realm child's and grandchild's members mirror with no special handling. Organization-group membership has no user-side read (org groups are a separate namespace, refused through the realm API), so it depends entirely on the group-side members endpoint, which is gated by the flat parent index above.

## Proposed shape

Either (implementer picks one, records the reasoning):

1. **A scoped members fetch.** Let `FetchChildren` (or a sibling) select the org-group members collection for an organization-parented group — `ScopedChildCollection("group","organization","user")` → `/organizations/{org-id}/groups/{group-id}/members` — so a consumer holding each nested org group (via ISSUE 0006) can read its members with one addressed GET and no override. Mirrors ISSUE 0006's fix, one endpoint over.
2. **Deepen the membership parent index.** Seed the org-group membership read's parent index with the org children discovered during the traversal, so it fires for nested org groups too. Server-driven, but a deeper change to index-building (and pairs O(orgs × groups) as the traversal already warns).

## Acceptance Criteria

- [ ] A member of a nested org group (child and grandchild) is reachable through a supported `admin.Service` API, keyed to that group — verified against the fixture above (a member seeded on the top-level org group, its child, and its grandchild; all three read back)
- [ ] Works at arbitrary depth, bounded only by the caller
- [ ] Top-level org-group membership and realm-group membership are unchanged — no regression
- [ ] `go vet ./...` clean, `go test ./...` green
- [ ] Tagged release so iga-dash can bump `syncengine/go.mod` (three vendor trees regenerated together)

## Out of Scope

- Organization-group **containment** at depth — done in ISSUE 0006 (v1.4.0).
- Realm-group membership at depth — already works via the user-side read.
- Org invitations, org identity-provider links.

## Notes

- Reproduction: seed `demo-a-user-1` as a member of an org top-level group, its child, and its grandchild; sync; observe only the top-level membership crossing. The downstream consumer (iga-dash) now seeds exactly this so the gap is visible on every clean stack.
- The org group items still advertise no `subGroupCount` and an empty inline `subGroups`, so any consumer recursion here terminates on an empty read, as with ISSUE 0006.
