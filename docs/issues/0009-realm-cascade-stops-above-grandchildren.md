---
type: Issue
title: Realm-rooted depth traversal stops one level below the realm, so grandchildren (identity-provider mappers, client roles, protocol mappers) are never reached
description: The depth traversal classifies any resource whose ParentType differs from its own Type as an organization-scoped parent. Resources fetched as children of the realm carry ParentType="realm", so at the next level they are misrouted through ScopedChildCollection, which strips the realm placeholder from the chain and finds no match. Descent halts one level below the realm. `fetch realm --depth N` therefore returns identity providers but not their mappers, clients but not their roles or protocol mappers, groups but not their subgroups. Directly naming the resource type (`fetch identityprovider`, `fetch client`) is unaffected because those seeds carry an empty ParentType.
tags: [issue, bug, admin, fetch, depth, identity-providers, clients]
timestamp: 2026-07-30T00:00:00Z
---

# ISSUE 0009: Realm-rooted depth traversal stops above grandchildren

- **Type**: bug
- **Status**: done
- **Priority**: medium
- **Labels**: [admin, fetch, depth, identity-providers, clients]
- **Assignee**: none
- **Related**: [ISSUE 0006](0006-org-scoped-group-deep-nesting.md) (introduced the org-scoped descent branch this bug over-triggers), [ISSUE 0007](0007-org-group-membership-deep-nesting.md)
- **Related code**: [`pkg/admin/fetch.go`](../../pkg/admin/fetch.go)
- **Closing commits**: `6573386` (isScopedParent guard + realm-cascade regression tests). Verified live on `demo123`.

## Summary

`fetch realm --depth N` descends exactly one level below the realm and then
stops. Identity providers are returned but **their mappers are not**; clients are
returned but **their roles and protocol mappers are not**; groups are returned but
**their subgroups are not**. The traversal misclassifies every realm-child as an
organization-scoped parent and routes it through a selector that cannot match
realm-rooted paths. Directly naming a resource type (`fetch identityprovider`,
`fetch client`) already works — those seeds carry an empty `ParentType` — so the
per-type export path is correct and only the realm-rooted cascade is broken.

## Details

Verified 2026-07-30 against the live Keycloak at `KEYCLOAK_BASE_URL` (realm
`demo123`, keycloak-cli built from `main` at `c16e186`).

The downward containment graph (spec-derived) is correct and already contains
every needed edge:

```
realm            -> identityprovider, client, group, role, clientscope,
                    organization, component, authenticationflow, user
identityprovider -> identityprovidermapper   (/identity-provider/instances/{alias}/mappers)
client           -> role, protocolmapper      (/clients/{client-uuid}/roles, .../protocol-mappers/models)
```

So the descent *should* reach every grandchild. Observed instead:

| Command | Result |
| --- | --- |
| `fetch identityprovider -r demo123` (Depth 1) | identityprovider **+ identityprovidermapper** ✅ |
| `fetch client -r demo123 --depth 2` | 7 client, **28 role, 9 protocolmapper** ✅ |
| `fetch realm -r demo123 --depth 3` | identityprovider ✅ but **0 identityprovidermapper**, client ✅ but **0 client-role, 0 protocolmapper** ❌ |

### Root cause — the scoped heuristic over-fires on `ParentType="realm"`

`fetchDepthLevels` (`pkg/admin/fetch.go`) decides how to descend from each
frontier parent with:

```go
scoped := parent.ParentType != "" && parent.ParentType != parent.Type
```

- A **seed** fetched by directly naming its type (`fetchRealmScopedResources`)
  carries `ParentType == ""` → `scoped == false` → the structural branch
  (`fetchNestedResourceCollection`) fires and resolves the child path from the
  parent's own data. This is why `fetch identityprovider --depth 1` reaches
  mappers and `fetch client --depth 2` reaches roles.
- A **realm-child** discovered at level 0 of a realm-rooted cascade is tagged
  `ParentType = parent.Type = "realm"` by `fetchNestedResourceCollection`
  (`pkg/admin/fetch.go:632`). At level 1, `"realm" != "" && "realm" != "identityprovider"`
  → `scoped == true` → it is routed through `ScopedChildCollection(parent.Type,
  "realm", childType)`.

`ScopedChildCollection` treats a non-empty, non-self `parentParentType` as an
org-style scope and requires the child's GET path to carry a **two-deep**
placeholder chain whose second-to-last element is the grandparent. But
`parentPlaceholderChain` **skips the `realm` placeholder**
(`pkg/catalog/dependencies.go`), so the mapper path
`/identity-provider/instances/{alias}/mappers` yields a one-element chain
`[{alias→identityprovider}]`. The `scoped` arm needs `len(chain) >= 2` with the
grandparent at `chain[-2]`, finds none, returns `ok == false`, and the child is
never fetched. Descent halts.

`realm` is the *only* placeholder `parentPlaceholderChain` strips, so `realm` is
the unique `ParentType` value that breaks the heuristic. Organization-scoped
descent (`ParentType == "organization"`) is unaffected: `org-id` is a real
placeholder that stays in the chain, so its two-deep match still succeeds.

### Why the per-type path is already correct

When a mapper is reached via the structural branch it is tagged
`ParentType = "identityprovider"` and the parent `alias` is injected into its
data (`ParentReferenceFields`), so it round-trips on apply. The fix makes the
realm-rooted cascade produce **the same** mapper/role resources the direct
`fetch identityprovider` / `fetch client` commands already produce today — it
adds no new resource shape, only restores the missing descent.

## Proposed fix

Exclude `realm` from the scoped classification so a realm-child descends through
the structural branch (identical to how a directly-named seed descends).
Centralize the rule so the two call sites cannot drift:

```go
// isScopedParent reports whether parent must descend through the scoped
// child-collection selector (org-style grandparent chain) rather than the
// structural downward path. The realm is never a scoping grandparent: the realm
// placeholder is stripped from every path's placeholder chain, so a realm-child
// (identity provider, client, group, ...) must descend structurally, exactly as
// a directly-named seed does. Without this guard the cascade halts one level
// below the realm (ISSUE 0009).
func isScopedParent(parent manifest.Resource) bool {
	return parent.ParentType != "" &&
		parent.ParentType != parent.Type &&
		parent.ParentType != "realm"
}
```

Replace the inline `scoped :=` at `pkg/admin/fetch.go:527` and the equivalent
`parent.ParentType != "" && parent.ParentType != parent.Type` guard in
`FetchChildren` (`pkg/admin/fetch.go:246`) with `isScopedParent(parent)`.

The change is additive for the realm cascade and inert for the org cascade and
the per-type path (their `ParentType` is never the literal `realm`).

## Acceptance Criteria

- [x] `fetch realm -r <realm> --depth 3` returns `identityprovidermapper`
      resources for every identity provider that has mappers, each tagged
      `ParentType="identityprovider"` with the parent `alias` present in its data.
      *Verified live on `demo123`: 0 → 1 mapper (`alias=idp-1` injected).*
- [x] The same cascade returns client `role` and `protocolmapper` resources for
      clients that have them (parity with `fetch client --depth 2`).
      *Verified live: 0 → 28 client roles, 0 → 43 protocol mappers.*
- [x] The same cascade returns nested realm `group` subgroups
      (parity with `fetch group --depth 2`) — same structural descent path.
- [x] Directly-named fetches (`fetch identityprovider`, `fetch client`) and the
      organization-scoped depth traversal are unchanged — full `pkg/admin` suite
      and the org depth/children tests still pass.
- [x] A new `admin_test` httptest-fake regression test pins realm → identity
      provider → mapper (and realm → client → role) descent, mirroring
      `fetch_depth_org_test.go`: the fake serves only the realm-rooted paths, so
      the test fails if the traversal reverts to the scoped (org) path.
      *Both tests confirmed to fail with the `!= "realm"` guard removed.*

## Implementation

- `pkg/admin/fetch.go`: added `isScopedParent(parent)` — the scoped test plus a
  `parent.ParentType != "realm"` clause — and routed both the depth loop
  (`fetchDepthLevels`) and `FetchChildren` through it, replacing the two inline
  `parent.ParentType != "" && parent.ParentType != parent.Type` checks.
- `pkg/admin/fetch_depth_realm_cascade_test.go`: two TDD regression tests
  (identity-provider → mapper, client → role), both watched fail before the fix.

## Out of Scope

- Adding `identityprovider` (or other types) to the default resource set
  `realm,user,client,group,role`. This issue only fixes descent depth, not
  default breadth; a full-realm export still requires either naming
  `identityprovider` or rooting at `realm`.
- Any change to `ScopedChildCollection` / `parentPlaceholderChain` in
  `pkg/catalog`. The misclassification is an admin-layer routing decision; the
  catalog selector is correct as written.
- Applying/writing identity-provider mappers. Read/round-trip already works via
  the injected `alias`; no apply gap was observed.

## Notes

- Manual verification recipe against a realm with an IdP:
  ```bash
  # seed a mapper
  curl -sX POST -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
    -d '{"name":"probe","identityProviderAlias":"idp-1","identityProviderMapper":"oidc-hardcoded-role-idp-mapper","config":{"role":"offline_access","syncMode":"INHERIT"}}' \
    "$BASE/admin/realms/<realm>/identity-provider/instances/idp-1/mappers"
  # before fix: 0 mappers; after fix: 1
  keycloak-cli fetch realm -r <realm> --depth 3 -f json | \
    jq '[.resources[]|select(.type=="identityprovidermapper")]|length'
  ```
- Edge case to keep in a test: a realm-child whose own `Type == ParentType`
  cannot occur for realm (`realm` is not itself a realm-child), so the
  `ParentType != parent.Type` clause and the new `!= "realm"` clause are
  independent; both must remain.
- Trace of the group case confirms the same bug and the same fix repairs it:
  realm → group (ParentType `realm`, currently scoped → fails) → subgroup. After
  the fix the subgroup descends structurally and, being `ParentType == Type ==
  "group"`, its own children continue structurally — matching direct
  `fetch group --depth N`.
