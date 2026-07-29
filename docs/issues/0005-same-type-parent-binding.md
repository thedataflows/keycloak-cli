---
type: Issue
title: Same-type parent binding on the resource channel (nested groups)
description: A group nested under a group cannot be created through the resource channel — the parent id left in the payload to render /groups/{group-id}/children is not stripped, because ParentReferenceFieldNames skips every placeholder whose type equals the resource's own type, not just the resource's own primary identifier. Consumers must hand-build a relationship operation instead.
tags: [issue, enhancement, catalog, resolver, groups]
timestamp: 2026-07-29T00:00:00Z
---

# ISSUE 0005: Same-type parent binding on the resource channel (nested groups)

- **Type**: feature
- **Status**: done
- **Priority**: medium
- **Labels**: [catalog, resolver, groups, enhancement]
- **Assignee**: none
- **Related**: [ISSUE 0002](0002-org-scoped-groups.md) (deferred subgroup hierarchy; that gap is the two-*parent* case, this is the same-*type* single-parent case), [ISSUE 0003](0003-parent-scoped-collection-fetch-public-api.md) (sibling parent-scoping work), downstream consumer: iga-dash `syncengine/docs/issues/0027-realm-group-subgroups-are-never-synced.md`
- **Related code**: [`pkg/catalog/resolver.go`](../../pkg/catalog/resolver.go), [`pkg/catalog/dependencies.go`](../../pkg/catalog/dependencies.go), [`pkg/catalog/relationship_registry.go`](../../pkg/catalog/relationship_registry.go)
- **Closing commits**: `db1ed41` (option 1 + Gap 3 + update-body strip, verified live). Released in `v1.3.0`.

## Summary

A realm group nested under another realm group (`POST /admin/realms/{realm}/groups/{group-id}/children`) cannot be created through the resource channel. Routing a nested resource means leaving the parent's id in the resource `Data` (as `groupId`) so the resolver can render the path param, and the library strips such parent-reference fields from the request body before sending — **except** when the placeholder's resource type equals the resource's own type. For a group under a group the parent placeholder is itself of type `group`, so the guard that exists to protect a resource's *own* identifier also swallows the parent binding: the field reaches Keycloak and the POST fails. A consumer therefore has to hand-build a relationship operation on that path instead of using the resource channel, and cannot even register it as a relationship kind (that would delete the structural-graph edge its *read* depends on). This is the single-parent, same-type sibling of the two-parent gap [ISSUE 0002](0002-org-scoped-groups.md) deferred.

## Details

Verified 2026-07-28/29 against a live Keycloak **26.6.3** instance with the `26.6.2` spec, and against the resolver/registry code at `main`. The downstream consumer (iga-dash syncengine ISSUE 0027) hit every constraint below live and works around them with a hand-built relationship operation; this issue is to remove that workaround at the source.

### What already works — no change needed

- **Read.** `BuildDownwardGraph` already carries `group → /groups/{group-id}/children` (`pkg/catalog/dependencies.go:108`), so a `Depth: 1` group fetch returns subgroups as resources with `ParentType: "group"` and the parent's UUID injected as `groupId`. `briefRepresentation=false` propagates to the nested fetch, so a realm subgroup carries its `description` and `attributes` (unlike the org-scoped children listing — ISSUE 0002 Gap 2). The read is fine; do not touch it.
- **Single-group CRUD.** `("group", "", PUT/DELETE, single)` → `/groups/{group-id}` resolves correctly and addresses the group by its own id.

### Gap 1 (P1) — the parent binding is not stripped from a same-type nested create

`ParentReferenceFieldNames` (`pkg/catalog/resolver.go:319`) decides which `Data` fields are parent path-params to render-and-strip rather than send in the body. It already skips the resource's own primary identifier by **position** (`skipPrimary`, `resolver.go:327-338`: the placeholder that is the path's trailing `/{id}` segment). But it then also skips **by type** (`resolver.go:340`):

```go
parentType := placeholderMap[placeholder]
if parentType == "" || parentType == resourceType {
    continue
}
```

For `POST /groups/{group-id}/children`:

- `primary = group-id`, but the path ends in `/children`, so `skipPrimary` is `false` — the position guard does *not* fire.
- `placeholderMap["group-id"] = "group"` and `resourceType = "group"`, so the **type** guard fires and `group-id`/`groupId` is *not* added to the strip list.

So `groupId` stays in the request body and Keycloak rejects it:

```
POST /admin/realms/sync-target/groups/{parent}/children
  -> 400 {"error":"Invalid json representation for GroupRepresentation.
           Unrecognized field \"groupId\" at line 1 column 113."}
```

The type guard is meant to protect a resource's *own* identifier from being stripped, but the *position* guard (`skipPrimary`) already does that job precisely. The type guard is broader than its intent: it also swallows a *non-primary* same-type placeholder, which is exactly a same-type parent binding.

> A fake `Apply` does not validate the body, so this passes every fake-backed unit test and only fails against a real Keycloak — worth a spec-validation test, not just a fake.

### Gap 2 (P1) — the workaround (register a relationship kind) is not available here

The obvious consumer workaround — register `/groups/{group-id}/children` as a relationship override so writes go through the relationship channel — is actively harmful, because `BuildDownwardGraph` **excludes registered relationship read paths** (`pkg/catalog/dependencies.go:135-137`):

```go
if _, isRelationship := relationshipPaths[normalized]; isRelationship {
    continue
}
```

Registering the children path as a kind would delete the very `group → /groups/{group-id}/children` graph edge the **read** relies on. So a consumer that needs both directions cannot register the kind; it must read through the structural graph and hand-build a bare `manifest.RelationshipOperation` for the write (which is exactly what iga-dash does today: `realmSubgroupWriteTemplate`). One channel each way, on one path, none of it going through the resource abstraction.

### Gap 3 (the trap any fix must respect) — path params resolve from data before the identifier

`resolvePathParamValue` (`pkg/catalog/resolver.go:187`) resolves a path parameter from `Data` (kebab then camel) *before* falling back to the resource's own identifier (`resolver.go:197-199`). On the realm children path the parent placeholder and the self placeholder are the **same** token `{group-id}`: it identifies the parent in `/groups/{group-id}/children` and the group itself in `/groups/{group-id}`. So if a fix lets `groupId` live in a child's `Data`, an **update or delete** of that child — which resolves `group-id` on the single-group path — would read `group-id` from `data["groupId"]` (the parent) and render the *parent's* path: a PUT that overwrites the parent with the child's representation, or a DELETE that removes the parent and its whole subtree. Any fix must keep the parent binding load-bearing on the **create** path only, and make update/delete address the group by its own id regardless of a stray `groupId` in the data.

### Proposed shape

Two candidate fixes; the implementer picks one and records the reasoning:

1. **Narrow the type guard to position (smallest change).** Drop the `parentType == resourceType` clause at `resolver.go:340` and rely on the existing `skipPrimary`/`primary` position check to protect the resource's own identifier. A non-primary same-type placeholder (a same-type parent) then correctly becomes a parent-reference field that is rendered into the path and stripped from the body. This must be proven not to regress single-resource same-type endpoints (`/groups/{group-id}`), where `skipPrimary` already fires, and not to regress any other same-type nesting the spec contains.
2. **An explicit `ParentID` on `manifest.Resource`.** Carry the parent id out-of-band instead of in `Data`, so it is never a body field and never collides with the self identifier via `resolvePathParamValue`. Larger surface, but it also disambiguates Gap 3 structurally rather than by precedence rules, and generalizes toward the two-parent case ISSUE 0002 deferred.

Whichever lands, the resolver must guarantee Gap 3's precedence: on a single-resource operation, the resource's own identifier wins over any same-typed value in `Data`.

### Downstream requirement (why this is P1)

iga-dash ISSUE 0027: realm-group subgroups. syncengine wants:

```go
// create a realm subgroup through the resource channel, parent bound by id, stripped from body
manifest.Resource{Type: "group", ParentType: "group",
    Data: map[string]any{"name": "platform", "groupId": parentUUID /* or via ParentID */}}
// should resolve to POST /groups/{group-id}/children with groupId NOT in the body
```

Today syncengine cannot use the resource channel for this and hand-builds a relationship operation on a hardcoded `{realm}/groups/{group-id}/children` template, precisely because of Gaps 1–3. Closing this issue lets that hardcode be retired.

## Acceptance Criteria

- [x] A `group` with `ParentType: "group"` creates via `POST /groups/{group-id}/children` with the parent id resolved from the parent binding and **stripped** from the request body (no `groupId` field on the wire) — `TestApplyNestedRealmGroupStripsParentBindingFromBody` validates the body against the embedded spec; verified live (create `201`)
- [x] Update and delete of a nested group resolve to `/groups/{group-id}` addressing the group **by its own id**, never the parent's, even when a parent id is present in the resource's `Data` (Gap 3) — guarded by `TestNestedGroupUpdateDeleteAddressOwnID` and the apply tests; verified live (update `204`, parent intact)
- [x] The structural read is unchanged: `group → /groups/{group-id}/children` is still in `BuildDownwardGraph` (`TestNestedGroupReadEdgeUnchanged`); the fix does not register the path as a relationship kind
- [x] Resolver behaviour pinned by tests: `("group", "group", POST, collection)` → `/groups/{group-id}/children` with `group-id` a rendered-and-stripped parent reference; `("group", "", …)` still → `/groups` and `/groups/{group-id}` (`same_type_parent_test.go`)
- [x] Realm-group, org-scoped-group and client-scoped-role behaviour unchanged — no regression; every same-type nesting in the spec is recorded by `TestSameTypeParentNestingsAreRecorded` (only the realm children path is resource-channel reachable)
- [x] `go vet ./...` clean, `go test ./...` green (except a pre-existing, unrelated `pkg/output` TOML failure)
- [x] Tagged release so iga-dash can bump `syncengine/go.mod`. Note for the consumer: iga-dash vendors three trees (`igadash/vendor`, `syncengine/vendor`, and the root `go work vendor`) — all three must be regenerated together — released in `v1.3.0`

## Out of Scope

- **Two-parent** nested paths (`/organizations/{org-id}/groups/{group-id}/children`, org-group membership) — that is ISSUE 0002 Gap 2/3, a different generalization (`ParentType` is a single string). If the `ParentID` route is chosen here, note whether it is a stepping stone toward that, but do not take it on.
- More than one level of realm nesting (a grandchild). The read is `Depth`-bounded and the consumer models a single parent name; deeper hierarchies are a separate concern.
- Moving a group between parents as a first-class move. The consumer treats a re-parent as delete+create; nothing here changes that.

## Notes

- **Why a fake test is not enough.** The downstream consumer's first live push failed exactly here while every fake-backed unit test passed — a fake `Apply` validates no body. Both children write paths in iga-dash now have a spec-validation test over the exact operation on the wire; that is the pattern to copy.
- The realm children path and the single-group path sharing the `{group-id}` token (Gap 3) is the sharpest hazard: a fix that is correct for create but sloppy about update/delete precedence is worse than the current hardcode, because it can delete a parent subtree. Lead with that test.
- The org-scoped children path (`/organizations/{org-id}/groups/{group-id}/children`) is *not* covered by this issue and would still need the two-parent model — but note that its parent placeholders are `organization` + `group`, so the same-type guard does not bite there the way it does on the realm path.

### Implementation notes (resolved)

- **Chose option 1** (narrow the type guard), the smallest change: dropped the `parentType == resourceType` clause in `ParentReferenceFieldNames` and rely on the existing `skipPrimary` position guard to protect the resource's own id. Recorded every same-type nesting in the 26.6.2 spec (`TestSameTypeParentNestingsAreRecorded`): only `group → /groups/{group-id}/children` is reachable via the resource channel; the two-parent org children path and the `client → registration-access-token` sub-action are not resolved by any canonical create, so the narrowing does not regress them.
- **Gap 3 was worse than the issue framed it — it also bites *create*.** The create-time existence probe (`resolveExistingResource` → single GET) resolves `{group-id}` with no own id yet; before the fix it fell through to the camel `groupId` (the parent) and would probe — then update — the parent. The fix makes the self placeholder on a single-resource path resolve **only** from the resource's own identifier (or an exact-kebab field), never the camel/parent-binding form. This covers update, delete, and the create probe. Exact-kebab overrides (`data["user-id"]`) still win, and the "primary is always present" contract is preserved.
- **Live testing found a second body leak the create-only view missed.** On the
  single-resource PUT the parent placeholder `{group-id}` *is* the self id, so
  `groupId` is invisible to the single contract and leaked into the update body —
  a live update failed `400 Unrecognized field "groupId"` (the parent stayed
  intact, so Gap 3's path precedence held; only the body was wrong). Fixed in
  `sanitizeResourcePayload`: when `ParentType` is set, also strip the
  *nested-create* contract's parent-reference fields on every write, since the
  binding is a property of the create endpoint, not the method. Guarded by
  `TestApplyNestedRealmGroupUpdateAddressesChildAndStripsParentBinding`.
- Tests: resolver pins + Gap 3 destructive-path guard (`same_type_parent_test.go`), the nesting record (`same_type_parent_internal_test.go`), and end-to-end applies that validate the create *and* update bodies against the embedded spec — `groupId` absent on the wire — and assert no write ever touches the parent (`apply_nested_group_test.go`).
- **Verified live** against Keycloak 26.6.3 (realm `claude-0005-test`, throwaway): create `201`, idempotent update `204`, parent group intact after the child update, child re-parented attributes applied. Full create→update cycle through the resource channel with no hand-built relationship operation.
