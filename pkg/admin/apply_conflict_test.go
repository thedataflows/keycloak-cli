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

func TestApplyFallsBackToUpdateOnConflict(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/clients/stale-id":
			http.Error(w, "not found", http.StatusNotFound)
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/clients":
			writeJSON(t, w, []map[string]interface{}{
				{"id": "server-id", "clientId": "client-1", "name": "Client 1"},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/admin/realms/demo/clients":
			http.Error(w, `{"errorMessage":"Client client-1 already exists"}`, http.StatusConflict)
		case r.Method == http.MethodPut && r.URL.Path == "/admin/realms/demo/clients/server-id":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.RequestURI())
		}
	}))
	defer server.Close()
	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{{
		Type:  "client",
		Realm: "demo",
		Data:  map[string]interface{}{"id": "stale-id", "clientId": "client-1", "name": "Client 1"},
	}}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "updated", report.Results[0].Action)
	assert.Equal(t, http.StatusNoContent, report.Results[0].Status)
	assert.Zero(t, report.Failed)
}
