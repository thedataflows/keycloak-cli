package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
)

const relationshipPathPrefix = "/admin/realms/"

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

type RelationshipManifest struct {
	Relationships []RelationshipOperation `json:"relationships"`
}

type ResourceMismatch struct {
	Expected Resource `json:"expected"`
	Actual   Resource `json:"actual"`
}

type ComparisonReport struct {
	Match                   bool                    `json:"match"`
	MissingResources        []Resource              `json:"missingResources,omitempty"`
	MismatchedResources     []ResourceMismatch      `json:"mismatchedResources,omitempty"`
	UnexpectedResources     []Resource              `json:"unexpectedResources,omitempty"`
	MissingRelationships    []RelationshipOperation `json:"missingRelationships,omitempty"`
	UnexpectedRelationships []RelationshipOperation `json:"unexpectedRelationships,omitempty"`
}

type LoadResult struct {
	Resources     []Resource
	Relationships []RelationshipOperation
	Skipped       []SkippedFile
}

type SkippedFile struct {
	Path   string
	Reason string
}

// ParseResources decodes a resource manifest from JSON or YAML.
func ParseResources(data []byte) ([]Resource, bool, error) {
	resources, ok, _, err := parseResourceList(data, json.Unmarshal)
	if ok || err != nil {
		return resources, ok, err
	}

	resources, ok, _, err = parseResourceList(data, yaml.Unmarshal)
	if ok || err != nil {
		return resources, ok, err
	}

	resources, ok, _, err = parseSingleResource(data, json.Unmarshal)
	if ok || err != nil {
		return resources, ok, err
	}

	resources, ok, _, err = parseSingleResource(data, yaml.Unmarshal)
	return resources, ok, err
}

func ParseRelationshipManifest(data []byte) (*RelationshipManifest, bool, error) {
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, false, err
	}

	if _, ok := envelope["relationships"]; !ok {
		return nil, false, nil
	}

	if len(envelope) > 1 {
		extra := make([]string, 0, len(envelope)-1)
		for key := range envelope {
			if key == "relationships" {
				continue
			}
			extra = append(extra, key)
		}
		sort.Strings(extra)
		return nil, false, fmt.Errorf("unexpected keys in relationship manifest: %s", strings.Join(extra, ", "))
	}

	var manifest RelationshipManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, false, err
	}

	for i := range manifest.Relationships {
		manifest.Relationships[i].Path = normalizeRelationshipPath(manifest.Relationships[i].Path)
	}

	return &manifest, true, nil
}

func LoadPaths(paths []string) (LoadResult, error) {
	files, err := collectManifestFiles(paths)
	if err != nil {
		return LoadResult{}, err
	}

	var result LoadResult

	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			return LoadResult{}, err
		}

		resources, relationships, err := parseManifestData(data)
		if err != nil {
			result.Skipped = append(result.Skipped, SkippedFile{Path: path, Reason: err.Error()})
			continue
		}

		result.Resources = append(result.Resources, resources...)
		result.Relationships = append(result.Relationships, relationships...)
	}

	return result, nil
}

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

func (r *RelationshipOperation) RebuildPath() {
	r.Path = buildActualRelationshipPath(r.Template, r.PathParams)
}

// ValidateResources ensures each resource has the minimum manifest fields required.
func ValidateResources(resources []Resource) error {
	var errs []error
	for index, resource := range resources {
		switch {
		case strings.TrimSpace(resource.Type) == "":
			errs = append(errs, fmt.Errorf("resource at index %d has empty type", index))
		case len(resource.Data) == 0:
			errs = append(errs, fmt.Errorf("resource '%s' missing data", resource.Type))
		}
	}
	return errors.Join(errs...)
}

// SortResources applies the current resource dependency ordering used by uploads.
// If priorityMap is nil, falls back to hardcoded priorities.
func SortResources(resources []Resource, priorityMap map[string]int) []Resource {
	if len(resources) == 0 {
		return resources
	}
	sorted := make([]Resource, len(resources))
	copy(sorted, resources)
	sort.SliceStable(sorted, func(left, right int) bool {
		return resolvePriority(sorted[left].Type, priorityMap) < resolvePriority(sorted[right].Type, priorityMap)
	})
	return sorted
}

func resolvePriority(resourceType string, priorityMap map[string]int) int {
	if priorityMap != nil {
		if p, ok := priorityMap[resourceType]; ok {
			return p
		}
	}
	return dependencyPriority(resourceType)
}

// NormalizeRoundTrip returns canonical resource and relationship slices for
// round-trip comparisons. It removes volatile server-managed resource fields,
// canonicalizes relationship payload JSON, and sorts both collections.
func NormalizeRoundTrip(resources []Resource, relationships []RelationshipOperation) ([]Resource, []RelationshipOperation) {
	index := buildResourceIndex(resources)
	return normalizeResourcesForRoundTrip(resources, nil), normalizeRelationshipsForRoundTrip(relationships, index)
}

// NormalizeForApply returns a canonical manifest suitable for re-applying to a
// different Keycloak server. Like NormalizeRoundTrip it strips volatile fields,
// but it preserves the id of any resource that is referenced by another
// resource in the manifest so that inline UUID references can be remapped.
func NormalizeForApply(resources []Resource, relationships []RelationshipOperation) ([]Resource, []RelationshipOperation) {
	preserveIDs := collectReferencedIDs(resources)
	index := buildResourceIndex(resources)
	return normalizeResourcesForRoundTrip(resources, preserveIDs), normalizeRelationshipsForRoundTrip(relationships, index)
}

// CompareRoundTrip compares expected manifests with actual fetched state using
// round-trip normalization. Resource comparison is subset-based per matching
// resource key, while relationship comparison is exact on normalized entries.
func CompareRoundTrip(expectedResources []Resource, expectedRelationships []RelationshipOperation, actualResources []Resource, actualRelationships []RelationshipOperation) ComparisonReport {
	normalizedExpectedResources, normalizedExpectedRelationships := NormalizeRoundTrip(expectedResources, expectedRelationships)
	normalizedActualResources, normalizedActualRelationships := NormalizeRoundTrip(actualResources, actualRelationships)

	report := ComparisonReport{
		Match:                   true,
		MissingResources:        make([]Resource, 0),
		MismatchedResources:     make([]ResourceMismatch, 0),
		UnexpectedResources:     make([]Resource, 0),
		MissingRelationships:    make([]RelationshipOperation, 0),
		UnexpectedRelationships: make([]RelationshipOperation, 0),
	}

	actualByKey := make(map[string]Resource, len(normalizedActualResources))
	expectedKeys := make(map[string]struct{}, len(normalizedExpectedResources))
	for _, resource := range normalizedActualResources {
		if isIgnoredActualResource(resource) {
			continue
		}
		actualByKey[roundTripResourceKey(resource)] = resource
	}
	for _, resource := range normalizedExpectedResources {
		key := roundTripResourceKey(resource)
		expectedKeys[key] = struct{}{}
		actual, ok := actualByKey[key]
		if !ok {
			report.MissingResources = append(report.MissingResources, resource)
			continue
		}
		if !isValueSubset(resource.Data, actual.Data) {
			report.MismatchedResources = append(report.MismatchedResources, ResourceMismatch{Expected: resource, Actual: actual})
		}
	}
	for _, resource := range normalizedActualResources {
		key := roundTripResourceKey(resource)
		if _, ok := expectedKeys[key]; ok {
			continue
		}
		if isIgnoredActualResource(resource) {
			continue
		}
		report.UnexpectedResources = append(report.UnexpectedResources, resource)
	}

	expectedRelationshipCounts := countRelationships(normalizedExpectedRelationships)
	actualRelationshipCounts := countRelationships(normalizedActualRelationships)
	for _, relationship := range normalizedExpectedRelationships {
		key := relationshipSortKey(relationship)
		if actualRelationshipCounts[key] == 0 {
			report.MissingRelationships = append(report.MissingRelationships, relationship)
			continue
		}
		actualRelationshipCounts[key]--
	}
	for _, relationship := range normalizedActualRelationships {
		key := relationshipSortKey(relationship)
		if expectedRelationshipCounts[key] == 0 {
			report.UnexpectedRelationships = append(report.UnexpectedRelationships, relationship)
			continue
		}
		expectedRelationshipCounts[key]--
	}
	if len(report.MissingResources) > 0 || len(report.MismatchedResources) > 0 || len(report.UnexpectedResources) > 0 || len(report.MissingRelationships) > 0 || len(report.UnexpectedRelationships) > 0 {
		report.Match = false
	}

	return report
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
		return firstStringField(r.Data, fields)
	}
	return stringField(r.Data, "id")
}

var nameFields = map[string][]string{
	"realm":  {"realm"},
	"user":   {"username"},
	"client": {"clientId"},
}

func (r Resource) Name() string {
	if fields, ok := nameFields[r.Type]; ok {
		return firstStringField(r.Data, fields)
	}
	return firstStringField(r.Data, []string{"name", "alias"})
}

var displayNameFields = []string{"name", "clientId", "username", "alias", "realm", "id"}

func (r Resource) DisplayName() string {
	return firstStringField(r.Data, displayNameFields)
}

func firstStringField(data map[string]interface{}, keys []string) string {
	for _, key := range keys {
		if value := stringField(data, key); value != "" {
			return value
		}
	}
	return ""
}

func dependencyPriority(resourceType string) int {
	switch resourceType {
	case "realm":
		return 1
	case "clientscope", "authenticationflow":
		return 2
	case "group", "role":
		return 3
	case "identityprovider":
		return 4
	case "client":
		return 5
	case "organization":
		return 6
	case "user":
		return 7
	default:
		return 10
	}
}

type unmarshalFunc func(data []byte, v interface{}) error

func parseResourceList(data []byte, unmarshal unmarshalFunc) ([]Resource, bool, string, error) {
	var resources []Resource
	if err := unmarshal(data, &resources); err != nil {
		return nil, false, err.Error(), nil
	}
	if len(resources) == 0 {
		return nil, false, "empty", nil
	}
	if err := ValidateResources(resources); err != nil {
		return nil, true, "", err
	}
	return resources, true, "", nil
}

func parseSingleResource(data []byte, unmarshal unmarshalFunc) ([]Resource, bool, string, error) {
	var resource Resource
	if err := unmarshal(data, &resource); err != nil {
		return nil, false, err.Error(), nil
	}
	if strings.TrimSpace(resource.Type) == "" {
		return nil, false, "missing 'type' field", nil
	}
	if err := ValidateResources([]Resource{resource}); err != nil {
		return nil, true, "", err
	}
	return []Resource{resource}, true, "", nil
}

func stringField(data map[string]interface{}, key string) string {
	if data == nil {
		return ""
	}

	value, _ := data[key].(string)
	return value
}

type resourceIndex map[string]map[string]string

func isBuiltinResource(resource Resource) bool {
	if builtin, ok := resource.Data["builtIn"].(bool); ok {
		return builtin
	}
	return false
}

func isIgnoredActualResource(resource Resource) bool {
	if isBuiltinResource(resource) {
		return true
	}
	if IsBuiltInResource != nil && IsBuiltInResource(resource) {
		return true
	}
	return false
}

func normalizeResourcesForRoundTrip(resources []Resource, preserveIDs map[string]struct{}) []Resource {
	if len(resources) == 0 {
		return nil
	}

	normalized := make([]Resource, 0, len(resources))
	for _, resource := range resources {
		normalized = append(normalized, cloneResource(resource, preserveIDs))
	}

	sort.Slice(normalized, func(left, right int) bool {
		return resourceSortKey(normalized[left]) < resourceSortKey(normalized[right])
	})

	return normalized
}

// cloneResource returns a copy of the resource with server-managed and volatile
// fields removed. If preserveIDs is non-nil, any resource whose "id" appears in
// the set keeps its id so that inline references from other resources can be
// remapped during apply.
func cloneResource(resource Resource, preserveIDs map[string]struct{}) Resource {
	clone := Resource{
		Type:       resource.Type,
		Realm:      strings.TrimSpace(resource.Realm),
		Delete:     resource.Delete,
		ParentType: resource.ParentType,
	}
	if len(resource.Data) > 0 {
		clone.Data = normalizeMapForRoundTrip(resource.Data, true)
	}
	if clone.Realm == "" {
		clone.Realm = stringField(clone.Data, "realm")
	}
	if clone.Data != nil {
		if fields, ok := writeOnlyResourceFields[clone.Type]; ok {
			deleteWriteOnlyFields(clone.Data, fields)
		}
		if preserveIDs != nil {
			if id := stringField(resource.Data, "id"); id != "" {
				if _, ok := preserveIDs[id]; ok {
					clone.Data["id"] = id
				}
			}
		}
	}
	return clone
}

// StripVolatileFields returns a copy of the resource with server-managed and
// write-only fields removed, suitable for validation before apply.
func StripVolatileFields(resource Resource) Resource {
	return cloneResource(resource, nil)
}

func deleteWriteOnlyFields(data map[string]interface{}, fields map[string]struct{}) {
	for field := range fields {
		delete(data, field)
	}
	for _, value := range data {
		switch typed := value.(type) {
		case map[string]interface{}:
			deleteWriteOnlyFields(typed, fields)
		case []interface{}:
			for _, item := range typed {
				if m, ok := item.(map[string]interface{}); ok {
					deleteWriteOnlyFields(m, fields)
				}
			}
		}
	}
}

func normalizeRelationshipsForRoundTrip(relationships []RelationshipOperation, index resourceIndex) []RelationshipOperation {
	if len(relationships) == 0 {
		return nil
	}

	normalized := make([]RelationshipOperation, 0, len(relationships))
	for _, relationship := range relationships {
		clone := RelationshipOperation{
			Kind:     strings.TrimSpace(relationship.Kind),
			Delete:   relationship.Delete,
			Method:   strings.ToUpper(strings.TrimSpace(relationship.Method)),
			Template: strings.TrimPrefix(strings.TrimSpace(relationship.Template), "/"),
			Path:     normalizeRelationshipPath(relationship.Path),
		}
		if len(relationship.PathParams) > 0 {
			clone.PathParams = normalizeRelationshipPathParams(relationship, index)
			if clone.Template != "" {
				clone.Path = buildActualRelationshipPath(clone.Template, clone.PathParams)
			}
		}
		if len(relationship.Data) > 0 {
			clone.Data = normalizeRelationshipDataForRoundTrip(relationship, index)
		}
		normalized = append(normalized, clone)
	}

	sort.Slice(normalized, func(left, right int) bool {
		leftKey := relationshipSortKey(normalized[left])
		rightKey := relationshipSortKey(normalized[right])
		return leftKey < rightKey
	})

	return normalized
}

func buildResourceIndex(resources []Resource) resourceIndex {
	index := make(resourceIndex)
	for _, resource := range resources {
		display := strings.TrimSpace(resource.Name())
		if display == "" {
			display = strings.TrimSpace(resource.DisplayName())
		}
		if display == "" {
			continue
		}
		byType, exists := index[resource.Type]
		if !exists {
			byType = make(map[string]string)
			index[resource.Type] = byType
		}
		for _, identifier := range resourceLookupKeys(resource) {
			byType[identifier] = display
		}
	}
	return index
}

// collectReferencedIDs scans every resource Data for UUID-shaped strings that
// may be inline references to other resources. The returned set is used by
// NormalizeForApply to decide which resource ids must be kept during
// canonicalization so they can be remapped on the target server.
func collectReferencedIDs(resources []Resource) map[string]struct{} {
	ids := make(map[string]struct{})
	for _, r := range resources {
		for key, value := range r.Data {
			// Skip each resource's own id so self-references are not treated
			// as cross-resource references.
			if strings.EqualFold(key, "id") {
				continue
			}
			collectReferencedIDsFromValue(value, ids)
		}
	}
	return ids
}

func collectReferencedIDsFromValue(value interface{}, ids map[string]struct{}) {
	switch typed := value.(type) {
	case string:
		if isUUID(typed) {
			ids[typed] = struct{}{}
		}
	case map[string]interface{}:
		for _, v := range typed {
			collectReferencedIDsFromValue(v, ids)
		}
	case []interface{}:
		for _, v := range typed {
			collectReferencedIDsFromValue(v, ids)
		}
	}
}

func isUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	return s[8] == '-' && s[13] == '-' && s[18] == '-' && s[23] == '-'
}

func resourceLookupKeys(resource Resource) []string {
	id := strings.TrimSpace(resource.Identifier())
	dataID := strings.TrimSpace(stringField(resource.Data, "id"))
	if id == "" {
		if dataID != "" {
			return []string{dataID}
		}
		return nil
	}
	if dataID != "" && dataID != id {
		return []string{id, dataID}
	}
	return []string{id}
}

func normalizeMapForRoundTrip(values map[string]interface{}, stripVolatile bool) map[string]interface{} {
	if len(values) == 0 {
		return nil
	}

	result := make(map[string]interface{}, len(values))
	for key, value := range values {
		if stripVolatile && isVolatileField(key) {
			continue
		}
		result[key] = normalizeValueForRoundTrip(value, stripVolatile)
	}
	return result
}

func isVolatileField(name string) bool {
	_, ok := volatileRoundTripFields[strings.ToLower(strings.TrimSpace(name))]
	return ok
}

func normalizeValueForRoundTrip(value interface{}, stripVolatile bool) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		return normalizeMapForRoundTrip(typed, stripVolatile)
	case []interface{}:
		normalized := make([]interface{}, len(typed))
		keys := make([]string, len(typed))
		canSort := true
		for idx, item := range typed {
			normalized[idx] = normalizeValueForRoundTrip(item, stripVolatile)
			key, ok := sortIdentity(normalized[idx])
			if !ok {
				canSort = false
				continue
			}
			keys[idx] = key
		}
		if canSort {
			sort.SliceStable(normalized, func(left, right int) bool {
				return keys[left] < keys[right]
			})
		}
		return normalized
	default:
		return value
	}
}

func normalizeRelationshipPathParams(relationship RelationshipOperation, index resourceIndex) map[string]string {
	params := maps.Clone(relationship.PathParams)
	if params == nil {
		params = make(map[string]string)
	}
	for key, resourceType := range RelationshipParamTypes(relationship.Kind) {
		if resolved := index.lookup(resourceType, relationship.PathParams[key]); resolved != "" {
			params[key] = resolved
		}
	}
	return params
}

func normalizeRelationshipDataForRoundTrip(relationship RelationshipOperation, index resourceIndex) json.RawMessage {
	var payload interface{}
	if err := json.Unmarshal(relationship.Data, &payload); err != nil {
		return append(json.RawMessage(nil), relationship.Data...)
	}

	payload = rewriteRelationshipPayloadForRoundTrip(relationship.Kind, payload, index)
	payload = normalizeValueForRoundTrip(payload, true)

	normalized, err := json.Marshal(payload)
	if err != nil {
		return append(json.RawMessage(nil), relationship.Data...)
	}
	return normalized
}

func rewriteRelationshipPayloadForRoundTrip(kind string, payload interface{}, index resourceIndex) interface{} {
	switch kind {
	case "organization-member":
		identifier, ok := payload.(string)
		if !ok {
			return payload
		}
		if resolved := index.lookup("user", identifier); resolved != "" {
			return resolved
		}
	case "organization-identity-provider":
		identifier, ok := payload.(string)
		if !ok {
			return payload
		}
		if resolved := index.lookup("identityprovider", identifier); resolved != "" {
			return resolved
		}
	case "user-federated-identity":
		return normalizeFederatedIdentityPayload(payload, index)
	case "user-realm-role-mapping", "group-realm-role-mapping", "role-composite-mapping", "client-scope-realm-role-mapping", "client-scope-client-role-mapping", "role-scope-mapping":
		return normalizeRoleRelationshipPayload(payload, index)
	}
	return payload
}

func normalizeFederatedIdentityPayload(payload interface{}, index resourceIndex) interface{} {
	mapped, ok := payload.(map[string]interface{})
	if !ok {
		return payload
	}

	normalized := make(map[string]interface{})
	if provider := coalesceString(mapped, "identityProvider", "provider"); provider != "" {
		if resolved := index.lookup("identityprovider", provider); resolved != "" {
			provider = resolved
		}
		normalized["identityProvider"] = provider
	}
	if userID := stringValue(mapped, "userId"); userID != "" {
		normalized["userId"] = userID
	}
	if userName := stringValue(mapped, "userName"); userName != "" {
		normalized["userName"] = userName
	}
	if len(normalized) == 0 {
		return mapped
	}
	return normalized
}

func normalizeRoleRelationshipPayload(payload interface{}, index resourceIndex) interface{} {
	items, ok := payload.([]interface{})
	if !ok {
		return payload
	}

	normalized := make([]interface{}, 0, len(items))
	for _, item := range items {
		mapped, ok := item.(map[string]interface{})
		if !ok {
			normalized = append(normalized, item)
			continue
		}

		role := make(map[string]interface{})
		if name := stringValue(mapped, "name"); name != "" {
			role["name"] = name
		}
		if clientRole, _ := mapped["clientRole"].(bool); clientRole {
			role["clientRole"] = true
			if resolved := index.lookup("client", stringValue(mapped, "containerId")); resolved != "" {
				role["client"] = resolved
			}
		}
		if len(role) == 0 {
			normalized = append(normalized, mapped)
			continue
		}
		normalized = append(normalized, role)
	}
	return normalized
}

// RelationshipParamTypes resolves the resource types for path parameters of a
// relationship kind. The default implementation is used when the catalog package
// has not yet installed a registry. Assigning a replacement function allows
// catalog-driven overrides without creating an import cycle.
var RelationshipParamTypes = defaultRelationshipParamTypes

// IsBuiltInResource allows the catalog package to inject knowledge of built-in
// resources that should be excluded from round-trip comparison.
var IsBuiltInResource = func(Resource) bool { return false }

func defaultRelationshipParamTypes(kind string) map[string]string {
	switch kind {
	case "user-group-membership":
		return map[string]string{"user-id": "user", "groupId": "group"}
	case "user-realm-role-mapping":
		return map[string]string{"user-id": "user"}
	case "group-realm-role-mapping":
		return map[string]string{"group-id": "group"}
	case "role-composite-mapping":
		return map[string]string{"role-id": "role"}
	case "default-group-membership":
		return map[string]string{"groupId": "group"}
	case "realm-default-client-scope", "realm-optional-client-scope":
		return map[string]string{"clientScopeId": "clientscope"}
	case "client-default-scope", "client-optional-scope":
		return map[string]string{"client-uuid": "client", "clientScopeId": "clientscope"}
	case "client-scope-realm-role-mapping":
		return map[string]string{"client-scope-id": "clientscope"}
	case "client-scope-client-role-mapping":
		return map[string]string{"client-scope-id": "clientscope", "client": "client"}
	case "user-federated-identity":
		return map[string]string{"user-id": "user", "provider": "identityprovider"}
	case "organization-member", "organization-identity-provider":
		return map[string]string{"org-id": "organization"}
	default:
		return nil
	}
}

func resourceSortKey(resource Resource) string {
	encoded, _ := json.Marshal(resource.Data)
	return strings.Join([]string{resource.Type, resource.Realm, resource.Name(), resource.DisplayName(), string(encoded)}, "|")
}

func roundTripResourceKey(resource Resource) string {
	name := strings.TrimSpace(resource.Name())
	if name == "" {
		name = strings.TrimSpace(resource.DisplayName())
	}
	return strings.Join([]string{resource.Type, strings.TrimSpace(resource.Realm), name}, "|")
}

func relationshipSortKey(relationship RelationshipOperation) string {
	return strings.Join([]string{relationship.Kind, relationship.Method, relationship.Path, string(relationship.Data)}, "|")
}

var volatileRoundTripFields = map[string]struct{}{
	"id":               {},
	"internalid":       {},
	"containerid":      {},
	"access":           {},
	"origin":           {},
	"createdtimestamp": {},
}

var writeOnlyResourceFields = map[string]map[string]struct{}{
	"user":   {"credentials": {}},
	"client": {"clientSecret": {}},
}

// InstallVolatileFields replaces the global set of volatile round-trip field
// names. It is used by the catalog package after loading field overrides.
func InstallVolatileFields(fields map[string]map[string]struct{}) {
	merged := make(map[string]struct{})
	for name := range volatileRoundTripFields {
		merged[name] = struct{}{}
	}
	for _, set := range fields {
		for name := range set {
			merged[strings.ToLower(name)] = struct{}{}
		}
	}
	volatileRoundTripFields = merged
}

// InstallWriteOnlyFields replaces the per-resource-type write-only field sets.
// It is used by the catalog package after loading field overrides.
func InstallWriteOnlyFields(fields map[string]map[string]struct{}) {
	writeOnlyResourceFields = fields
}

func sortIdentity(value interface{}) (string, bool) {
	mapped, ok := value.(map[string]interface{})
	if !ok {
		return "", false
	}
	for _, key := range []string{"name", "username", "clientId", "alias", "realm", "provider", "identityProvider", "path", "type"} {
		if s, ok := mapped[key].(string); ok && strings.TrimSpace(s) != "" {
			return key + ":" + strings.TrimSpace(s), true
		}
	}
	return "", false
}

func coalesceString(m map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if v := stringValue(m, key); v != "" {
			return v
		}
	}
	return ""
}

func stringValue(m map[string]interface{}, key string) string {
	if s, ok := m[key].(string); ok {
		return strings.TrimSpace(s)
	}
	return ""
}

func (i resourceIndex) lookup(resourceType, identifier string) string {
	if i == nil {
		return ""
	}
	items, ok := i[resourceType]
	if !ok {
		return ""
	}
	return items[strings.TrimSpace(identifier)]
}

func isValueSubset(expected interface{}, actual interface{}) bool {
	switch expectedTyped := expected.(type) {
	case map[string]interface{}:
		actualTyped, ok := actual.(map[string]interface{})
		if !ok {
			return false
		}
		for key, expectedValue := range expectedTyped {
			actualValue, ok := actualTyped[key]
			if !ok || !isValueSubset(expectedValue, actualValue) {
				return false
			}
		}
		return true
	case []interface{}:
		actualTyped, ok := actual.([]interface{})
		if !ok || len(actualTyped) < len(expectedTyped) {
			return false
		}
		for idx := range expectedTyped {
			if !isValueSubset(expectedTyped[idx], actualTyped[idx]) {
				return false
			}
		}
		return true
	default:
		return fmt.Sprintf("%v", expected) == fmt.Sprintf("%v", actual)
	}
}

func countRelationships(relationships []RelationshipOperation) map[string]int {
	counts := make(map[string]int, len(relationships))
	for _, relationship := range relationships {
		counts[relationshipSortKey(relationship)]++
	}
	return counts
}

func normalizeRelationshipPath(path string) string {
	trimmed := strings.TrimSpace(path)
	trimmed = strings.TrimPrefix(trimmed, relationshipPathPrefix)
	trimmed = strings.TrimPrefix(trimmed, "/")
	return trimmed
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

func collectManifestFiles(paths []string) ([]string, error) {
	var files []string
	for _, input := range paths {
		info, err := os.Stat(input)
		if err != nil {
			return nil, err
		}

		if info.IsDir() {
			if err := filepath.WalkDir(input, func(path string, d fs.DirEntry, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}
				if !d.IsDir() && supportedManifestPath(path) {
					files = append(files, path)
				}
				return nil
			}); err != nil {
				return nil, err
			}
			continue
		}

		if supportedManifestPath(input) {
			files = append(files, input)
		}
	}

	sort.Strings(files)
	return files, nil
}

func supportedManifestPath(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".json" || ext == ".yaml" || ext == ".yml"
}

type resourceParser func([]byte) ([]Resource, bool, string, error)

type resourceParserAttempt struct {
	name   string
	parser resourceParser
}

func tryEnvelope(data []byte, unmarshal unmarshalFunc) ([]Resource, []RelationshipOperation, bool, string, error) {
	var envelope struct {
		Resources     []Resource              `json:"resources"`
		Relationships []RelationshipOperation `json:"relationships"`
	}
	if err := unmarshal(data, &envelope); err != nil {
		return nil, nil, false, err.Error(), nil
	}
	if len(envelope.Resources) == 0 && len(envelope.Relationships) == 0 {
		return nil, nil, false, "missing 'resources' or 'relationships' key", nil
	}
	for i := range envelope.Relationships {
		envelope.Relationships[i].Path = normalizeRelationshipPath(envelope.Relationships[i].Path)
	}
	return envelope.Resources, envelope.Relationships, true, "", nil
}

func parseManifestData(data []byte) ([]Resource, []RelationshipOperation, error) {
	trimmed := strings.TrimSpace(string(data))
	var reasons []string

	if strings.HasPrefix(trimmed, "{") {
		resources, relationships, ok, reason, err := tryEnvelope(data, json.Unmarshal)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			return resources, relationships, nil
		}
		if reason != "" {
			reasons = append(reasons, "JSON envelope: "+reason)
		}

		relationshipManifest, ok, err := ParseRelationshipManifest([]byte(trimmed))
		if err != nil {
			return nil, nil, fmt.Errorf("relationship manifest: %w", err)
		}
		if ok {
			return nil, relationshipManifest.Relationships, nil
		}
		reasons = append(reasons, "relationship manifest: missing 'relationships' key")
	}

	if !strings.HasPrefix(trimmed, "{") {
		resources, relationships, ok, reason, err := tryEnvelope(data, yaml.Unmarshal)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			return resources, relationships, nil
		}
		if reason != "" {
			reasons = append(reasons, "YAML envelope: "+reason)
		}
	}

	for _, attempt := range resourceParsers() {
		resources, ok, reason, err := attempt.parser(data)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", attempt.name, err)
		}
		if ok {
			return resources, nil, nil
		}
		if reason != "" {
			reasons = append(reasons, attempt.name+": "+reason)
		}
	}

	return nil, nil, fmt.Errorf("unsupported manifest format; attempted the following parsers:\n- %s", strings.Join(reasons, "\n- "))
}

func resourceParsers() []resourceParserAttempt {
	unmarshalers := []struct {
		name string
		fn   unmarshalFunc
	}{
		{"JSON", json.Unmarshal},
		{"YAML", yaml.Unmarshal},
	}

	parsers := make([]resourceParserAttempt, 0, len(unmarshalers)*2)
	for _, unmarshaler := range unmarshalers {
		fn := unmarshaler.fn
		parsers = append(parsers,
			resourceParserAttempt{name: unmarshaler.name + " resource list", parser: func(data []byte) ([]Resource, bool, string, error) {
				return parseResourceList(data, fn)
			}},
			resourceParserAttempt{name: unmarshaler.name + " single resource", parser: func(data []byte) ([]Resource, bool, string, error) {
				return parseSingleResource(data, fn)
			}},
		)
	}
	return parsers
}
