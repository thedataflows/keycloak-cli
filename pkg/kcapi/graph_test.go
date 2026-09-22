package kcapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// queryRecorder records the query strings that hit the fake server, so tests
// can pin what went on the wire.
type queryRecorder struct {
	mu     sync.Mutex
	byPath map[string]string // request path -> raw query
}

func (q *queryRecorder) record(path, rawQuery string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.byPath[path] = rawQuery
}

func (q *queryRecorder) queryFor(path string) string {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.byPath[path]
}

func (q *queryRecorder) hit(path string) bool {
	_, ok := q.byPath[path]
	return ok
}

// fakeKC is a minimal in-memory Keycloak whose routes mirror the vendored
// spec's real path shapes (26.6.2.spec.json: {user-id}-style placeholders,
// no operationIds). The client loads that same committed spec, so Resolve
// searches through the real users collection contract and Neighbors walks the
// real users->groups and users->consents edges:
//
//	users->groups   /admin/realms/{realm}/users/{user-id}/groups/{groupId}
//	users->consents /admin/realms/{realm}/users/{user-id}/consents/{client}
//
// (There is no users->roles edge in the real spec: /users/{user-id}/role-mappings
// carries no second placeholder and .../role-mappings/clients/{client-id}'s
// prefix is not a spec path, so the brief's role-mappings traversal is adapted
// to the membership edges that actually exist.)
func fakeKC(t *testing.T) (*httptest.Server, *Client, *queryRecorder) {
	t.Helper()
	rec := &queryRecorder{byPath: map[string]string{}}
	mux := http.NewServeMux()
	// GET /admin/realms/{realm}/users — collection search; the real spec
	// declares username as a query param on this endpoint.
	mux.HandleFunc("GET /admin/realms/demo/users", func(w http.ResponseWriter, r *http.Request) {
		rec.record("/users", r.URL.RawQuery)
		switch r.URL.Query().Get("username") {
		case "alice":
			_, _ = w.Write([]byte(`[{"id":"u-1","username":"alice"}]`))
		case "":
			_, _ = w.Write([]byte(`[{"id":"u-1","username":"alice"},{"id":"u-2","username":"bob"}]`))
		default:
			_, _ = w.Write([]byte(`[]`))
		}
	})
	// GET /admin/realms/{realm}/users/{user-id} — the single-resource item.
	mux.HandleFunc("GET /admin/realms/demo/users/u-1", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"id":"u-1","username":"alice","firstName":"Alice"}`))
	})
	// GET /admin/realms/{realm}/users/{user-id}/groups — the users->groups
	// edge's collection (the edge path hangs under it).
	mux.HandleFunc("GET /admin/realms/demo/users/u-1/groups", func(w http.ResponseWriter, r *http.Request) {
		rec.record("/groups", r.URL.RawQuery)
		_, _ = w.Write([]byte(`[{"id":"g-1","name":"admins","path":"/admins"}]`))
	})
	// GET /admin/realms/{realm}/users/{user-id}/consents — answers a bare
	// object, exercising the single-object fallback.
	mux.HandleFunc("GET /admin/realms/demo/users/u-1/consents", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"grantCount":1}`))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	c, err := New(Config{
		BaseURL: srv.URL,
		Spec:    SpecSource{Raw: realSpecBytes(t)},
		Auth:    staticTokenProvider("test-token"),
	})
	require.NoError(t, err)
	return srv, c, rec
}

func TestResolveByName(t *testing.T) {
	_, c, rec := fakeKC(t)
	node, err := c.Resolve(context.Background(), Ref{Type: "users", Name: "alice", Realm: "demo"})
	require.NoError(t, err)
	assert.Equal(t, "u-1", node.Resource.Data["id"])
	// The search goes over the wire with the identity field the real spec
	// declares as a query param (username), not a client-side-only filter.
	assert.Equal(t, "username=alice", rec.queryFor("/users"))
	assert.Equal(t, "users", node.Ref.Type)
	assert.Equal(t, "alice", node.Ref.Name)
	assert.Equal(t, "u-1", node.Ref.ID)
	assert.Equal(t, "demo", node.Ref.Realm)
}

func TestResolveByID(t *testing.T) {
	_, c, _ := fakeKC(t)
	node, err := c.Resolve(context.Background(), Ref{Type: "users", ID: "u-1", Realm: "demo"})
	require.NoError(t, err)
	// ID set: the single-resource item path is hit directly, no search.
	assert.Equal(t, "alice", node.Resource.Data["username"])
	assert.Equal(t, "alice", node.Ref.Name, "the fetched representation names the node")
}

func TestResolveNotFoundAndAmbiguous(t *testing.T) {
	_, c, _ := fakeKC(t)
	_, err := c.Resolve(context.Background(), Ref{Type: "users", Name: "nobody", Realm: "demo"})
	assert.ErrorIs(t, err, ErrNotFound)
	_, err = c.Resolve(context.Background(), Ref{Type: "users", Realm: "demo"}) // no name -> lists all -> 2 hits
	require.ErrorIs(t, err, ErrValidation)                                      // ambiguous
	assert.Contains(t, err.Error(), "alice")
	assert.Contains(t, err.Error(), "bob")
}

func TestNeighborsFetchesChildCollection(t *testing.T) {
	_, c, rec := fakeKC(t)
	node, err := c.Resolve(context.Background(), Ref{Type: "users", Name: "alice", Realm: "demo"})
	require.NoError(t, err)

	nodes, edges, err := c.Neighbors(context.Background(), node, EdgeFilter{Child: "groups"})
	require.NoError(t, err)
	require.Len(t, nodes, 1)
	assert.Equal(t, "admins", nodes[0].Resource.Data["name"])
	assert.Equal(t, "groups", nodes[0].Ref.Type)
	assert.Equal(t, "admins", nodes[0].Ref.Name)
	assert.Equal(t, "demo", nodes[0].Ref.Realm)
	require.NotEmpty(t, edges)
	assert.Equal(t, "users", edges[0].Parent)
	assert.Equal(t, "groups", edges[0].Child)
	// The walk instantiates the edge's collection prefix (the real spec path
	// the edge hangs under), not the edge's own item path.
	assert.True(t, rec.hit("/groups"), "the users->groups collection must be fetched")
}

func TestNeighborsSingleObjectFallback(t *testing.T) {
	_, c, _ := fakeKC(t)
	node, err := c.Resolve(context.Background(), Ref{Type: "users", Name: "alice", Realm: "demo"})
	require.NoError(t, err)

	// The consents edge's collection answers a bare JSON object (the edge
	// itself is cardinality one); it still wraps as one node.
	nodes, _, err := c.Neighbors(context.Background(), node, EdgeFilter{Child: "consents"})
	require.NoError(t, err)
	require.Len(t, nodes, 1)
	assert.Equal(t, "consents", nodes[0].Ref.Type)
}

func TestNeighborsFilter(t *testing.T) {
	_, c, _ := fakeKC(t)
	node, err := c.Resolve(context.Background(), Ref{Type: "users", Name: "alice", Realm: "demo"})
	require.NoError(t, err)

	// A Parent filter that does not match the node's own type as parent
	// selects nothing without error.
	nodes, edges, err := c.Neighbors(context.Background(), node, EdgeFilter{Parent: "clients"})
	require.NoError(t, err)
	assert.Empty(t, nodes)
	assert.Empty(t, edges)

	// Parent and Child combine with AND.
	nodes, edges, err = c.Neighbors(context.Background(), node, EdgeFilter{Parent: "users", Child: "groups"})
	require.NoError(t, err)
	require.Len(t, nodes, 1)
	require.Len(t, edges, 1)
}
