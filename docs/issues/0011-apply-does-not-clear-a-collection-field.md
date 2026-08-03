---
type: Issue
title: apply does not clear a collection field (an explicit empty is not removed on the target)
description: Applying a resource whose data carries an explicit-empty collection (e.g. attributes {}) to remove it does not clear that collection on the target. The resource updates but the collection keeps its previous values. A direct full-representation PUT with the same explicit-empty collection does clear it, so Keycloak supports removal; the gap is in the library apply path. Distinct from ISSUE 0010 (which was attributes not persisting, a read-side brief-representation issue) — this is the removal/clear case on the write side, observed for groups.
tags: [issue, bug, apply, attributes, group, collection, data-loss]
---

# ISSUE 0011: apply does not clear a collection field (an explicit empty is not removed on the target)

- **Type**: bug
- **Status**: done
- **Priority**: medium
- **Labels**: [apply, attributes, group, collection]
- **Assignee**: none
- **Related**: [ISSUE 0010](0010-organization-apply-drops-attributes.md) (attributes not persisting — read side; this is the write-side removal/clear case)
- **Related code**: [`pkg/admin/apply.go`](../../pkg/admin/apply.go), [`pkg/admin/internal/client.go`](../../pkg/admin/internal/client.go)
- **Closing commits**: `c65f8b0` (`sanitizeResourceData` preserves explicit-empty collections via `isExplicitEmptyCollection` + group clear regression tests). Released in `v1.6.3`.

## Resolution

Confirmed as the **empty-value omission** hypothesis, not the create-on-existing
route. `sanitizeResourceData` (`pkg/admin/apply.go`) ran once over every
resource at the top of `Apply` and dropped any value that sanitized to an empty
map via `isEmptyValue`. An explicit `"attributes": {}` was therefore stripped
before the request body was ever built, so both the create (`POST`) and update
(`PUT`) paths sent a body that omitted the collection — and Keycloak preserves a
collection it is not sent, so the previous values survived and the leg never
converged. Empty **lists** (`[]`) already survived, because `isEmptyValue` only
flags empty maps.

**Fix:** `sanitizeResourceData` and `sanitizeMap` now keep a collection that was
**explicitly empty in the input** (`isExplicitEmptyCollection` — an already-empty
`{}` or `[]`), while still dropping a map that only *becomes* empty after
nil-stripping (e.g. `{"stale": null}`) and still dropping plain nulls. This
distinguishes "the caller wants to clear this field" from "Keycloak returned a
null-only object we must not send." The write path already marshals the raw
`resource.Data` directly (`pkg/admin/internal/client.go`), so once the empty map
survives sanitization it reaches Keycloak unchanged on both paths.

Regression tests capture the actual HTTP body and assert the collection is
**sent** as `{}` (not omitted) and the resource ends up cleared, on both create
and update: `pkg/admin/apply_clear_collection_test.go`. The
`sanitizeResourceData` unit contract for explicit-empty vs collapsed-empty maps
is in `pkg/admin/apply_internal_test.go`.

## Summary

`Apply` on a resource whose `Data` carries an **explicit empty** collection
(e.g. `"attributes": {}`), intending to **remove** that collection from the
target, does not clear it: the resource updates but the collection keeps its
previous values. A direct full-representation `PUT` with the same explicit-empty
collection **does** clear it, so Keycloak supports removal — the gap is inside
the library apply path. Observed for groups.

## Details

Reported downstream by the syncengine A→B Keycloak sync (thedataflows/iga-dash,
its ISSUE 0033), using keycloak-cli `v1.6.2`. The connector implements a
"drop this field" feature: to remove a group's `attributes`, it builds an apply
whose data has `attributes` set to an explicit empty map and calls
`svc.Apply(..., manifest.Resource{Type: <group>, Data: groupMap})`.

Evidence the input is correct and the drop is inside the library:

1. The manifest `Data` at the `Apply` call had `attributes` present and empty
   (`{}`) — confirmed by logging `groupToMap` output in the connector:
   `attrs={} hasAttrsKey=true` for the affected groups.
2. `adapter.UpdateGroup` was invoked with the **real target group id** (the same
   id used in step 3's manual PUT).
3. A direct `PUT /admin/realms/{realm}/groups/{id}` with
   `{"id": "...", "name": "...", "attributes": {}}` **clears** the group's
   attributes (HTTP 204), verified against the same server/version.
4. Despite (1)–(3), after the sync applies, the target group's `attributes`
   remain populated (previous values), and the leg keeps re-issuing the update
   every cycle (never converges).

So the library, given an apply whose data carries an explicit-empty collection,
does not send a request that Keycloak treats as "clear this collection."

## Reproduction

1. Apply a group with a non-empty `attributes` map; confirm it lands.
2. Apply the same group again with `attributes` set to an explicit empty map
   (`{}`) — the "remove the attributes" intent.
3. Read the group back. Observe `attributes` is unchanged (still populated).
   Expected: `attributes` is empty/removed, matching a direct
   `PUT .../groups/{id}` with `attributes: {}`.

## Acceptance Criteria

- [x] Applying a resource whose `Data` carries an explicit-empty collection
      (`{}` for a map, `[]` for a list) clears that collection on the target,
      for both the create and update paths. *(Fixed: `sanitizeResourceData`
      preserves explicitly-empty collections via `isExplicitEmptyCollection`.)*
- [x] A regression test captures the actual HTTP request body for an
      explicit-empty-collection apply and asserts the collection is **sent**
      (as `{}`/`[]`, not omitted) and that the resource ends up with the
      collection cleared. *(`pkg/admin/apply_clear_collection_test.go` — group
      create + update; unit contract in `pkg/admin/apply_internal_test.go`.)*

## Out of Scope

- Persisting non-empty attributes (ISSUE 0010 — already working).
- Per-attribute (partial) removal / merge semantics; this is whole-collection
  clear via full-representation replacement.

## Notes

Hypotheses for where the clear is lost (unconfirmed — the maintainer's
unredacted view of `pkg/admin` will pin it):

- **Create-on-existing / 409 path.** The downstream apply arrived classified as
  a create for an already-existing group (the connector resolves the existing id
  and calls the update, but the delta op was `create`). If `applyResource`
  treats these as create → `POST` → 409, the 409-fallback update may not carry
  the explicit-empty collection body. Instrumentation in the vendored copy did
  not observe a `PUT /groups/{id}` during the group apply, which points here —
  though ISSUE 0010's resolution states the update path PUTs the raw `Data`, so
  the divergence may be specific to the create-on-existing route and/or to an
  **empty** collection value.
- **Empty-value omission.** An empty map/list may be dropped before the request
  is built (marshaling, `sanitizeResourceData`, or request validation), so the
  body Keycloak receives omits the collection and Keycloak preserves it.

Because ISSUE 0010 established the non-empty write path is correct, the likely
delta here is the **empty-collection** case and/or the **create-on-existing**
route for groups. A body-capturing test for the empty case (create + update)
would confirm which.

Reported by: thedataflows/iga-dash ADR 0020 `target_field_drop` (its ISSUE 0033).
