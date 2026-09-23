package kcapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The manifest fetchers emit briefRepresentation=false (and search/max/exact)
// for collection endpoints. Only the operations that declare a query parameter
// accept it; the fetch layer must not put undeclared parameters on the wire
// and rely on the server ignoring them. Declared parameters must survive.

// newFetchClient builds the RuntimeClient the manifest service uses, pointed
// at a fake Keycloak and loaded with the real vendored spec.
func newFetchClient(t *testing.T, baseURL string) *RuntimeClient {
	t.Helper()
	spec, err := NewSpecFromBytes(realSpecBytes(t))
	require.NoError(t, err)
	rc, err := NewRuntimeClient(RuntimeConfig{BaseURL: baseURL, Spec: spec}, staticTokenProvider("test-token"))
	require.NoError(t, err)
	return rc
}

// GET /admin/realms declares briefRepresentation but not search, so a realm
// collection fetch keeps the former and drops the latter.
func TestFetchResourcesSendsOnlyDeclaredQueryParams(t *testing.T) {
	var gotQuery map[string][]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		_, _ = w.Write([]byte(`[{"id":"m","realm":"demo"}]`))
	}))
	defer srv.Close()

	c := newFetchClient(t, srv.URL)
	_, err := c.FetchResources(context.Background(), "realm", nil,
		map[string]string{"briefRepresentation": "false", "search": "demo"})
	require.NoError(t, err)

	assert.Equal(t, []string{"false"}, gotQuery["briefRepresentation"],
		"declared query parameter must reach the wire")
	assert.NotContains(t, gotQuery, "search",
		"undeclared query parameter must not reach the wire")
}

// GET .../organizations/{org-id}/members declares briefRepresentation, max
// and search but not q (the ${KEYCLOAK_VERSION} spec): a child collection fetch keeps the
// declared parameters and drops q.
func TestFetchPathCollectionSendsOnlyDeclaredQueryParams(t *testing.T) {
	var gotQuery map[string][]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		_, _ = w.Write([]byte(`[{"id":"m1"}]`))
	}))
	defer srv.Close()

	c := newFetchClient(t, srv.URL)
	_, err := c.FetchPathCollection(context.Background(),
		"/admin/realms/{realm}/organizations/{org-id}/members",
		map[string]string{"realm": "demo", "org-id": "o1"},
		map[string]string{"briefRepresentation": "false", "max": "5", "q": "acme"})
	require.NoError(t, err)

	assert.Equal(t, []string{"false"}, gotQuery["briefRepresentation"],
		"declared query parameter must reach the wire")
	assert.Equal(t, []string{"5"}, gotQuery["max"],
		"declared query parameter must reach the wire")
	assert.NotContains(t, gotQuery, "q",
		"undeclared query parameter must not reach the wire")
}
