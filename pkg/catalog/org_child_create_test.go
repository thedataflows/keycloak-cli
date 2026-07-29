package catalog_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/catalog"
)

// TestOrgChildrenCreateNestsUnderChild is ISSUE 0006 AC #3: nesting under an org
// CHILD is not special-cased to top-level org groups. The create path is
// spec-valid with {group-id} bound to a child id and a GroupRepresentation body.
// The read is the blocker; this confirms the write is not a second one. Full
// resource-channel writes of grandchildren remain out of scope (the consumer
// writes via its relationship template).
func TestOrgChildrenCreateNestsUnderChild(t *testing.T) {
	spec := loadOrgGroupSpec(t)
	require.NoError(t, spec.ValidateOperationRequest(
		"/admin/realms/{realm}/organizations/{org-id}/groups/{group-id}/children",
		http.MethodPost,
		catalog.RequestValidation{
			PathParams: map[string]string{"realm": "demo", "org-id": "org-1", "group-id": "child-1"},
			Body:       map[string]interface{}{"name": "grandchild"},
		}))
}
