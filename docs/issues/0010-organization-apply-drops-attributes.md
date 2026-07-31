---
type: Issue
title: organization apply does not persist attributes
description: Applying an organization resource whose data carries an attributes map does not write those attributes to Keycloak. The organization is created/updated but its attributes remain unset in the target realm, even though the OrganizationRepresentation model has an Attributes field and Keycloak's organization PUT accepts attributes.
tags: [issue, bug, organization, attributes, apply]
---

# ISSUE 0010: organization apply does not persist attributes

- **Type**: bug
- **Status**: open
- **Priority**: high
- **Labels**: [organization, attributes, apply, data-loss]
- **Assignee**: none
- **Related**: none
- **Related code**: [`pkg/admin/apply.go`](../../pkg/admin/apply.go), [`pkg/models/models.gen.go`](../../pkg/models/models.gen.go)

## Summary

`Apply` on an organization resource whose `Data` includes an `attributes` map
(`map[string][]string`) creates/updates the organization but does not persist
its attributes. The attributes are silently dropped: the target organization
ends up with `attributes: null`.

## Details

Observed downstream in the syncengine Keycloak connector (v1.6.1). The connector
builds an organization apply whose data map contains attributes and calls
`svc.Apply(..., manifest.Resource{Type: <organization>, Data: orgMap})`. The
organization syncs (name, description, domains, enabled all land correctly) but
its attributes never appear in the target realm.

Evidence that the input is correct and the drop is inside the library:

1. The caller's decoded apply payload carried the attributes, e.g. (values
   base64-decoded from the wire form):
   ```json
   {"department":["ENGINEERING"],"costcenter":["OPS-42"],"region":["EMEA-WEST"],"tier":["GOLD"]}
   ```
   These were present in `Resource.Data["attributes"]` at the `Apply` call.
2. Keycloak itself accepts organization attributes: a direct
   `PUT /admin/realms/{realm}/organizations/{id}` with an `attributes` object
   persists them (verified against the same server/version).
3. The generated model `OrganizationRepresentation` in
   `pkg/models/models.gen.go` already has an `Attributes *map[string][]string`
   field with a JSON tag, so the type can carry them.

Given (1)-(3), the organization create/update write path in the library is not
including `attributes` in the request body it sends to Keycloak (or is sending
a representation that omits them). Other resource types (users, groups) carry
attributes; organizations appear to have been implemented for
name/description/domains/enabled only, so attributes were never exercised.

Note: this report was drafted from a terminal that redacts the literal tokens
"organization"/"attributes" in source excerpts, so the exact function/line in
the apply path is not cited; the maintainer's unredacted view will show the
organization create/update body construction in `pkg/admin/apply.go` (and the
organization resource definition in `pkg/catalog`).

## Reproduction

1. Apply an organization resource with a non-empty `attributes` map to a realm
   that has organizations enabled.
2. Read the organization back (list or get) from the target realm.
3. Observe `attributes` is null/absent, while name/description/domains persisted.

## Acceptance Criteria

- [ ] Applying an organization with an `attributes` map persists those
      attributes to Keycloak (create and update paths).
- [ ] Reading the organization back returns the applied attributes.
- [ ] A test covers organization apply with attributes (create + update).

## Out of Scope

- User/group/role attribute handling (already working).
- Per-attribute merge semantics beyond full-representation apply.

## Notes

- Reported by the syncengine A->B Keycloak sync (thedataflows/iga-dash), where a
  Mapper transform_cel uppercases org attribute values; the transformed values
  reach `Apply` correctly but never land, isolating the drop to this library.
