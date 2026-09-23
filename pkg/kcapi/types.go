// kcapi must not import pkg/manifest (manifest aliases kcapi core types).

package kcapi

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"strings"
)

// Resource is the shared manifest model for Keycloak resources.
type Resource struct {
	Type       string                 `json:"type"`
	Realm      string                 `json:"realm"`
	Delete     bool                   `json:"delete,omitempty"`
	Data       map[string]interface{} `json:"data"`
	ParentType string                 `json:"parentType,omitempty"` // disambiguates multi-location types (e.g. protocolmapper under clientscope vs client)
}

type RelationshipOperation struct {
	Kind   string          `json:"kind,omitempty"`
	Path   string          `json:"path"`
	Delete bool            `json:"delete,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`

	Method     string            `json:"-"`
	Template   string            `json:"-"`
	PathParams map[string]string `json:"-"`
}

func (r *RelationshipOperation) RebuildPath() {
	r.Path = buildActualRelationshipPath(r.Template, r.PathParams)
}

var identifierFields = map[string][]string{
	"realm":            {"realm"},
	"client":           {"id", "clientId"},
	"role":             {"name", "alias"},
	"identityprovider": {"alias"},
	"user":             {"id", "username"},
}

func (r Resource) Identifier() string {
	if fields, ok := identifierFields[r.Type]; ok {
		return rawFirstStringField(r.Data, fields)
	}
	return rawStringField(r.Data, "id")
}

var nameFields = map[string][]string{
	"realm":  {"realm"},
	"user":   {"username"},
	"client": {"clientId"},
}

func (r Resource) Name() string {
	if fields, ok := nameFields[r.Type]; ok {
		return rawFirstStringField(r.Data, fields)
	}
	return rawFirstStringField(r.Data, []string{"name", "alias"})
}

var displayNameFields = []string{"name", "clientId", "username", "alias", "realm", "id"}

func (r Resource) DisplayName() string {
	return rawFirstStringField(r.Data, displayNameFields)
}

// rawStringField and rawFirstStringField are the untrimmed string helpers moved
// from pkg/manifest alongside the Resource methods; package kcapi's existing
// stringField/firstStringField trim whitespace and must not rebind these reads.
func rawStringField(data map[string]interface{}, key string) string {
	if data == nil {
		return ""
	}

	value, _ := data[key].(string)
	return value
}

func rawFirstStringField(data map[string]interface{}, keys []string) string {
	for _, key := range keys {
		if value := rawStringField(data, key); value != "" {
			return value
		}
	}
	return ""
}

func buildActualRelationshipPath(template string, params map[string]string) string {
	resolved := template
	for key, value := range params {
		placeholder := fmt.Sprintf("{%s}", key)
		resolved = strings.ReplaceAll(resolved, placeholder, value)
	}
	resolved = strings.TrimPrefix(resolved, "/")
	return resolved
}

// NewRelationshipOperation is moved verbatim from pkg/manifest/manifest.go (it
// constructs the kcapi-owned RelationshipOperation). relationshipPathPrefix and
// buildActualRelationshipPath resolve to the byte-identical kcapi copies.
func NewRelationshipOperation(template, method string, params map[string]string, payload interface{}) (RelationshipOperation, error) {
	trimmedTemplate := strings.TrimSpace(template)
	trimmedTemplate = strings.TrimPrefix(trimmedTemplate, relationshipPathPrefix)
	trimmedTemplate = strings.TrimPrefix(trimmedTemplate, "/")
	if trimmedTemplate == "" {
		return RelationshipOperation{}, fmt.Errorf("relationship template cannot be empty")
	}

	resolvedMethod := strings.ToUpper(strings.TrimSpace(method))
	if resolvedMethod == "" {
		return RelationshipOperation{}, fmt.Errorf("relationship method cannot be empty")
	}

	operation := RelationshipOperation{
		Template:   trimmedTemplate,
		Method:     resolvedMethod,
		PathParams: maps.Clone(params),
	}
	if resolvedMethod == http.MethodDelete {
		operation.Delete = true
	}

	operation.Path = buildActualRelationshipPath(operation.Template, operation.PathParams)

	if payload != nil {
		switch typed := payload.(type) {
		case json.RawMessage:
			operation.Data = append(json.RawMessage(nil), typed...)
		default:
			data, err := json.Marshal(payload)
			if err != nil {
				return RelationshipOperation{}, fmt.Errorf("marshal relationship payload: %w", err)
			}
			operation.Data = data
		}
	}

	return operation, nil
}
