package catalog

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

var fallbackPlaceholderToResourceType = map[string]string{
	"realm":           "realm",
	"client-uuid":     "client",
	"client":          "client",
	"client-scope-id": "clientscope",
	"clientScopeId":   "clientscope",
	"group-id":        "group",
	"groupId":         "group",
	"user-id":         "user",
	"alias":           "identityprovider",
	"provider":        "identityprovider",
	"role-name":       "role",
	"role-id":         "role",
	"org-id":          "organization",
	"flowAlias":       "authenticationflow",
	"executionId":     "authenticationexecution",
}

// BuildDependencyGraph constructs a dependency graph from the spec contracts.
// An edge from A to B means "A depends on B" (B must be created before A).
func (s *Spec) BuildDependencyGraph() (map[string][]string, error) {
	contracts, err := s.ResourceContracts()
	if err != nil {
		return nil, err
	}

	placeholderMap, err := s.PlaceholderToResourceType()
	if err != nil {
		return nil, err
	}
	for k, v := range fallbackPlaceholderToResourceType {
		if _, ok := placeholderMap[k]; !ok {
			placeholderMap[k] = v
		}
	}

	graph := make(map[string][]string)

	for resourceType, contract := range contracts {
		postOp, hasPost := contract.Operations[http.MethodPost]
		if !hasPost {
			continue
		}

		path := postOp.Path
		segments := strings.Split(path, "/")

		for _, seg := range segments {
			if !strings.HasPrefix(seg, "{") || !strings.HasSuffix(seg, "}") {
				continue
			}
			placeholder := seg[1 : len(seg)-1]
			if placeholder == "realm" {
				continue
			}

			parentType, known := placeholderMap[placeholder]
			if !known {
				continue
			}
			if parentType == resourceType {
				continue
			}
			if !slices.Contains(graph[resourceType], parentType) {
				graph[resourceType] = append(graph[resourceType], parentType)
			}
		}
	}

	// Ensure all contract types are in the graph (top-level resources depend on realm).
	for resourceType := range contracts {
		if _, exists := graph[resourceType]; exists {
			continue
		}
		if resourceType == "realm" {
			graph[resourceType] = []string{}
			continue
		}
		graph[resourceType] = []string{"realm"}
	}

	return graph, nil
}

// DownwardChild describes a structural child resource type and the spec path
// used to list instances of that child under a parent instance.
type DownwardChild struct {
	ChildType string
	Path      string // GET collection path, e.g. /admin/realms/{realm}/clients/{client-uuid}/roles
}

// BuildDownwardGraph returns a map from parent resource type to the resource
// types that can be created directly underneath it in the Keycloak object
// hierarchy. Edges represent structural containment discovered from the spec
// POST paths. Relationship/mapping endpoints are excluded so the graph only
// describes object ownership, not associations.
func (s *Spec) BuildDownwardGraph() (map[string][]DownwardChild, error) {
	contracts, err := s.ResourceContracts()
	if err != nil {
		return nil, err
	}

	placeholderMap, err := s.PlaceholderToResourceType()
	if err != nil {
		return nil, err
	}
	for k, v := range fallbackPlaceholderToResourceType {
		if _, ok := placeholderMap[k]; !ok {
			placeholderMap[k] = v
		}
	}

	relationshipPaths := relationshipReadPaths()
	collectionPaths := s.collectionGetPaths()
	graph := make(map[string][]DownwardChild)

	for childType, contract := range contracts {
		for _, post := range contract.AllOperations[http.MethodPost] {
			path := post.Path
			if !isCollectionEndpoint(path) {
				continue
			}
			normalized := normalizeReadPath(path)
			if _, isRelationship := relationshipPaths[normalized]; isRelationship {
				continue
			}

			parentType := parentTypeForPath(path, placeholderMap)
			if parentType == "" {
				continue
			}

			getPath := findCollectionGetPath(contract.AllOperations[http.MethodGet], path)
			if getPath == "" {
				continue
			}
			if _, ok := collectionPaths[getPath]; !ok {
				continue
			}

			if !containsDownwardChild(graph[parentType], childType) {
				graph[parentType] = append(graph[parentType], DownwardChild{ChildType: childType, Path: getPath})
			}
		}
	}

	return graph, nil
}

// collectionGetPaths returns the set of GET paths that return an array response.
func (s *Spec) collectionGetPaths() map[string]struct{} {
	paths := make(map[string]struct{})
	s.ForEachOperation(func(path, method string, op *v3.Operation, item *v3.PathItem) {
		if method != http.MethodGet {
			return
		}
		if !strings.Contains(path, "{") {
			return
		}
		if isCollectionResponse(op) {
			paths[path] = struct{}{}
		}
	})
	return paths
}

// findCollectionGetPath returns the GET operation path that corresponds to the
// supplied POST path. It prefers an exact path match, then a path with the same
// parent placeholders and a trailing collection segment.
func findCollectionGetPath(getOps []OperationContract, postPath string) string {
	for _, get := range getOps {
		if get.Path == postPath {
			return get.Path
		}
	}
	postParams := extractPathParams(postPath)
	postLast := lastNonPlaceholderSegment(postPath)
	for _, get := range getOps {
		if !isCollectionEndpoint(get.Path) {
			continue
		}
		getParams := extractPathParams(get.Path)
		if !slices.Equal(postParams, getParams) {
			continue
		}
		getLast := lastNonPlaceholderSegment(get.Path)
		if getLast == postLast {
			return get.Path
		}
	}
	return ""
}

func lastNonPlaceholderSegment(path string) string {
	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		if part == "" {
			continue
		}
		if !strings.HasPrefix(part, "{") || !strings.HasSuffix(part, "}") {
			return part
		}
	}
	return ""
}

func containsDownwardChild(children []DownwardChild, childType string) bool {
	for _, c := range children {
		if c.ChildType == childType {
			return true
		}
	}
	return false
}

// parentTypeForPath returns the resource type of the closest non-realm path
// placeholder. A top-level resource under the realm returns "realm".
func parentTypeForPath(path string, placeholderMap map[string]string) string {
	params := extractPathParams(path)
	for i := len(params) - 1; i >= 0; i-- {
		param := params[i]
		if param == "realm" {
			continue
		}
		parentType, ok := placeholderMap[param]
		if !ok {
			return ""
		}
		return parentType
	}
	return "realm"
}

// relationshipReadPaths returns the normalized read paths of all registered
// relationship kinds. These are used to keep relationship/mapping endpoints out
// of the structural downward graph.
func relationshipReadPaths() map[string]struct{} {
	paths := make(map[string]struct{})
	for _, kind := range DefaultRegistry().Kinds() {
		paths[normalizeReadPath(kind.ReadPath)] = struct{}{}
	}
	return paths
}

// InvertDependencyGraph returns a new graph where each edge A -> B is reversed
// to B -> A. The input graph is interpreted as "child depends on parent".
func InvertDependencyGraph(graph map[string][]string) map[string][]string {
	downward := make(map[string][]string, len(graph))
	for node := range graph {
		if _, ok := downward[node]; !ok {
			downward[node] = []string{}
		}
	}
	for child, parents := range graph {
		for _, parent := range parents {
			if !slices.Contains(downward[parent], child) {
				downward[parent] = append(downward[parent], child)
			}
		}
	}
	return downward
}

// DownwardLevels returns the children found at each successive level below the
// supplied roots in the downward graph. The returned slice index corresponds to
// depth - 1 (index 0 is one level below the roots). Cycles are broken by
// tracking visited types. A non-positive maxDepth returns nil.
func DownwardLevels(graph map[string][]DownwardChild, roots []string, maxDepth int) [][]DownwardChild {
	if maxDepth <= 0 {
		return nil
	}
	seen := make(map[string]struct{})
	for _, root := range roots {
		seen[root] = struct{}{}
	}

	levels := make([][]DownwardChild, 0, maxDepth)
	frontier := roots
	for depth := 0; depth < maxDepth; depth++ {
		var next []DownwardChild
		for _, node := range frontier {
			for _, child := range graph[node] {
				if _, ok := seen[child.ChildType]; ok {
					continue
				}
				seen[child.ChildType] = struct{}{}
				next = append(next, child)
			}
		}
		if len(next) == 0 {
			break
		}
		levels = append(levels, next)
		frontier = make([]string, len(next))
		for i, child := range next {
			frontier[i] = child.ChildType
		}
	}
	return levels
}

// DependencyPriorityMap computes a topological sort of the dependency graph
// and returns a map from resource type to priority (lower = created first).
func (s *Spec) DependencyPriorityMap() (map[string]int, error) {
	graph, err := s.BuildDependencyGraph()
	if err != nil {
		return nil, err
	}

	return topologicalSortPriorityMap(graph)
}

func topologicalSortPriorityMap(graph map[string][]string) (map[string]int, error) {
	allNodes := make(map[string]struct{})
	for node, parents := range graph {
		allNodes[node] = struct{}{}
		for _, parent := range parents {
			allNodes[parent] = struct{}{}
		}
	}

	inDegree := make(map[string]int, len(allNodes))
	children := make(map[string][]string)
	for node := range allNodes {
		inDegree[node] = 0
	}
	for node, parents := range graph {
		for _, parent := range parents {
			if parent == node {
				continue
			}
			if _, ok := allNodes[parent]; !ok {
				continue
			}
			inDegree[node]++
			children[parent] = append(children[parent], node)
		}
	}

	queue := make([]string, 0, len(allNodes))
	for node, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, node)
		}
	}

	priority := 1
	result := make(map[string]int, len(allNodes))
	for len(queue) > 0 {
		nextQueue := make([]string, 0)
		for _, node := range queue {
			result[node] = priority
			for _, child := range children[node] {
				inDegree[child]--
				if inDegree[child] == 0 {
					nextQueue = append(nextQueue, child)
				}
			}
		}
		queue = nextQueue
		priority++
	}

	if len(result) != len(allNodes) {
		return nil, fmt.Errorf("dependency graph contains a cycle")
	}
	return result, nil
}

// DiscoverRelationshipOperations scans the spec for GET endpoints that return
// relationship/mapping collections and generates corresponding RelationshipOperation
// patterns. These patterns can be used to fetch and apply relationships.
type RelationshipOperationPattern struct {
	Kind                 string
	Method               string
	Path                 string
	PathParams           []string
	ResourceA            string // First resource type in relationship
	ResourceB            string // Second resource type in relationship
	Direction            string // "one-way" or "bidirectional"
	RelationshipTemplate string // Template for the relationship operation
	RelationshipMethod   string // HTTP method for the relationship operation
	ItemParamName        string // Item field to extract for leaf path param; empty if none
	PayloadField         string // Item field to use as payload; empty means use the whole item
	BulkPayload          bool   // True if the entire collection is sent as one operation
}

// ParentResourceTypes maps each path placeholder to its resource type using
// the placeholderToResourceType map. It returns an error if a placeholder is
// unknown.
func (p RelationshipOperationPattern) ParentResourceTypes(placeholderMap map[string]string) ([]string, error) {
	if placeholderMap == nil {
		placeholderMap = fallbackPlaceholderToResourceType
	}
	types := make([]string, 0, len(p.PathParams))
	for _, param := range p.PathParams {
		if param == "realm" {
			continue
		}
		resourceType, ok := placeholderMap[param]
		if !ok {
			resourceType = fallbackPlaceholderToResourceType[param]
		}
		if resourceType == "" {
			return nil, fmt.Errorf("unknown path parameter '%s' in pattern %s", param, p.Path)
		}
		types = append(types, resourceType)
	}
	return types, nil
}

// DiscoverRelationshipPatterns scans the spec for relationship endpoints.
func (s *Spec) DiscoverRelationshipPatterns() ([]RelationshipOperationPattern, error) {
	if _, err := s.ResourceContracts(); err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})
	var patterns []RelationshipOperationPattern

	s.ForEachOperation(func(path, method string, operation *v3.Operation, item *v3.PathItem) {
		if method != http.MethodGet {
			return
		}
		if !strings.Contains(path, "{") {
			return
		}
		if !isCollectionResponse(operation) {
			return
		}

		pattern := classifyRelationshipEndpoint(path)
		if pattern == nil {
			return
		}
		pattern.PathParams = extractPathParams(path)

		key := pattern.Kind + ":" + pattern.Path
		if _, exists := seen[key]; exists {
			return
		}
		seen[key] = struct{}{}
		patterns = append(patterns, *pattern)
	})

	return patterns, nil
}

// isCollectionResponse reports whether an operation returns an array schema.
func isCollectionResponse(operation *v3.Operation) bool {
	if operation == nil || operation.Responses == nil {
		return false
	}
	for pair := operation.Responses.Codes.First(); pair != nil; pair = pair.Next() {
		response := pair.Value()
		if response == nil || response.Content == nil {
			continue
		}
		for item := response.Content.First(); item != nil; item = item.Next() {
			media := item.Value()
			if media == nil || media.Schema == nil {
				continue
			}
			schema, err := media.Schema.BuildSchema()
			if err != nil {
				continue
			}
			if schema != nil && schema.Items != nil && schema.Items.IsA() {
				return true
			}
		}
	}
	return false
}
