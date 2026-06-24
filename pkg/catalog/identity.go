package catalog

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/thedataflows/keycloak-cli/pkg/catalog/internal"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

// ResourceIdentity describes how a resource type is identified in manifests and URLs.
type ResourceIdentity struct {
	Type             string
	IDParam          string
	PrimaryKey       string
	IdentifierFields []string
	DisplayField     string
}

var defaultPlaceholderResource = map[string]string{
	"realm":           "realm",
	"user-id":         "user",
	"client-uuid":     "client",
	"client":          "client",
	"group-id":        "group",
	"groupId":         "group",
	"role-name":       "role",
	"role-id":         "role",
	"roleName":        "role",
	"alias":           "identityprovider",
	"org-id":          "organization",
	"client-scope-id": "clientscope",
	"flowAlias":       "authenticationflow",
	"executionId":     "authenticationexecution",
}

var defaultPlaceholderFields = map[string][]string{
	"realm":           {"realm"},
	"user-id":         {"id", "username"},
	"client-uuid":     {"id", "clientId"},
	"client":          {"clientId", "id"},
	"group-id":        {"id", "name", "path"},
	"groupId":         {"id", "name", "path"},
	"role-name":       {"name", "id"},
	"role-id":         {"id", "name"},
	"roleName":        {"name", "id"},
	"alias":           {"alias", "id"},
	"org-id":          {"alias", "id"},
	"client-scope-id": {"id", "name"},
	"flowAlias":       {"alias", "id"},
	"executionId":     {"id", "alias"},
}

var defaultDisplayFields = map[string]string{
	"realm":                   "realm",
	"user":                    "username",
	"client":                  "clientId",
	"role":                    "name",
	"group":                   "name",
	"identityprovider":        "alias",
	"organization":            "alias",
	"clientscope":             "name",
	"authenticationflow":      "alias",
	"authenticationexecution": "id",
	"protocolmapper":          "name",
}

var defaultVolatileFields = map[string][]string{
	"realm":                   {"id"},
	"user":                    {"id", "createdTimestamp", "access", "origin"},
	"client":                  {"id", "containerId", "access"},
	"role":                    {"id", "containerId"},
	"group":                   {"id", "access"},
	"identityprovider":        {"id", "types"},
	"organization":            {"id"},
	"clientscope":             {"id"},
	"authenticationflow":      {"id"},
	"authenticationexecution": {"id"},
	"protocolmapper":          {"id"},
}

var defaultWriteOnlyFields = map[string][]string{
	"user":             {"credentials"},
	"client":           {"clientSecret"},
	"identityprovider": {"clientSecret"},
}

// FieldOverride describes per-resource-type additions or removals to the
// volatile and write-only field sets used for normalization and apply.
type FieldOverride struct {
	Type            string   `yaml:"type"`
	AddVolatile     []string `yaml:"addVolatile,omitempty"`
	RemoveVolatile  []string `yaml:"removeVolatile,omitempty"`
	AddWriteOnly    []string `yaml:"addWriteOnly,omitempty"`
	RemoveWriteOnly []string `yaml:"removeWriteOnly,omitempty"`
}

// FieldOverridesFile is the top-level shape of a field overrides YAML file.
type FieldOverridesFile struct {
	Overrides []FieldOverride `yaml:"overrides"`
}

// LoadFieldOverrides reads field overrides from a YAML file.
func LoadFieldOverrides(path string) ([]FieldOverride, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read field overrides: %w", err)
	}

	var file FieldOverridesFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse field overrides: %w", err)
	}
	return file.Overrides, nil
}

// ApplyFieldOverrides mutates the default volatile and write-only field maps.
func ApplyFieldOverrides(overrides []FieldOverride) error {
	for _, override := range overrides {
		if override.Type == "" {
			return fmt.Errorf("field override missing required field 'type'")
		}
		defaultVolatileFields[override.Type] = patchFieldList(
			defaultVolatileFields[override.Type],
			override.AddVolatile,
			override.RemoveVolatile,
		)
		defaultWriteOnlyFields[override.Type] = patchFieldList(
			defaultWriteOnlyFields[override.Type],
			override.AddWriteOnly,
			override.RemoveWriteOnly,
		)
	}
	return nil
}

func patchFieldList(current, add, remove []string) []string {
	seen := make(map[string]struct{})
	for _, f := range current {
		seen[f] = struct{}{}
	}
	for _, f := range add {
		seen[f] = struct{}{}
	}
	for _, f := range remove {
		delete(seen, f)
	}
	result := make([]string, 0, len(seen))
	for f := range seen {
		result = append(result, f)
	}
	sort.Strings(result)
	return result
}

// InstallDefaultFieldOverrides loads field overrides from the directory
// containing the spec. It is safe to call when no override file exists.
func InstallDefaultFieldOverrides(specPath string) error {
	specDir := filepath.Dir(specPath)
	overridePath := filepath.Join(specDir, "field-overrides.yaml")

	// Install catalog defaults unconditionally so that spec-derived volatile and
	// write-only fields are available even when no override file exists.
	manifest.InstallVolatileFields(volatileFieldMap())
	manifest.InstallWriteOnlyFields(writeOnlyFieldMap())

	overrides, err := LoadFieldOverrides(overridePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if err := ApplyFieldOverrides(overrides); err != nil {
		return err
	}
	manifest.InstallVolatileFields(volatileFieldMap())
	manifest.InstallWriteOnlyFields(writeOnlyFieldMap())
	return nil
}

func volatileFieldMap() map[string]map[string]struct{} {
	result := make(map[string]map[string]struct{})
	for resourceType, fields := range defaultVolatileFields {
		set := make(map[string]struct{}, len(fields))
		for _, f := range fields {
			set[f] = struct{}{}
		}
		result[resourceType] = set
	}
	return result
}

func writeOnlyFieldMap() map[string]map[string]struct{} {
	result := make(map[string]map[string]struct{})
	for resourceType, fields := range defaultWriteOnlyFields {
		set := make(map[string]struct{}, len(fields))
		for _, f := range fields {
			set[f] = struct{}{}
		}
		result[resourceType] = set
	}
	return result
}

// PlaceholderToResourceType maps path-parameter placeholders to their resource types.
func (s *Spec) PlaceholderToResourceType() (map[string]string, error) {
	if s == nil {
		return nil, internal.ErrSpecNotInitialized
	}
	contracts, err := s.ResourceContracts()
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, len(defaultPlaceholderResource))
	for k, v := range defaultPlaceholderResource {
		result[k] = v
	}
	s.ForEachOperation(func(path, method string, op *v3.Operation, item *v3.PathItem) {
		if method != http.MethodPost {
			return
		}
		parts := strings.Split(path, "/")
		for i, part := range parts {
			if !strings.HasPrefix(part, "{") || !strings.HasSuffix(part, "}") {
				continue
			}
			placeholder := part[1 : len(part)-1]
			if placeholder == "realm" || result[placeholder] != "" {
				continue
			}
			resourceType := inferResourceTypeFromPlaceholder(placeholder, parts, i, contracts)
			if resourceType != "" {
				result[placeholder] = resourceType
			}
		}
	})
	return result, nil
}

func inferResourceTypeFromPlaceholder(placeholder string, parts []string, idx int, contracts map[string]ResourceContract) string {
	if resourceType, ok := defaultPlaceholderResource[placeholder]; ok {
		return resourceType
	}
	base := stripSuffixes(placeholder, []string{"-id", "-uuid", "-name", "Id", "ID", "Uuid", "UUID", "Name"})
	candidate := normalizeSchemaResourceType(base)
	if _, ok := contracts[candidate]; ok {
		return candidate
	}
	if idx > 0 {
		prev := parts[idx-1]
		prev = strings.TrimPrefix(prev, "{")
		prev = strings.TrimSuffix(prev, "}")
		candidate = singularOf(prev)
		if _, ok := contracts[candidate]; ok {
			return candidate
		}
	}
	return ""
}

// ResourceIdentities derives identity metadata for every discovered resource type.
func (s *Spec) ResourceIdentities() (map[string]ResourceIdentity, error) {
	if s == nil {
		return nil, internal.ErrSpecNotInitialized
	}
	contracts, err := s.ResourceContracts()
	if err != nil {
		return nil, err
	}
	identities := make(map[string]ResourceIdentity)
	for resourceType, contract := range contracts {
		idParam, primaryKey, fields := s.inferIdentity(resourceType, contract)
		display := defaultDisplayFields[resourceType]
		if display == "" {
			display = primaryKey
		}
		identities[resourceType] = ResourceIdentity{
			Type:             resourceType,
			IDParam:          idParam,
			PrimaryKey:       primaryKey,
			IdentifierFields: fields,
			DisplayField:     display,
		}
	}
	return identities, nil
}

func (s *Spec) inferIdentity(resourceType string, contract ResourceContract) (idParam, primaryKey string, fields []string) {
	idParam = ""
	primaryKey = "id"
	fields = []string{"id"}

	var singleGets []OperationContract
	for _, op := range contract.AllOperations[http.MethodGet] {
		if !isCollectionEndpoint(op.Path) {
			singleGets = append(singleGets, op)
		}
	}
	sort.Slice(singleGets, func(i, j int) bool {
		return len(strings.Split(singleGets[i].Path, "/")) < len(strings.Split(singleGets[j].Path, "/"))
	})

	if len(singleGets) > 0 {
		idParam = primaryIdentifierParam(singleGets[0].Path)
	}
	if idParam == "" {
		if op, ok := contract.Operations[http.MethodGet]; ok && !isCollectionEndpoint(op.Path) {
			idParam = primaryIdentifierParam(op.Path)
		}
	}
	if idParam == "" {
		if op, ok := contract.Operations[http.MethodPost]; ok {
			idParam = primaryIdentifierParam(op.Path)
		}
	}

	if candidateFields, ok := defaultPlaceholderFields[idParam]; ok {
		fields = candidateFields
		primaryKey = candidateFields[0]
	} else if idParam != "" {
		fields = []string{idParam, "id"}
		primaryKey = idParam
	} else {
		fields = inferIdentifierFieldsFromSchema(contract.Schema)
		if len(fields) > 0 {
			primaryKey = fields[0]
		}
	}

	if preferredDisplay := defaultDisplayFields[resourceType]; preferredDisplay != "" {
		primaryKey = preferredDisplay
	}

	return idParam, primaryKey, fields
}

func inferIdentifierFieldsFromSchema(proxy *base.SchemaProxy) []string {
	if proxy == nil {
		return []string{"id"}
	}
	schema, err := proxy.BuildSchema()
	if err != nil {
		return []string{"id"}
	}
	required := schema.Required
	var candidates []string
	for _, key := range []string{"id", "name", "alias", "clientId", "username", "realm"} {
		if containsString(required, key) || schemaPropertyExists(schema, key) {
			candidates = append(candidates, key)
		}
	}
	if len(candidates) == 0 {
		return []string{"id"}
	}
	return candidates
}

// VolatileFields returns server-managed fields that should be stripped before applying or comparing.
func (s *Spec) VolatileFields(resourceType string) ([]string, error) {
	if s == nil {
		return nil, internal.ErrSpecNotInitialized
	}
	fields := map[string]struct{}{}
	for _, f := range defaultVolatileFields[resourceType] {
		fields[f] = struct{}{}
	}
	contracts, err := s.ResourceContracts()
	if err != nil {
		return nil, err
	}
	contract, ok := contracts[resourceType]
	if !ok {
		return sortedStringSet(fields), nil
	}
	for _, f := range schemaReadOnlyFields(contract.Schema) {
		fields[f] = struct{}{}
	}
	return sortedStringSet(fields), nil
}

// WriteOnlyFields returns fields that should never be exposed in fetched output.
func (s *Spec) WriteOnlyFields(resourceType string) ([]string, error) {
	if s == nil {
		return nil, internal.ErrSpecNotInitialized
	}
	fields := map[string]struct{}{}
	for _, f := range defaultWriteOnlyFields[resourceType] {
		fields[f] = struct{}{}
	}
	contracts, err := s.ResourceContracts()
	if err != nil {
		return nil, err
	}
	contract, ok := contracts[resourceType]
	if !ok {
		return sortedStringSet(fields), nil
	}
	for _, f := range schemaWriteOnlyFields(contract.Schema) {
		fields[f] = struct{}{}
	}
	return sortedStringSet(fields), nil
}

// IdentifierOf returns the best identifier value for a resource using its identity model.
func IdentifierOf(resource manifest.Resource, identity ResourceIdentity) string {
	return firstStringField(resource.Data, identity.IdentifierFields)
}

// NameOf returns the human-readable name for a resource.
func NameOf(resource manifest.Resource, identity ResourceIdentity) string {
	if identity.DisplayField != "" {
		if v := stringField(resource.Data, identity.DisplayField); v != "" {
			return v
		}
	}
	return firstStringField(resource.Data, []string{"name", "alias", "clientId", "username", "realm", "id"})
}

// DisplayNameOf returns a display name, falling back to the identifier.
func DisplayNameOf(resource manifest.Resource, identity ResourceIdentity) string {
	if name := NameOf(resource, identity); name != "" {
		return name
	}
	return IdentifierOf(resource, identity)
}

func stripSuffixes(s string, suffixes []string) string {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return s[:len(s)-len(suffix)]
		}
	}
	return s
}

func singularOf(s string) string {
	if strings.HasSuffix(s, "ies") && len(s) > 3 {
		return s[:len(s)-3] + "y"
	}
	s = strings.TrimSuffix(s, "es")
	s = strings.TrimSuffix(s, "s")
	return s
}

func containsString(list []string, want string) bool {
	for _, item := range list {
		if item == want {
			return true
		}
	}
	return false
}

func schemaPropertyExists(schema *base.Schema, key string) bool {
	if schema == nil || schema.Properties == nil {
		return false
	}
	for pair := schema.Properties.First(); pair != nil; pair = pair.Next() {
		if pair.Key() == key {
			return true
		}
	}
	return false
}

func schemaReadOnlyFields(proxy *base.SchemaProxy) []string {
	schema, err := proxy.BuildSchema()
	if err != nil {
		return nil
	}
	var fields []string
	if schema.Properties != nil {
		for pair := schema.Properties.First(); pair != nil; pair = pair.Next() {
			prop, err := pair.Value().BuildSchema()
			if err != nil {
				continue
			}
			if prop.ReadOnly != nil && *prop.ReadOnly {
				fields = append(fields, pair.Key())
			}
		}
	}
	return fields
}

func schemaWriteOnlyFields(proxy *base.SchemaProxy) []string {
	schema, err := proxy.BuildSchema()
	if err != nil {
		return nil
	}
	var fields []string
	if schema.Properties != nil {
		for pair := schema.Properties.First(); pair != nil; pair = pair.Next() {
			prop, err := pair.Value().BuildSchema()
			if err != nil {
				continue
			}
			if prop.WriteOnly != nil && *prop.WriteOnly {
				fields = append(fields, pair.Key())
			}
		}
	}
	return fields
}

func sortedStringSet(set map[string]struct{}) []string {
	result := make([]string, 0, len(set))
	for k := range set {
		result = append(result, k)
	}
	sort.Strings(result)
	return result
}

func stringField(data map[string]interface{}, key string) string {
	if data == nil {
		return ""
	}
	if s, ok := data[key].(string); ok {
		return strings.TrimSpace(s)
	}
	return ""
}

func firstStringField(data map[string]interface{}, keys []string) string {
	for _, key := range keys {
		if v := stringField(data, key); v != "" {
			return v
		}
	}
	return ""
}
