# Design: Organization-scoped group hierarchies at arbitrary depth (ISSUE 0006)

- **Date**: 2026-07-29
- **Issue**: [ISSUE 0006](../../issues/0006-org-scoped-group-deep-nesting.md)
- **Approach chosen**: Native + caller-driven (option 1 in the issue)

## Problem

A Keycloak 26 organization owns a group hierarchy under
`/organizations/{org-id}/groups`, and those groups nest arbitrarily via
`/organizations/{org-id}/groups/{group-id}/children`. The library reads only the
first level of that hierarchy (the flat groups of an organization); a grandchild
is never returned at any `Depth`, and no targeted API reaches it.

Two facts, verified against the 26.6.2 spec, cause this:

1. `BuildDownwardGraph` dedupes edges by `(parentType, childType)` and keeps a
   single path per pair. `group → group` therefore resolves only to the **realm**
   children path `/groups/{group-id}/children`; the org children path is dropped.
2. The resolver selects an operation by `(resourceType, parentType)` where
   `parentType` matches the *first* parent placeholder. Both the flat org groups
   path and the org children path lead with `{org-id}` (organization), so
   `("group","organization")` resolves to the flat `/organizations/{org-id}/groups`,
   never the children path. The org children path needs **two** parent
   placeholders (`org-id` + `group-id`), which the single-`ParentType` model
   cannot select — the two-parent gap ISSUE 0002 deferred.

The realm half (ISSUE 0005) is reachable because realm groups sit in the downward
graph and `FetchChildren` (ISSUE 0003) targets a specific group's children; the
consumer recurses per parent. The org half has neither property.

## Approach

Native + caller-driven: make the **library** reach org group children with no
consumer-side relationship override, and let the caller drive depth via
`FetchChildren` recursion — exactly as it already does for realm groups.

### 1. Scoped child-collection selection (`catalog`)

Add a `Spec` helper that selects the child collection path by matching the
parent's **placeholder chain**, rather than the deduped downward graph. Given
`(parentType, parentParentType, childType)`, enumerate the child's POST
collection paths that have a matching GET collection (the existing
`BuildDownwardGraph` predicates, without the dedupe) and pick the one whose
ordered non-realm parent placeholder *types* match the parent's chain:

- **Top-level parent** (`parentParentType == ""`): exactly one non-realm parent
  placeholder, of type `parentType`.
- **Nested parent** (`parentParentType != ""`): the trailing parent placeholder
  is `parentType` and the placeholder before it is `parentParentType`.

Results:

| Parent (`type`, `parentType`) | child | Selected path |
|---|---|---|
| `client`, `""` | role | `/clients/{client-uuid}/roles` |
| `group`, `""` | group | `/groups/{group-id}/children` |
| `group`, `organization` | group | `/organizations/{org-id}/groups/{group-id}/children` |

The helper also returns the **inherited parent-reference fields**: the camel-case
field names of the *leading* (grandparent-chain) placeholders — `["orgId"]` for
the org children path, empty for single-parent paths.

This reproduces every current `FetchChildren` result and, as a bonus, removes the
ISSUE 0003 composites divergence: for `client → role` the composites path's
trailing placeholder is `role-name` (type role), not `client`, so it is correctly
excluded — `FetchChildren` no longer needs to route through `BuildDownwardGraph`
to avoid it.

### 2. `FetchChildren` routing + child shape (`admin`)

`FetchChildren(parent, childType, query)`:

1. Resolve `(path, inheritedFields)` from the scoped selector. If none, return an
   error (unchanged contract).
2. Render path params from `parent` and issue exactly one GET (unchanged: reuse
   the single-GET fetch + `FullRepresentation`, 404 → `FetchFailure{NotFound}`).
3. Tag each child: `Type = childType`, `Realm = parent.Realm`,
   `ParentType = parent.ParentType`, `Data = raw`, then **copy each
   `inheritedField` from `parent.Data` into the child** (the grandparent chain is
   constant down the subtree). For single-parent paths `inheritedFields` is empty
   and the existing immediate parent-reference injection (`clientUuid`, etc.)
   applies unchanged.

Key semantic decision: **`ParentType: "organization"` is a scope marker, not
"immediate parent is an organization".** Every group in an org subtree stays
tagged `ParentType: organization` and carries `{orgId, id}`. That lets the
selector pick the org path at every level and lets a returned grandchild be
passed straight back into `FetchChildren` — recursion to arbitrary depth,
terminating on an empty read (org groups expose no `subGroupCount`).

### 3. Recursion safety — no colliding `groupId`

A child is deliberately **not** given a `groupId` field set to its immediate
parent's id. On the org children path `{group-id}` is a collection placeholder
rendered from the parent's *own* identifier; a `groupId` in the child's data
would win via camel-case lookup and fetch the *parent's* children instead (the
ISSUE 0005 Gap 3 hazard, on a collection path where the single-resource
precedence fix does not apply). Children are identified by `(orgId, own id,
ParentType: organization)`; the immediate parent is caller-tracked, which is how
the consumer's per-parent recursion already works.

### 4. Write path (AC #3) — confirmation only

The resource-channel *write* selection has the same two-parent limitation.
Extending `Apply`'s write-path resolution is a materially larger change and is
**out of scope**. AC #3 is satisfied by a focused resolver/spec test that the org
children `POST` renders correctly with `{group-id}` bound to a child id (proving
nesting under a child is not special-cased to top-level org groups). Full
resource-channel creation of org grandchildren is left to the consumer's existing
relationship-template write; this issue unblocks the **read**.

## Components

- `pkg/catalog/dependencies.go` — new `ScopedChildCollection` helper (and any
  shared enumeration extracted from `BuildDownwardGraph`).
- `pkg/admin/fetch.go` — `FetchChildren` routes through the selector and copies
  inherited fields; no public signature change.
- No change to `manifest.Resource` (single `ParentType` retained; used as a scope
  marker for org groups).

## Testing

- **catalog unit**: `ScopedChildCollection` for all three parent shapes; that an
  org grandchild's `(group, organization)` shape still selects the org children
  path (arbitrary-depth selection); realm/client paths unchanged.
- **admin `FetchChildren`**: fake server serving org → group → child → grandchild;
  assert one GET per level, correct path per level, children tagged
  `ParentType: organization` + `orgId` propagated, recursion reaches the
  grandchild, and no `groupId` collision (a child's next fetch targets its own id).
- **AC #3**: resolver/spec test that org children `POST` renders under a child id.
- **No regression**: realm group nesting (0005), client roles (0003), one-level
  org groups (0002).
- **Live**: seed org → child → grandchild in a throwaway realm; verify
  `FetchChildren` recursion reaches the grandchild with **no** consumer override.

## Acceptance criteria mapping

- AC1 (grandchild reachable, tagged) → §2 + admin test + live.
- AC2 (arbitrary depth, caller-bounded) → §1–§3 recursion + unit test.
- AC3 (create nests under a child) → §4 confirmation test.
- AC4 (no regression) → regression tests.
- AC5 (`go vet` / `go test`) → CI gate.
- AC6 (tagged release) → post-merge `v1.4.0` (user-driven).

## Out of scope

- Full resource-channel *writes* of org grandchildren via `Apply` (confirmation
  only — §4).
- Org-group membership at depth, invitations, identity-provider links.
- Generalizing `manifest.Resource.ParentType` into a multi-level chain — this
  design uses the single string as an org-scope marker and does not need it.
