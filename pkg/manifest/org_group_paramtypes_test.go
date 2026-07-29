package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestOrgGroupKindParamTypesFromDefault pins ISSUE 0008 AC#2: the built-in
// org-group kinds' param types resolve through manifest.RelationshipParamTypes
// from the package default (no admin.New / registry install), so fake-backed
// consumer tests resolve them.
func TestOrgGroupKindParamTypesFromDefault(t *testing.T) {
	assert.Equal(t, map[string]string{"org-id": "organization", "group-id": "group", "userId": "user"},
		RelationshipParamTypes("organization-group-member"))
	assert.Equal(t, map[string]string{"org-id": "organization", "group-id": "group"},
		RelationshipParamTypes("organization-group-child"))
}
