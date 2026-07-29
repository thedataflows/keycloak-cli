---
type: Issue
title: Ship organization-group-member and organization-group-child as built-in relationship kinds
description: The organization group membership and containment relationship kinds are not in the built-in registry, so every consumer has to register them as overrides. keycloak-cli already implements everything they need (org-scoped fetch, FetchChildren member/child, depth traversal); ship the two kinds built-in so consumers stop carrying an identical override shim.
tags: [issue, enhancement, catalog, groups, organizations]
timestamp: 2026-07-29T00:00:00Z
---

# ISSUE 0008: Ship organization-group-member and organization-group-child as built-in relationship kinds

- **Type**: feature
- **Status**: open
- **Priority**: low
- **Labels**: [catalog, groups, organizations, enhancement]
- **Assignee**: none
- **Related**: [ISSUE 0002](0002-org-scoped-groups.md) (established these kinds as the deferred org-group work), [ISSUE 0006](0006-org-scoped-group-deep-nesting.md) + [ISSUE 0007](0007-org-group-membership-deep-nesting.md) (the reads that made these kinds fully usable are now in the library), downstream consumer: iga-dash `syncengine/connectors/keycloak/relationship-overrides.yaml` (the shim to delete)
- **Related code**: [`pkg/catalog/relationship_registry.go`](../../pkg/catalog/relationship_registry.go)
- **Closing commits**: none

## Summary

The `relationship_registry.go` built-in set has no entry for an organization's group membership (`/organizations/{org-id}/groups/{group-id}/members`) or its group containment (`/organizations/{org-id}/groups/{group-id}/children`), so every consumer that syncs org groups must register those two kinds itself via `catalog.RelationshipOverridesFile`. The library now implements everything those kinds need — org-scoped resolution (ISSUE 0002), `FetchChildren` selection of the org children and members collections at any depth (ISSUE 0006/0007). The only thing missing is that the two kinds ship as *overrides in the consumer* rather than *built-ins in the library*. Ship them built-in and the consumer deletes an identical shim.

## Details

iga-dash (syncengine) carries this embedded override file today and applies it at connector construction. It is byte-for-byte the two kinds below — copy them into the built-in registry verbatim (same names, paths, methods, param types) and the consumer's shim becomes deletable with **no behaviour change**:

```yaml
- name: organization-group-member
  resourceA: group
  resourceB: user
  readPath: /organizations/{org-id}/groups/{group-id}/members
  writeTemplate: "{realm}/organizations/{org-id}/groups/{group-id}/members/{userId}"
  writeMethod: PUT
  deleteTemplate: "{realm}/organizations/{org-id}/groups/{group-id}/members/{userId}"
  deleteMethod: DELETE
  itemParamName: id
  paramTypes: {org-id: organization, group-id: group, userId: user}

- name: organization-group-child
  resourceA: group
  resourceB: group
  readPath: /organizations/{org-id}/groups/{group-id}/children
  writeTemplate: "{realm}/organizations/{org-id}/groups/{group-id}/children"
  writeMethod: POST
  # no itemParamName, no payloadField — the fetched Data must be the whole child item
  paramTypes: {org-id: organization, group-id: group}
```

Two properties the built-ins must preserve, because consumers depend on them:

- **`organization-group-child` has no `ItemParamName` and no `PayloadField`.** That is what makes `admin.relationshipPayload` hand back the whole fetched child item as the operation `Data` (a child group is an entity, not a bare edge). An item param or payload field would silently change the shape the consumer decodes.
- **`organization-group-member`'s `{org-id}` is a path *scope*, not an endpoint of the edge** (the edge is group↔user). Its `paramTypes` must keep `org-id: organization` so a consumer resolves the owning org separately from the edge endpoints (the same shape as a client-role mapping's `{client-id}`).

### Why this is safe / a no-op for consumers

A consumer references these kinds only by name through `catalog.DefaultRegistry().ByName("organization-group-member" | "organization-group-child")` — for the write templates, the param types, and (for the child) the edge-leg classification. It does not care whether the kind was registered by an override or shipped built-in, as long as the registered `RelationshipKind` is identical. So shipping them built-in lets the consumer delete its override file + registration with zero behaviour change.

### The one requirement beyond "add to the registry"

The built-ins must be present in `DefaultRegistry()` from package initialization (as the other built-in kinds are), not only after `admin.New`. Consumers' unit tests use fake clients and never call `admin.New`; today they register the override in a test `init()` to populate the registry. Once these are built-in, `DefaultRegistry()` must already contain them so those tests keep resolving the kinds with no registration step.

## Acceptance Criteria

- [ ] `organization-group-member` and `organization-group-child` are in the built-in relationship registry with exactly the definitions above (names, read/write/delete paths, methods, param types; child has no item param / payload field)
- [ ] Both resolve via `DefaultRegistry().ByName(...)` and `ByPath(...)`, and `manifest.RelationshipParamTypes(...)` returns the param types, from package init — no `admin.New` required (so fake-backed consumer tests resolve them)
- [ ] A guard test asserts each kind's read/write paths validate against the embedded spec (the pattern ISSUE 0005/0006 use)
- [ ] No regression to existing built-in kinds or to org-scoped fetch/apply
- [ ] `go vet ./...` clean, `go test ./...` green
- [ ] Tagged release so iga-dash can bump `syncengine/go.mod` (three vendor trees regenerated together) and then delete `relationship-overrides.yaml` + `installRelationshipOverrides`

## Out of Scope

- Any change to the org-scoped **read** paths (ISSUE 0006/0007, done) or to how members/children are fetched.
- The consumer-side deletion itself — that is iga-dash's follow-up once this releases.

## Notes

- This is purely a "promote the override to a built-in" change; the endpoints, methods, and param semantics are already proven in production by the consumer's override. Low priority — the override works fine; this just removes duplication.
- Keep the two kinds' names stable (`organization-group-member`, `organization-group-child`) — the consumer's edge-leg skip and edge-builder key on those exact strings.
