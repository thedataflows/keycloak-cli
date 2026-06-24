package catalog

import (
	"net/http"
)

// classifyRelationshipEndpoint matches a spec path against the relationship
// registry and returns the corresponding pattern metadata. It returns nil for
// unrecognized paths.
func classifyRelationshipEndpoint(path string) *RelationshipOperationPattern {
	if !isRelationshipPath(path) {
		return nil
	}
	kind, ok := DefaultRegistry().ByPath(path)
	if !ok {
		return nil
	}
	pattern := kind.Pattern()
	pattern.Path = path
	pattern.Method = http.MethodGet
	return &pattern
}
