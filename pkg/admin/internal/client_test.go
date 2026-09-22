package internal

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thedataflows/keycloak-cli/pkg/kcapi"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

func TestSanitizeResourcePayloadStripsRealmAndParentReferences(t *testing.T) {
	spec, err := kcapi.NewSpec(filepath.Join("..", "..", "..", "keycloak-oapi", "26.6.2.spec.json"))
	require.NoError(t, err)
	client := &RuntimeClient{spec: spec}

	tests := []struct {
		name     string
		resource manifest.Resource
		contract kcapi.OperationContract
		want     map[string]interface{}
	}{
		{
			name: "role under client strips realm and clientUuid",
			resource: manifest.Resource{
				Type:       "role",
				Realm:      "demo",
				ParentType: "client",
				Data: map[string]interface{}{
					"name":       "admin",
					"realm":      "demo",
					"clientUuid": "target-client-uuid",
				},
			},
			contract: kcapi.OperationContract{Path: "/admin/realms/{realm}/clients/{client-uuid}/roles/{role-name}", Method: http.MethodPost},
			want: map[string]interface{}{
				"name": "admin",
			},
		},
		{
			name: "protocolmapper under client scope strips clientScopeId",
			resource: manifest.Resource{
				Type:       "protocolmapper",
				Realm:      "demo",
				ParentType: "clientscope",
				Data: map[string]interface{}{
					"name":          "email",
					"clientScopeId": "scope-1",
					"protocol":      "openid-connect",
				},
			},
			contract: kcapi.OperationContract{Path: "/admin/realms/{realm}/client-scopes/{client-scope-id}/protocol-mappers/models/{id}", Method: http.MethodPost},
			want: map[string]interface{}{
				"name":     "email",
				"protocol": "openid-connect",
			},
		},
		{
			name: "realm keeps realm field",
			resource: manifest.Resource{
				Type:  "realm",
				Realm: "demo",
				Data:  map[string]interface{}{"realm": "demo", "enabled": true},
			},
			contract: kcapi.OperationContract{Path: "/admin/realms", Method: http.MethodPost},
			want: map[string]interface{}{
				"realm":   "demo",
				"enabled": true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.sanitizeResourcePayload(&tt.resource, tt.contract.Method, tt.contract)
			assert.Equal(t, tt.want, tt.resource.Data)
		})
	}
}
