package admin_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/admin"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

// groupClearServer records the body of the group create/update write so the
// test can assert whether an explicit-empty attributes map reaches Keycloak.
// existing controls whether the group is already present (drives create vs
// update).
func groupClearServer(t *testing.T, capturedBody *map[string]interface{}, existing bool) *httptest.Server {
	t.Helper()
	groupList := []map[string]interface{}{{
		"id":         "group-1",
		"name":       "platform",
		"attributes": map[string]interface{}{"department": []interface{}{"ENGINEERING"}},
	}}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.TrimRight(r.URL.Path, "/") == "/admin/realms/demo/groups":
			if !existing {
				writeJSON(t, w, []map[string]interface{}{})
				return
			}
			writeJSON(t, w, groupList)
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/groups/group-1":
			if !existing {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			writeJSON(t, w, groupList[0])
		case r.Method == http.MethodPost && r.URL.Path == "/admin/realms/demo/groups":
			body, _ := io.ReadAll(r.Body)
			require.NoError(t, json.Unmarshal(body, capturedBody))
			w.Header().Set("Location", r.URL.Path+"/group-1")
			w.WriteHeader(http.StatusCreated)
		case r.Method == http.MethodPut && r.URL.Path == "/admin/realms/demo/groups/group-1":
			body, _ := io.ReadAll(r.Body)
			require.NoError(t, json.Unmarshal(body, capturedBody))
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
}

// TestApplyGroupUpdateClearsAttributes is the ISSUE 0011 regression: applying a
// group whose Data carries an explicit-empty attributes map ({}) to remove it
// must send that empty map on the wire (not omit it), so Keycloak clears the
// collection. Omitting it makes Keycloak preserve the previous values and the
// leg never converges.
func TestApplyGroupUpdateClearsAttributes(t *testing.T) {
	var body map[string]interface{}
	server := groupClearServer(t, &body, true)
	defer server.Close()

	resource := manifest.Resource{
		Type:  "group",
		Realm: "demo",
		Data: map[string]interface{}{
			"id":         "group-1",
			"name":       "platform",
			"attributes": map[string]interface{}{},
		},
	}

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{resource}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "updated", report.Results[0].Action)

	require.Contains(t, body, "attributes", "update body must carry the explicit-empty attributes so Keycloak clears it")
	attrs, ok := body["attributes"].(map[string]interface{})
	require.True(t, ok, "attributes must be an object, got %T", body["attributes"])
	assert.Empty(t, attrs, "attributes must be sent as an empty object")
}

// TestApplyGroupCreateSendsExplicitEmptyAttributes covers the create side of the
// same acceptance criterion: an explicit-empty collection must survive to the
// create body too, not be dropped before the request is built.
func TestApplyGroupCreateSendsExplicitEmptyAttributes(t *testing.T) {
	var body map[string]interface{}
	server := groupClearServer(t, &body, false)
	defer server.Close()

	resource := manifest.Resource{
		Type:  "group",
		Realm: "demo",
		Data: map[string]interface{}{
			"name":       "platform",
			"attributes": map[string]interface{}{},
		},
	}

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{resource}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Equal(t, "created", report.Results[0].Action)

	require.Contains(t, body, "attributes", "create body must carry the explicit-empty attributes")
	attrs, ok := body["attributes"].(map[string]interface{})
	require.True(t, ok, "attributes must be an object, got %T", body["attributes"])
	assert.Empty(t, attrs, "attributes must be sent as an empty object")
}
