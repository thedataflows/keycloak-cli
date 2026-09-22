package kcapi

import (
	"net/http"
	"sort"
	"strings"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// Edge is one parent->child relationship implied by the spec's path structure.
type Edge struct {
	Parent       string // parent resource in the plural path vocabulary, e.g. "realms"
	Child        string // child resource in the same vocabulary, e.g. "roles"
	Path         string // spec path template spanning both placeholders
	Cardinality  string // "one" when the path ends in an item placeholder, "many" when it ends in a collection segment
	CascadeOwner bool   // true when the child has no standalone top-level DELETE, so removing it goes through the parent
	Display      string // optional curated name from the registry display overrides
}

const (
	edgeCardinalityOne  = "one"
	edgeCardinalityMany = "many"
)

// Edges returns every relationship edge implied by the live spec, sorted by
// (Parent, Child, Path). Edges are structural: they are inferred from the path
// templates alone, then decorated with the registry's curated display names.
func (c *Client) Edges() []Edge {
	spec, _ := c.snapshot() // capture the pair once: in-flight calls keep it across Reload
	edges := inferEdges(spec)
	for i := range edges {
		if name, ok := displayOverride(edges[i].Parent, edges[i].Child); ok {
			edges[i].Display = name
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		a, b := edges[i], edges[j]
		return a.Parent+a.Child+a.Path < b.Parent+b.Child+b.Path
	})
	return edges
}

// edgePlaceholder is one typed placeholder occurrence in a path template.
type edgePlaceholder struct {
	index    int    // segment index of the placeholder
	resource string // resource type from the placeholder map
}

// inferEdges walks every spec path and records an edge per (parent placeholder,
// child placeholder) pair of consecutive typed placeholders. A pair yields an
// edge when the two placeholders name distinct resource types and the path
// prefix up to the child placeholder is itself a path in the spec — the parent
// collection the child hangs off. Paths of the roles-by-id family therefore
// yield nothing: roles are listed at /roles, so {role-id} has no parent
// collection of its own.
//
// Vocabulary: Parent/Child carry the plural path-segment names of
// operationResource — the same vocabulary Call{Resource} and Ref{Type} resolve
// by — so Edges() output feeds Invoke/Resolve directly. The placeholder map's
// singular catalog types ("user", "group", …) are used only to decide
// distinctness and never surface in the Edge fields.
func inferEdges(spec *Spec) []Edge {
	if spec == nil {
		return nil
	}

	placeholderMap, err := spec.PlaceholderToResourceType()
	if err != nil {
		placeholderMap = make(map[string]string) // degrade to the curated fallback below
	}
	for k, v := range fallbackPlaceholderToResourceType {
		if _, ok := placeholderMap[k]; !ok {
			placeholderMap[k] = v
		}
	}

	// First pass: the set of spec paths (for the parent-collection check) and
	// the collections that have a standalone top-level DELETE item
	// (/admin/realms/{realm}/<collection>/{id}).
	allPaths := make(map[string]struct{})
	standaloneDeletes := make(map[string]struct{})
	spec.ForEachOperation(func(path, method string, _ *v3.Operation, _ *v3.PathItem) {
		allPaths[path] = struct{}{}
		if method != http.MethodDelete {
			return
		}
		segments := strings.Split(path, "/")
		if len(segments) == 6 && segments[3] == "{realm}" && isPathPlaceholderSegment(segments[5]) {
			standaloneDeletes[segments[4]] = struct{}{}
		}
	})

	pathSeen := make(map[string]struct{})
	edgeSeen := make(map[string]struct{})
	var edges []Edge
	spec.ForEachOperation(func(path, method string, _ *v3.Operation, _ *v3.PathItem) {
		if _, ok := pathSeen[path]; ok { // one edge set per path template, not per method
			return
		}
		pathSeen[path] = struct{}{}

		segments := strings.Split(path, "/")
		var placeholders []edgePlaceholder
		for i, segment := range segments {
			if !isPathPlaceholderSegment(segment) {
				continue
			}
			resource := placeholderMap[segment[1:len(segment)-1]]
			if resource == "" {
				continue // an untyped placeholder cannot anchor an edge
			}
			placeholders = append(placeholders, edgePlaceholder{index: i, resource: resource})
		}

		for i := 0; i+1 < len(placeholders); i++ {
			parent, child := placeholders[i], placeholders[i+1]
			if parent.resource == child.resource {
				continue // self-edges are noise (e.g. {client-uuid} -> {client})
			}
			parentName := collectionSegment(segments, parent.index)
			childName := collectionSegment(segments, child.index)
			if parentName == "" || childName == "" || parentName == childName {
				continue
			}
			if _, ok := allPaths[strings.Join(segments[:child.index], "/")]; !ok {
				continue // the child has no parent-collection path in the spec
			}

			cardinality := edgeCardinalityOne
			if isCollectionEndpoint(path) {
				cardinality = edgeCardinalityMany
			}
			_, standalone := standaloneDeletes[childName]

			key := parentName + "->" + childName + " " + path
			if _, ok := edgeSeen[key]; ok {
				continue
			}
			edgeSeen[key] = struct{}{}
			edges = append(edges, Edge{
				Parent:       parentName,
				Child:        childName,
				Path:         path,
				Cardinality:  cardinality,
				CascadeOwner: !standalone,
			})
		}
	})
	return edges
}

// isPathPlaceholderSegment reports whether a path segment is a {name}
// placeholder.
func isPathPlaceholderSegment(segment string) bool {
	return len(segment) > 2 && strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}")
}

// collectionSegment returns the lowercased literal segment naming the
// collection the placeholder at segments[index] hangs under — the segment
// right before it. It returns "" when that segment is missing or is itself a
// placeholder, since neither names a collection.
func collectionSegment(segments []string, index int) string {
	if index == 0 {
		return ""
	}
	name := strings.ToLower(segments[index-1])
	if name == "" || isPathPlaceholderSegment(name) {
		return ""
	}
	return name
}
