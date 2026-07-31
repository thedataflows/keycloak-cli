package admin_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/admin"
)

// Keycloak's organizations list returns a brief representation that omits the
// attributes map unless briefRepresentation=false is sent. Fetching
// organizations must always request the full representation so an
// export→apply round-trip preserves attributes, even without
// --full-representation (ISSUE 0010).
func TestFetchOrganizationAlwaysRequestsFullRepresentation(t *testing.T) {
	var orgQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms":
			writeJSON(t, w, []map[string]interface{}{{"realm": "demo"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/organizations":
			orgQuery = r.URL.RawQuery
			// Echo attributes only when the full representation is requested,
			// mirroring Keycloak's brief/full behavior.
			org := map[string]interface{}{"id": "org-1", "alias": "acme", "name": "acme"}
			if r.URL.Query().Get("briefRepresentation") == "false" {
				org["attributes"] = map[string]interface{}{"department": []interface{}{"ENGINEERING"}}
			}
			writeJSON(t, w, []map[string]interface{}{org})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	// Note: FullRepresentation deliberately left false — organizations must opt
	// into the full form on their own.
	report, err := service.Fetch(context.Background(), admin.FetchQuery{Resources: "organization"})
	require.NoError(t, err)
	require.Len(t, report.Resources, 1)

	assert.Contains(t, orgQuery, "briefRepresentation=false", "org fetch must send briefRepresentation=false")
	attrs, ok := report.Resources[0].Data["attributes"].(map[string]interface{})
	require.True(t, ok, "fetched organization must carry attributes, got %T", report.Resources[0].Data["attributes"])
	assert.Equal(t, []interface{}{"ENGINEERING"}, attrs["department"])
}
