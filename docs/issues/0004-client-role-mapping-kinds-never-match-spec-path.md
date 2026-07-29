---
type: Issue
title: Client-scoped role-mapping kinds never match the spec path ({client} vs {client-id})
description: Four relationship kinds are registered with a {client} placeholder while the Keycloak 26.6.2 spec declares {client-id}. Registry.ByPath is an exact lookup, so those endpoints are never classified as relationships and their edges are never fetched — silently, with no failure reported.
tags: [issue, bug, catalog, relationships, roles, clients]
timestamp: 2026-07-28T00:00:00Z
---

# ISSUE 0004: Client-scoped role-mapping kinds never match the spec path ({client} vs {client-id})

- **Type**: bug
- **Status**: open
- **Priority**: high
- **Labels**: [catalog, relationships, bug]
- **Assignee**: none
- **Related**: [ISSUE 0003](0003-parent-scoped-collection-fetch-public-api.md) (same consumer, different gap), downstream consumer: iga-dash `syncengine/docs/issues/0022-client-roles-never-synced-to-target-realm.md`
- **Related code**: [`pkg/catalog/relationship_registry.go`](../../pkg/catalog/relationship_registry.go), [`pkg/catalog/relationship_patterns.go`](../../pkg/catalog/relationship_patterns.go)
- **Closing commits**: none

## Summary

`user-client-role-mapping`, `group-client-role-mapping`, `client-scope-client-role-mapping` and `client-client-scope-mapping` are registered with a `{client}` path placeholder, but the Keycloak 26.6.2 spec declares those segments as `{client-id}`. `Registry.ByPath` matches by exact string after prefix normalisation, so `classifyRelationshipEndpoint` returns nil for those spec paths, `DiscoverRelationshipPatterns` never yields a pattern for them, and a `--relationships` fetch silently returns no edges of those kinds — no failure, no warning, no `FetchFailure`.

A user or group that has client roles assigned therefore looks, to any library consumer, exactly like one that has none.

## Details

The registry read paths (`pkg/catalog/relationship_registry.go`):

```
/users/{user-id}/role-mappings/clients/{client}                   user-client-role-mapping        (:339)
/groups/{group-id}/role-mappings/clients/{client}                 group-client-role-mapping       (:342)
/client-scopes/{client-scope-id}/scope-mappings/clients/{client}  client-scope-client-role-mapping(:346)
/clients/{client-uuid}/scope-mappings/clients/{client}            client-client-scope-mapping     (:348)
```

The spec paths (`keycloak-oapi/26.6.2.spec.json`):

```
/admin/realms/{realm}/users/{user-id}/role-mappings/clients/{client-id}
/admin/realms/{realm}/groups/{group-id}/role-mappings/clients/{client-id}
```

`normalizeReadPath` only strips the `/admin/realms/` prefix and the `{realm}/` segment (`relationship_registry.go:124-127` → `normalizeRelationshipPath`, `relationships.go:293-298`); it does not rewrite placeholder names. So the lookup key is `users/{user-id}/role-mappings/clients/{client-id}` against a map holding `users/{user-id}/role-mappings/clients/{client}` — a miss.

Chain of consequences: `Registry.ByPath` miss → `classifyRelationshipEndpoint` returns nil (`relationship_patterns.go:14-17`) → the path is dropped by `DiscoverRelationshipPatterns` → `fetchRelationshipsForRealm` never iterates it → zero edges, and because nothing errored, zero failures.

### Verified live

Against Keycloak **26.6.3** with the `26.6.2` spec, realm `sync-source`, which has exactly two client-role mappings (`demo-a-user-1` → `demo-a-clientrole-1` on client `demo-a-client-1`, and the same role on group `demo-a-group-1`, both confirmed through the raw admin API):

```bash
keycloak-cli fetch user,role,client,group --realm sync-source --relationships -f json -o rels.json
jq -r '[.relationships[].kind] | unique' rels.json
```

123 relationships are returned. The kind breakdown:

```
55 client-default-scope        5 user-realm-role-mapping     1 group-realm-role-mapping
45 client-optional-scope       5 realm-optional-client-scope 1 organization-member
 8 realm-default-client-scope  1 user-group-membership       1 role-composite-mapping
                                                             1 client-scope-realm-role-mapping
```

Not one edge of the four `{client}`-keyed kinds. The correlation is exact: **every** registry path whose placeholders match the spec appears in the results (`user-group-membership` → `/users/{user-id}/groups`, `client-default-scope` → `/clients/{client-uuid}/default-client-scopes`, …), and **only** the four `{client}` paths are missing.

Note the traversal machinery itself is not the problem: `fetchRelationshipsForPattern` already recurses over multiple parent path params (`fetch_relationships.go:175-193`), and `paramNameByResourceType` maps `client` → `{"client-uuid", "client"}`. Once the path matches, a (user × client) iteration is what these kinds need — this is not the deferred two-parent-param work from [ISSUE 0002](0002-org-scoped-groups.md) Gap 2, because both params here are ordinary independent parents rather than a nested `(org-id, group-id)` pair.

### Why the silence is the worse half of the bug

An unclassified relationship path is indistinguishable from "this realm has no such edges". The iga-dash consumer spent a full investigation cycle attributing the missing edges to its own ref-resolution filter (which reports a `skipped_unresolved_endpoint` counter) before the mismatch was found here. A registry entry whose path does not exist in the loaded spec is a startup-time detectable defect.

## Acceptance Criteria

- [ ] The four `{client}` registry paths match the loaded spec, so `Registry.ByPath` resolves them (rename to `{client-id}`, or normalise client placeholder aliases in `normalizeReadPath` — the alias route also covers a future spec that renames the segment again)
- [ ] A `--relationships` fetch of a realm where a user and a group each hold a client role returns `user-client-role-mapping` and `group-client-role-mapping` edges, with the owning client in `PathParams` and the role in the bulk payload
- [ ] `client-scope-client-role-mapping` and `client-client-scope-mapping` resolve too, or are explicitly scoped out with a reason
- [ ] Apply round-trips the recovered edges: link and unlink against `/users/{user-id}/role-mappings/clients/{client-id}` (bulk POST / DELETE)
- [ ] A test asserts every registry `ReadPath` resolves against the embedded spec, so a placeholder that exists in no spec path fails the suite instead of silently disabling a kind — this is the criterion that prevents the next instance of this bug
- [ ] Realm-role mappings and the currently-working kinds are unchanged (no regression in the 123-edge baseline above, plus the newly recovered edges)
- [ ] `go vet ./...` clean, `go test ./...` green
- [ ] Tagged release so iga-dash can bump. Reminder: iga-dash vendors three trees (`igadash/vendor`, `syncengine/vendor`, root `go work vendor`) and all three must be regenerated together

## Out of Scope

- The parent-scoped collection fetch requested in [ISSUE 0003](0003-parent-scoped-collection-fetch-public-api.md) — independent gap, same consumer.
- Org-scoped group hierarchy/membership (ISSUE 0002 Gap 2/3) — genuinely needs the two-*nested*-parent model, unlike these kinds.
- Composite role membership semantics (`role-composite-mapping` already works).

## Notes

- Suggested fix shape, if the alias route is taken: teach `normalizeReadPath` a placeholder-alias table (`client-id` → `client`, applied to both the registry key and the incoming spec path) rather than editing four string literals. The registry then keeps one canonical placeholder name per resource while tolerating spec renames — which is what bit here.
- Worth checking the same class of mismatch across the whole registry as part of the fix: the acceptance criterion above (assert every `ReadPath` resolves against the spec) will surface any other silently-dead kind in one run.
