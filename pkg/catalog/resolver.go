package catalog

import (
	"fmt"
	"strings"

	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

// OperationShape controls which kind of endpoint the resolver should pick.
type OperationShape int

const (
	OperationAny OperationShape = iota
	OperationCollection
	OperationSingle
)

// Resolver resolves spec operations and path parameters for resources.
type Resolver struct {
	spec      *Spec
	contracts map[string]ResourceContract
}

// Resolver returns a new Resolver backed by the spec's resource contracts.
func (s *Spec) Resolver() *Resolver {
	contracts, _ := s.ResourceContracts()
	return &Resolver{spec: s, contracts: contracts}
}

// ResolveResourceOperation finds the best operation for a resource type and method.
// parentType disambiguates multi-location types (e.g. protocolmapper under
// clientscope vs client) by matching the parent segment in the path.
func (r *Resolver) ResolveResourceOperation(resourceType, parentType, method string, shape OperationShape) (OperationContract, error) {
	method = strings.ToUpper(method)
	contract, ok := r.contracts[resourceType]
	if !ok {
		return OperationContract{}, fmt.Errorf("no operations found for resource type '%s'", resourceType)
	}

	var candidates []OperationContract
	if ops, exists := contract.AllOperations[method]; exists {
		candidates = ops
	} else if op, e := contract.Operations[method]; e {
		candidates = []OperationContract{op}
	}
	if len(candidates) == 0 {
		return OperationContract{}, fmt.Errorf("no '%s' operation found for resource type '%s'", method, resourceType)
	}

	// Filter by parent type if specified.
	if parentType != "" {
		filtered := r.filterByParentType(candidates, parentType)
		if len(filtered) > 0 {
			candidates = filtered
		}
	}

	if len(candidates) == 1 {
		return candidates[0], nil
	}

	// Filter candidates by shape when possible.
	filtered := filterByShape(candidates, shape)
	if len(filtered) > 0 {
		candidates = filtered
	}

	if len(candidates) == 1 {
		return candidates[0], nil
	}

	var best OperationContract
	var bestScore int
	for _, candidate := range candidates {
		score := operationPriority(candidate, resourceType)
		if best.Path == "" || score > bestScore {
			best = candidate
			bestScore = score
		}
	}
	return best, nil
}

func filterByShape(candidates []OperationContract, shape OperationShape) []OperationContract {
	var single []OperationContract
	var collection []OperationContract
	for _, c := range candidates {
		if isCollectionEndpoint(c.Path) {
			collection = append(collection, c)
		} else {
			single = append(single, c)
		}
	}
	switch shape {
	case OperationSingle:
		if len(single) > 0 {
			return single
		}
	case OperationCollection:
		if len(collection) > 0 {
			return collection
		}
	}
	return nil
}

func (r *Resolver) filterByParentType(candidates []OperationContract, parentType string) []OperationContract {
	placeholderMap, err := r.spec.PlaceholderToResourceType()
	if err != nil {
		return candidates
	}

	var filtered []OperationContract
	for _, c := range candidates {
		firstParent := firstParentPlaceholderType(c.Path, placeholderMap)
		if firstParent == "" {
			continue
		}
		if firstParent == parentType {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

// firstParentPlaceholderType returns the resource type of the first path
// placeholder that appears after the realm segment. For nested resources this
// is the immediate parent type, which disambiguates endpoints where the parent
// collection segment appears later in the path (e.g. role composites).
func firstParentPlaceholderType(path string, placeholderMap map[string]string) string {
	parts := strings.Split(path, "/")
	seenRealm := false
	for _, part := range parts {
		if part == "{realm}" {
			seenRealm = true
			continue
		}
		if !seenRealm {
			continue
		}
		if !strings.HasPrefix(part, "{") || !strings.HasSuffix(part, "}") {
			continue
		}
		placeholder := part[1 : len(part)-1]
		if resourceType, ok := placeholderMap[placeholder]; ok {
			return resourceType
		}
	}
	return ""
}

// ResolveResourcePath resolves the operation for the resource and renders the
// path by substituting path parameters from resource.Data and resource.Realm.
func (r *Resolver) ResolveResourcePath(resource manifest.Resource, method string, shape OperationShape) (string, OperationContract, map[string]string, error) {
	contract, err := r.ResolveResourceOperation(resource.Type, resource.ParentType, method, shape)
	if err != nil {
		return "", OperationContract{}, nil, err
	}
	params, err := r.PathParams(resource, contract)
	if err != nil {
		return "", OperationContract{}, nil, err
	}
	path := RenderPath(contract.Path, params)
	return path, contract, params, nil
}

// PathParams extracts path parameters for a resource given an operation contract.
func (r *Resolver) PathParams(resource manifest.Resource, op OperationContract) (map[string]string, error) {
	params := map[string]string{"realm": resource.Realm}
	if op.Path == "" {
		return params, nil
	}
	primary := primaryIdentifierParam(op.Path)
	// On a single-resource path (…/{primary}) the primary placeholder addresses
	// the resource itself, so its own identifier must win over any same-named
	// value left in Data. A nested group carries its parent's id under "groupId",
	// which is also the camel form of the {group-id} self param; without this an
	// update or delete of a nested child would resolve to and clobber its parent
	// (ISSUE 0005 Gap 3, a destructive failure).
	selfPrimary := strings.HasSuffix(op.Path, "/{"+primary+"}")
	placeholders := extractPathParams(op.Path)
	for _, name := range placeholders {
		if name == "realm" {
			continue
		}
		preferIdentifier := selfPrimary && name == primary
		if value, ok := resolvePathParamValue(name, resource.Data, primary, resource.Identifier(), preferIdentifier); ok {
			params[name] = value
		}
	}
	return params, nil
}

func resolvePathParamValue(name string, data map[string]interface{}, primary, identifier string, preferIdentifier bool) (string, bool) {
	// An exact placeholder-name field is always an explicit override and wins.
	if value, ok := data[name].(string); ok && value != "" {
		return value, true
	}
	// The self placeholder on a single-resource path (…/{primary}) resolves only
	// from the resource's own identifier — never from the camel-form field. For a
	// nested child that field holds the PARENT's id (e.g. groupId), so consulting
	// it would make an update, a delete, or a create's existence probe address the
	// parent and clobber it (ISSUE 0005 Gap 3, a destructive failure). Falling
	// through to the own identifier here also preserves the "primary is always
	// present" contract even when the identifier is still empty on a create.
	if preferIdentifier {
		return identifier, true
	}
	camel := kebabToCamelCase(name)
	if camel != name {
		if value, ok := data[camel].(string); ok && value != "" {
			return value, true
		}
	}
	if name == primary {
		return identifier, true
	}
	return "", false
}

// RenderPath substitutes path placeholders with their values.
func RenderPath(path string, params map[string]string) string {
	result := path
	for key, value := range params {
		result = strings.ReplaceAll(result, "{"+key+"}", value)
	}
	return result
}

func primaryIdentifierParam(path string) string {
	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		if !strings.HasPrefix(part, "{") || !strings.HasSuffix(part, "}") {
			continue
		}
		name := part[1 : len(part)-1]
		if name != "realm" {
			return name
		}
	}
	return ""
}

func extractPathParams(path string) []string {
	var params []string
	for {
		start := strings.Index(path, "{")
		if start == -1 {
			break
		}
		end := strings.Index(path[start:], "}")
		if end == -1 {
			break
		}
		params = append(params, path[start+1:start+end])
		path = path[start+end+1:]
	}
	return params
}

func kebabToCamelCase(s string) string {
	parts := strings.Split(s, "-")
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

// ParentReferenceFieldTypes returns the data field names that identify parent
// resources for a nested operation path, mapped to the parent resource type.
// These fields are needed to render path parameters but should not be sent in
// the request body.
func (r *Resolver) ParentReferenceFieldTypes(resourceType string, op OperationContract) map[string]string {
	if op.Path == "" {
		return nil
	}
	placeholderMap, err := r.spec.PlaceholderToResourceType()
	if err != nil {
		return nil
	}
	result := make(map[string]string)
	for _, placeholder := range extractPathParams(op.Path) {
		if placeholder == "realm" {
			continue
		}
		parentType := placeholderMap[placeholder]
		if parentType == "" || parentType == resourceType {
			continue
		}
		result[kebabToCamelCase(placeholder)] = parentType
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// ParentReferenceFields returns child data field names and values that identify
// the parent resource for a nested operation path. When a child is fetched from
// a path such as /admin/realms/{realm}/client-scopes/{client-scope-id}/protocol-mappers/models,
// the returned map contains {"clientScopeId": "<parent-scope-id>"} so the child
// can be applied independently as a top-level resource.
func (r *Resolver) ParentReferenceFields(childPath, parentType string, parent manifest.Resource) map[string]string {
	if childPath == "" || parentType == "" {
		return nil
	}
	placeholderMap, err := r.spec.PlaceholderToResourceType()
	if err != nil {
		return nil
	}

	result := make(map[string]string)
	for _, param := range extractPathParams(childPath) {
		paramType, ok := placeholderMap[param]
		if !ok || paramType != parentType {
			continue
		}
		value := parentReferenceValue(parent, param)
		if value == "" {
			continue
		}
		field := kebabToCamelCase(param)
		result[field] = value
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// ParentReferenceFieldNames returns the data field names that identify parent
// resources for a nested operation path. These fields are needed to render path
// parameters but should not be sent in the request body.
func (r *Resolver) ParentReferenceFieldNames(resourceType string, op OperationContract) []string {
	if op.Path == "" {
		return nil
	}
	placeholderMap, err := r.spec.PlaceholderToResourceType()
	if err != nil {
		return nil
	}
	primary := primaryIdentifierParam(op.Path)
	// For single-resource endpoints the primary path parameter identifies the
	// resource itself (e.g. {id}) and must not be stripped from the body.
	skipPrimary := strings.HasSuffix(op.Path, "/{"+primary+"}")
	var fields []string
	for _, placeholder := range extractPathParams(op.Path) {
		if placeholder == "realm" {
			continue
		}
		if skipPrimary && placeholder == primary {
			continue
		}
		// The position guard (skipPrimary) already protects the resource's own
		// identifier. A non-primary same-type placeholder is a same-type parent
		// binding — a group nested under a group renders /groups/{group-id}/children
		// — and must also be stripped from the body. So we intentionally do not
		// skip parentType == resourceType here; that clause was broader than its
		// intent and swallowed the parent binding (ISSUE 0005 Gap 1).
		if placeholderMap[placeholder] == "" {
			continue
		}
		fields = append(fields, kebabToCamelCase(placeholder))
	}
	if len(fields) == 0 {
		return nil
	}
	return fields
}

func parentReferenceValue(parent manifest.Resource, placeholder string) string {
	if data := parent.Data; data != nil {
		if id, ok := data["id"].(string); ok && id != "" {
			return id
		}
	}
	return parent.Identifier()
}
