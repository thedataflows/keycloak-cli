package kcapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperationsListsAllSpecOps(t *testing.T) {
	c := testClient(t) // helper below loads the real repo spec
	ops, err := c.Operations(OpFilter{})
	require.NoError(t, err)
	assert.Greater(t, len(ops), 100, "real Keycloak spec has hundreds of operations")

	users := filterOps(ops, func(o Operation) bool { return o.Resource == "users" })
	assert.NotEmpty(t, users)
}

func TestOperationsFilterByVerbAndSearch(t *testing.T) {
	c := testClient(t)
	ops, err := c.Operations(OpFilter{Resource: "users", Method: Get})
	require.NoError(t, err)
	for _, o := range ops {
		assert.Equal(t, Get, o.Verb)
	}
	bySearch, err := c.Operations(OpFilter{Search: "users/{user-id}"})
	require.NoError(t, err)
	// The vendored 26.6.2 spec has no operationId fields at all, so search
	// runs against Path/Summary; "users/{user-id}" pins the single-user path.
	found := false
	for _, o := range bySearch {
		if o.Verb == Get && o.Path == "/admin/realms/{realm}/users/{user-id}" {
			found = true
			assert.Equal(t, "users", o.Resource)
		}
	}
	assert.True(t, found, "search by path substring should find the single-user GET")
}
