package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

// RelationshipKind describes a family of Keycloak relationships discovered from
// the OpenAPI spec or supplied through an override file.
type RelationshipKind struct {
	Name                 string
	ResourceA            string
	ResourceB            string
	Direction            string
	ReadMethod           string
	ReadPath             string
	WriteMethod          string
	WriteTemplate        string
	DeleteMethod         string
	DeleteTemplate       string
	DeleteItemParam      string
	DeletePayloadField   string
	ItemParamName        string
	PayloadField         string
	BulkPayload          bool
	ParamTypes           map[string]string
	relationshipTemplate string
}

// Pattern returns the read-side pattern used by the relationship fetcher.
func (k RelationshipKind) Pattern() RelationshipOperationPattern {
	return RelationshipOperationPattern{
		Kind:                 k.Name,
		Method:               k.ReadMethod,
		Path:                 k.ReadPath,
		PathParams:           extractPathParams(k.ReadPath),
		ResourceA:            k.ResourceA,
		ResourceB:            k.ResourceB,
		Direction:            k.Direction,
		RelationshipTemplate: k.WriteTemplate,
		RelationshipMethod:   k.WriteMethod,
		ItemParamName:        k.ItemParamName,
		PayloadField:         k.PayloadField,
		BulkPayload:          k.BulkPayload,
	}
}

// ParamTypesFor returns the resource types corresponding to path parameters for
// this relationship kind.
func (k RelationshipKind) ParamTypesFor(param string) string {
	if k.ParamTypes != nil {
		if t, ok := k.ParamTypes[param]; ok {
			return t
		}
	}
	return fallbackPlaceholderToResourceType[param]
}

// Registry holds known relationship kinds indexed by name, read path, and write
// template.
type Registry struct {
	kinds           map[string]RelationshipKind
	byPath          map[string]RelationshipKind
	byWriteTemplate map[string]RelationshipKind
}

// NewRegistry creates an empty relationship registry.
func NewRegistry() *Registry {
	return &Registry{
		kinds:           make(map[string]RelationshipKind),
		byPath:          make(map[string]RelationshipKind),
		byWriteTemplate: make(map[string]RelationshipKind),
	}
}

// Register adds a relationship kind to the registry.
func (r *Registry) Register(kind RelationshipKind) {
	if r.kinds == nil {
		r.kinds = make(map[string]RelationshipKind)
	}
	if r.byPath == nil {
		r.byPath = make(map[string]RelationshipKind)
	}
	if r.byWriteTemplate == nil {
		r.byWriteTemplate = make(map[string]RelationshipKind)
	}
	readPath := normalizeReadPath(kind.ReadPath)
	if readPath != "" {
		r.byPath[readPath] = kind
	}
	writeTemplate := normalizeRelationshipPath(kind.WriteTemplate)
	if writeTemplate != "" {
		r.byWriteTemplate[writeTemplate] = kind
	}
	r.kinds[kind.Name] = kind
}

// ByName returns a kind by its canonical name.
func (r *Registry) ByName(name string) (RelationshipKind, bool) {
	kind, ok := r.kinds[name]
	return kind, ok
}

// ByPath looks up a relationship kind by its read path.
func (r *Registry) ByPath(path string) (RelationshipKind, bool) {
	path = normalizeReadPath(path)
	kind, ok := r.byPath[path]
	return kind, ok
}

// normalizeReadPath normalizes a spec read path by removing the admin/realms
// prefix and the realm placeholder segment so that both
// "/admin/realms/{realm}/users/{user-id}/groups" and "users/{user-id}/groups"
// resolve to the same registry key.
func normalizeReadPath(path string) string {
	path = normalizeRelationshipPath(path)
	return strings.TrimPrefix(path, "{realm}/")
}

// ByWriteTemplate looks up a relationship kind by its write template.
func (r *Registry) ByWriteTemplate(template string) (RelationshipKind, bool) {
	template = normalizeRelationshipPath(template)
	kind, ok := r.byWriteTemplate[template]
	return kind, ok
}

// Kinds returns all registered relationship kinds sorted by name.
func (r *Registry) Kinds() []RelationshipKind {
	names := make([]string, 0, len(r.kinds))
	for name := range r.kinds {
		names = append(names, name)
	}
	sort.Strings(names)
	result := make([]RelationshipKind, 0, len(names))
	for _, name := range names {
		result = append(result, r.kinds[name])
	}
	return result
}

// ParamTypes returns the path-parameter to resource-type mapping for a kind.
func (r *Registry) ParamTypes(kindName string) map[string]string {
	kind, ok := r.kinds[kindName]
	if !ok {
		return nil
	}
	return kind.ParamTypes
}

var defaultRelationshipRegistry = func() *Registry {
	r := NewRegistry()
	for _, kind := range defaultRelationshipKinds() {
		r.Register(kind)
	}
	return r
}()

// DefaultRegistry returns the built-in relationship registry bundled with the
// CLI. It is safe for read-only use; callers that need overrides should build a
// registry with NewRegistry and merge overrides on top.
func DefaultRegistry() *Registry {
	return defaultRelationshipRegistry
}

// InstallDefaultRegistry loads relationship overrides from the directory
// containing the spec and wires the default registry into the manifest package.
// It is safe to call when no override file exists.
func InstallDefaultRegistry(specPath string) error {
	specDir := filepath.Dir(specPath)
	overridePath := filepath.Join(specDir, "relationship-overrides.yaml")

	overrides, err := LoadRelationshipOverrides(overridePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			InstallManifestRegistry()
			return nil
		}
		return err
	}

	if err := ApplyRelationshipOverrides(defaultRelationshipRegistry, overrides); err != nil {
		return err
	}
	InstallManifestRegistry()
	return nil
}

// InstallManifestRegistry wires the default registry into the manifest package
// so relationship normalization can resolve parameter types from the registry
// instead of a hardcoded switch.
func InstallManifestRegistry() {
	manifest.RelationshipParamTypes = defaultRelationshipRegistry.paramTypesFunc()
}

func (r *Registry) paramTypesFunc() func(string) map[string]string {
	return func(kindName string) map[string]string {
		return r.ParamTypes(kindName)
	}
}

// RelationshipOverride describes a single override entry in a relationship
// overrides YAML file.
type RelationshipOverride struct {
	Name           string            `yaml:"name"`
	ResourceA      string            `yaml:"resourceA"`
	ResourceB      string            `yaml:"resourceB"`
	Direction      string            `yaml:"direction"`
	ReadPath       string            `yaml:"readPath"`
	WriteTemplate  string            `yaml:"writeTemplate"`
	WriteMethod    string            `yaml:"writeMethod"`
	DeleteTemplate string            `yaml:"deleteTemplate,omitempty"`
	DeleteMethod   string            `yaml:"deleteMethod,omitempty"`
	ItemParamName  string            `yaml:"itemParamName,omitempty"`
	PayloadField   string            `yaml:"payloadField,omitempty"`
	BulkPayload    bool              `yaml:"bulkPayload,omitempty"`
	ParamTypes     map[string]string `yaml:"paramTypes,omitempty"`
	Disabled       bool              `yaml:"disabled,omitempty"`
}

// RelationshipOverridesFile is the top-level shape of an override file.
type RelationshipOverridesFile struct {
	Overrides []RelationshipOverride `yaml:"overrides"`
}

// LoadRelationshipOverrides reads relationship overrides from a YAML file.
func LoadRelationshipOverrides(path string) ([]RelationshipOverride, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read relationship overrides: %w", err)
	}

	var file RelationshipOverridesFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse relationship overrides: %w", err)
	}
	return file.Overrides, nil
}

// ApplyRelationshipOverrides merges override entries into the registry. Entries
// with Disabled=true are removed; otherwise the entry is registered or updated.
func ApplyRelationshipOverrides(r *Registry, overrides []RelationshipOverride) error {
	for _, override := range overrides {
		if override.Name == "" {
			return fmt.Errorf("relationship override missing required field 'name'")
		}
		if override.Disabled {
			delete(r.kinds, override.Name)
			if override.ReadPath != "" {
				delete(r.byPath, normalizeReadPath(override.ReadPath))
			}
			if override.WriteTemplate != "" {
				delete(r.byWriteTemplate, normalizeRelationshipPath(override.WriteTemplate))
			}
			continue
		}

		if override.ReadPath == "" {
			return fmt.Errorf("relationship override %q missing required field 'readPath'", override.Name)
		}
		if override.WriteTemplate == "" {
			return fmt.Errorf("relationship override %q missing required field 'writeTemplate'", override.Name)
		}
		if override.WriteMethod == "" {
			override.WriteMethod = http.MethodPut
		}

		kind := RelationshipKind{
			Name:           override.Name,
			ResourceA:      override.ResourceA,
			ResourceB:      override.ResourceB,
			Direction:      override.Direction,
			ReadMethod:     http.MethodGet,
			ReadPath:       override.ReadPath,
			WriteMethod:    override.WriteMethod,
			WriteTemplate:  override.WriteTemplate,
			DeleteMethod:   override.DeleteMethod,
			DeleteTemplate: override.DeleteTemplate,
			ItemParamName:  override.ItemParamName,
			PayloadField:   override.PayloadField,
			BulkPayload:    override.BulkPayload,
			ParamTypes:     override.ParamTypes,
		}
		if kind.Direction == "" {
			kind.Direction = "one-way"
		}
		if kind.DeleteMethod == "" {
			kind.DeleteMethod = http.MethodDelete
		}
		if kind.DeleteTemplate == "" {
			kind.DeleteTemplate = kind.WriteTemplate
		}
		if kind.ParamTypes == nil {
			kind.ParamTypes = inferParamTypesForTemplate(kind.WriteTemplate)
		}
		r.Register(kind)
	}
	return nil
}

func inferParamTypesForTemplate(template string) map[string]string {
	params := extractPathParams(template)
	result := make(map[string]string, len(params))
	for _, param := range params {
		if param == "realm" {
			continue
		}
		if resourceType, ok := fallbackPlaceholderToResourceType[param]; ok {
			result[param] = resourceType
		}
	}
	return result
}

func defaultRelationshipKinds() []RelationshipKind {
	patterns := []struct {
		match                string
		name                 string
		resourceA            string
		resourceB            string
		relationshipTemplate string
		relationshipMethod   string
		itemParamName        string
		payloadField         string
		bulkPayload          bool
		deleteTemplate       string
		deleteItemParam      string
	}{
		{"/users/{user-id}/groups", "user-group-membership", "user", "group", "{realm}/users/{user-id}/groups/{groupId}", "PUT", "id", "", false, "", ""},
		{"/users/{user-id}/role-mappings/realm", "user-realm-role-mapping", "user", "role", "{realm}/users/{user-id}/role-mappings/realm", "POST", "", "", true, "", ""},
		{"/users/{user-id}/role-mappings/clients/{client}", "user-client-role-mapping", "user", "role", "{realm}/users/{user-id}/role-mappings/clients/{client}", "POST", "", "", true, "", ""},
		{"/users/{user-id}/federated-identity", "user-federated-identity", "user", "identityprovider", "{realm}/users/{user-id}/federated-identity/{provider}", "POST", "identityProvider", "", false, "", ""},
		{"/groups/{group-id}/role-mappings/realm", "group-realm-role-mapping", "group", "role", "{realm}/groups/{group-id}/role-mappings/realm", "POST", "", "", true, "", ""},
		{"/groups/{group-id}/role-mappings/clients/{client}", "group-client-role-mapping", "group", "role", "{realm}/groups/{group-id}/role-mappings/clients/{client}", "POST", "", "", true, "", ""},
		{"/roles-by-id/{role-id}/composites", "role-composite-mapping", "role", "role", "{realm}/roles-by-id/{role-id}/composites", "POST", "", "", true, "", ""},
		{"/clients/{client-uuid}/roles/{role-name}/composites", "client-role-composite", "role", "role", "{realm}/clients/{client-uuid}/roles/{role-name}/composites", "POST", "", "", true, "", ""},
		{"/client-scopes/{client-scope-id}/scope-mappings/realm", "client-scope-realm-role-mapping", "clientscope", "role", "{realm}/client-scopes/{client-scope-id}/scope-mappings/realm", "POST", "", "", true, "", ""},
		{"/client-scopes/{client-scope-id}/scope-mappings/clients/{client}", "client-scope-client-role-mapping", "clientscope", "role", "{realm}/client-scopes/{client-scope-id}/scope-mappings/clients/{client}", "POST", "", "", true, "", ""},
		{"/clients/{client-uuid}/scope-mappings/realm", "client-realm-scope-mapping", "client", "role", "{realm}/clients/{client-uuid}/scope-mappings/realm", "POST", "", "", true, "", ""},
		{"/clients/{client-uuid}/scope-mappings/clients/{client}", "client-client-scope-mapping", "client", "role", "{realm}/clients/{client-uuid}/scope-mappings/clients/{client}", "POST", "", "", true, "", ""},
		{"/organizations/{org-id}/members", "organization-member", "organization", "user", "{realm}/organizations/{org-id}/members", "POST", "", "id", false, "{realm}/organizations/{org-id}/members/{member-id}", "member-id"},
		{"/organizations/{org-id}/identity-providers", "organization-identity-provider", "organization", "identityprovider", "{realm}/organizations/{org-id}/identity-providers", "POST", "", "alias", false, "{realm}/organizations/{org-id}/identity-providers/{alias}", "alias"},
		{"/default-groups", "default-group-membership", "realm", "group", "{realm}/default-groups/{groupId}", "PUT", "id", "", false, "", ""},
		{"/default-default-client-scopes", "realm-default-client-scope", "realm", "clientscope", "{realm}/default-default-client-scopes/{clientScopeId}", "PUT", "id", "", false, "", ""},
		{"/default-optional-client-scopes", "realm-optional-client-scope", "realm", "clientscope", "{realm}/default-optional-client-scopes/{clientScopeId}", "PUT", "id", "", false, "", ""},
		{"/clients/{client-uuid}/default-client-scopes", "client-default-scope", "client", "clientscope", "{realm}/clients/{client-uuid}/default-client-scopes/{clientScopeId}", "PUT", "id", "", false, "", ""},
		{"/clients/{client-uuid}/optional-client-scopes", "client-optional-scope", "client", "clientscope", "{realm}/clients/{client-uuid}/optional-client-scopes/{clientScopeId}", "PUT", "id", "", false, "", ""},
	}

	kinds := make([]RelationshipKind, 0, len(patterns))
	for _, p := range patterns {
		deleteTemplate := p.deleteTemplate
		if deleteTemplate == "" {
			deleteTemplate = p.relationshipTemplate
		}
		kind := RelationshipKind{
			Name:               p.name,
			ResourceA:          p.resourceA,
			ResourceB:          p.resourceB,
			Direction:          "one-way",
			ReadMethod:         http.MethodGet,
			ReadPath:           p.match,
			WriteMethod:        p.relationshipMethod,
			WriteTemplate:      p.relationshipTemplate,
			DeleteMethod:       http.MethodDelete,
			DeleteTemplate:     deleteTemplate,
			DeleteItemParam:    p.deleteItemParam,
			DeletePayloadField: p.payloadField,
			ItemParamName:      p.itemParamName,
			PayloadField:       p.payloadField,
			BulkPayload:        p.bulkPayload,
			ParamTypes:         inferParamTypesForTemplate(p.relationshipTemplate),
		}
		kinds = append(kinds, kind)
	}
	return kinds
}

// extractStringValue returns a string value from an interface{} when the payload
// is a string or a map with the given key.
func extractStringValue(payload interface{}, key string) string {
	if payload == nil {
		return ""
	}
	if s, ok := payload.(string); ok {
		return strings.TrimSpace(s)
	}
	if m, ok := payload.(map[string]interface{}); ok {
		if s, ok := m[key].(string); ok {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

// BuildDeleteOperation creates a delete operation for an existing relationship.
func BuildDeleteOperation(actual manifest.RelationshipOperation, kind RelationshipKind) (manifest.RelationshipOperation, error) {
	params := make(map[string]string, len(actual.PathParams))
	for k, v := range actual.PathParams {
		params[k] = v
	}

	var parsedPayload interface{}
	if len(actual.Data) > 0 {
		_ = json.Unmarshal(actual.Data, &parsedPayload)
	}

	deletePayload := actual.Data
	if kind.DeleteItemParam != "" {
		value := extractStringValue(parsedPayload, kind.DeletePayloadField)
		if value == "" {
			return manifest.RelationshipOperation{}, fmt.Errorf("cannot build delete for %s: missing %s in payload", kind.Name, kind.DeleteItemParam)
		}
		params[kind.DeleteItemParam] = value
		deletePayload = nil
	} else if !kind.BulkPayload {
		deletePayload = nil
	}

	return manifest.NewRelationshipOperation(kind.DeleteTemplate, kind.DeleteMethod, params, deletePayload)
}
