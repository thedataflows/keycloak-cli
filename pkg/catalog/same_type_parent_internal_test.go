package catalog

import (
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSameTypeParentNestingsAreRecorded enumerates every POST collection path in
// the spec that carries a non-primary placeholder of the resource's own type — a
// "same-type parent nesting". ISSUE 0005 narrows ParentReferenceFieldNames so
// such a placeholder is stripped from the body; this test records the full set so
// the effect of that narrowing is visible and cannot silently grow. Of these,
// only the realm-group children path is reachable through the resource channel;
// the others are documented here with why they are unaffected.
func TestSameTypeParentNestingsAreRecorded(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)
	contracts, err := spec.ResourceContracts()
	require.NoError(t, err)
	placeholderMap, err := spec.PlaceholderToResourceType()
	require.NoError(t, err)
	for k, v := range fallbackPlaceholderToResourceType {
		if placeholderMap[k] == "" {
			placeholderMap[k] = v
		}
	}

	type nesting struct{ resourceType, path string }
	found := make(map[nesting]struct{})
	for resourceType, contract := range contracts {
		for _, op := range contract.AllOperations[http.MethodPost] {
			if !isCollectionEndpoint(op.Path) {
				continue
			}
			primary := primaryIdentifierParam(op.Path)
			skipPrimary := strings.HasSuffix(op.Path, "/{"+primary+"}")
			for _, ph := range extractPathParams(op.Path) {
				if ph == "realm" || (skipPrimary && ph == primary) {
					continue
				}
				if placeholderMap[ph] == resourceType {
					found[nesting{resourceType, op.Path}] = struct{}{}
				}
			}
		}
	}

	// The recorded set as of the 26.6.2 spec. A new entry here means the guard
	// narrowing now strips a body field on another endpoint — review it before
	// updating this expectation.
	expected := map[nesting]struct{}{
		// Reachable via the resource channel — this is the case ISSUE 0005 fixes.
		{"group", "/admin/realms/{realm}/groups/{group-id}/children"}: {},
		// Two-parent org children (ISSUE 0002 Gap 2, deferred). The org parent
		// resolves to /organizations/{org-id}/groups, never this children path,
		// so the resource channel never reaches it.
		{"group", "/admin/realms/{realm}/organizations/{org-id}/groups/{group-id}/children"}: {},
		// A client sub-action, not the canonical client create (POST /clients).
		// {client-uuid} is the client's own id and clientUuid is not a
		// ClientRepresentation field, so stripping-if-present is harmless.
		{"client", "/admin/realms/{realm}/clients/{client-uuid}/registration-access-token"}: {},
	}
	assert.Equal(t, expected, found)
}
