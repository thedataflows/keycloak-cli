---
type: Issue
title: organization apply does not persist attributes
description: Applying an organization resource whose data carries an attributes map does not write those attributes to Keycloak. The organization is created/updated but its attributes remain unset in the target realm, even though the OrganizationRepresentation model has an Attributes field and Keycloak's organization PUT accepts attributes.
tags: [issue, bug, organization, attributes, apply]
---

# ISSUE 0010: organization apply does not persist attributes

- **Type**: bug
- **Status**: done
- **Priority**: high
- **Labels**: [organization, attributes, apply, fetch, data-loss]
- **Assignee**: none
- **Related**: [ISSUE 0001](0001-full-representation-collection-fetch.md) (added the `--full-representation` flag this bug shows was insufficient for organizations)
- **Related code**: [`pkg/admin/fetch.go`](../../pkg/admin/fetch.go), [`pkg/admin/apply.go`](../../pkg/admin/apply.go)
- **Closing commits**: organization fetch forces full representation (`alwaysFullRepresentation`/`forceFullRepresentation`) + apply/fetch attribute regression tests. Released in `v1.6.2`.

## Resolution

The reported root cause — the **apply** path dropping organization attributes —
**does not reproduce**. The write path marshals the raw `resource.Data` map
directly (`pkg/admin/internal/client.go`, unchanged since `v1.6.1`), no typed
`OrganizationRepresentation` round-trip is involved, `sanitizeResourceData`
strips attributes only for `client`, and the OpenAPI schema declares
`attributes` so validation accepts it. Regression tests capturing the actual
HTTP body confirm attributes reach Keycloak on both create (`POST`) and update
(`PUT`).

The real drop is on the **read** side. Keycloak's organizations list returns a
brief representation that omits the `attributes` map unless
`briefRepresentation=false` is sent, and that was gated behind the opt-in
`--full-representation` flag (default off). So a `fetch organization` →
`apply` round-trip silently lost attributes even though `apply` itself is
correct — matching the observed "attributes end up null in the target" symptom.

**Fix:** organization collection fetches now always request the full
representation (`alwaysFullRepresentation`/`forceFullRepresentation` in
`pkg/admin/fetch.go`), independent of `--full-representation`. Organizations are
few and their attributes are pure data, so the cost is negligible; users and
groups stay flag-gated because their lists can be large.

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

- [x] Applying an organization with an `attributes` map persists those
      attributes to Keycloak (create and update paths). *(Already worked;
      verified by regression tests, not a code change.)*
- [x] Reading the organization back returns the applied attributes. *(Fixed:
      org fetch now always requests the full representation.)*
- [x] A test covers organization apply with attributes (create + update).
      *(`pkg/admin/apply_org_attributes_test.go`; read side covered by
      `pkg/admin/fetch_org_attributes_test.go`.)*

## Out of Scope

- User/group/role attribute handling (already working).
- Per-attribute merge semantics beyond full-representation apply.

## Notes

- Reported by the syncengine A->B Keycloak sync (thedataflows/iga-dash), where a
  Mapper transform_cel uppercases org attribute values; the transformed values
  reach `Apply` correctly but never land, isolating the drop to this library.
