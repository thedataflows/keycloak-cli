package kcapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEdgesInferredFromSpec(t *testing.T) {
	c := testClient(t)
	edges := c.Edges()
	require.NotEmpty(t, edges)

	// realms -> users: /admin/realms/{realm}/users is a child collection,
	// but that path has only {realm}; the pair must come from paths that
	// carry a second placeholder, like /admin/realms/{realm}/users/{user-id}
	// or /admin/realms/{realm}/users/{user-id}/role-mappings/realm.
	byPair := map[string]Edge{}
	for _, e := range edges {
		byPair[e.Parent+"->"+e.Child] = e
	}
	assert.Contains(t, byPair, "realms->users")
	// Pairs beyond the realm star: the placeholder-to-placeholder legs that
	// Neighbors walks. /admin/realms/{realm}/users/{user-id}/groups/{groupId}
	// spans {user-id} -> {groupId}, and /clients/{client-uuid}/roles/{role-name}
	// spans {client-uuid} -> {role-name}. (The user side has no matching pair:
	// /users/{user-id}/role-mappings/clients exists only with a {client-id}
	// suffix, so its prefix is not a spec path.)
	assert.Contains(t, byPair, "users->groups")
	assert.Contains(t, byPair, "clients->roles")

	seen := map[string]struct{}{}
	for _, e := range edges {
		assert.NotEqual(t, e.Parent, e.Child, "self-edges are noise: %v", e)
		assert.Contains(t, e.Path, "{realm}")
		key := e.Parent + "->" + e.Child + " " + e.Path
		assert.NotContains(t, seen, key, "duplicate edge: %v", e)
		seen[key] = struct{}{}

		assert.True(t, e.Cardinality == "one" || e.Cardinality == "many",
			"cardinality must be one|many, got %q on %v", e.Cardinality, e)
	}

	for i := 1; i < len(edges); i++ {
		prev := edges[i-1].Parent + edges[i-1].Child + edges[i-1].Path
		curr := edges[i].Parent + edges[i].Child + edges[i].Path
		assert.LessOrEqual(t, prev, curr, "edges must be sorted by (Parent, Child, Path)")
	}

	// The prefix-up-to-the-second-placeholder rule: /admin/realms/{realm}/roles-by-id
	// is not a path in this spec (roles are listed at /roles), so no edge may be
	// anchored on it.
	for _, e := range edges {
		assert.NotContains(t, e.Path, "/roles-by-id/",
			"roles-by-id has no parent-collection path in the spec; edge %v must not exist", e)
	}
}

func TestEdgeCardinalityAndCascade(t *testing.T) {
	c := testClient(t)

	// One path can carry several edges (one per placeholder pair), so the
	// lookup key is the (Parent, Child, Path) triple.
	byKey := map[string]Edge{}
	for _, e := range c.Edges() {
		byKey[e.Parent+"->"+e.Child+" "+e.Path] = e
	}

	// Trailing placeholder -> single item.
	e, ok := byKey["realms->users /admin/realms/{realm}/users/{user-id}"]
	require.True(t, ok, "expected an edge on the user item path")
	assert.Equal(t, "realms", e.Parent)
	assert.Equal(t, "users", e.Child)
	assert.Equal(t, "one", e.Cardinality)
	// users have a standalone top-level DELETE (/admin/realms/{realm}/users/{user-id}).
	assert.False(t, e.CascadeOwner)

	// Trailing collection -> many.
	e, ok = byKey["realms->users /admin/realms/{realm}/users/{user-id}/groups"]
	require.True(t, ok, "expected an edge on the user groups collection path")
	assert.Equal(t, "many", e.Cardinality)

	// A client's default scopes are removed only through the client: no
	// top-level DELETE /admin/realms/{realm}/default-client-scopes/{id} exists.
	e, ok = byKey["clients->default-client-scopes /admin/realms/{realm}/clients/{client-uuid}/default-client-scopes/{clientScopeId}"]
	require.True(t, ok, "expected an edge on the client default-client-scopes item path")
	assert.Equal(t, "clients", e.Parent)
	assert.Equal(t, "default-client-scopes", e.Child)
	assert.Equal(t, "one", e.Cardinality)
	assert.True(t, e.CascadeOwner, "child %q has no standalone top-level DELETE", e.Child)
}

func TestEdgeDisplayOverride(t *testing.T) {
	c := testClient(t)
	for _, e := range c.Edges() {
		if e.Parent == "users" && e.Child == "groups" {
			// The registry's curated name for the membership edge flows into Display.
			assert.Equal(t, "user-group-membership", e.Display,
				"registry override should name this common edge")
			return
		}
	}
	t.Fatalf("no users->groups edge inferred from the committed spec; the display override has no anchor")
}
