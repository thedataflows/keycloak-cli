package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thedataflows/keycloak-cli/pkg/kcapi"
)

func TestGraphFilterFromFlags(t *testing.T) {
	assert.Equal(t, kcapi.EdgeFilter{Child: "roles"}, graphFilter("", "roles"))
	assert.Equal(t, kcapi.EdgeFilter{Parent: "realms"}, graphFilter("realms", ""))
	assert.Equal(t, kcapi.EdgeFilter{}, graphFilter("", ""))
}

func TestGraphRefFromArgs(t *testing.T) {
	ref := graphRef("users", "alice", "demo", "")
	assert.Equal(t, "users", ref.Type)
	assert.Equal(t, "alice", ref.Name)
	assert.Equal(t, "demo", ref.Realm)
	assert.Equal(t, "", ref.ID)
}

func TestFilterGraphEdges(t *testing.T) {
	edges := []kcapi.Edge{
		{Parent: "users", Child: "groups", Path: "/a"},
		{Parent: "realms", Child: "roles", Path: "/b"},
		{Parent: "users", Child: "roles", Path: "/c"},
	}
	assert.Equal(t, edges, filterGraphEdges(edges, graphFilter("", "")), "empty filter keeps every edge")
	assert.Equal(t, []kcapi.Edge{edges[2]}, filterGraphEdges(edges, graphFilter("users", "roles")), "set fields AND")
	assert.Equal(t, []kcapi.Edge{edges[0], edges[2]}, filterGraphEdges(edges, graphFilter("users", "")))
	assert.Empty(t, filterGraphEdges(nil, graphFilter("users", "")), "nil stays nil, never panics")
}
