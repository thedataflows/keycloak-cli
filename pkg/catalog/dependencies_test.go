package catalog

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func childTypes(children []DownwardChild) []string {
	types := make([]string, len(children))
	for i, c := range children {
		types[i] = c.ChildType
	}
	return types
}

func TestBuildDependencyGraph(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	graph, err := spec.BuildDependencyGraph()
	require.NoError(t, err)

	// Verify key dependencies exist
	assert.Contains(t, graph["protocolmapper"], "clientscope")
	assert.Contains(t, graph["identityprovidermapper"], "identityprovider")
	assert.Contains(t, graph["group"], "realm")
	assert.Contains(t, graph["client"], "realm")
	assert.Contains(t, graph["user"], "realm")
	assert.Contains(t, graph["role"], "realm")
	assert.Contains(t, graph["clientscope"], "realm")
	assert.Contains(t, graph["identityprovider"], "realm")
	assert.Contains(t, graph["organization"], "realm")
	assert.Contains(t, graph["authenticationflow"], "realm")
}

func TestTopologicalSortPriorityMap(t *testing.T) {
	graph := map[string][]string{
		"realm":          {},
		"group":          {"realm"},
		"client":         {"realm"},
		"user":           {"realm"},
		"role":           {"realm"},
		"clientscope":    {"realm"},
		"protocolmapper": {"clientscope", "client"},
	}

	priorityMap, err := topologicalSortPriorityMap(graph)
	require.NoError(t, err)

	// Verify realm comes before everything else
	assert.Greater(t, priorityMap["group"], priorityMap["realm"])
	assert.Greater(t, priorityMap["client"], priorityMap["realm"])
	assert.Greater(t, priorityMap["protocolmapper"], priorityMap["clientscope"])
	assert.Greater(t, priorityMap["protocolmapper"], priorityMap["client"])
}

func TestInvertDependencyGraph(t *testing.T) {
	graph := map[string][]string{
		"realm":          {},
		"client":         {"realm"},
		"clientscope":    {"realm"},
		"protocolmapper": {"clientscope", "client"},
	}

	downward := InvertDependencyGraph(graph)

	assert.ElementsMatch(t, downward["realm"], []string{"client", "clientscope"})
	assert.ElementsMatch(t, downward["client"], []string{"protocolmapper"})
	assert.ElementsMatch(t, downward["clientscope"], []string{"protocolmapper"})
	assert.Empty(t, downward["protocolmapper"])
}

func TestDownwardLevels(t *testing.T) {
	graph := map[string][]DownwardChild{
		"realm":          {{ChildType: "client"}, {ChildType: "clientscope"}},
		"client":         {{ChildType: "protocolmapper"}},
		"clientscope":    {{ChildType: "protocolmapper"}},
		"protocolmapper": {},
	}

	levels := DownwardLevels(graph, []string{"realm"}, 2)
	require.Len(t, levels, 2)
	assert.ElementsMatch(t, []string{"client", "clientscope"}, childTypes(levels[0]))
	assert.ElementsMatch(t, []string{"protocolmapper"}, childTypes(levels[1]))
}

func TestDownwardLevelsRespectsMaxDepth(t *testing.T) {
	graph := map[string][]DownwardChild{
		"a": {{ChildType: "b"}},
		"b": {{ChildType: "c"}},
		"c": {{ChildType: "d"}},
	}

	levels := DownwardLevels(graph, []string{"a"}, 2)
	require.Len(t, levels, 2)
	assert.Equal(t, []string{"b"}, childTypes(levels[0]))
	assert.Equal(t, []string{"c"}, childTypes(levels[1]))
}

func TestDownwardLevelsBreaksCycles(t *testing.T) {
	graph := map[string][]DownwardChild{
		"a": {{ChildType: "b"}},
		"b": {{ChildType: "c"}},
		"c": {{ChildType: "b"}},
	}

	levels := DownwardLevels(graph, []string{"a"}, 5)
	require.Len(t, levels, 2)
	assert.Equal(t, []string{"b"}, childTypes(levels[0]))
	assert.Equal(t, []string{"c"}, childTypes(levels[1]))
}

func TestBuildDownwardGraphFromSpec(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	downward, err := spec.BuildDownwardGraph()
	require.NoError(t, err)

	assert.Contains(t, childTypes(downward["clientscope"]), "protocolmapper")
	assert.Contains(t, childTypes(downward["client"]), "protocolmapper")
	assert.Contains(t, childTypes(downward["realm"]), "user")
	assert.Contains(t, childTypes(downward["realm"]), "client")
}

func TestTopologicalSortDetectsCycle(t *testing.T) {
	graph := map[string][]string{
		"a": {"b"},
		"b": {"a"},
	}
	_, err := topologicalSortPriorityMap(graph)
	assert.Error(t, err)
}

func TestDependencyPriorityMapFromSpec(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	priorityMap, err := spec.DependencyPriorityMap()
	require.NoError(t, err)
	require.NotNil(t, priorityMap)

	// Verify realm has the lowest priority (created first)
	for rt, p := range priorityMap {
		if rt == "realm" {
			continue
		}
		assert.Greater(t, p, priorityMap["realm"],
			"%s priority %d should be greater than realm priority %d", rt, p, priorityMap["realm"])
	}
}

func TestDiscoverRelationshipPatterns(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	patterns, err := spec.DiscoverRelationshipPatterns()
	require.NoError(t, err)
	require.NotEmpty(t, patterns)

	hasUserGroup := false
	for _, p := range patterns {
		if p.Kind == "user-group-membership" {
			hasUserGroup = true
		}
	}
	assert.True(t, hasUserGroup, "expected user-group-membership pattern")
}

func TestClassifyRelationshipEndpoint(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/admin/realms/{realm}/users/{user-id}/groups", "user-group-membership"},
		{"/admin/realms/{realm}/roles-by-id/{role-id}/composites", "role-composite-mapping"},
		{"/admin/realms/{realm}/users/{user-id}/role-mappings/realm", "user-realm-role-mapping"},
		{"/admin/realms/{realm}/unknown/path", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := classifyRelationshipEndpoint(tt.path)
			if tt.want == "" {
				assert.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			assert.Equal(t, tt.want, got.Kind)
		})
	}
}

func TestBuildDependencyGraphFiltersOutSelfReferences(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	graph, err := spec.BuildDependencyGraph()
	require.NoError(t, err)

	for child, parents := range graph {
		for _, parent := range parents {
			assert.NotEqual(t, child, parent, "self-reference detected for %s", child)
		}
	}
}

func TestBuildDependencyGraphRealmIsRoot(t *testing.T) {
	spec, err := NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)

	graph, err := spec.BuildDependencyGraph()
	require.NoError(t, err)

	// Realm should only depend on itself or nothing
	realmDeps, exists := graph["realm"]
	require.True(t, exists, "realm should be in graph")
	assert.Empty(t, realmDeps, "realm should have no external dependencies")
}

func TestDependencyPriorityMapIgnoresUnknownParents(t *testing.T) {
	graph := map[string][]string{
		"realm": {},
		"child": {"realm", "unknown-parent"},
	}

	priorityMap, err := topologicalSortPriorityMap(graph)
	require.NoError(t, err)

	assert.Greater(t, priorityMap["child"], priorityMap["realm"])
}

func TestDependencyPriorityMapHandlesDisconnectedNodes(t *testing.T) {
	graph := map[string][]string{
		"a": {},
		"b": {},
	}

	priorityMap, err := topologicalSortPriorityMap(graph)
	require.NoError(t, err)

	assert.Contains(t, priorityMap, "a")
	assert.Contains(t, priorityMap, "b")
}

func TestTopologicalSortPriorityMapLevelGrouping(t *testing.T) {
	graph := map[string][]string{
		"root":  {},
		"a":     {"root"},
		"b":     {"root"},
		"child": {"a", "b"},
	}

	priorityMap, err := topologicalSortPriorityMap(graph)
	require.NoError(t, err)

	// All nodes at same level should have same priority
	assert.Equal(t, priorityMap["a"], priorityMap["b"], "a and b should be at same level")
	assert.Greater(t, priorityMap["child"], priorityMap["a"])
	assert.Greater(t, priorityMap["child"], priorityMap["b"])
}
