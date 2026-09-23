// graph.go is the graph-query surface of kcapi: Resolve turns a Ref (a
// resource type plus a human name or explicit id) into a fetched Node, and
// Neighbors walks the structural Edges() outward from a resolved node,
// listing the child collections the spec hangs under it.
package kcapi

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// Ref names a resource to resolve. Type speaks the plural path-segment
// vocabulary of Invoke's Resource field and Edge.Parent/Child (e.g. "users");
// Name is the human identity (username, client id, role name); ID, when set,
// skips the name search and fetches the single-resource path directly.
type Ref struct {
	Type  string
	Name  string
	ID    string
	Realm string
}

// Node is a resolved resource: the ref that found it plus the fetched
// object. Ref.Type keeps the caller's plural vocabulary while
// Node.Resource.Type carries the catalog's singular resource type, matching
// the identity helpers.
type Node struct {
	Ref      Ref
	Resource Resource
}

// EdgeFilter narrows a Neighbors walk. Zero fields match everything; set
// fields are ANDed with the node's own parent type.
type EdgeFilter struct {
	Child  string // only edges with this child type
	Parent string // only edges with this parent type
}

// Resolve fetches the resource a Ref names. With ID set it issues the
// single-resource GET for the type; otherwise it searches the type's
// collection (query param chosen from the identity fields when the spec
// declares one), filters hits by exact identity-field match, and requires
// exactly one hit: zero is ErrNotFound, several is a validation error listing
// the candidates.
func (c *Client) Resolve(ctx context.Context, ref Ref) (Node, error) {
	spec, _ := c.snapshot() // capture the pair once: in-flight calls keep it across Reload
	identities, err := spec.ResourceIdentities()
	if err != nil {
		return Node{}, &Error{Kind: KindValidation, Op: resolveOpLabel(ref), Err: err}
	}
	identity, ok := lookupIdentity(identities, ref.Type)
	if !ok {
		return Node{}, &Error{Kind: KindValidation, Op: resolveOpLabel(ref), Err: fmt.Errorf("unknown resource type %q", ref.Type)}
	}
	if ref.ID == "" && identity.IDParam == "realm" {
		// A realm's identity field IS the {realm} path parameter: the name
		// is the address, so resolve straight onto GET /admin/realms/{name}.
		ref.ID = ref.Name
	}
	if ref.ID != "" {
		return c.resolveSingle(ctx, spec, ref, identity)
	}
	return c.resolveByName(ctx, spec, ref, identity)
}

// resolveSingle fetches the single-resource GET for an explicitly identified
// ref. The id fills the identity's path parameter (e.g. "user-id" for users);
// the resolution must land on an item endpoint, never a collection, so a
// whole-collection body can never masquerade as the requested object.
func (c *Client) resolveSingle(ctx context.Context, spec *Spec, ref Ref, identity ResourceIdentity) (Node, error) {
	if identity.IDParam == "" {
		return Node{}, &Error{Kind: KindValidation, Op: resolveOpLabel(ref), Err: fmt.Errorf("resource type %q has no single-resource identity parameter", ref.Type)}
	}
	call := Call{Resource: ref.Type, Verb: Get, Realm: ref.Realm, Params: P{identity.IDParam: ref.ID}}
	if err := requireItemOperation(spec, call, resolveOpLabel(ref), ref, identity); err != nil {
		return Node{}, err
	}

	raw, err := c.Invoke(ctx, call)
	if err != nil {
		return Node{}, err
	}
	if len(raw) == 0 {
		return Node{}, &Error{Kind: KindNotFound, Op: resolveOpLabel(ref), Err: fmt.Errorf("%s %q not found in realm %q", ref.Type, ref.ID, ref.Realm)}
	}
	var data map[string]interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return Node{}, &Error{Kind: KindServer, Op: resolveOpLabel(ref), Err: fmt.Errorf("decode response: %w", err)}
	}
	return wrapNode(ref.Type, ref.Realm, identity, data), nil
}

// resolveMaxResults is the page size requested on the name-search collection
// GET. Keycloak collection endpoints default to small server-side pages, and
// the exact-match filtering here is client-side — without a cap it would only
// see the first default page and answer false ErrNotFound in large realms.
// This mirrors the manifest-side conflict-resolution cap semantics (the
// constant there lives in pkg/manifest, which kcapi must not import, so the
// value is restated locally).
const resolveMaxResults = "10000"

// resolveByName searches the type's collection GET and requires exactly one
// exact identity-field match. The search param is the first identity field
// the operation contract declares as a query parameter (username for users,
// clientId for clients); when none is declared, the walk lists the collection
// and filters client-side, which is always correct, just less efficient.
// The fetch is capped with max=resolveMaxResults (never first: the search is
// single-shot with client-side exact filtering), and a param named max that
// is already present in the call — e.g. the search param itself being "max" —
// wins over the cap.
func (c *Client) resolveByName(ctx context.Context, spec *Spec, ref Ref, identity ResourceIdentity) (Node, error) {
	call := Call{Resource: ref.Type, Verb: Get, Realm: ref.Realm}
	op, _, err := resolveCall(spec, call)
	if err != nil {
		return Node{}, &Error{Kind: KindValidation, Op: resolveOpLabel(ref), Err: err}
	}
	if !isCollectionEndpoint(op.path) {
		return Node{}, &Error{Kind: KindValidation, Op: resolveOpLabel(ref), Err: fmt.Errorf("resource type %q has no collection GET to search by name", ref.Type)}
	}
	if ref.Name != "" {
		if param := searchQueryParam(op.contract, identity); param != "" {
			call.Params = P{param: ref.Name}
		}
	}
	if _, exists := call.Params["max"]; !exists {
		if call.Params == nil {
			call.Params = P{}
		}
		call.Params["max"] = resolveMaxResults
	}

	raw, err := c.Invoke(ctx, call)
	if err != nil {
		return Node{}, err
	}
	var hits []map[string]interface{}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &hits); err != nil {
			return Node{}, &Error{Kind: KindServer, Op: resolveOpLabel(ref), Err: fmt.Errorf("decode response: %w", err)}
		}
	}

	var matches []map[string]interface{}
	for _, hit := range hits {
		if ref.Name == "" || matchesIdentity(hit, identity, ref.Name) {
			matches = append(matches, hit)
		}
	}
	switch len(matches) {
	case 0:
		if ref.Name == "" {
			return Node{}, &Error{Kind: KindNotFound, Op: resolveOpLabel(ref), Err: fmt.Errorf("no %s found in realm %q", ref.Type, ref.Realm)}
		}
		return Node{}, &Error{Kind: KindNotFound, Op: resolveOpLabel(ref), Err: fmt.Errorf("no %s named %q in realm %q", ref.Type, ref.Name, ref.Realm)}
	case 1:
		return wrapNode(ref.Type, ref.Realm, identity, matches[0]), nil
	default:
		var candidates []string
		for _, hit := range matches {
			candidates = append(candidates, NameOf(Resource{Type: identity.Type, Data: hit}, identity))
		}
		return Node{}, &Error{Kind: KindValidation, Op: resolveOpLabel(ref), Err: fmt.Errorf("ambiguous %s in realm %q: %d matches (%s); resolve by ID instead",
			ref.Type, ref.Realm, len(matches), strings.Join(candidates, ", "))}
	}
}

// Neighbors walks every edge whose parent is the node's type and that the
// filter admits, fetching each reachable child collection once. It returns
// the fetched child nodes plus the edges that produced them (one entry per
// fetched collection, in walk order); edges that cannot be instantiated from
// the node are skipped silently, never reported.
//
// Traversal instantiates the edge's collection prefix — the spec path the
// edge hangs under (Task 8's own parent-collection check guarantees that
// prefix exists) — rather than the edge's Path itself: an edge path ends at
// or after the child's own placeholder (…/users/{user-id}/groups/{groupId}),
// which identifies a single child Neighbors cannot know. The prefix
// (…/users/{user-id}/groups) is the listable collection. Its placeholders
// are filled from the node: {realm} from the node's realm and the anchor —
// the last placeholder before the child collection — from the node's
// identifier. An edge whose prefix still carries any other placeholder has
// nothing to fill it from and is skipped.
func (c *Client) Neighbors(ctx context.Context, node Node, filter EdgeFilter) ([]Node, []Edge, error) {
	spec, _ := c.snapshot() // capture the pair once: in-flight calls keep it across Reload
	edges := inferEdges(spec)
	sort.Slice(edges, func(i, j int) bool {
		a, b := edges[i], edges[j]
		return a.Parent+a.Child+a.Path < b.Parent+b.Child+b.Path
	})
	specPaths := make(map[string]struct{})
	spec.ForEachOperation(func(path, _ string, _ *v3.Operation, _ *v3.PathItem) {
		specPaths[path] = struct{}{}
	})
	identities, err := spec.ResourceIdentities()
	if err != nil {
		return nil, nil, &Error{Kind: KindValidation, Op: neighborsOpLabel(node), Err: err}
	}

	var nodes []Node
	var walked []Edge
	fetched := make(map[string]struct{}) // collection prefixes already requested
	for _, edge := range edges {
		if !edgeMatches(edge, node, filter) {
			continue
		}
		prefix, anchor, ok := edgeCollectionPrefix(edge, specPaths)
		if !ok {
			continue // the edge hangs under no spec path: nothing listable
		}
		realm := nodeRealm(node)
		values := map[string]string{}
		if realm != "" {
			values["realm"] = realm
		}
		if anchor == "realm" {
			// A realm node anchors its own edges: the identifier fills {realm}.
			values["realm"] = nodeIdentifier(node, identities)
		} else {
			values[anchor] = nodeIdentifier(node, identities)
		}
		if !pathPlaceholdersFilled(prefix, values) {
			continue // unfillable placeholder in the collection prefix: skip
		}

		call := Call{Resource: edge.Child, Verb: Get, Realm: realm, Params: values}
		// Invoke resolves Resource+Verb specificity-first, which can land on a
		// different operation than the edge's own collection (shared literal
		// segments, competing paths). Only fetch when resolution is exact.
		op, _, err := resolveCall(spec, call)
		if err != nil || op.path != prefix {
			continue
		}
		key := strings.ToLower(op.path)
		if _, seen := fetched[key]; seen {
			continue
		}

		raw, err := c.Invoke(ctx, call)
		if err != nil {
			return nil, nil, err
		}
		items, err := decodeNeighborItems(raw)
		if err != nil {
			return nil, nil, err
		}
		fetched[key] = struct{}{}
		walked = append(walked, edge)
		childIdentity, _ := lookupIdentity(identities, edge.Child)
		for _, item := range items {
			nodes = append(nodes, wrapNode(edge.Child, realm, childIdentity, item))
		}
	}
	return nodes, walked, nil
}

// edgeMatches applies the Neighbors selection: the edge must start from the
// node's type, and every set filter field must match too.
func edgeMatches(edge Edge, node Node, filter EdgeFilter) bool {
	if edge.Parent != node.Ref.Type {
		return false
	}
	if filter.Child != "" && edge.Child != filter.Child {
		return false
	}
	if filter.Parent != "" && edge.Parent != filter.Parent {
		return false
	}
	return true
}

// edgeCollectionPrefix returns the collection prefix of an edge's path — the
// spec path ending directly before the edge's child placeholder — together
// with the anchor placeholder that the parent node's identifier fills (the
// last placeholder before the child collection). candidate child positions
// are scanned in order and the first whose prefix is a real spec path wins,
// mirroring how inferEdges records the first (parent, child) pair per path.
func edgeCollectionPrefix(edge Edge, specPaths map[string]struct{}) (prefix, anchor string, ok bool) {
	segments := strings.Split(edge.Path, "/")
	wantChild := strings.ToLower(edge.Child)
	for j, segment := range segments {
		if j == 0 || !isPathPlaceholderSegment(segment) {
			continue
		}
		if strings.ToLower(segments[j-1]) != wantChild {
			continue
		}
		candidate := strings.Join(segments[:j], "/")
		if _, exists := specPaths[candidate]; !exists {
			continue
		}
		return candidate, lastPlaceholderName(segments[:j]), true
	}
	return "", "", false
}

// lastPlaceholderName returns the name of the last {placeholder} in a path
// prefix, or "" when the prefix carries none.
func lastPlaceholderName(prefixSegments []string) string {
	for i := len(prefixSegments) - 1; i >= 0; i-- {
		segment := prefixSegments[i]
		if isPathPlaceholderSegment(segment) {
			return segment[1 : len(segment)-1]
		}
	}
	return ""
}

// pathPlaceholdersFilled reports whether every placeholder of the template
// has a non-empty entry in values.
func pathPlaceholdersFilled(template string, values map[string]string) bool {
	for _, name := range pathPlaceholderNames(template) {
		if values[name] == "" {
			return false
		}
	}
	return true
}

// decodeNeighborItems unmarshals an Invoke response as a list of objects,
// falling back to a single wrapped object: collection endpoints of
// cardinality-one edges can answer a bare JSON object, and dropping it would
// hide the neighbor.
func decodeNeighborItems(raw json.RawMessage) ([]map[string]interface{}, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var items []map[string]interface{}
	if err := json.Unmarshal(raw, &items); err == nil {
		return items, nil
	}
	var single map[string]interface{}
	if err := json.Unmarshal(raw, &single); err != nil {
		return nil, &Error{Kind: KindServer, Err: fmt.Errorf("decode response: %w", err)}
	}
	return []map[string]interface{}{single}, nil
}

// wrapNode builds a Node around a fetched representation: the ref keeps the
// caller's plural type and gains the fetched identity and display name.
func wrapNode(typeName, realm string, identity ResourceIdentity, data map[string]interface{}) Node {
	resource := Resource{
		Type:  identity.Type,
		Realm: extractRealm(data, map[string]string{"realm": realm}),
		Data:  data,
	}
	if resource.Type == "" {
		resource.Type = typeName
	}
	return Node{
		Ref: Ref{
			Type:  typeName,
			Name:  NameOf(resource, identity),
			ID:    IdentifierOf(resource, identity),
			Realm: realm,
		},
		Resource: resource,
	}
}

// lookupIdentity finds the identity for a type in either vocabulary: the
// caller's plural path segment ("users") or the catalog's singular resource
// type ("user").
func lookupIdentity(identities map[string]ResourceIdentity, typeName string) (ResourceIdentity, bool) {
	if identity, ok := identities[typeName]; ok {
		return identity, true
	}
	identity, ok := identities[singularOf(typeName)]
	return identity, ok
}

// searchQueryParam returns the first identity field the operation declares as
// a query parameter, or "" when the search must happen client-side.
func searchQueryParam(contract OperationContract, identity ResourceIdentity) string {
	declared := make(map[string]struct{}, len(contract.Parameters))
	for _, parameter := range contract.Parameters {
		if parameter.In == "query" {
			declared[parameter.Name] = struct{}{}
		}
	}
	for _, field := range identity.IdentifierFields {
		if _, ok := declared[field]; ok {
			return field
		}
	}
	return ""
}

// matchesIdentity reports whether a fetched representation carries the wanted
// name in any of its identity fields, compared case-insensitively like the
// server-side search it supplements.
func matchesIdentity(data map[string]interface{}, identity ResourceIdentity, name string) bool {
	for _, field := range identity.IdentifierFields {
		if value := stringField(data, field); value != "" && strings.EqualFold(value, name) {
			return true
		}
	}
	return false
}

// nodeIdentifier returns the identifier a Neighbors walk fills the anchor
// placeholder with: the resource's identity-field identifier, falling back to
// the ref's explicit id and then to a bare "id" field.
func nodeIdentifier(node Node, identities map[string]ResourceIdentity) string {
	identity, ok := lookupIdentity(identities, node.Ref.Type)
	if ok {
		if id := IdentifierOf(node.Resource, identity); id != "" {
			return id
		}
	}
	if node.Ref.ID != "" {
		return node.Ref.ID
	}
	return stringField(node.Resource.Data, "id")
}

// nodeRealm returns the realm a Neighbors walk requests under: the ref's
// realm, falling back to a realm field on the fetched representation.
func nodeRealm(node Node) string {
	if node.Ref.Realm != "" {
		return node.Ref.Realm
	}
	return stringField(node.Resource.Data, "realm")
}

// requireItemOperation checks that a by-ID call resolves to a single-resource
// GET, so a collection body can never stand in for the named object.
func requireItemOperation(spec *Spec, call Call, label string, ref Ref, identity ResourceIdentity) error {
	op, _, err := resolveCall(spec, call)
	if err != nil {
		return &Error{Kind: KindValidation, Op: label, Err: err}
	}
	if isCollectionEndpoint(op.path) {
		return &Error{Kind: KindValidation, Op: label, Err: fmt.Errorf("resource type %q has no single-resource GET for %s %q", ref.Type, identity.IDParam, ref.ID)}
	}
	return nil
}

func resolveOpLabel(ref Ref) string {
	return "resolve " + ref.Type
}

func neighborsOpLabel(node Node) string {
	return "neighbors " + node.Ref.Type
}
