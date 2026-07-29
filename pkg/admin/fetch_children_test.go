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

const demoClientUUID = "0bd3b65d-40c6-46ed-8ffe-81610c031ce1"

// TestFetchChildrenReturnsOneClientRoleCollection pins ISSUE 0003: a single
// (parent, childType) fetch issues exactly one GET to the nested endpoint, with
// no sibling child types and no realm-wide reference sweep, and returns children
// shaped like Depth: 1 does — Type/Realm/ParentType plus the injected clientUuid.
func TestFetchChildrenReturnsOneClientRoleCollection(t *testing.T) {
	var gets []string
	var briefParam string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gets = append(gets, r.Method+" "+r.URL.Path)
		if r.Method == http.MethodGet && r.URL.Path == "/admin/realms/sync-source/clients/"+demoClientUUID+"/roles" {
			briefParam = r.URL.Query().Get("briefRepresentation")
			writeJSON(t, w, []map[string]interface{}{{
				"id":          "e6f79b96-9640-4e88-bd5d-4e6488452eac",
				"name":        "demo-a-clientrole-1",
				"description": "Seeded demo client role",
				"clientRole":  true,
				"composite":   false,
				"attributes":  map[string]interface{}{"tier": []interface{}{"gold"}},
			}})
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	parent := manifest.Resource{Type: "client", Realm: "sync-source", Data: map[string]interface{}{"id": demoClientUUID}}

	report, err := service.FetchChildren(context.Background(), parent, "role", admin.ChildFetchQuery{FullRepresentation: true})
	require.NoError(t, err)
	require.Empty(t, report.Failures)

	// Exactly one HTTP GET — no sibling child types, no reference sweep.
	assert.Equal(t, []string{"GET /admin/realms/sync-source/clients/" + demoClientUUID + "/roles"}, gets)
	// FullRepresentation honoured (ISSUE 0001 parity).
	assert.Equal(t, "false", briefParam)

	require.Len(t, report.Resources, 1)
	role := report.Resources[0]
	assert.Equal(t, "role", role.Type)
	assert.Equal(t, "sync-source", role.Realm)
	assert.Equal(t, "client", role.ParentType)
	assert.Equal(t, demoClientUUID, role.Data["clientUuid"], "parent reference must be injected so the child can be applied standalone")
	assert.Equal(t, "demo-a-clientrole-1", role.Data["name"])
	assert.Contains(t, role.Data, "attributes", "attributes must be present under FullRepresentation")
}

// TestFetchChildrenClassifies404AsNotFound pins that an absent optional child
// collection is a benign FetchFailure{NotFound: true}, not a hard error.
func TestFetchChildrenClassifies404AsNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	parent := manifest.Resource{Type: "client", Realm: "sync-source", Data: map[string]interface{}{"id": demoClientUUID}}

	report, err := service.FetchChildren(context.Background(), parent, "role", admin.ChildFetchQuery{})
	require.NoError(t, err, "a 404 must not be a hard error")
	assert.Empty(t, report.Resources)
	require.Len(t, report.Failures, 1)
	assert.True(t, report.Failures[0].NotFound, "404 must be classified as NotFound")
}
