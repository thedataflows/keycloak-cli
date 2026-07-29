package admin_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/admin"
	"github.com/thedataflows/keycloak-cli/pkg/catalog"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

// A realm group nested under another realm group is created via
// POST /groups/{group-id}/children. The parent id is carried in Data as groupId
// so the resolver can render {group-id}, and it must be stripped from the body —
// Keycloak rejects "Unrecognized field groupId". ISSUE 0005 insists this is
// proven against the embedded spec, not only a fake Apply, because a fake
// validates no body.
func TestApplyNestedRealmGroupStripsParentBindingFromBody(t *testing.T) {
	var createBody map[string]interface{}
	var writePaths []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		// Existence probes: the child does not exist yet.
		case r.Method == http.MethodGet:
			// Single GET on the self path (empty own id) and the children listing
			// both report "not found" so Apply proceeds to create.
			if strings.HasSuffix(r.URL.Path, "/children") {
				writeJSON(t, w, []map[string]interface{}{})
				return
			}
			http.Error(w, "not found", http.StatusNotFound)
		case r.Method == http.MethodPost && r.URL.Path == "/admin/realms/demo/groups/parent-1/children":
			writePaths = append(writePaths, r.Method+" "+r.URL.Path)
			body, _ := io.ReadAll(r.Body)
			require.NoError(t, json.Unmarshal(body, &createBody))
			w.WriteHeader(http.StatusCreated)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{{
		Type:       "group",
		Realm:      "demo",
		ParentType: "group",
		Data:       map[string]interface{}{"name": "platform", "groupId": "parent-1"},
	}}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Zero(t, report.Failed)
	assert.Equal(t, "created", report.Results[0].Action)

	// Routed to the children endpoint under the parent, never the realm root.
	assert.Equal(t, []string{"POST /admin/realms/demo/groups/parent-1/children"}, writePaths)

	// The parent binding must not be on the wire, and the create body must be a
	// valid GroupRepresentation per the embedded spec.
	require.NotNil(t, createBody)
	assert.NotContains(t, createBody, "groupId", "parent binding must be stripped from the body")
	assert.Equal(t, "platform", createBody["name"])

	spec, err := catalog.NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)
	require.NoError(t, spec.ValidateOperationRequest(
		"/admin/realms/{realm}/groups/{group-id}/children", http.MethodPost,
		catalog.RequestValidation{
			PathParams: map[string]string{"realm": "demo", "group-id": "parent-1"},
			Body:       createBody,
		}),
		"the create body must validate against the spec GroupRepresentation")
}

// TestApplyNestedRealmGroupUpdateAddressesChildAndStripsParentBinding is the
// regression a fake nearly missed and a live push caught: on the single-resource
// PUT the parent placeholder {group-id} IS the self id, so the parent binding
// groupId is invisible to the single contract and would leak into the update
// body — Keycloak rejects "Unrecognized field groupId". The update must (a)
// address the child by its own id, never the parent's (Gap 3), and (b) send a
// body with no groupId that validates against the spec GroupRepresentation.
func TestApplyNestedRealmGroupUpdateAddressesChildAndStripsParentBinding(t *testing.T) {
	var updateBody map[string]interface{}
	var writes []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/groups/child-1":
			writeJSON(t, w, map[string]interface{}{"id": "child-1", "name": "platform"})
		case r.Method == http.MethodPut && r.URL.Path == "/admin/realms/demo/groups/child-1":
			writes = append(writes, r.Method+" "+r.URL.Path)
			body, _ := io.ReadAll(r.Body)
			require.NoError(t, json.Unmarshal(body, &updateBody))
			w.WriteHeader(http.StatusNoContent)
		default:
			// Any GET/PUT on the PARENT id is a destructive Gap 3 failure.
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Apply(context.Background(), []manifest.Resource{{
		Type:       "group",
		Realm:      "demo",
		ParentType: "group",
		Data:       map[string]interface{}{"id": "child-1", "name": "platform", "groupId": "parent-1"},
	}}, nil, admin.ApplyOptions{})
	require.NoError(t, err)
	require.Len(t, report.Results, 1)
	assert.Zero(t, report.Failed)

	// Addressed the child by its own id — never the parent (Gap 3).
	assert.Equal(t, []string{"PUT /admin/realms/demo/groups/child-1"}, writes)

	require.NotNil(t, updateBody)
	assert.NotContains(t, updateBody, "groupId", "parent binding must be stripped from the update body")

	spec, err := catalog.NewSpec(filepath.Join("..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)
	require.NoError(t, spec.ValidateOperationRequest(
		"/admin/realms/{realm}/groups/{group-id}", http.MethodPut,
		catalog.RequestValidation{
			PathParams: map[string]string{"realm": "demo", "group-id": "child-1"},
			Body:       updateBody,
		}),
		"the update body must validate against the spec GroupRepresentation")
}
