---
type: Issue
title: Organization-group membership is unreadable below the top level
description: The organization-group-member read returns members only for an organization's TOP-LEVEL groups, because the traversal that fires it builds its parent index from the flat /organizations/{org-id}/groups collection. A member of a nested org group (a child or grandchild) is never read, so it cannot be mirrored. This is the membership sibling of ISSUE 0006, which fixed only the containment read (FetchChildren).
tags: [issue, enhancement, admin, catalog, groups, organizations, membership]
timestamp: 2026-07-29T00:00:00Z
---

# ISSUE 0007: Organization-group membership is unreadable below the top level

- **Type**: feature
- **Status**: done
- **Priority**: medium
- **Labels**: [admin, catalog, groups, organizations, membership, enhancement]
- **Assignee**: none
- **Related**: [ISSUE 0006](0006-org-scoped-group-deep-nesting.md) (the containment sibling — fixed org group *reads* to arbitrary depth via FetchChildren; this is the same gap for *membership*), [ISSUE 0002](0002-org-scoped-groups.md), [ISSUE 0005](0005-same-type-parent-binding.md), downstream consumer: iga-dash `syncengine/docs/issues/0026-org-group-membership-and-child-groups.md` (G1) and the nested-membership fixtures in its seeder
- **Related code**: [`pkg/admin/fetch.go`](../../pkg/admin/fetch.go), [`pkg/catalog/dependencies.go`](../../pkg/catalog/dependencies.go)
- **Closing commits**: `3222fb0` (GET-only selector), `4642746` (FetchChildren members), `6e570a3` (depth exposure + read-only apply). Verified live. Released in `v1.5.0`.

## Summary

An organization's groups can nest, and a user can be a member of a group at any level (`/organizations/{org-id}/groups/{group-id}/members`). The library reads members only for an organization's **top-level** groups: the traversal that fires the org-group membership read binds `{group-id}` from the flat `/organizations/{org-id}/groups` collection, which lists top-level org groups only. A member of a nested org group — a child or a grandchild — is never returned, so a consumer cannot mirror it. This is the membership sibling of [ISSUE 0006](0006-org-scoped-group-deep-nesting.md), which fixed the *containment* read (`FetchChildren`) but not the *membership* read.

**Update (2026-07-29, v1.4.1) — partially resolved, one level.** v1.4.1's
`feat(admin): depth traversal descends org group hierarchies` (ISSUE 0006 follow-up)
also deepened the membership read: because it seeds the org relationship parent
index with depth-fetched org children, the member read now fires for an org group's
**direct children** too. Verified live end to end by the iga-dash consumer: a member
of a top-level org group **and** of its child now sync; a member of a **grandchild**
(a child of a child) still does not — the traversal descends one level, not
arbitrarily. So the remaining gap is org membership at **depth ≥ 2** (grandchild and
below); the acceptance criteria below still stand for that. A targeted per-group
members fetch (proposed shape #1) would close it at any depth without relying on
traversal depth.

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

- [x] A member of a nested org group (child and grandchild) is reachable through a supported `admin.Service` API (`FetchChildren(orgGroup, "member")`), keyed to that group — **verified live** on `sync-source`: `demo-a-user-1` read back at `demo-a-orggroup-1` (top), `demo-a-orgchild-1` (child), and `demo-a-orggrandchild-1` (grandchild)
- [x] Works at arbitrary depth, bounded only by the caller — `FetchChildren` selects the members path at every level via the scoped selector; also surfaced through `fetch --depth`
- [x] Top-level org-group membership and realm-group membership are unchanged — no regression (relationship pass untouched; the 125-edge baseline is unchanged; realm membership still reads user-side)
- [x] `go vet ./...` clean, `go test ./...` green (except a pre-existing, unrelated `pkg/output` TOML failure)
- [x] Tagged release so iga-dash can bump `syncengine/go.mod` (three vendor trees regenerated together) — released in `v1.5.0`

## Out of Scope

- Organization-group **containment** at depth — done in ISSUE 0006 (v1.4.0).
- Realm-group membership at depth — already works via the user-side read.
- Org invitations, org identity-provider links.

## Notes

- Reproduction: seed `demo-a-user-1` as a member of an org top-level group, its child, and its grandchild; sync; observe only the top-level membership crossing. The downstream consumer (iga-dash) now seeds exactly this so the gap is visible on every clean stack.
- The org group items still advertise no `subGroupCount` and an empty inline `subGroups`, so any consumer recursion here terminates on an empty read, as with ISSUE 0006.

### Implementation notes (resolved)

- **Chose the read channel (proposed shape 1), not the relationship pass (shape 2).**
  The members endpoint is **GET-only** (no POST/DELETE) and returns
  `MemberRepresentation`; org-group membership is a read view. The relationship
  subsystem — as designed in the initial commit — is uniformly *writable* and
  generic (no per-kind branches). Forcing a read-only, two-nested-parent read
  through it would have required inventing read-only relationship kinds and a
  per-shape special-case in the generic engine — against the original
  architecture. The read channel (`ScopedChildCollection`/`FetchChildren`, the
  ISSUE 0006 machinery) already fits a GET-only scoped collection exactly.
- **`ScopedChildCollection` gained a GET-only fallback.** It preferred POST
  (structural containment); now, when no POST create exists, it matches the
  child type's GET collection by the same placeholder chain. So
  `("group","organization","member")` → `/organizations/{org-id}/groups/{group-id}/members`.
- **`FetchChildren(orgGroup, "member")`** returns members tagged
  `member`/`organization` with `orgId` propagated and the immediate `groupId`
  injected to key each member to its group — safe because a member is never
  recursed as a group parent (the ISSUE 0005 Gap 3 collision applies only to
  same-type nesting, so nested groups still get no `groupId`).
- **CLI via `fetch --depth`.** The depth traversal reads each org group's members
  and emits them as `member` resources. `resourceKey` keys a member by
  (user, group) so the *same* user's memberships across groups are not collapsed
  by dedup. Apply filters read-only `member` resources (no write path) and records
  them skipped, so a `fetch → upload` round-trip does not error.
- **Verified live** against Keycloak 26.6.3, `sync-source`, no override:
  `fetch organization demo-a-org-1 --depth 4` returns `demo-a-user-1` as five
  distinct membership edges across the org group tree (top-level, child, and
  grandchild levels), where the default path read zero.
