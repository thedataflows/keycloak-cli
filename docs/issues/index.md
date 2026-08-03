# Issues

Issues track units of work: features, bugs, tasks, and chores. They
describe *what* needs to happen; ADRs describe *why* things are the way
they are. When an issue and an ADR disagree, the ADR wins.

## Open

*No open issues.*

## In Progress

*No issues in progress.*

## Done

* [ISSUE 0001](0001-full-representation-collection-fetch.md) - Full-representation collection fetch (`briefRepresentation=false`) so `attributes` are returned
* [ISSUE 0002](0002-org-scoped-groups.md) - Org-scoped groups as a parent-scoped resource; subgroup hierarchy and membership deferred
* [ISSUE 0003](0003-parent-scoped-collection-fetch-public-api.md) - Parent-scoped collection fetch on the public `admin.Service` (one child collection of one parent, without the `Depth` fan-out)
* [ISSUE 0004](0004-client-role-mapping-kinds-never-match-spec-path.md) - Client-scoped role-mapping kinds never match the spec path (`{client}` vs `{client-id}`), so their edges are silently never fetched
* [ISSUE 0005](0005-same-type-parent-binding.md) - Same-type parent binding on the resource channel: a group nested under a group can't be created via the resource channel because the parent id is not stripped from the body
* [ISSUE 0006](0006-org-scoped-group-deep-nesting.md) - Organization-scoped group hierarchies reachable at arbitrary depth via `FetchChildren` (scoped child-collection selection), no consumer override needed
* [ISSUE 0007](0007-org-group-membership-deep-nesting.md) - Organization-group membership readable at any depth via the scoped members read (`FetchChildren`/`fetch --depth`), keyed to each group, read-only on apply
* [ISSUE 0008](0008-ship-org-group-relationship-kinds-builtin.md) - `organization-group-member` / `organization-group-child` shipped as built-in relationship kinds so consumers can delete their override shim
* [ISSUE 0009](0009-realm-cascade-stops-above-grandchildren.md) - Realm-rooted depth traversal reaches grandchildren (identity-provider mappers, client roles, protocol mappers, nested groups); realm-children no longer misclassified as org-scoped parents
* [ISSUE 0010](0010-organization-apply-drops-attributes.md) - organization attributes dropped on a fetch→apply round-trip: apply was already correct; the drop was on the read side (org list returned brief representation), fixed by always requesting full representation for organizations
* [ISSUE 0011](0011-apply-does-not-clear-a-collection-field.md) - apply does not clear a collection field: an explicit-empty collection (attributes {}) meant to remove it is not cleared on the target, though a direct full-representation PUT clears it; write-side removal case, distinct from 0010's persist case, observed for groups

## Conventions

* [Issue Conventions](issues-conventions.md) - Numbering, header fields, body sections, and lifecycle rules
* [Issue Template](TEMPLATE.md) - Copy-and-fill template for new issues
