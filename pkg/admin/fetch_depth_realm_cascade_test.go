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

// TestFetchDepthRealmCascadeReachesIdentityProviderMappers pins that a
// realm-rooted depth traversal descends past the realm's direct children into
// their children — here realm -> identityprovider -> mapper. A resource fetched
// as a child of the realm carries ParentType="realm"; the traversal must NOT
// treat that as an organization-style scope (which would route the descent
// through the scoped child-collection selector and, because the realm
// placeholder is stripped from every path chain, match nothing and halt one
// level below the realm). The fake serves only the structural realm-rooted
// paths, so the test passes only if the mapper is reached structurally
// (ISSUE 0009).
func TestFetchDepthRealmCascadeReachesIdentityProviderMappers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms":
			writeJSON(t, w, []map[string]interface{}{{"realm": "demo"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/identity-provider/instances":
			writeJSON(t, w, []map[string]interface{}{{"alias": "idp-1", "providerId": "oidc"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/identity-provider/instances/idp-1/mappers":
			writeJSON(t, w, []map[string]interface{}{{
				"id":                     "m1",
				"name":                   "mapper-1",
				"identityProviderAlias":  "idp-1",
				"identityProviderMapper": "oidc-hardcoded-role-idp-mapper",
			}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Fetch(context.Background(), admin.FetchQuery{
		Realm: "demo", Resources: "realm", Depth: 2,
	})
	require.NoError(t, err)

	var mapper *struct {
		parentType string
		alias      interface{}
	}
	for _, res := range report.Resources {
		if res.Type == "identityprovidermapper" && res.Data["name"] == "mapper-1" {
			mapper = &struct {
				parentType string
				alias      interface{}
			}{res.ParentType, res.Data["alias"]}
		}
	}

	require.NotNil(t, mapper, "realm-rooted depth traversal must reach the identity provider's mapper")
	assert.Equal(t, "identityprovider", mapper.parentType, "mapper must be keyed to its identity-provider parent for re-apply")
	assert.Equal(t, "idp-1", mapper.alias, "parent alias must be injected so the mapper round-trips")
}

// TestFetchDepthRealmCascadeReachesClientRoles pins the same realm-rooted
// descent for a second grandchild leg — realm -> client -> role — so the fix is
// not narrowly scoped to identity providers. A client fetched as a realm-child
// carries ParentType="realm" and must descend structurally into its roles
// (ISSUE 0009).
func TestFetchDepthRealmCascadeReachesClientRoles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms":
			writeJSON(t, w, []map[string]interface{}{{"realm": "demo"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/clients":
			writeJSON(t, w, []map[string]interface{}{{"id": "client-1", "clientId": "app"}})
		case r.Method == http.MethodGet && r.URL.Path == "/admin/realms/demo/clients/client-1/roles":
			writeJSON(t, w, []map[string]interface{}{{"id": "r1", "name": "role-1", "clientRole": true}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	service := newServiceForTest(t, server.URL)
	report, err := service.Fetch(context.Background(), admin.FetchQuery{
		Realm: "demo", Resources: "realm", Depth: 2,
	})
	require.NoError(t, err)

	var role *string
	for _, res := range report.Resources {
		if res.Type == "role" && res.Data["name"] == "role-1" {
			pt := res.ParentType
			role = &pt
		}
	}

	require.NotNil(t, role, "realm-rooted depth traversal must reach the client's role")
	assert.Equal(t, "client", *role, "client role must be keyed to its client parent for re-apply")
}
