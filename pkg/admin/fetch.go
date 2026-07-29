package admin

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/thedataflows/keycloak-cli/pkg/catalog"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

const defaultFetchResources = "realm,user,client,group,role"

type FetchQuery struct {
	Realm                string
	Resources            string
	Search               string
	Max                  int
	Parent               string
	IncludeRelationships bool
	Depth                int
	Filter               string
	ExactMatch           bool
	// FullRepresentation requests Keycloak's complete resource representation on
	// collection fetches by sending briefRepresentation=false. Keycloak's list
	// endpoints for users, groups and organizations otherwise return a brief
	// representation that omits the attributes map. The false default preserves
	// the existing brief behavior.
	FullRepresentation bool
}

// ChildFetchQuery configures a parent-scoped collection fetch (FetchChildren).
type ChildFetchQuery struct {
	// FullRepresentation requests briefRepresentation=false so children that
	// carry an attributes map (e.g. client roles) return it, matching the
	// representation flag on the realm-level fetch path (ISSUE 0001 parity).
	FullRepresentation bool
}

type FetchReport struct {
	Resources     []manifest.Resource
	Relationships []manifest.RelationshipOperation
	Failures      []FetchFailure
}

// FetchFailure describes a single resource that could not be fetched.
//
// NotFound is true when the underlying error was an HTTP 404, meaning an
// optional resource is simply absent rather than genuinely broken — for
// example authorization resources/scopes on a client that does not have
// authorization services enabled, or organizations on a realm where the
// feature is disabled. The command layer aggregates these separately from
// real failures so the output distinguishes "absent" from "broken".
//
// Resource is the resource type (e.g. "scope", "organization") and is used
// as the aggregation key. Detail carries additional context such as the
// parent identifier or realm.
type FetchFailure struct {
	Resource string
	Detail   string
	NotFound bool
	Err      error
}

func (f FetchFailure) String() string {
	s := f.Resource
	if f.Detail != "" {
		s += " " + f.Detail
	}
	if f.Err != nil {
		if s != "" {
			s += ": "
		}
		s += f.Err.Error()
	}
	return s
}

func fetchFailure(resource, detail string, err error) FetchFailure {
	return FetchFailure{
		Resource: resource,
		Detail:   detail,
		NotFound: isNotFound(err),
		Err:      err,
	}
}

func logFetchError(label string, err error) {
	evt := log.Logger.Error().Str("pkg", "admin").Err(err)
	if isNotFound(err) {
		evt = log.Logger.Debug().Str("pkg", "admin").Err(err)
	}
	evt.Msgf("fetch %s", label)
}

func (s *service) Fetch(ctx context.Context, query FetchQuery) (FetchReport, error) {
	realms, err := s.fetchRealms(ctx, query.Realm)
	if err != nil {
		return FetchReport{}, err
	}

	resourceList := strings.TrimSpace(query.Resources)
	if resourceList == "" {
		resourceList = defaultFetchResources
	}

	reqs, includeRelationships := requestedResources(resourceList)
	includeRelationships = includeRelationships || query.IncludeRelationships || query.Depth > 1
	var results []manifest.Resource
	var failures []FetchFailure
	realmNames := realmNamesFromResources(realms)
	queryParams := buildQueryParams(query)

	for _, rawResource := range reqs {
		resource := strings.TrimSpace(rawResource)
		if resource == "" {
			continue
		}

		if resource == "realm" {
			fetched, fetchErr := s.fetchResourceCollection(ctx, resource, nil, query.Realm)
			if fetchErr != nil {
				logFetchError(resource, fetchErr)
				failures = append(failures, fetchFailure(resource, "", fetchErr))
				continue
			}
			results = append(results, fetched...)
			continue
		}

		if len(realmNames) == 0 {
			failures = append(failures, FetchFailure{Resource: resource, Detail: "no realms available"})
			continue
		}

		if resource == "authenticationexecution" || resource == "authenticationexecutioninfo" {
			for _, realm := range realmNames {
				fetched, fetchErr := s.fetchAuthenticationExecutions(ctx, realm, query.Parent)
				if fetchErr != nil {
					logFetchError(resource+" for realm "+realm, fetchErr)
					failures = append(failures, fetchFailure(resource, realm, fetchErr))
					continue
				}
				results = append(results, fetched...)
			}
			continue
		}

		for _, realm := range realmNames {
			fetched, fetchErr := s.fetchRealmScopedResources(ctx, resource, realm, queryParams)
			if fetchErr != nil {
				logFetchError(resource+" for realm "+realm, fetchErr)
				failures = append(failures, fetchFailure(resource, realm, fetchErr))
				continue
			}
			results = append(results, fetched...)
		}
	}

	seeds := results
	if query.Filter != "" {
		seeds = filterResources(results, query.Filter)
		results = seeds
	}

	var childTypes map[string]struct{}
	if query.Depth > 0 && len(realmNames) > 0 {
		var depthResources []manifest.Resource
		var depthFailures []FetchFailure
		depthResources, childTypes, depthFailures = s.fetchDepthLevels(ctx, query.Depth, realmNames, seeds, buildNestedQueryParams(query))
		results = append(results, depthResources...)
		failures = append(failures, depthFailures...)
	}

	relationships := make([]manifest.RelationshipOperation, 0)
	if includeRelationships && len(realmNames) > 0 {
		if query.IncludeRelationships {
			fetchedRelationships, relationshipFailures := s.fetchRelationships(ctx, realmNames, nil)
			relationships = append(relationships, fetchedRelationships...)
			failures = append(failures, relationshipFailures...)
		}
		if query.Depth > 1 {
			parentTypes := make(map[string]struct{})
			for _, r := range seeds {
				parentTypes[r.Type] = struct{}{}
			}
			for t := range childTypes {
				parentTypes[t] = struct{}{}
			}
			scopedResources := make([]manifest.Resource, 0, len(seeds)+len(results))
			scopedResources = append(scopedResources, seeds...)
			for _, r := range results {
				if _, ok := parentTypes[r.Type]; ok {
					scopedResources = append(scopedResources, r)
				}
			}
			fetchedRelationships, relationshipFailures := s.fetchRelationshipsForResources(ctx, realmNames, parentTypes, scopedResources)
			relationships = append(relationships, fetchedRelationships...)
			failures = append(failures, relationshipFailures...)
		}
	}

	return FetchReport{Resources: results, Relationships: relationships, Failures: failures}, nil
}

// FetchChildren returns one child collection of one parent resource — e.g. the
// roles of a single client — with exactly one HTTP GET and none of Fetch's depth
// fan-out or realm-wide reference-resolution sweep. parent must carry enough
// identity to render the nested path (typically Realm + Data["id"]).
//
// Returned children carry Type=childType, Realm=parent.Realm, the parent's scope
// in ParentType, and the parent reference field(s) needed to re-fetch or apply
// them, matching what Fetch with Depth: 1 produces. A 404 on the child collection
// is reported as a FetchFailure{NotFound: true} rather than a hard error,
// consistent with the depth traversal, so an absent optional collection is
// benign for the caller.
//
// The child collection path is selected by matching the parent's placeholder
// chain (Spec.ScopedChildCollection), which reaches two-parent paths the deduped
// downward graph collapses — an organization group's children live at
// /organizations/{org-id}/groups/{group-id}/children. A caller recurses by
// passing each returned child straight back in; org groups keep
// ParentType=organization (a scope marker) at every level so descent terminates
// only on an empty read.
func (s *service) FetchChildren(ctx context.Context, parent manifest.Resource, childType string, query ChildFetchQuery) (FetchReport, error) {
	path, inherited, ok := s.Spec().ScopedChildCollection(parent.Type, parent.ParentType, childType)
	if !ok {
		return FetchReport{}, fmt.Errorf("no child collection %q under parent %q", childType, parent.Type)
	}

	var params []map[string]string
	if query.FullRepresentation {
		params = []map[string]string{{"briefRepresentation": "false"}}
	}

	// An org-scoped parent (its ParentType is a distinct grandparent type, e.g.
	// organization) needs the grandparent chain (orgId) propagated to children and
	// the scope marker preserved so recursion keeps selecting the org path. The
	// single-parent case keeps the existing immediate parent-reference injection.
	var (
		fetched  []manifest.Resource
		fetchErr error
	)
	if parent.ParentType != "" && parent.ParentType != parent.Type {
		fetched, fetchErr = s.fetchScopedChildren(ctx, childType, path, parent, inherited, params...)
	} else {
		fetched, fetchErr = s.fetchNestedResourceCollection(ctx, childType, path, parent.Type, parent, params...)
	}
	if fetchErr != nil {
		logFetchError(childType+" under "+parent.Type, fetchErr)
		return FetchReport{
			Failures: []FetchFailure{fetchFailure(childType, parent.Type+":"+parent.Identifier(), fetchErr)},
		}, nil
	}
	return FetchReport{Resources: fetched}, nil
}

// fetchScopedChildren fetches a child collection under an org-scoped parent. It
// tags children with the parent's scope (ParentType) and copies the inherited
// grandparent-chain fields (orgId) from the parent, but deliberately injects no
// immediate parent id: on the org children path {group-id} is a collection
// placeholder rendered from the parent's own id, and a groupId in the child's
// data would win via camel-case lookup and re-fetch the parent's children
// (the ISSUE 0005 Gap 3 hazard, on a collection path). Children are identified by
// (orgId, own id, ParentType).
func (s *service) fetchScopedChildren(ctx context.Context, childType, path string, parent manifest.Resource, inherited []string, params ...map[string]string) ([]manifest.Resource, error) {
	contract := catalog.OperationContract{Path: path, Method: http.MethodGet}
	scope, err := s.Spec().Resolver().PathParams(parent, contract)
	if err != nil {
		return nil, err
	}

	operationCtx, cancel := s.operationContext(ctx)
	defer cancel()

	raw, err := s.specClient.FetchPathCollection(operationCtx, path, scope, params...)
	if err != nil {
		return nil, classifyError(err, 0, "fetch", childType)
	}

	inheritedVals := make(map[string]string, len(inherited))
	for _, field := range inherited {
		if v, ok := parent.Data[field].(string); ok && v != "" {
			inheritedVals[field] = v
		}
	}

	resources := make([]manifest.Resource, len(raw))
	for i, item := range raw {
		for field, value := range inheritedVals {
			if _, ok := item[field]; !ok {
				item[field] = value
			}
		}
		resources[i] = manifest.Resource{
			Type:       childType,
			Realm:      parent.Realm,
			ParentType: parent.ParentType,
			Data:       item,
		}
	}
	return resources, nil
}

func (s *service) fetchRealms(ctx context.Context, realm string) ([]manifest.Resource, error) {
	realm = strings.TrimSpace(realm)
	if realm != "" {
		return []manifest.Resource{{
			Type:  "realm",
			Realm: realm,
			Data:  map[string]interface{}{"realm": realm},
		}}, nil
	}

	operationCtx, cancel := s.operationContext(ctx)
	defer cancel()

	realms, err := s.specClient.FetchResources(operationCtx, "realm", nil)
	if err != nil {
		return nil, classifyError(err, 0, "fetch", "realm")
	}
	return realms, nil
}

func (s *service) fetchRealmScopedResources(ctx context.Context, resource, realm string, params []map[string]string) ([]manifest.Resource, error) {
	fetched, err := s.fetchResourceCollection(ctx, resource, map[string]string{"realm": realm}, "", params...)
	if err != nil {
		return nil, err
	}
	for i := range fetched {
		fetched[i].Realm = realm
	}
	return fetched, nil
}

func (s *service) fetchAuthenticationExecutions(ctx context.Context, realm, parent string) ([]manifest.Resource, error) {
	flows, err := s.fetchRealmScopedResources(ctx, "authenticationflow", realm, nil)
	if err != nil {
		return nil, err
	}

	var results []manifest.Resource
	for _, flow := range flows {
		alias := ""
		if s, ok := flow.Data["alias"].(string); ok {
			alias = strings.TrimSpace(s)
		}
		if alias == "" {
			continue
		}
		if parent != "" && alias != parent {
			continue
		}

		operationCtx, cancel := s.operationContext(ctx)
		fetched, fetchErr := s.specClient.FetchPathCollection(operationCtx, "/admin/realms/{realm}/authentication/flows/{flowAlias}/executions", map[string]string{
			"realm":     realm,
			"flowAlias": alias,
		})
		cancel()
		if fetchErr != nil {
			return nil, classifyError(fetchErr, 0, "fetch", "authenticationexecution")
		}

		for _, raw := range fetched {
			raw["flowAlias"] = alias
			results = append(results, manifest.Resource{
				Type:  "authenticationexecution",
				Realm: realm,
				Data:  raw,
			})
		}
	}
	return results, nil
}

func (s *service) fetchResourceCollection(ctx context.Context, resource string, scope map[string]string, realmFilter string, params ...map[string]string) ([]manifest.Resource, error) {
	operationCtx, cancel := s.operationContext(ctx)
	defer cancel()

	fetched, err := s.specClient.FetchResources(operationCtx, resource, scope, params...)
	if err != nil {
		return nil, classifyError(err, 0, "fetch", resource)
	}

	if realmFilter == "" || resource != "realm" {
		return fetched, nil
	}

	filtered := make([]manifest.Resource, 0)
	for _, item := range fetched {
		if item.Realm == realmFilter {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func buildQueryParams(query FetchQuery) []map[string]string {
	params := make(map[string]string)
	if query.Search != "" {
		params["search"] = query.Search
	}
	if query.Max > 0 {
		params["max"] = fmt.Sprintf("%d", query.Max)
	}
	if query.ExactMatch {
		params["exact"] = "true"
	}
	// ponytail: emitted for every collection endpoint, not just the users, groups
	// and organizations operations that declare it. Keycloak ignores unrecognized
	// query parameters and validateOperationInput only checks spec-declared ones,
	// so this is inert elsewhere. Upgrade path: gate on the resolved operation's
	// parameter list, which is available in RuntimeClient.FetchResources.
	if query.FullRepresentation {
		params["briefRepresentation"] = "false"
	}

	if len(params) == 0 {
		return nil
	}

	return []map[string]string{params}
}

// buildNestedQueryParams returns the query params that apply to structural child
// collections. Only the representation flag propagates: search and max scope the
// requested resources, and forwarding them would silently filter children too.
func buildNestedQueryParams(query FetchQuery) []map[string]string {
	if !query.FullRepresentation {
		return nil
	}
	return []map[string]string{{"briefRepresentation": "false"}}
}

func realmNamesFromResources(realms []manifest.Resource) []string {
	names := make([]string, 0, len(realms))
	for _, realm := range realms {
		if name := realmName(realm); name != "" {
			names = append(names, name)
		}
	}
	return names
}

func realmName(realm manifest.Resource) string {
	if name := strings.TrimSpace(realm.Realm); name != "" {
		return name
	}
	if name := strings.TrimSpace(realm.DisplayName()); name != "" {
		return name
	}
	if value, ok := realm.Data["realm"].(string); ok {
		return strings.TrimSpace(value)
	}
	return ""
}

func requestedResources(resourceList string) ([]string, bool) {
	parts := strings.Split(resourceList, ",")
	resources := make([]string, 0, len(parts))
	includeRelationships := false
	for _, raw := range parts {
		resource := strings.TrimSpace(raw)
		if resource == "" {
			continue
		}
		if resource == "relationship" || resource == "relationships" {
			includeRelationships = true
			continue
		}
		resources = append(resources, resource)
	}
	return resources, includeRelationships
}

func (s *service) fetchDepthLevels(ctx context.Context, depth int, realmNames []string, seeds []manifest.Resource, nestedParams []map[string]string) ([]manifest.Resource, map[string]struct{}, []FetchFailure) {
	downward, err := s.Spec().BuildDownwardGraph()
	if err != nil {
		return nil, nil, []FetchFailure{{Resource: "depth-graph", Err: err}}
	}

	childParents := make(map[string][]string)
	childTypes := make(map[string]struct{})
	for parentType, children := range downward {
		for _, child := range children {
			childParents[child.ChildType] = append(childParents[child.ChildType], parentType)
			childTypes[child.ChildType] = struct{}{}
		}
	}

	resourcesByTypeRealm := indexResourcesByTypeRealm(seeds)
	var results []manifest.Resource
	var failures []FetchFailure
	seen := make(map[string]struct{})
	for _, r := range seeds {
		seen[resourceKey(r)] = struct{}{}
	}

	frontier := seeds
	for level := 0; level < depth; level++ {
		var levelResources []manifest.Resource

		for _, parent := range frontier {
			children := downward[parent.Type]
			// An org-scoped parent (its ParentType is a distinct grandparent type,
			// e.g. a group under an organization) descends through the scoped child
			// collection — the same mechanism FetchChildren uses — so the traversal
			// follows /organizations/{org-id}/groups/{group-id}/children instead of
			// the realm children path Keycloak rejects for org groups, and keeps the
			// org scope marker so descent continues past one level (ISSUE 0006).
			scoped := parent.ParentType != "" && parent.ParentType != parent.Type
			for _, child := range children {
				var (
					fetched  []manifest.Resource
					fetchErr error
				)
				if scoped {
					path, inherited, ok := s.Spec().ScopedChildCollection(parent.Type, parent.ParentType, child.ChildType)
					if !ok {
						continue
					}
					fetched, fetchErr = s.fetchScopedChildren(ctx, child.ChildType, path, parent, inherited, nestedParams...)
				} else {
					fetched, fetchErr = s.fetchNestedResourceCollection(ctx, child.ChildType, child.Path, parent.Type, parent, nestedParams...)
				}
				if fetchErr != nil {
					failures = append(failures, fetchFailure(child.ChildType, parent.Type+":"+parent.Identifier(), fetchErr))
					continue
				}
				// The fetch helpers already tag Realm and ParentType (parent.Type for
				// the structural case, the org scope for the scoped case); do not
				// overwrite, or a scoped child would lose its organization marker and
				// the next level would resolve to the realm children path.
				for i := range fetched {
					key := resourceKey(fetched[i])
					if _, ok := seen[key]; ok {
						continue
					}
					seen[key] = struct{}{}
					levelResources = append(levelResources, fetched[i])
				}
			}
		}

		references, refFailures := s.resolveReferences(ctx, realmNames, append(seeds, append(results, levelResources...)...))
		failures = append(failures, refFailures...)
		for _, r := range references {
			key := resourceKey(r)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			levelResources = append(levelResources, r)
		}

		for _, r := range levelResources {
			byRealm, ok := resourcesByTypeRealm[r.Type]
			if !ok {
				byRealm = make(map[string][]manifest.Resource)
				resourcesByTypeRealm[r.Type] = byRealm
			}
			byRealm[r.Realm] = append(byRealm[r.Realm], r)
		}

		results = append(results, levelResources...)
		frontier = levelResources
		if len(frontier) == 0 {
			break
		}
	}

	return results, childTypes, failures
}

func (s *service) fetchNestedResourceCollection(ctx context.Context, childType, path, parentType string, parent manifest.Resource, params ...map[string]string) ([]manifest.Resource, error) {
	contract := catalog.OperationContract{Path: path, Method: http.MethodGet}
	scope, err := s.Spec().Resolver().PathParams(parent, contract)
	if err != nil {
		return nil, err
	}

	operationCtx, cancel := s.operationContext(ctx)
	defer cancel()

	fetched, err := s.specClient.FetchPathCollection(operationCtx, path, scope, params...)
	if err != nil {
		return nil, classifyError(err, 0, "fetch", childType)
	}

	parentFields := s.Spec().Resolver().ParentReferenceFields(path, parentType, parent)

	resources := make([]manifest.Resource, len(fetched))
	for i, raw := range fetched {
		for field, value := range parentFields {
			if _, ok := raw[field]; !ok {
				raw[field] = value
			}
		}
		resources[i] = manifest.Resource{
			Type:       childType,
			Realm:      parent.Realm,
			ParentType: parentType,
			Data:       raw,
		}
	}
	return resources, nil
}

func indexResourcesByTypeRealm(resources []manifest.Resource) map[string]map[string][]manifest.Resource {
	index := make(map[string]map[string][]manifest.Resource)
	for _, r := range resources {
		byRealm, ok := index[r.Type]
		if !ok {
			byRealm = make(map[string][]manifest.Resource)
			index[r.Type] = byRealm
		}
		byRealm[r.Realm] = append(byRealm[r.Realm], r)
	}
	return index
}

func resourceKey(r manifest.Resource) string {
	return strings.Join([]string{r.Type, r.Realm, r.Identifier()}, "|")
}

func filterResources(resources []manifest.Resource, filter string) []manifest.Resource {
	needle := strings.ToLower(strings.TrimSpace(filter))
	out := make([]manifest.Resource, 0, len(resources))
	for _, r := range resources {
		if strings.ToLower(r.Name()) == needle || strings.ToLower(r.Identifier()) == needle {
			out = append(out, r)
		}
	}
	return out
}
