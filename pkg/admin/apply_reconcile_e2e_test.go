package admin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/admin"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

func TestApplyReconcileRemovesUnexpectedRelationships(t *testing.T) {
	var requests []struct {
		method string
		path   string
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, struct {
			method string
			path   string
		}{method: r.Method, path: r.URL.Path})
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/users/alice/groups":
			writeJSON(t, w, []map[string]interface{}{{"id": "old-group", "name": "old"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/users":
			writeJSON(t, w, []map[string]interface{}{{"id": "alice", "username": "alice"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/groups":
			writeJSON(t, w, []map[string]interface{}{{"id": "old-group", "name": "old"}})
		case r.Method == http.MethodDelete && r.URL.Path == "/admin/realms/demo/users/alice/groups/old-group":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodPut && r.URL.Path == "/admin/realms/demo/users/alice/groups/new-group":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/admin/realms/demo/"):
			writeJSON(t, w, []map[string]interface{}{})
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	relationships := []manifest.RelationshipOperation{{
		Kind:   "user-group-membership",
		Path:   "demo/users/alice/groups/new-group",
		Method: "PUT",
	}}

	report, err := service.Apply(context.Background(), nil, relationships, admin.ApplyOptions{Reconcile: true})
	require.NoError(t, err)
	assert.Zero(t, report.Failed)

	var sawDelete, sawPut bool
	for _, req := range requests {
		if req.method == http.MethodDelete && req.path == "/admin/realms/demo/users/alice/groups/old-group" {
			sawDelete = true
		}
		if req.method == http.MethodPut && req.path == "/admin/realms/demo/users/alice/groups/new-group" {
			sawPut = true
		}
	}
	assert.True(t, sawDelete, "expected DELETE for unexpected relationship")
	assert.True(t, sawPut, "expected PUT for desired relationship")
}
