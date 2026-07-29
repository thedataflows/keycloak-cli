# Org-scoped group deep nesting — Implementation Plan

> **For agentic workers:** Use superpowers:executing-plans to implement task-by-task. Steps use checkbox (`- [ ]`) syntax.

**Goal:** Make an organization's group hierarchy reachable at arbitrary depth through `admin.Service.FetchChildren`, with no consumer-side relationship override.

**Architecture:** Add a placeholder-chain-aware child-collection selector in `catalog` (`Spec.ScopedChildCollection`) and route `FetchChildren` through it. For an org-scoped parent it selects `/organizations/{org-id}/groups/{group-id}/children`, tags children `ParentType: organization` (a scope marker), and propagates the constant `orgId` down the subtree so a returned child is a valid parent for the next call. Realm-group and client-role behaviour is unchanged.

**Tech Stack:** Go 1.26, testify, embedded Keycloak 26.6.2 OpenAPI spec, `httptest` fakes.

## Global Constraints

- Go 1.26; no new dependencies.
- `go vet ./...` clean; `go test ./...` green (the pre-existing unrelated `pkg/output` TOML failure excepted).
- Follow existing package patterns; no public signature change to `FetchChildren`.
- Spec fixture path in tests: `filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json")`.

---

### Task 1: `ScopedChildCollection` selector in catalog

**Files:**
- Modify: `pkg/catalog/dependencies.go` (add `ScopedChildCollection` + `parentPlaceholderChain`)
- Test: `pkg/catalog/scoped_child_collection_test.go` (create)

**Interfaces:**
- Produces: `func (s *Spec) ScopedChildCollection(parentType, parentParentType, childType string) (path string, inheritedFields []string, ok bool)` — returns the GET collection path nested under a parent of the given chain, plus the camel-case grandparent-chain fields the child inherits (`["orgId"]` for org groups, `nil` for single-parent paths). `ok` is false when no such path exists.

- [ ] **Step 1: Write the failing test**

```go
package catalog_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScopedChildCollection(t *testing.T) {
	spec := loadOrgGroupSpec(t) // helper in org_scoped_groups_test.go (package catalog_test)

	tests := []struct {
		name             string
		parentType       string
		parentParentType string
		childType        string
		wantPath         string
		wantInherited    []string
	}{
		{"client role", "client", "", "role",
			"/admin/realms/{realm}/clients/{client-uuid}/roles", nil},
		{"realm group top", "group", "", "group",
			"/admin/realms/{realm}/groups/{group-id}/children", nil},
		{"realm subgroup", "group", "group", "group",
			"/admin/realms/{realm}/groups/{group-id}/children", nil},
		{"org group", "group", "organization", "group",
			"/admin/realms/{realm}/organizations/{org-id}/groups/{group-id}/children", []string{"orgId"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, inherited, ok := spec.ScopedChildCollection(tt.parentType, tt.parentParentType, tt.childType)
			require.True(t, ok)
			assert.Equal(t, tt.wantPath, path)
			assert.Equal(t, tt.wantInherited, inherited)
		})
	}

	_, _, ok := spec.ScopedChildCollection("client", "organization", "role")
	assert.False(t, ok, "no org-scoped client role collection exists")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/catalog/ -run TestScopedChildCollection`
Expected: FAIL — `spec.ScopedChildCollection` undefined.

- [ ] **Step 3: Implement the selector**

Add to `pkg/catalog/dependencies.go`:

```go
type placeholderRef struct {
	param string
	ptype string
}

// parentPlaceholderChain returns the ordered non-realm path placeholders of a
// collection path with their resource types.
func parentPlaceholderChain(path string, placeholderMap map[string]string) []placeholderRef {
	var chain []placeholderRef
	for _, param := range extractPathParams(path) {
		if param == "realm" {
			continue
		}
		chain = append(chain, placeholderRef{param: param, ptype: placeholderMap[param]})
	}
	return chain
}

// ScopedChildCollection selects the GET collection path for childType nested
// directly under a parent whose own type is parentType and whose parent type is
// parentParentType. It matches by placeholder chain rather than the deduped
// downward graph, so it can reach two-parent paths the graph collapses (an
// organization group's children). It returns the path and the camel-case
// grandparent-chain fields the child inherits from the parent (constant down the
// subtree). ok is false when no matching collection exists.
func (s *Spec) ScopedChildCollection(parentType, parentParentType, childType string) (string, []string, bool) {
	contracts, err := s.ResourceContracts()
	if err != nil {
		return "", nil, false
	}
	placeholderMap, err := s.PlaceholderToResourceType()
	if err != nil {
		return "", nil, false
	}
	for k, v := range fallbackPlaceholderToResourceType {
		if _, ok := placeholderMap[k]; !ok {
			placeholderMap[k] = v
		}
	}
	collectionPaths := s.collectionGetPaths()

	contract, ok := contracts[childType]
	if !ok {
		return "", nil, false
	}
	scoped := parentParentType != "" && parentParentType != parentType

	for _, post := range contract.AllOperations[http.MethodPost] {
		if !isCollectionEndpoint(post.Path) {
			continue
		}
		getPath := findCollectionGetPath(contract.AllOperations[http.MethodGet], post.Path)
		if getPath == "" {
			continue
		}
		if _, ok := collectionPaths[getPath]; !ok {
			continue
		}
		chain := parentPlaceholderChain(getPath, placeholderMap)
		if len(chain) == 0 {
			continue
		}
		if chain[len(chain)-1].ptype != parentType {
			continue // immediate container must be the parent's own type
		}
		if scoped {
			if len(chain) < 2 || chain[len(chain)-2].ptype != parentParentType {
				continue
			}
			inherited := make([]string, 0, len(chain)-1)
			for _, p := range chain[:len(chain)-1] {
				inherited = append(inherited, kebabToCamelCase(p.param))
			}
			return getPath, inherited, true
		}
		if len(chain) != 1 {
			continue // single-parent path only
		}
		return getPath, nil, true
	}
	return "", nil, false
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./pkg/catalog/ -run TestScopedChildCollection -v`
Expected: PASS (all four sub-cases + the false case).

- [ ] **Step 5: Commit**

```bash
git add pkg/catalog/dependencies.go pkg/catalog/scoped_child_collection_test.go
git commit -m "feat(catalog): scoped child-collection selector for two-parent nesting (ISSUE 0006)"
```

---

### Task 2: Route `FetchChildren` through the selector + org child tagging

**Files:**
- Modify: `pkg/admin/fetch.go` (`FetchChildren`; add `fetchScopedChildren`)
- Test: `pkg/admin/fetch_children_org_test.go` (create)

**Interfaces:**
- Consumes: `Spec.ScopedChildCollection` (Task 1); existing `fetchNestedResourceCollection`, `FetchPathCollection`, `classifyError`, `operationContext`, `fetchFailure`, `logFetchError`.
- Produces: unchanged public `FetchChildren` signature; new private `func (s *service) fetchScopedChildren(ctx context.Context, childType, path string, parent manifest.Resource, inherited []string, params ...map[string]string) ([]manifest.Resource, error)`.

- [ ] **Step 1: Write the failing test**

```go
package admin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/admin"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

// Walk org -> child -> grandchild via repeated FetchChildren, passing each
// returned child straight back in as the next parent (recursion must work out of
// the box: correct path per level, ParentType=organization, orgId propagated, no
// groupId that would collide and re-fetch the parent's children).
func TestFetchChildrenWalksOrgGroupHierarchy(t *testing.T) {
	const org = "org-1"
	var gets []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gets = append(gets, r.URL.Path)
		switch r.URL.Path {
		case "/admin/realms/demo/organizations/" + org + "/groups/parent-1/children":
			writeJSON(t, w, []map[string]interface{}{{"id": "child-1", "name": "child"}})
		case "/admin/realms/demo/organizations/" + org + "/groups/child-1/children":
			writeJSON(t, w, []map[string]interface{}{{"id": "grand-1", "name": "grandchild"}})
		case "/admin/realms/demo/organizations/" + org + "/groups/grand-1/children":
			writeJSON(t, w, []map[string]interface{}{})
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	// A top-level org group as ISSUE 0002 produces it: ParentType=organization + orgId.
	parent := manifest.Resource{Type: "group", Realm: "demo", ParentType: "organization",
		Data: map[string]interface{}{"id": "parent-1", "orgId": org, "name": "parent"}}

	names := []string{}
	cur := parent
	for {
		rep, err := service.FetchChildren(context.Background(), cur, "group", admin.ChildFetchQuery{})
		require.NoError(t, err)
		require.Empty(t, rep.Failures)
		if len(rep.Resources) == 0 {
			break
		}
		require.Len(t, rep.Resources, 1)
		child := rep.Resources[0]
		assert.Equal(t, "group", child.Type)
		assert.Equal(t, "organization", child.ParentType, "org scope marker must propagate for recursion")
		assert.Equal(t, org, child.Data["orgId"], "orgId must be propagated down the subtree")
		_, hasGroupID := child.Data["groupId"]
		assert.False(t, hasGroupID, "must not inject a colliding parent groupId")
		names = append(names, child.Data["name"].(string))
		cur = child
	}

	assert.Equal(t, []string{"child", "grandchild"}, names, "recursion must reach the grandchild")
	assert.Equal(t, []string{
		"/admin/realms/demo/organizations/" + org + "/groups/parent-1/children",
		"/admin/realms/demo/organizations/" + org + "/groups/child-1/children",
		"/admin/realms/demo/organizations/" + org + "/groups/grand-1/children",
	}, gets, "each level addresses the group by its own id")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/admin/ -run TestFetchChildrenWalksOrgGroupHierarchy`
Expected: FAIL — current `FetchChildren` resolves via `BuildDownwardGraph` (realm children path), so it never hits the org path; the first GET 404s / mis-routes.

- [ ] **Step 3: Implement routing + scoped fetch**

In `pkg/admin/fetch.go`, replace the body of `FetchChildren` (the `BuildDownwardGraph` lookup) with:

```go
func (s *service) FetchChildren(ctx context.Context, parent manifest.Resource, childType string, query ChildFetchQuery) (FetchReport, error) {
	path, inherited, ok := s.Spec().ScopedChildCollection(parent.Type, parent.ParentType, childType)
	if !ok {
		return FetchReport{}, fmt.Errorf("no child collection %q under parent %q", childType, parent.Type)
	}

	var params []map[string]string
	if query.FullRepresentation {
		params = []map[string]string{{"briefRepresentation": "false"}}
	}

	// An org-scoped parent (ParentType is a distinct grandparent type, e.g.
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
// (ISSUE 0005 Gap 3 on a collection path). Children are identified by
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
```

Remove the now-unused `BuildDownwardGraph` call from `FetchChildren`. Keep the `catalog` import (still used by `catalog.OperationContract`).

- [ ] **Step 4: Run the org test + the existing FetchChildren/regression suites**

Run: `go test ./pkg/admin/ -run 'TestFetchChildren' -v`
Expected: PASS — new org-walk test plus the existing `TestFetchChildrenReturnsOneClientRoleCollection` and `TestFetchChildrenClassifies404AsNotFound` (client-role path unchanged via the single-parent branch).

Run: `go test ./pkg/admin/... ./pkg/catalog/...`
Expected: PASS (realm-group apply tests, org-group apply tests, resolver pins unchanged).

- [ ] **Step 5: Commit**

```bash
git add pkg/admin/fetch.go pkg/admin/fetch_children_org_test.go
git commit -m "feat(admin): FetchChildren reaches org group hierarchies at any depth (ISSUE 0006)"
```

---

### Task 3: AC #3 — confirm the org-children create nests under a child

**Files:**
- Test: `pkg/catalog/org_child_create_test.go` (create)

**Interfaces:**
- Consumes: existing `Spec.ValidateOperationRequest`, `catalog.RequestValidation`.

- [ ] **Step 1: Write the test**

```go
package catalog_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/catalog"
)

// AC #3: nesting under an org CHILD is not special-cased to top-level org groups.
// The create path is spec-valid with {group-id} bound to a child id and a
// GroupRepresentation body. (The read is the blocker; this confirms the write is
// not a second one. Full resource-channel writes of grandchildren are out of scope.)
func TestOrgChildrenCreateNestsUnderChild(t *testing.T) {
	spec := loadOrgGroupSpec(t)
	require.NoError(t, spec.ValidateOperationRequest(
		"/admin/realms/{realm}/organizations/{org-id}/groups/{group-id}/children",
		http.MethodPost,
		catalog.RequestValidation{
			PathParams: map[string]string{"realm": "demo", "org-id": "org-1", "group-id": "child-1"},
			Body:       map[string]interface{}{"name": "grandchild"},
		}))
}
```

- [ ] **Step 2: Run test**

Run: `go test ./pkg/catalog/ -run TestOrgChildrenCreateNestsUnderChild -v`
Expected: PASS (POST on the org children path accepts a GroupRepresentation with a child id as `{group-id}`).

- [ ] **Step 3: Commit**

```bash
git add pkg/catalog/org_child_create_test.go
git commit -m "test(catalog): confirm org children create nests under a child (ISSUE 0006 AC3)"
```

---

### Task 4: Full verification, live check, close issue

**Files:**
- Modify: `docs/issues/0006-org-scoped-group-deep-nesting.md` (status, ACs, closing commits, live note)
- Modify: `docs/issues/index.md` (move 0006 to Done)

- [ ] **Step 1: Full suite + vet**

Run: `go vet ./... && go test ./pkg/catalog/... ./pkg/admin/...`
Expected: green.

- [ ] **Step 2: Live verification** — write a throwaway Go harness in the job tmp dir (NOT committed) that, against the running Keycloak: creates a realm + organization + org group + child + grandchild via the raw admin API, then uses `admin.New` + `FetchChildren` recursion starting from the org group, asserts the grandchild is reached with `ParentType: organization` and `orgId` set, and deletes the realm. Run with the base URL + a fresh admin token. Confirms the library reaches the grandchild with **no** relationship override registered.

- [ ] **Step 3: Close the issue doc** — set `Status: done`, tick ACs 1–5 (AC 6 release pending), fill `Closing commits`, add a live-verification note; move 0006 from Open to Done in `index.md`. Commit:

```bash
git add docs/issues/0006-org-scoped-group-deep-nesting.md docs/issues/index.md
git commit -m "docs: close ISSUE 0006 (org group deep nesting)"
```

- [ ] **Step 4: Push branch** `git push -u origin feat/org-group-deep-nesting-0006` and report. Fast-forward main + `v1.4.0` tag are user-driven (as in the v1.3.0 flow).

---

## Self-Review

- **Spec coverage:** §1 selector → Task 1; §2 routing/child shape → Task 2; §3 recursion safety (no colliding groupId) → Task 2 test assertion; §4 write confirmation → Task 3; testing/live → Task 4. AC1/AC2 → Task 2 + live; AC3 → Task 3; AC4 → Task 2 regression runs; AC5 → Task 4; AC6 → user-driven release.
- **Placeholder scan:** none — all steps carry real code and commands.
- **Type consistency:** `ScopedChildCollection(parentType, parentParentType, childType) (string, []string, bool)` defined in Task 1 and consumed identically in Task 2; `fetchScopedChildren(ctx, childType, path, parent, inherited, params...)` defined and called consistently.
