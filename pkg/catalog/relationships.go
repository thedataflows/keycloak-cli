package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

const relationshipPathPrefix = "/admin/realms/"

var relationshipPathIndicators = []string{
	"/role-mappings",
	"/scope-mappings",
	"/users/{user-id}/groups",
	"/groups/{group-id}/members",
	"/groups/{group-id}/role-mappings",
	"/roles-by-id/{role-id}/composites",
	"/federated-identity",
	"/default-groups",
	"/default-default-client-scopes",
	"/default-optional-client-scopes",
	"/default-client-scopes",
	"/optional-client-scopes",
	"/organizations",
}

var relationshipSchemaNames = map[string]struct{}{
	"RoleRepresentation":              {},
	"GroupRepresentation":             {},
	"MappingsRepresentation":          {},
	"FederatedIdentityRepresentation": {},
	"UserRepresentation":              {},
}

type RelationshipTemplate struct {
	Template       string
	Method         string
	Summary        string
	Body           map[string]interface{}
	RequiresBody   bool
	PathParamNames []string
}

func RelationshipTemplatePattern(template string) (string, error) {
	normalized := normalizeRelationshipPath(template)
	if normalized == "" {
		return "", fmt.Errorf("relationship template cannot be empty")
	}

	paramNames, err := extractTemplateParamNames(normalized)
	if err != nil {
		return "", err
	}
	return buildTemplatePattern(normalized, paramNames)
}

func ValidateRelationshipOperations(spec *Spec, ops []manifest.RelationshipOperation) error {
	if spec == nil {
		return fmt.Errorf("spec is required")
	}

	templates, err := CollectRelationshipTemplates(spec)
	if err != nil {
		return err
	}

	matchers := make([]templateMatcher, len(templates))
	for i := range templates {
		matcher, matcherErr := newTemplateMatcher(templates[i])
		if matcherErr != nil {
			return matcherErr
		}
		matchers[i] = matcher
	}

	validationErrors := make([]error, 0)
	for idx := range ops {
		if err := validateRelationshipOperation(spec, &ops[idx], matchers); err != nil {
			validationErrors = append(validationErrors, fmt.Errorf("relationships[%d]: %w", idx, err))
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}
	return errors.Join(validationErrors...)
}

func validateRelationshipOperation(spec *Spec, rel *manifest.RelationshipOperation, matchers []templateMatcher) error {
	normalized := normalizeRelationshipPath(rel.Path)
	if normalized == "" {
		return fmt.Errorf("path is required")
	}

	selected, err := selectRelationshipTemplate(rel.Delete, normalized, matchers)
	if err != nil {
		return err
	}

	rel.Template = selected.template.Template
	rel.Method = selected.template.Method
	rel.Delete = strings.EqualFold(selected.template.Method, http.MethodDelete)
	rel.Kind = inferRelationshipKind(selected.template.Template, selected.template.Method)
	if len(selected.params) > 0 {
		rel.PathParams = maps.Clone(selected.params)
	} else {
		rel.PathParams = nil
	}
	rel.RebuildPath()

	if selected.template.RequiresBody && len(rel.Data) == 0 {
		return fmt.Errorf("data is required")
	}
	if len(rel.Data) > 0 && !json.Valid(rel.Data) {
		return fmt.Errorf("data is not valid JSON")
	}

	request := RequestValidation{PathParams: rel.PathParams}
	if len(rel.Data) > 0 {
		var payload interface{}
		if err := json.Unmarshal(rel.Data, &payload); err != nil {
			return fmt.Errorf("decode payload: %w", err)
		}
		request.Body = payload
	}

	return spec.ValidateOperationRequest(relationshipPathPrefix+selected.template.Template, selected.template.Method, request)
}

func CollectRelationshipTemplates(spec *Spec) ([]RelationshipTemplate, error) {
	entries := make([]RelationshipTemplate, 0)
	seen := make(map[string]struct{})
	var collectedErrs []error

	spec.ForEachOperation(func(path, method string, operation *v3.Operation, item *v3.PathItem) {
		if !shouldConsiderMethod(method) {
			return
		}

		bodySchema, bodyErr := extractRequestSchema(operation)
		if bodyErr != nil {
			collectedErrs = append(collectedErrs, fmt.Errorf("%s %s: %w", method, path, bodyErr))
			return
		}

		key := method + " " + path
		if _, exists := seen[key]; exists {
			return
		}

		if !isRelationshipPath(path) {
			refs := collectRefsFromSchema(bodySchema)
			if !hasKnownRelationshipRef(refs) {
				return
			}
		}

		template := strings.TrimPrefix(path, relationshipPathPrefix)
		requiresBody := operation.RequestBody != nil && operation.RequestBody.Required != nil && *operation.RequestBody.Required
		entries = append(entries, RelationshipTemplate{
			Template:       template,
			Method:         strings.ToUpper(method),
			Summary:        operationSummary(operation),
			Body:           bodySchema,
			RequiresBody:   requiresBody,
			PathParamNames: pathParamNamesForTemplate(template),
		})
		seen[key] = struct{}{}
	})

	sort.Slice(entries, func(left, right int) bool {
		if entries[left].Template == entries[right].Template {
			return entries[left].Method < entries[right].Method
		}
		return entries[left].Template < entries[right].Template
	})

	if len(collectedErrs) == 0 {
		return entries, nil
	}
	return entries, errors.Join(collectedErrs...)
}

type relationshipSelection struct {
	template RelationshipTemplate
	params   map[string]string
}

type templateMatcher struct {
	template RelationshipTemplate
	regex    *regexp.Regexp
}

func selectRelationshipTemplate(deleteRequested bool, path string, matchers []templateMatcher) (relationshipSelection, error) {
	candidates := make([]relationshipSelection, 0)
	for _, matcher := range matchers {
		params, ok := matcher.match(path)
		if !ok {
			continue
		}
		candidates = append(candidates, relationshipSelection{template: matcher.template, params: params})
	}

	if len(candidates) == 0 {
		return relationshipSelection{}, fmt.Errorf("unknown relationship path '%s'", path)
	}

	if deleteRequested {
		for _, candidate := range candidates {
			if candidate.template.Method == http.MethodDelete {
				return candidate, nil
			}
		}
		return relationshipSelection{}, fmt.Errorf("delete flag set but no DELETE operation available for '%s'", path)
	}

	bestIdx := -1
	bestScore := 0
	for idx, candidate := range candidates {
		if candidate.template.Method == http.MethodDelete {
			continue
		}

		score := methodPreference[candidate.template.Method]
		if bestIdx == -1 || score < bestScore {
			bestIdx = idx
			bestScore = score
		}
	}

	if bestIdx == -1 {
		return relationshipSelection{}, fmt.Errorf("relationship path '%s' only supports DELETE operations; set delete: true to remove it", path)
	}

	return candidates[bestIdx], nil
}

// methodPreference ranks HTTP methods for selecting the best relationship
// template when multiple candidates match a path.
var methodPreference = map[string]int{
	http.MethodPost:   0,
	http.MethodPut:    1,
	http.MethodPatch:  2,
	http.MethodGet:    3,
	http.MethodDelete: 4,
}

func newTemplateMatcher(template RelationshipTemplate) (templateMatcher, error) {
	paramNames, err := extractTemplateParamNames(template.Template)
	if err != nil {
		return templateMatcher{}, err
	}

	pattern, err := buildTemplatePattern(template.Template, paramNames)
	if err != nil {
		return templateMatcher{}, err
	}

	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return templateMatcher{}, err
	}

	return templateMatcher{template: template, regex: compiled}, nil
}

func (m templateMatcher) match(path string) (map[string]string, bool) {
	matches := m.regex.FindStringSubmatch(path)
	if matches == nil {
		return nil, false
	}

	params := make(map[string]string, len(m.template.PathParamNames))
	for i, name := range m.template.PathParamNames {
		idx := i + 1
		if idx < len(matches) {
			params[name] = matches[idx]
		}
	}
	return params, true
}

func normalizeRelationshipPath(path string) string {
	trimmed := strings.TrimSpace(path)
	trimmed = strings.TrimPrefix(trimmed, relationshipPathPrefix)
	trimmed = strings.TrimPrefix(trimmed, "/")
	return trimmed
}

func buildTemplatePattern(template string, paramNames []string) (string, error) {
	pattern := regexp.QuoteMeta(template)
	for _, name := range paramNames {
		placeholder := "{" + name + "}"
		pattern = strings.Replace(pattern, regexp.QuoteMeta(placeholder), "([^/]+)", 1)
	}
	return "^" + pattern + "$", nil
}

func extractTemplateParamNames(template string) ([]string, error) {
	paramNames, err := extractParamNames(template)
	if err != nil {
		return nil, err
	}
	return paramNames, nil
}

func pathParamNamesForTemplate(template string) []string {
	paramNames, _ := extractParamNames(template)
	return paramNames
}

func extractParamNames(template string) ([]string, error) {
	params := make([]string, 0)
	cursor := template
	for {
		start := strings.Index(cursor, "{")
		if start == -1 {
			break
		}

		end := strings.Index(cursor[start:], "}")
		if end == -1 {
			return nil, fmt.Errorf("invalid template '%s'", template)
		}

		params = append(params, cursor[start+1:start+end])
		cursor = cursor[start+end+1:]
	}
	return params, nil
}

var consideredMethods = map[string]struct{}{
	http.MethodPost:   {},
	http.MethodPut:    {},
	http.MethodDelete: {},
	http.MethodPatch:  {},
}

func shouldConsiderMethod(method string) bool {
	_, ok := consideredMethods[method]
	return ok
}

func operationSummary(op *v3.Operation) string {
	if op == nil {
		return ""
	}
	if op.Summary != "" {
		return op.Summary
	}
	return op.Description
}

func inferRelationshipKind(template, method string) string {
	template = normalizeRelationshipPath(template)
	if kind, ok := DefaultRegistry().ByWriteTemplate(template); ok {
		return kind.Name
	}
	if strings.EqualFold(method, http.MethodDelete) {
		return "relationship-delete"
	}
	return "relationship"
}

func isRelationshipPath(path string) bool {
	lower := strings.ToLower(path)
	for _, indicator := range relationshipPathIndicators {
		if strings.Contains(lower, indicator) {
			return true
		}
	}
	return false
}

func hasKnownRelationshipRef(refs []string) bool {
	for _, ref := range refs {
		name := refName(ref)
		if _, ok := relationshipSchemaNames[name]; ok {
			return true
		}
	}
	return false
}

func refName(ref string) string {
	if ref == "" {
		return ""
	}
	if idx := strings.LastIndex(ref, "/"); idx >= 0 && idx+1 < len(ref) {
		return ref[idx+1:]
	}
	return ref
}

func collectRefsFromSchema(schema map[string]interface{}) []string {
	if schema == nil {
		return nil
	}

	refs := make(map[string]struct{})
	walkRefs(schema, refs)

	list := make([]string, 0, len(refs))
	for ref := range refs {
		list = append(list, ref)
	}
	return list
}

func walkRefs(node interface{}, refs map[string]struct{}) {
	switch value := node.(type) {
	case map[string]interface{}:
		if ref, ok := value["$ref"].(string); ok && ref != "" {
			refs[ref] = struct{}{}
		}
		for _, nested := range value {
			walkRefs(nested, refs)
		}
	case []interface{}:
		for _, item := range value {
			walkRefs(item, refs)
		}
	}
}

func extractRequestSchema(operation *v3.Operation) (map[string]interface{}, error) {
	if operation == nil || operation.RequestBody == nil || operation.RequestBody.Content == nil {
		return nil, nil
	}

	media := firstJSONMedia(operation.RequestBody.Content)
	if media == nil || media.Schema == nil {
		return nil, nil
	}

	return schemaProxyToJSON(media.Schema)
}

func firstJSONMedia(content *orderedmap.Map[string, *v3.MediaType]) *v3.MediaType {
	if content == nil {
		return nil
	}

	if media := content.GetOrZero("application/json"); media != nil {
		return media
	}

	for item := content.First(); item != nil; item = item.Next() {
		if strings.Contains(item.Key(), "json") {
			return item.Value()
		}
	}

	return nil
}

func schemaProxyToJSON(proxy *base.SchemaProxy) (map[string]interface{}, error) {
	if proxy == nil {
		return nil, nil
	}

	if ref := proxy.GetReference(); ref != "" {
		return map[string]interface{}{"$ref": ref}, nil
	}

	schema, err := proxy.BuildSchema()
	if err != nil {
		return nil, err
	}

	data, err := schema.MarshalJSON()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}

	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}

	return out, nil
}
