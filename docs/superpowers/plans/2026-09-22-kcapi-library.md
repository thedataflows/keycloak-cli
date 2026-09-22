# kcapi Library + CLI Migration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build `pkg/kcapi`, a first-class spec-driven Keycloak client library (generic invoke over every OpenAPI operation, plus object-relationship graph queries), migrate the CLI onto it, and dissolve `pkg/admin`.

**Architecture:** Absorb `pkg/catalog` into `pkg/kcapi`, hoist the core data types out of `pkg/manifest` into `kcapi` (type aliases keep compatibility), move the runtime HTTP client from `pkg/admin/internal` into `kcapi`, then add the new generic surface (`Operations`, `Invoke`, `Reload`, `Edges`, `Resolve`, `Neighbors`, `ListAll`) and two new CLI commands (`invoke`, `graph`). Finally move the declarative fetch/apply flows into `pkg/manifest` and delete `pkg/admin`.

**Tech Stack:** Go 1.26, pb33f/libopenapi (spec parsing, already vendored), stretchr/testify, kong (CLI), `httptest` fakes for tests.

**Spec:** `docs/design/kcapi-library.md` (read it first; this plan argues from it).

**Out of scope (separate plan, written after this one lands):** the MCP server (design section of the spec). Its tools map 1:1 onto the `kcapi` API created here.

## Global Constraints

- Module path: `github.com/thedataflows/keycloak-cli`. New code lives in package `kcapi` at `pkg/kcapi/`.
- No generated typed models: `pkg/models/models.gen.go` is deleted in Task 1 and stays deleted.
- Every new behavior is written test-first (TDD). Existing tests are migrated, not rewritten, wherever possible.
- All commands run from the repo root. Build gate: `go build ./...`; test gate: `go test ./...` (integration tests env-gated, skipped by default — keep it that way).
- `pkg/manifest` may import `pkg/kcapi`; `pkg/kcapi` must never import `pkg/manifest` (cycle guard, enforced from Task 3).
- Backward compatibility only matters for CLI behavior, not for Go API of `pkg/admin` (it is deleted in Task 13).
- Commit after every task, message style: existing git log (plain imperative, no scope prefixes — verify with `git log --oneline -5` before first commit).

---

### Task 1: Phase-0 deletions (audit fixes)

**Files:**
- Delete: `pkg/models/models.gen.go`, `cmd/gen-rel-schema/main.go`, `cmd/gen-rel-schema/main_test.go`, `generated/relationships.json`
- Verify no other file references the deleted paths.

**Interfaces:**
- Consumes: nothing.
- Produces: a repo with no `pkg/models` package, no `cmd/gen-rel-schema` binary, and no stale `generated/relationships.json`. Later tasks never reference them.

- [ ] **Step 1: Confirm nothing imports the doomed packages**

Run:
```bash
grep -rn 'pkg/models' --include='*.go' . | grep -v vendor/ | grep -v '^./pkg/models/'
grep -rn 'gen-rel-schema' --include='*.go' . | grep -v vendor/ | grep -v '^./cmd/gen-rel-schema/'
grep -rn 'relationships.json' --include='*.go' . | grep -v vendor/
```
Expected: no hits (or only hits inside the files being deleted). If a test references `generated/relationships.json`, update that test in this task to use inline data instead.

- [ ] **Step 2: Delete the files**

```bash
git rm -r pkg/models cmd/gen-rel-schema
git rm generated/relationships.json
```
Keep `generated/realm.json` (it is the sample output of the `generate` command).

- [ ] **Step 3: Build and test**

Run: `go build ./... && go test ./...`
Expected: PASS (zero failures; vet clean).

- [ ] **Step 4: Commit**

```bash
git commit -m "delete dead code: generated models, gen-rel-schema tool, stale relationships.json"
```

---

### Task 2: Absorb `pkg/catalog` into `pkg/kcapi` (mechanical move)

**Files:**
- Move: every file in `pkg/catalog/` (including `pkg/catalog/internal/`) → `pkg/kcapi/`
- Modify: all importers of `github.com/thedataflows/keycloak-cli/pkg/catalog` (16 files across `cmd/`, `pkg/admin/`, `pkg/admin/internal/`, `pkg/realmgen/internal/` — `cmd/gen-rel-schema` is gone after Task 1).

**Interfaces:**
- Consumes: the catalog API unchanged (`catalog.Spec`, `NewSpec`, `NewSpecFromBytes`, `OperationContract`, `ResourceContract`, `ResourceIdentity`, `ForEachOperation`, `OperationContract(path, method)`, `ResourceContracts()`, `ValidateSchema`, `ValidateOperationRequest/Response`, `ResourceIdentities`, `PlaceholderToResourceType`, `VolatileFields`, `WriteOnlyFields`).
- Produces: identical API with package name `kcapi`, import path `github.com/thedataflows/keycloak-cli/pkg/kcapi`. All later tasks build in this package. `pkg/catalog` no longer exists.

- [ ] **Step 1: Move files and rename the package**

```bash
git mv pkg/catalog pkg/kcapi
# rewrite package clauses (all files are `package catalog` or `package catalog_internal` style)
grep -rl '^package catalog' pkg/kcapi | xargs sed -i 's/^package catalog$/package kcapi/'
grep -rl 'catalog/internal' pkg/kcapi | xargs sed -i 's|catalog/internal|kcapi/internal|g'
git mv pkg/kcapi/internal/testdata 2>/dev/null || true   # only if it exists as its own dir path reference
```
Then rewrite the import path in all importers:
```bash
grep -rl 'thedataflows/keycloak-cli/pkg/catalog' --include='*.go' . | grep -v vendor/ | xargs sed -i 's|thedataflows/keycloak-cli/pkg/catalog|thedataflows/keycloak-cli/pkg/kcapi|g'
```
Fix any remaining `catalog.` qualified identifiers in moved/importer files:
```bash
grep -rl 'pkg/kcapi' --include='*.go' . | grep -v vendor/ | xargs sed -i 's/\bcatalog\./kcapi./g'
```
(Some importer files may alias the import as `catalog "…/pkg/catalog"`; remove those aliases so the default `kcapi` name applies.)

- [ ] **Step 2: Build and test**

Run: `go build ./... && go test ./...`
Expected: PASS. The moved test files (`catalog_test.go` etc.) now run as `package kcapi`.

- [ ] **Step 3: Rename test files for clarity**

```bash
cd pkg/kcapi && for f in *_test.go; do :; done   # optional: keep names as-is; names are historical but harmless
```
Keep original file names (cheaper diff; renaming is cosmetic).

- [ ] **Step 4: Commit**

```bash
git add -A && git commit -m "absorb pkg/catalog into pkg/kcapi (mechanical move, API unchanged)"
```

---

### Task 3: Hoist `manifest.Resource` and `manifest.RelationshipOperation` into `kcapi`

`pkg/kcapi/contracts.go` (ex-catalog) imports `pkg/manifest` for `Resource`/`RelationshipOperation`. Task 13 needs `pkg/manifest` to import `kcapi` — that would cycle. Fix now: own the two core types in `kcapi`, alias them in `manifest`.

**Files:**
- Create: `pkg/kcapi/types.go`
- Modify: `pkg/manifest/manifest.go` (replace the two struct declarations with aliases), all `pkg/kcapi` files referencing `manifest.Resource` / `manifest.RelationshipOperation`
- Test: `pkg/kcapi/types_test.go`, `pkg/manifest/manifest_test.go` (existing, must stay green)

**Interfaces:**
- Consumes: exact current declarations of `manifest.Resource`, `manifest.RelationshipOperation` (read `pkg/manifest/manifest.go:21-40` first — move them verbatim, fields unchanged).
- Produces: `kcapi.Resource`, `kcapi.RelationshipOperation` (struct definitions), and in manifest: `type Resource = kcapi.Resource`, `type RelationshipOperation = kcapi.RelationshipOperation`. Type identities are preserved across the alias, so all existing call sites compile unchanged.

- [ ] **Step 1: Write the failing test**

Create `pkg/kcapi/types_test.go`:
```go
package kcapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceRoundTrip(t *testing.T) {
	r := Resource{Type: "users", Data: map[string]interface{}{"username": "alice"}}
	assert.Equal(t, "users", r.Type)
	assert.Equal(t, "alice", r.Data["username"])
}
```
(The struct does not exist yet, so this fails to compile — that is the failure.)

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/kcapi/ -run TestResourceRoundTrip`
Expected: FAIL — `undefined: Resource`.

- [ ] **Step 3: Move the declarations**

Cut `Resource` and `RelationshipOperation` struct declarations verbatim from `pkg/manifest/manifest.go` into new `pkg/kcapi/types.go` (package `kcapi`). In `pkg/manifest/manifest.go` replace them with:
```go
// Resource and RelationshipOperation are owned by pkg/kcapi since the
// kcapi library absorbed the catalog; aliased here for API stability.
type Resource = kcapi.Resource

type RelationshipOperation = kcapi.RelationshipOperation
```
(add the `kcapi` import; manifest already parses these types so no further changes). In `pkg/kcapi`, rewrite `manifest.Resource` → `Resource` and `manifest.RelationshipOperation` → `RelationshipOperation`, then delete the now-unused `manifest` import from every `pkg/kcapi` file (`goimports -w pkg/kcapi` or by hand).

- [ ] **Step 4: Run tests to verify they pass**

Run: `go build ./... && go test ./pkg/kcapi/ ./pkg/manifest/ ./cmd/`
Expected: PASS everywhere.

- [ ] **Step 5: Cycle guard**

Add to the top of `pkg/kcapi/types.go` a comment: `// kcapi must not import pkg/manifest (manifest aliases kcapi core types).`
Run: `go list -deps ./pkg/manifest | grep keycloak-cli/pkg/kcapi` — expected: prints the kcapi path (manifest depends on kcapi), and `go list -deps ./pkg/kcapi | grep keycloak-cli/pkg/manifest` — expected: empty.

- [ ] **Step 6: Commit**

```bash
git add -A && git commit -m "own Resource/RelationshipOperation in kcapi, alias in manifest (breaks import cycle)"
```

---

### Task 4: Move the runtime client into `kcapi`; unify errors as `kcapi.Error`

**Files:**
- Move: `pkg/admin/internal/client.go` → `pkg/kcapi/runtime.go`; `pkg/admin/internal/http_error.go` → merged into `pkg/kcapi/errors.go`; `pkg/admin/errors.go` → merged into `pkg/kcapi/errors.go`
- Modify: `pkg/admin/service_impl.go`, `pkg/admin/apply.go`, `pkg/admin/fetch.go` (switch from `internal.client` to `kcapi.RuntimeClient`)
- Delete: `pkg/admin/internal/` (empty after the move), `pkg/admin/errors.go`
- Test: `pkg/admin/internal/client_test.go` → `pkg/kcapi/runtime_test.go` (move, keep cases); `pkg/admin/apply_typed_error_test.go` stays green (assertions switch to `kcapi.Error`)

**Interfaces:**
- Consumes: `RuntimeClient` API as-is (`NewRuntimeClient(config Config, tokens TokenProvider)`, `Spec()`, `FetchResources`, `FetchResourcesWithParent`, `FetchPathCollection`, `FetchResource`, `CreateResource`, `UpdateResource`, `DeleteResource`, `ExecuteRelationship`), admin `ErrorKind`/`classifyError`/`kindFromStatus`.
- Produces (later tasks rely on these exact names):
```go
package kcapi

type ErrorKind string

const (
	KindNotFound   ErrorKind = "not_found"
	KindConflict   ErrorKind = "conflict"
	KindValidation ErrorKind = "validation"
	KindAuth       ErrorKind = "auth"
	KindServer     ErrorKind = "server"
	KindNetwork    ErrorKind = "network"
)

var ErrNotFound = errors.New("kcapi: not found")
var ErrConflict = errors.New("kcapi: conflict")
var ErrValidation = errors.New("kcapi: validation failed")
var ErrAuth = errors.New("kcapi: authentication failed")

type Error struct {
	Kind   ErrorKind
	Op     string // e.g. "invoke getUser", "fetch users"
	Status int    // HTTP status, 0 for non-HTTP
	Body   string // truncated response body
	Err    error  // wrapped cause
}

func (e *Error) Error() string
func (e *Error) Unwrap() error
func (e *Error) Is(target error) bool // maps Kind→sentinel: KindNotFound→ErrNotFound etc.
func classifyError(err error, statusCode int, op string) *Error
```

- [ ] **Step 1: Move the files**

```bash
git mv pkg/admin/internal/client.go pkg/kcapi/runtime.go
git mv pkg/admin/internal/client_test.go pkg/kcapi/runtime_test.go
```
In `pkg/kcapi/runtime.go` + `runtime_test.go`: change `package client` → `package kcapi`, drop unexported-ness friction (nothing to change — `RuntimeClient` is already exported), and make the local `Config` type `RuntimeConfig` (it clashes with the future public `Config`):
```bash
sed -i 's/\btype Config struct {/type RuntimeConfig struct {/; s/NewRuntimeClient(config Config,/NewRuntimeClient(config RuntimeConfig,/' pkg/kcapi/runtime.go pkg/kcapi/runtime_test.go
```

- [ ] **Step 2: Write the failing error-unification test**

Create `pkg/kcapi/errors_test.go`:
```go
package kcapi

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClassifyErrorKinds(t *testing.T) {
	tests := []struct {
		status int
		want   ErrorKind
		sent   error
	}{
		{http.StatusNotFound, KindNotFound, ErrNotFound},
		{http.StatusConflict, KindConflict, ErrConflict},
		{http.StatusUnauthorized, KindAuth, ErrAuth},
		{http.StatusBadRequest, KindValidation, nil},
		{http.StatusInternalServerError, KindServer, nil},
	}
	for _, tt := range tests {
		err := classifyError(errors.New("boom"), tt.status, "invoke getUser")
		assert.Equal(t, tt.want, err.Kind)
		assert.Equal(t, "invoke getUser", err.Op)
		assert.Equal(t, tt.status, err.Status)
		if tt.sent != nil {
			assert.ErrorIs(t, err, tt.sent)
		}
	}
}

func TestErrorIsSentinel(t *testing.T) {
	err := &Error{Kind: KindNotFound, Op: "op"}
	assert.ErrorIs(t, err, ErrNotFound)
	assert.NotErrorIs(t, err, ErrConflict)
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./pkg/kcapi/ -run 'TestClassify|TestErrorIs'`
Expected: FAIL — `undefined: classifyError` / `undefined: Error`.

- [ ] **Step 4: Implement `kcapi/errors.go`**

Create `pkg/kcapi/errors.go` — port the classification logic from the old `pkg/admin/errors.go` (`classifyError`, `kindFromStatus`, `formatHTTPMessage`, `Error.message`) but rename its `ErrorKind` constants to the six `Kind*` values above and add:
```go
func (e *Error) Is(target error) bool {
	switch target {
	case ErrNotFound:
		return e.Kind == KindNotFound
	case ErrConflict:
		return e.Kind == KindConflict
	case ErrValidation:
		return e.Kind == KindValidation
	case ErrAuth:
		return e.Kind == KindAuth
	}
	return errors.Is(e.Err, target)
}
```
Delete `pkg/admin/errors.go` and `pkg/admin/internal/http_error.go`; in `pkg/admin`, replace references to `admin.Error`/`admin.ErrorKind` with `kcapi.Error`/`kcapi.ErrorKind` (the admin package already imports `kcapi` after Task 2).

- [ ] **Step 5: Run all tests**

Run: `go build ./... && go test ./...`
Expected: PASS, including the moved `runtime_test.go` and `apply_typed_error_test.go` (update its imports only, not assertions).

- [ ] **Step 6: Commit**

```bash
git add -A && git commit -m "move runtime client into kcapi, unify error taxonomy as kcapi.Error"
```

---

### Task 5: `kcapi.Client` + discovery (`Operations`)

**Files:**
- Create: `pkg/kcapi/kcapi.go` (Client, Config, New), `pkg/kcapi/discovery.go` (Operation, OpFilter, Operations)
- Test: `pkg/kcapi/kcapi_test.go`, `pkg/kcapi/discovery_test.go`

**Interfaces:**
- Consumes: `kcapi.NewSpec`/`NewSpecFromBytes` (Task 2), `kcapi.NewRuntimeClient` + `RuntimeClient` (Task 4), `auth.Service` (existing package).
- Produces:
```go
type SpecSource struct {
	Path string // file path
	URL  string // fetched once at New; Reload refetches
	Raw  []byte // inline bytes
}

type Credentials struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
}

type Config struct {
	BaseURL     string
	Spec        SpecSource
	Credentials Credentials
	Auth        auth.Service // optional; overrides Credentials
	Timeout     time.Duration
	HTTP        *http.Client // optional override
}

func New(cfg Config) (*Client, error)

type Verb string

const (
	Get    Verb = "GET"
	Post   Verb = "POST"
	Put    Verb = "PUT"
	Delete Verb = "DELETE"
	Patch  Verb = "PATCH"
)

type Param struct {
	Name     string
	In       string // "path" | "query" | "body"
	Required bool
}

type Operation struct {
	ID       string // operationId, may be empty
	Resource string // inferred from path, e.g. "users"
	Verb     Verb
	Path     string // raw template, e.g. "/admin/realms/{realm}/users/{id}"
	Summary  string
	Tags     []string
	Params   []Param
}

type OpFilter struct {
	Resource string
	Method   Verb
	Tag      string
	Search   string // substring match on ID/Path/Summary, case-insensitive
}

func (c *Client) Operations(filter OpFilter) ([]Operation, error)
func (c *Client) Spec() *Spec // re-exported for manifest/CLI compatibility
```

- [ ] **Step 1: Write the failing test**

Create `pkg/kcapi/discovery_test.go`:
```go
package kcapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperationsListsAllSpecOps(t *testing.T) {
	c := testClient(t) // helper below loads the real repo spec
	ops, err := c.Operations(OpFilter{})
	require.NoError(t, err)
	assert.Greater(t, len(ops), 100, "real Keycloak spec has hundreds of operations")

	users := filterOps(ops, func(o Operation) bool { return o.Resource == "users" })
	assert.NotEmpty(t, users)
}

func TestOperationsFilterByVerbAndSearch(t *testing.T) {
	c := testClient(t)
	ops, err := c.Operations(OpFilter{Resource: "users", Method: Get})
	require.NoError(t, err)
	for _, o := range ops {
		assert.Equal(t, Get, o.Verb)
	}
	bySearch, err := c.Operations(OpFilter{Search: "getuser"})
	require.NoError(t, err)
	found := false
	for _, o := range bySearch {
		if o.ID == "getUser" {
			found = true
			assert.Equal(t, "/admin/realms/{realm}/users/{id}", o.Path)
		}
	}
	assert.True(t, found, "search by operationId substring should find getUser")
}
```
Create the helper in `pkg/kcapi/kcapi_test.go` (package `kcapi`):
```go
package kcapi

import (
	"os"
	"testing"
)

func testClient(t *testing.T) *Client {
	t.Helper()
	spec, err := os.ReadFile("../../keycloak-oapi/26.6.2.spec.json")
	if err != nil {
		t.Skipf("repo spec not available: %v", err)
	}
	c, err := New(Config{BaseURL: "http://test", Spec: SpecSource{Raw: spec}})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/kcapi/ -run TestOperations`
Expected: FAIL — `undefined: New` / `Operations`.

- [ ] **Step 3: Implement `kcapi.go` and `discovery.go`**

`pkg/kcapi/kcapi.go`:
```go
// Package kcapi is a spec-driven Keycloak client: every operation in the
// supplied OpenAPI spec is discoverable and invocable generically, plus
// relationship-graph queries over Keycloak objects.
package kcapi

import (
	"net/http"
	"time"

	"github.com/thedataflows/keycloak-cli/pkg/auth"
)

type Client struct {
	spec   *Spec
	runtime *RuntimeClient
	http   *http.Client
	specFn func() ([]byte, error) // current raw-spec loader, swapped by Reload
	mu     sync.RWMutex
}

func New(cfg Config) (*Client, error) {
	spec, err := loadSpec(cfg.Spec) // dispatch on Path/URL/Raw; URL uses cfg.HTTP or http.DefaultClient
	if err != nil {
		return nil, err
	}
	tokens := cfg.Auth
	if tokens == nil {
		tokens, err = authFromCredentials(cfg.Credentials) // thin wrapper picking password-grant or client-credentials from auth package
		if err != nil {
			return nil, err
		}
	}
	rt, err := NewRuntimeClient(RuntimeConfig{
		BaseURL: cfg.BaseURL,
		Timeout: cfg.Timeout,
		Spec:    spec,
	}, tokens)
	if err != nil {
		return nil, err
	}
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	httpClient := cfg.HTTP
	if httpClient == nil {
		httpClient = &http.Client{Timeout: timeout}
	}
	return &Client{spec: spec, runtime: rt, http: httpClient}, nil
}

func (c *Client) Spec() *Spec { return c.spec }
```
Note: `NewRuntimeClient` currently owns its own HTTP client — adapt `RuntimeConfig` minimally so the token-provider plumbing is reused but the `*http.Client` can be injected; keep the diff small.

`pkg/kcapi/discovery.go`:
```go
package kcapi

import (
	"strings"
)

func (c *Client) Operations(filter OpFilter) ([]Operation, error) {
	var out []Operation
	c.spec.ForEachOperation(func(path, method string, op *v3.Operation, item *v3.PathItem) {
		o := buildOperation(path, method, op)
		if matchOp(filter, o) {
			out = append(out, o)
		}
	})
	return out, nil
}

func buildOperation(path, method string, op *v3.Operation) Operation {
	o := Operation{
		ID:       op.OperationID,
		Verb:     Verb(method),
		Path:     path,
		Summary:  op.Summary,
		Resource: inferResourceTypeFromPath(path), // already in kcapi (ex-catalog contracts.go)
	}
	for _, tag := range op.Tags {
		o.Tags = append(o.Tags, tag)
	}
	for _, p := range append(op.Parameters, item.Parameters...) {
		if p == nil {
			continue
		}
		o.Params = append(o.Params, Param{Name: p.Name, In: p.In, Required: p.Required})
	}
	return o
}

func matchOp(f OpFilter, o Operation) bool {
	if f.Resource != "" && o.Resource != f.Resource {
		return false
	}
	if f.Method != "" && o.Verb != f.Method {
		return false
	}
	if f.Tag != "" && !containsString(o.Tags, f.Tag) {
		return false
	}
	if f.Search != "" {
		hay := strings.ToLower(o.ID + " " + o.Path + " " + o.Summary)
		if !strings.Contains(hay, strings.ToLower(f.Search)) {
			return false
		}
	}
	return true
}
```
(Use the libopenapi v3 types already imported by `contracts.go`; drop the placeholder `if` block in `buildOperation` — params come solely from the two slices.)

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./pkg/kcapi/ -run TestOperations -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add -A && git commit -m "kcapi: Client construction and generic operation discovery"
```

---

### Task 6: Generic `Invoke`

**Files:**
- Create: `pkg/kcapi/invoke.go`
- Test: `pkg/kcapi/invoke_test.go`

**Interfaces:**
- Consumes: `Client` (Task 5), `RuntimeClient` request plumbing + `buildPathWithOperation` (Task 4), `OperationContract` validation (`validateOperationInput` is unexported — call the exported `ValidateOperationRequest` wrapper), `classifyError` (Task 4).
- Produces:
```go
type P map[string]string

type Call struct {
	Op       string // operationId; XOR with Resource+Verb
	Resource string
	Verb     Verb
	Realm    string // convenience: fills the {realm} path param if present and not in Params
	Params   P      // path + query params
	Body     interface{} // JSON-serializable; nil for GET/DELETE
}

func (c *Client) Invoke(ctx context.Context, call Call) (json.RawMessage, error)
```

- [ ] **Step 1: Write the failing test**

Create `pkg/kcapi/invoke_test.go`:
```go
package kcapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvokeByOperationID(t *testing.T) {
	var gotPath, gotMethod string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath, gotMethod = r.URL.Path, r.Method
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"u-1","username":"alice"}`))
	}))
	defer srv.Close()

	c, err := New(Config{
		BaseURL: srv.URL,
		Spec:    SpecSource{Raw: testSpec(t)},
		Auth:    staticTokenProvider("test-token"), // helper: TokenProvider returning a fixed token
	})
	require.NoError(t, err)

	out, err := c.Invoke(context.Background(), Call{Op: "getUser", Realm: "demo", Params: P{"id": "u-1"}})
	require.NoError(t, err)
	assert.Equal(t, http.MethodGet, gotMethod)
	assert.Contains(t, gotPath, "/admin/realms/demo/users/u-1")

	var user map[string]interface{}
	require.NoError(t, json.Unmarshal(out, &user))
	assert.Equal(t, "alice", user["username"])
}

func TestInvokeValidationErrors(t *testing.T) {
	c, err := New(Config{BaseURL: "http://test", Spec: SpecSource{Raw: testSpec(t)}})
	require.NoError(t, err)

	_, err = c.Invoke(context.Background(), Call{Op: "no-such-op"})
	assert.ErrorIs(t, err, ErrValidation)

	_, err = c.Invoke(context.Background(), Call{Op: "getUser"}) // missing required id
	assert.ErrorIs(t, err, ErrValidation)

	_, err = c.Invoke(context.Background(), Call{Resource: "users", Verb: Get, Realm: "demo", Params: P{"id": "x"}})
	assert.ErrorIs(t, err, ErrValidation) // mutually exclusive modes not the failure here; getUser has no such overload → unknown op; see note
}
```
Note on the third case: `Resource+Verb` resolution must find the users-GET-by-id operation; with `id` present it should succeed against the fake server. Split that into `TestInvokeByResourceAndVerb` hitting an `httptest` server (mirror the first test; assert `POST /admin/realms/demo/users` works for `Call{Resource: "users", Verb: Post, Realm: "demo", Body: map[string]any{"username":"bob"}}`). Keep exactly one resolution mode per test.

Also add `testSpec(t)` helper in `kcapi_test.go` (reads the repo spec once per test via `testClient`'s loader) and `staticTokenProvider` helper in `runtime_test.go` (wraps the existing test token fake from the moved `client_test.go`).

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/kcapi/ -run TestInvoke`
Expected: FAIL — `undefined: Call` / `Invoke`.

- [ ] **Step 3: Implement `invoke.go`**

```go
package kcapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) Invoke(ctx context.Context, call Call) (json.RawMessage, error) {
	op, path, err := c.resolveCall(call)
	if err != nil {
		return nil, &Error{Kind: KindValidation, Op: callLabel(call), Err: err}
	}
	if err := c.validateCallParams(op, call); err != nil {
		return nil, &Error{Kind: KindValidation, Op: callLabel(call), Err: err}
	}
	body, err := c.marshalBody(op, call)
	if err != nil {
		return nil, &Error{Kind: KindValidation, Op: callLabel(call), Err: err}
	}
	req, err := http.NewRequestWithContext(ctx, string(op.verb), c.runtime.BaseURL()+path, bodyReader(body))
	if err != nil {
		return nil, &Error{Kind: KindNetwork, Op: callLabel(call), Err: err}
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if err := c.runtime.Authorize(ctx, req); err != nil {
		return nil, classifyError(err, 0, callLabel(call))
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, &Error{Kind: KindNetwork, Op: callLabel(call), Err: err}
	}
	defer resp.Body.Close()
	raw, readErr := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if resp.StatusCode >= 400 {
		return nil, classifyError(errors.New(string(raw)), resp.StatusCode, callLabel(call))
	}
	if readErr != nil {
		return nil, &Error{Kind: KindNetwork, Op: callLabel(call), Err: readErr}
	}
	if len(bytes.TrimSpace(raw)) == 0 {
		return nil, nil // 204-style empty bodies are a valid result
	}
	return json.RawMessage(raw), nil
}

// resolveCall finds the operation by ID or Resource+Verb and instantiates
// path placeholders from Params (+Realm), erroring on unknown or missing ones.
func (c *Client) resolveCall(call Call) (resolvedOp, string, error)
```
Implementation notes (concrete requirements — define these local helpers exactly like this):
0. `type resolvedOp struct { verb Verb; path string; contract OperationContract }`; `func callLabel(call Call) string` returns `"invoke <opID>"` or `"invoke <verb> <resource>"`; `func marshalBody(op resolvedOp, call Call) ([]byte, error)` returns `nil, nil` for a nil Body and otherwise `json.Marshal(call.Body)`, erroring when the operation contract has no request schema; `func bodyReader(b []byte) io.Reader` returns `bytes.NewReader(b)` or `http.NoBody` when empty.
1. `resolveCall`: iterate `c.Spec().ForEachOperation`; match `op.OperationID == call.Op` (when `call.Op != ""` and `Resource==""`), else match `inferResourceTypeFromPath(path) == call.Resource && method == string(call.Verb)` and the path has exactly the placeholders fillable from `Params`+`Realm` (choose the shortest path on ties). Zero matches or >1 distinct operation for the ID → validation error.
2. Path building reuses `RuntimeClient.buildPathWithOperation` semantics — extract that logic into a shared helper `instantiatePath(template string, values map[string]string) (string, error)` in `runtime.go` so both manifest flows and `Invoke` use one implementation (no duplicate placeholder logic).
3. `Realm` fills the `{realm}` placeholder when present and not already in `Params`.
4. `validateCallParams`: every `path` placeholder must be non-empty; required query params present; use the existing `OperationContract` (via `c.Spec().OperationContract(path, method)`) + `validateOperationInput` logic — if it is unexported in contracts.go, export a small wrapper `ValidateCall(spec *Spec, path, method string, call Call) error` in `contracts.go` reusing the internals.
5. `Body` → `json.Marshal`; forbidden when the operation contract has no request schema (return validation error instead of silently dropping).
6. Expose `func (r *RuntimeClient) BaseURL() string` and `func (r *RuntimeClient) Authorize(ctx context.Context, req *http.Request) error` from `runtime.go` (thin extractions of existing private logic in `newAuthRequest`).

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./pkg/kcapi/ -run TestInvoke -v && go test ./...`
Expected: PASS (new tests + no regressions).

- [ ] **Step 5: Commit**

```bash
git add -A && git commit -m "kcapi: generic Invoke by operationId or resource+verb"
```

---

### Task 7: `Reload` — atomic spec swap

**Files:**
- Modify: `pkg/kcapi/kcapi.go` (Reload), `pkg/kcapi/discovery.go`, `pkg/kcapi/invoke.go` (read spec under RLock)
- Test: `pkg/kcapi/reload_test.go`

**Interfaces:**
- Consumes: `SpecSource.URL` loader from Task 5; `sync.RWMutex` on `Client`.
- Produces: `func (c *Client) Reload(ctx context.Context) error` — refetches the spec (URL/Path source; `Raw` reloads the same bytes), parses it, builds a fresh `*Spec` + `*RuntimeClient`, swaps both under one write lock. In-flight calls keep the old pair (they read under RLock and use the values they captured).

- [ ] **Step 1: Write the failing test**

Create `pkg/kcapi/reload_test.go`:
```go
package kcapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReloadSwapsCatalog(t *testing.T) {
	var specBody atomic.Value
	specBody.Store(mustReadRepoSpec(t)) // start from the full repo spec
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(specBody.Load().([]byte))
	}))
	defer srv.Close()

	c, err := New(Config{BaseURL: "http://test", Spec: SpecSource{URL: srv.URL}})
	require.NoError(t, err)

	before, err := c.Operations(OpFilter{Search: "getUser"})
	require.NoError(t, err)
	require.NotEmpty(t, before)

	specBody.Store([]byte(`{"openapi":"3.0.0","paths":{}}`)) // spec update drops everything
	require.NoError(t, c.Reload(context.Background()))
	after, err := c.Operations(OpFilter{})
	require.NoError(t, err)
	assert.Empty(t, after)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/kcapi/ -run TestReload`
Expected: FAIL — `undefined: Reload` (or `c.specFn` unused path).

- [ ] **Step 3: Implement**

In `kcapi.go` add:
```go
func (c *Client) Reload(ctx context.Context) error {
	raw, err := c.loadSpecBytes(ctx)
	if err != nil {
		return err
	}
	spec, err := NewSpecFromBytes(raw)
	if err != nil {
		return err
	}
	rt, err := NewRuntimeClient(RuntimeConfig{BaseURL: c.runtime.BaseURL(), Timeout: c.timeout, Spec: spec}, c.runtime.tokens)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.spec, c.runtime = spec, rt
	return nil
}
```
Convert `Operations` and `Invoke` to snapshot `spec, runtime := c.snapshot()` (new RLock helper) at entry and use the snapshot throughout. Store `timeout` and `tokens` on `Client` at `New` (fields `timeout time.Duration`, `tokens auth.Service`) so Reload can rebuild. If `RuntimeClient.tokens` is not accessible, keep a `Client.tokens` field instead and pass it to `NewRuntimeClient`.

- [ ] **Step 4: Run tests**

Run: `go test ./pkg/kcapi/ -run TestReload -v && go test ./... && go test -race ./pkg/kcapi/`
Expected: PASS, race detector clean.

- [ ] **Step 5: Commit**

```bash
git add -A && git commit -m "kcapi: atomic spec Reload"
```

---

### Task 8: Generalized relationship edges

**Files:**
- Create: `pkg/kcapi/edges.go`
- Modify: `pkg/kcapi/relationship_registry.go` (demote to display-name overrides only; drop any pattern data that duplicates inference)
- Test: `pkg/kcapi/edges_test.go`

**Interfaces:**
- Consumes: `Spec.ForEachOperation`, `PlaceholderToResourceType` (identity.go), the existing `relationship_patterns.go` inference (generalize it).
- Produces:
```go
type Edge struct {
	Parent       string // resource type, e.g. "realms"
	Child        string // resource type, e.g. "roles"
	Path         string // template with both placeholders, e.g. "/admin/realms/{realm}/roles/{role-keepers}/users"
	Cardinality  string // "one" if child path has a single trailing identifier, "many" if collection
	CascadeOwner bool   // true when deleting the parent is the documented way to remove children (heuristic: child has no standalone DELETE at top level)
	Display      string // optional curated name from the registry overrides
}

func (c *Client) Edges() []Edge // all edges implied by the live spec, sorted by (Parent, Child, Path)
```
Inference rule (port + generalize `relationship_patterns.go`): a path is an edge when it contains ≥2 distinct resource placeholders AND the path prefix up to the second placeholder is a valid parent-collection path. Cardinality from the trailing segment's shape (`isCollectionEndpoint`, already in contracts.go).

- [ ] **Step 1: Write the failing test**

Create `pkg/kcapi/edges_test.go`:
```go
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
	// but that path has only {realm}; the pair must come from paths like
	// /admin/realms/{realm}/users/{user-identifier}/role-mappings/...
	byPair := map[string]Edge{}
	for _, e := range edges {
		byPair[e.Parent+"->"+e.Child] = e
	}
	assert.Contains(t, byPair, "realms->users")

	for _, e := range edges {
		assert.NotEqual(t, e.Parent, e.Child, "self-edges are noise: %v", e)
		assert.Contains(t, e.Path, "{realm}")
	}
}

func TestEdgeDisplayOverride(t *testing.T) {
	c := testClient(t)
	for _, e := range c.Edges() {
		if e.Parent == "realms" && e.Child == "clients" {
			assert.NotEmpty(t, e.Display, "registry override should name this common edge")
			return
		}
	}
	t.Skip("no realms->clients edge in this spec version")
}
```
(If the second test's assumption about the override table contents is wrong, adjust the expected pair to one actually present in `relationship_registry.go` after reading it — the point is that overrides flow into `Edge.Display`.)

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/kcapi/ -run TestEdge`
Expected: FAIL — `undefined: Edges`.

- [ ] **Step 3: Implement `edges.go`**

```go
package kcapi

import "sort"

func (c *Client) Edges() []Edge {
	spec, _ := c.snapshot() // snapshot helper from Task 7
	edges := inferEdges(spec) // port/generalize relationship_patterns.go walk over ForEachOperation
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
```
`inferEdges` walks every path template, splits it into segments, tracks placeholder→resource-type via the existing `PlaceholderToResourceType` map, and records an edge per (parent placeholder, child placeholder) pair where the prefix path exists as an operation in the spec. Delete from `relationship_patterns.go` whatever this replaces; move genuinely-curated names into a small `displayOverrides map[string]string` keyed `"parent->child"` in `relationship_registry.go`, deleting pattern logic there.

- [ ] **Step 4: Run tests**

Run: `go test ./pkg/kcapi/ -run 'TestEdge' -v && go test ./...`
Expected: PASS including all migrated relationship tests (`relationship_registry_test.go` may shrink — that is intended; keep its cases that assert display names).

- [ ] **Step 5: Commit**

```bash
git add -A && git commit -m "kcapi: structural relationship-edge inference, registry demoted to display overrides"
```

---

### Task 9: `Resolve` and `Neighbors`

**Files:**
- Create: `pkg/kcapi/graph.go`
- Test: `pkg/kcapi/graph_test.go`

**Interfaces:**
- Consumes: `Edges()` (Task 8), `Invoke` (Task 6), identity resolution helpers from `identity.go` (`ResourceIdentities`, `IdentifierOf`, `NameOf`), the search-or-list strategy in `pkg/admin/references.go` (`locateExistingResource` — move that helper into `graph.go` or a new `resolve.go` inside kcapi, since admin dissolves in Task 13 anyway).
- Produces:
```go
type Ref struct {
	Type  string // resource type, e.g. "users"
	Name  string // human name (username, client-id, role name…)
	ID    string // explicit id; skips name lookup
	Realm string
}

type Node struct {
	Ref      Ref
	Resource Resource // the fetched object (Type + Data)
}

func (c *Client) Resolve(ctx context.Context, ref Ref) (Node, error)

type EdgeFilter struct {
	Child string // only edges with this child type
	Parent string // only edges with this parent type
}

func (c *Client) Neighbors(ctx context.Context, node Node, filter EdgeFilter) ([]Node, []Edge, error)
```
`Resolve` strategy: if `ID` set → GET the single-resource path for that type; else search (query param `search`/`name`/`username` per identity fields) and require exactly one hit; zero → `ErrNotFound`, ambiguous → validation error listing candidates. `Neighbors` instantiates each matching edge's path from the node's identifier, calls `Invoke` on the collection, and wraps each element as a `Node` (child ref `Name` filled via `NameOf`).

- [ ] **Step 1: Write the failing test**

Create `pkg/kcapi/graph_test.go`:
```go
package kcapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeKC is a minimal in-memory Keycloak: users list + role-mappings edge.
func fakeKC(t *testing.T) (*httptest.Server, *Client) {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /admin/realms/demo/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("username") {
		case "alice":
			_, _ = w.Write([]byte(`[{"id":"u-1","username":"alice"}]`))
		case "":
			_, _ = w.Write([]byte(`[{"id":"u-1","username":"alice"},{"id":"u-2","username":"bob"}]`))
		default:
			_, _ = w.Write([]byte(`[]`))
		}
	})
	mux.HandleFunc("GET /admin/realms/demo/users/u-1/role-mappings/realm", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`[{"id":"r-9","name":"offline_access"}]`))
	})
	srv := httptest.NewServer(mux)
	c, err := New(Config{BaseURL: srv.URL, Spec: SpecSource{Raw: testSpec(t)}})
	require.NoError(t, err)
	return srv, c
}

func TestResolveByName(t *testing.T) {
	srv, c := fakeKC(t)
	defer srv.Close()
	node, err := c.Resolve(context.Background(), Ref{Type: "users", Name: "alice", Realm: "demo"})
	require.NoError(t, err)
	assert.Equal(t, "u-1", node.Resource.Data["id"])
}

func TestResolveNotFoundAndAmbiguous(t *testing.T) {
	srv, c := fakeKC(t)
	defer srv.Close()
	_, err := c.Resolve(context.Background(), Ref{Type: "users", Name: "nobody", Realm: "demo"})
	assert.ErrorIs(t, err, ErrNotFound)
	_, err = c.Resolve(context.Background(), Ref{Type: "users", Realm: "demo"}) // no name -> lists all -> 2 hits
	assert.Error(t, err)                                                        // ambiguous
}

func TestNeighborsRoleMappings(t *testing.T) {
	srv, c := fakeKC(t)
	defer srv.Close()
	node, err := c.Resolve(context.Background(), Ref{Type: "users", Name: "alice", Realm: "demo"})
	require.NoError(t, err)
	nodes, edges, err := c.Neighbors(context.Background(), node, EdgeFilter{Child: "roles"})
	require.NoError(t, err)
	require.Len(t, nodes, 1)
	assert.Equal(t, "offline_access", nodes[0].Resource.Data["name"])
	assert.NotEmpty(t, edges)
}
```
Adjust route/param details to the real spec shapes (path placeholders in `26.6.2.spec.json` may use `{user-identifier}`-style names; `Neighbors` derives them from the edge `Path` template, so the test data must match the actual placeholder names — read them from the spec when writing the test).

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/kcapi/ -run 'TestResolve|TestNeighbors'`
Expected: FAIL — `undefined: Resolve` / `Neighbors`.

- [ ] **Step 3: Implement `graph.go`**

Implement per the Produces contract:
1. `resolveSingle`: build the by-id call from the resource identity's id path (existing identity inference: `ResourceIdentities()[type]` gives idParam + identity fields).
2. `resolveByName`: `Invoke` the collection path with a query param chosen from the identity fields (`username` for users, `clientId` for clients, `name` otherwise), exact-match filter client-side on the identity fields.
3. `Neighbors`: for each `Edge` where `edge.Parent == node.Ref.Type` and the filter matches, `instantiatePath(edge.Path, values)` with the node id + realm, `Invoke` it, unmarshal as `[]map[string]interface{}` (fall back to a single object for `Cardinality: "one"`), wrap each element in `Node{Ref: Ref{Type: edge.Child, Name: NameOf(...), Realm: node.RefRealm}, Resource: ...}`.
Move `locateExistingResource` (from `pkg/admin/references.go`) into `graph.go` as the shared name→id resolution core, adapted to `Invoke` instead of the manifest-specific client calls; `pkg/admin/references.go` keeps a thin delegator until Task 13 deletes it.

- [ ] **Step 4: Run tests**

Run: `go test ./pkg/kcapi/ -run 'TestResolve|TestNeighbors' -v && go test ./...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add -A && git commit -m "kcapi: Resolve and Neighbors graph queries"
```

---

### Task 10: `ListAll` pagination helper

**Files:**
- Create: `pkg/kcapi/paging.go`
- Test: `pkg/kcapi/paging_test.go`

**Interfaces:**
- Consumes: `Invoke` (Task 6); Keycloak list paging (`first`, `max` query params).
- Produces:
```go
const DefaultPageSize = 100

// ListAll pages through a list operation, following first/max, until a
// short page is returned. Call must target a collection operation.
func (c *Client) ListAll(ctx context.Context, call Call) ([]json.RawMessage, error)
```

- [ ] **Step 1: Write the failing test**

Create `pkg/kcapi/paging_test.go`:
```go
package kcapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListAllPages(t *testing.T) {
	var seenFirst []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		first := r.URL.Query().Get("first")
		seenFirst = append(seenFirst, first)
		f := 0
		fmt.Sscanf(first, "%d", &f)
		start := f
		end := f + 100
		if end > 250 {
			end = 250
		}
		users := []string{}
		for i := start; i < end; i++ {
			users = append(users, fmt.Sprintf(`{"id":"u-%d"}`, i))
		}
		_, _ = w.Write([]byte("[" + join(users) + "]"))
	}))
	defer srv.Close()

	c, err := New(Config{BaseURL: srv.URL, Spec: SpecSource{Raw: testSpec(t)}})
	require.NoError(t, err)
	out, err := c.ListAll(context.Background(), Call{Resource: "users", Verb: Get, Realm: "demo"})
	require.NoError(t, err)
	assert.Len(t, out, 250)
	assert.Equal(t, []string{"", "100", "200"}, seenFirst)
}

func join(parts []string) string { s := ""; for i, p := range parts { if i > 0 { s += "," }; s += p }; return s }
```
(Move `join` into the test file's plain helper or use `strings.Join` — prefer `strings.Join`.)

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./pkg/kcapi/ -run TestListAll`
Expected: FAIL — `undefined: ListAll`.

- [ ] **Step 3: Implement `paging.go`**

```go
package kcapi

import (
	"context"
	"encoding/json"
)

func (c *Client) ListAll(ctx context.Context, call Call) ([]json.RawMessage, error) {
	var all []json.RawMessage
	for first := 0; ; first += DefaultPageSize {
		page := call
		page.Params = cloneParams(call.Params)
		page.Params["first"] = itoa(first)
		page.Params["max"] = itoa(DefaultPageSize)
		raw, err := c.Invoke(ctx, page)
		if err != nil {
			return nil, err
		}
		var items []json.RawMessage
		if raw == nil || json.Unmarshal(raw, &items) != nil || items == nil {
			return all, nil // empty body or non-array result: nothing more to page
		}
		all = append(all, items...)
		if len(items) < DefaultPageSize {
			return all, nil
		}
	}
}
```

- [ ] **Step 4: Run tests**

Run: `go test ./pkg/kcapi/ -run TestListAll -v && go test ./...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add -A && git commit -m "kcapi: ListAll pagination helper"
```

---

### Task 11: CLI `invoke` command

**Files:**
- Create: `cmd/invoke.go`, `cmd/invoke_test.go`
- Modify: `cmd/root.go` (register `InvokeCmd` in the CLI commands struct), wherever the shared `admin.Service`/`kcapi.Client` is constructed for commands (follow the existing pattern in `cmd/fetch.go`)

**Interfaces:**
- Consumes: `kcapi.Client.Operations`, `kcapi.Client.Invoke`, `kcapi.Call`, `pkg/output` (existing JSON output helper).
- Produces:
```
keycloak-cli invoke <operation-id> --realm R [--param k=v]... [--body @file.json | --body '{"json":true}'] [--op-users X --op-id Y]
```
Use `map[string]string` kong flag (`param:`). Output: raw JSON via `pkg/output`, exit code 1 with the error message on failure.

- [ ] **Step 1: Write the failing test**

Create `cmd/invoke_test.go` (follow the command-test pattern used by `cmd/fetch_test.go` — read it first and mirror its client-injection approach):
```go
package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvokeCmdBuildsCall(t *testing.T) {
	// Unit-level: assert the command turns args into a kcapi.Call correctly.
	got := invokeCall("getUser", "demo", map[string]string{"id": "u-1"}, json.RawMessage(`{"a":1}`))
	assert.Equal(t, "getUser", got.Op)
	assert.Equal(t, "demo", got.Realm)
	assert.Equal(t, "u-1", got.Params["id"])
	require.NotNil(t, got.Body)
}
```
(`invokeCall` is the small pure function the command delegates to; the HTTP-level behavior is already covered by `pkg/kcapi` tests, so the command test only covers arg→Call mapping and output formatting.)

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./cmd/ -run TestInvokeCmd`
Expected: FAIL — `undefined: invokeCall`.

- [ ] **Step 3: Implement `cmd/invoke.go`**

```go
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

type InvokeCmd struct {
	OpID     string            `arg:"" optional:"" help:"OpenAPI operationId (use --list to search)"`
	Realm    string            `help:"Realm scope" short:"r"`
	Param    map[string]string `help:"Path/query parameters (k=v, repeatable)"`
	Body     string            `help:"Request body: inline JSON or @file"`
	Resource string            `help:"Resolve by resource type + --verb instead of operationId"`
	Verb     string            `help:"HTTP verb when resolving by resource" default:"GET"`
	List     string            `help:"List operations matching this search substring and exit"`
}

func (c *InvokeCmd) Run(_ *kong.Context, cli *CLI) error {
	client, err := cli.Kcapi() // shared constructor added to CLI in this task; mirrors how FetchCmd gets its service
	if err != nil {
		return err
	}
	if c.List != "" {
		ops, err := client.Operations(kcapi.OpFilter{Search: c.List})
		if err != nil {
			return err
		}
		return printJSON(ops)
	}
	var body json.RawMessage
	if c.Body != "" {
		raw, err := readBodyArg(c.Body) // "@" prefix -> os.ReadFile, else inline
		if err != nil {
			return err
		}
		body = raw
	}
	out, err := client.Invoke(cmdContext(), invokeCall(c.OpID, c.Realm, c.Param, body, c.Resource, c.Verb))
	if err != nil {
		return err
	}
	return printJSON(out)
}

func invokeCall(opID, realm string, params map[string]string, body json.RawMessage, resource, verb string) kcapi.Call {
	call := kcapi.Call{Op: opID, Realm: realm, Params: kcapi.P(params), Body: body, Resource: resource}
	if resource != "" {
		call.Op = ""
		call.Verb = kcapi.Verb(verb)
	}
	if body == nil {
		call.Body = nil
	}
	return call
}

func readBodyArg(arg string) (json.RawMessage, error) {
	raw := []byte(arg)
	if len(arg) > 0 && arg[0] == '@' {
		var err error
		raw, err = os.ReadFile(arg[1:])
		if err != nil {
			return nil, fmt.Errorf("read body file: %w", err)
		}
	}
	if !json.Valid(raw) {
		return nil, fmt.Errorf("body is not valid JSON")
	}
	return json.RawMessage(raw), nil
}
```
Register in `cmd/root.go`'s CLI struct: `Invoke InvokeCmd `cmd:"" help:"Invoke any Keycloak API operation from the loaded spec"`. `cli.Kcapi()` (add to root.go next to the existing service constructor): builds a `kcapi.Config` from the existing CLI flags (`SpecPath`, base URL, credentials — reuse exactly the env/flag plumbing `FetchCmd` uses today). Use the existing context/output helpers (`cmdContext`, `printJSON` — check their real names in `cmd/output.go`/`cmd/root.go` and use those).

- [ ] **Step 4: Run tests**

Run: `go test ./cmd/ -run TestInvokeCmd -v && go build ./... && go test ./...`
Expected: PASS.

- [ ] **Step 5: Smoke-test against a fake server (manual verification step)**

```bash
go run . invoke --help
```
Expected: help text lists flags; no panic.

- [ ] **Step 6: Commit**

```bash
git add -A && git commit -m "cli: invoke command for generic spec-driven API calls"
```

---

### Task 12: CLI `graph` command

**Files:**
- Create: `cmd/graph.go`, `cmd/graph_test.go`
- Modify: `cmd/root.go` (register)

**Interfaces:**
- Consumes: `kcapi.Client.Resolve`, `kcapi.Client.Neighbors`, `kcapi.Client.Edges`, `kcapi.Ref`, `kcapi.EdgeFilter`, `pkg/output`.
- Produces:
```
keycloak-cli graph edges [--parent realms] [--child roles]
keycloak-cli graph resolve <type> <name> --realm R [--id]
keycloak-cli graph neighbors <type> <name> --realm R [--child roles] [--parent realms]
```

- [ ] **Step 1: Write the failing test**

Create `cmd/graph_test.go`:
```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./cmd/ -run TestGraph`
Expected: FAIL — `undefined: graphFilter` / `graphRef`.

- [ ] **Step 3: Implement `cmd/graph.go`**

```go
package cmd

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/thedataflows/keycloak-cli/pkg/kcapi"
)

type GraphCmd struct {
	Edges struct {
		Parent string `help:"Filter by parent resource type"`
		Child  string `help:"Filter by child resource type"`
	} `cmd:"" help:"List relationship edges implied by the spec"`
	Resolve struct {
		Type  string `arg:""`
		Name  string `arg:""`
		ID    string `help:"Skip name lookup, use this id"`
		Realm string `required:""`
	} `cmd:"" help:"Resolve a Keycloak object by type + name"`
	Neighbors struct {
		Type   string `arg:""`
		Name   string `arg:""`
		Realm  string `required:""`
		Child  string `help:"Only edges with this child type"`
		Parent string `help:"Only edges with this parent type"`
	} `cmd:"" help:"List objects related to a resolved object"`
}

func (c *GraphCmd) Run(_ *kong.Context, cli *CLI) error {
	client, err := cli.Kcapi()
	if err != nil {
		return err
	}
	switch kong.Must(cli).Selected() /* or use the nested Run pattern from cmd/fetch.go */ {
	// Implement the three subcommands exactly like other multi-verb commands
	// in this repo do (read cmd/compare.go for the nesting pattern).
	}
	return nil
}

func graphFilter(parent, child string) kcapi.EdgeFilter {
	return kcapi.EdgeFilter{Parent: parent, Child: child}
}

func graphRef(typ, name, realm, id string) kcapi.Ref {
	return kcapi.Ref{Type: typ, Name: name, Realm: realm, ID: id}
}
```
Follow the repo's actual kong nesting style (mirror `cmd/compare.go`); output JSON via the shared printer. `fmt` import only if actually used.

- [ ] **Step 4: Run tests**

Run: `go test ./cmd/ -run TestGraph -v && go build ./... && go test ./...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add -A && git commit -m "cli: graph edges/resolve/neighbors commands"
```

---

### Task 13: Dissolve `pkg/admin` — manifest flows onto `kcapi`

**Files:**
- Move: `pkg/admin/fetch.go`, `fetch_relationships.go`, `apply.go`, `references.go`, `service_impl.go` (Service plumbing) → `pkg/manifest/` (as `service.go`, `fetch.go`, `apply.go`, `references.go`), `pkg/admin/apply_*.go` test files → `pkg/manifest/`
- Delete: `pkg/admin/` entirely (after moves)
- Modify: `cmd/fetch.go`, `cmd/upload.go`, `cmd/compare.go` (construct `manifest.Service` instead of `admin.Service`), `pkg/output` (switch `admin` import to `manifest`/`kcapi` as needed)

**Interfaces:**
- Consumes: `kcapi.Client` (Tasks 5–9), existing `FetchQuery`/`FetchReport`/`ApplyOptions`/`ApplyReport` types (move them with their code, renamed to `manifest.*` only by package clause — type names unchanged).
- Produces:
```go
package manifest

type Service interface {
	Spec() *kcapi.Spec
	Fetch(ctx context.Context, query FetchQuery) (FetchReport, error)
	Apply(ctx context.Context, resources []Resource, relationships []RelationshipOperation, options ApplyOptions) (ApplyReport, error)
}

func NewService(cfg Config) (Service, error) // Config mirrors old admin.Config; Auth field kept
```
The service internally holds a `*kcapi.Client`; all `RuntimeClient` calls stay (RuntimeClient lives in kcapi now). `RealmCache` (wherever it lives in admin — locate with `grep -rn 'RealmCache\|realmCache' pkg/admin`) moves into `pkg/manifest`.

- [ ] **Step 1: Move files and rename package**

```bash
git mv pkg/admin/fetch.go pkg/admin/fetch_relationships.go pkg/admin/apply.go pkg/admin/references.go pkg/admin/service_impl.go pkg/manifest/
git mv pkg/admin/fetch_test.go pkg/admin/apply_conflict_test.go pkg/admin/apply_created_id_test.go pkg/admin/apply_reconcile_test.go pkg/admin/apply_reconcile_e2e_test.go pkg/admin/apply_typed_error_test.go pkg/admin/auth_injection_test.go pkg/admin/admin_test.go pkg/admin/admin_integration_test.go pkg/admin/references_test.go pkg/manifest/ 2>/dev/null || true
```
In moved files: `package admin` → `package manifest`; rewrite `admin.` qualifiers; where moved code calls `s.runtime`, keep the field but type it `*kcapi.RuntimeClient` (already true after Task 4). Merge old `admin.Config` + `newService` into `pkg/manifest/service.go` building a `kcapi.Client` first, then a thin service around it.

- [ ] **Step 2: Fix importers**

`cmd/fetch.go`, `cmd/upload.go`, `cmd/compare.go`, `pkg/output`: replace `admin.New`/`admin.Service` with `manifest.NewService`/`manifest.Service` (or `kcapi` where output helpers only need error kinds). Delete `pkg/admin/admin.go` and the emptied `pkg/admin/` directory.

- [ ] **Step 3: Full build + tests**

Run: `go build ./... && go test ./... && go vet ./...`
Expected: PASS. The migrated apply/fetch tests must pass **unchanged in assertions** — if one needs a logic change to pass, stop and fix the production code, not the test.

- [ ] **Step 4: Commit**

```bash
git add -A && git commit -m "dissolve pkg/admin into pkg/manifest on top of kcapi"
```

---

### Task 14: Docs + final verification

**Files:**
- Modify: `README.md` (library section: kcapi quickstart, new `invoke`/`graph` command docs), `docs/keycloak-admin.md` (note the new architecture; remove references to deleted assets), `docs/design/kcapi-library.md` (mark design shipped, note deviations)

**Interfaces:**
- Consumes: everything above.
- Produces: docs matching reality; green gates.

- [ ] **Step 1: Update README**

Add a "Library usage" section:
```markdown
## Library usage (pkg/kcapi)

`pkg/kcapi` is a spec-driven Keycloak client: every operation in
`keycloak-oapi/*.spec.json` is discoverable and invocable, plus graph queries
over object relationships.

```go
client, _ := kcapi.New(kcapi.Config{
    BaseURL: "https://kc.example.com",
    Spec:    kcapi.SpecSource{Path: "keycloak-oapi/26.6.2.spec.json"},
    Credentials: kcapi.Credentials{ClientID: "admin-cli", Username: "admin", Password: "..."},
})

ops, _ := client.Operations(kcapi.OpFilter{Resource: "users", Method: kcapi.Get})
user, _ := client.Invoke(ctx, kcapi.Call{Op: "getUser", Realm: "demo", Params: kcapi.P{"id": "..."}})
node, _ := client.Resolve(ctx, kcapi.Ref{Type: "users", Name: "alice", Realm: "demo"})
neighbors, edges, _ := client.Neighbors(ctx, node, kcapi.EdgeFilter{Child: "roles"})
```
```
Document `invoke` and `graph` in the CLI command list.

- [ ] **Step 2: Sweep stale references**

```bash
grep -rn 'pkg/admin\|pkg/catalog\|pkg/models\|gen-rel-schema' README.md docs/ --include='*.md'
```
Expected: no hits outside historical/changelog notes (update or delete them).

- [ ] **Step 3: Full verification (evidence before done)**

Run: `go build ./... && go vet ./... && go test ./... && go test -race ./pkg/kcapi/`
Expected: all PASS; record the output in the commit message body.

- [ ] **Step 4: Commit**

```bash
git add -A && git commit -m "docs: kcapi library usage, invoke/graph commands, architecture notes"
```

---

## Self-review notes

- Spec coverage: phase-0 fixes (T1), kcapi core absorb (T2–T4), Operations (T5), Invoke (T6), Reload (T7), relationships/graph (T8–T9), pagination (T10), CLI invoke/graph (T11–T12), admin dissolution (T13), docs (T14). MCP server intentionally deferred to its own plan per spec ("implementation as final phase", separate subsystem).
- Known risk called out for executors: Tasks 2–3 and 13 are mechanical-but-wide; run the full test suite after every step, and prefer `git mv` + `sed` over retyping.
- Where the plan says "check the real name in file X", that is an instruction to read the named file first — the repo's exact helper names win over this plan's sketches.
