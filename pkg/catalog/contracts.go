package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
	"go.yaml.in/yaml/v4"
)

type ParameterContract struct {
	Name        string
	In          string
	Description string
	Required    bool
	SchemaName  string
	Schema      *base.SchemaProxy
}

type OperationContract struct {
	Path                  string
	Method                string
	Summary               string
	Deprecated            bool
	Parameters            []ParameterContract
	RequestBodyRequired   bool
	RequestBodySchema     *base.SchemaProxy
	RequestBodySchemaName string
	ResponseSchema        *base.SchemaProxy
	ResponseSchemaName    string
}

type ResourceContract struct {
	Type          string
	SchemaName    string
	Schema        *base.SchemaProxy
	Operations    map[string]OperationContract   // highest-priority endpoint per method
	AllOperations map[string][]OperationContract // all endpoints per method for multi-location types
}

type RequestValidation struct {
	PathParams   map[string]string
	QueryParams  map[string]string
	HeaderParams map[string]string
	CookieParams map[string]string
	Body         interface{}
}

func (s *Spec) OperationContract(path, method string) (OperationContract, error) {
	op, item, err := s.Operation(path, method)
	if err != nil {
		return OperationContract{}, err
	}
	return buildOperationContract(path, method, op, item), nil
}

func (s *Spec) ResourceContracts() (map[string]ResourceContract, error) {
	contracts := make(map[string]ResourceContract)
	if s == nil {
		return contracts, nil
	}
	s.ForEachOperation(func(path, method string, operation *v3.Operation, item *v3.PathItem) {
		resourceType, schemaName, schema := extractResourceContract(operation)
		if resourceType == "" {
			resourceType = inferResourceTypeFromPath(path)
		}
		if resourceType == "" {
			return
		}
		contract := contracts[resourceType]
		if contract.Type == "" {
			contract = ResourceContract{
				Type:          resourceType,
				SchemaName:    schemaName,
				Schema:        schema,
				Operations:    make(map[string]OperationContract),
				AllOperations: make(map[string][]OperationContract),
			}
		}
		if contract.Schema == nil && schema != nil {
			contract.Schema = schema
			contract.SchemaName = schemaName
		}
		method = strings.ToUpper(method)
		candidate := buildOperationContract(path, method, operation, item)
		// Store in AllOperations
		contract.AllOperations[method] = append(contract.AllOperations[method], candidate)
		// Store highest-priority in Operations
		current, exists := contract.Operations[method]
		if !exists || operationPriority(candidate, resourceType) > operationPriority(current, resourceType) {
			contract.Operations[method] = candidate
		}
		contracts[resourceType] = contract
	})
	return contracts, nil
}

func (s *Spec) ValidateSchema(name string, value interface{}) error {
	if s == nil {
		return fmt.Errorf("spec not initialized")
	}
	schemas, err := s.GetSchemas()
	if err != nil {
		return err
	}
	proxy, ok := schemas[name]
	if !ok || proxy == nil {
		return fmt.Errorf("schema %s not found", name)
	}
	return validateSchemaProxy(proxy, value, name)
}

func (s *Spec) ValidateOperationRequest(path, method string, input RequestValidation) error {
	contract, err := s.OperationContract(path, method)
	if err != nil {
		return err
	}
	return validateOperationInput(contract, input)
}

func (s *Spec) ValidateOperationResponse(path, method string, body interface{}) error {
	contract, err := s.OperationContract(path, method)
	if err != nil {
		return err
	}
	if contract.ResponseSchema == nil {
		return nil
	}
	label := path + " " + strings.ToUpper(method) + " response"
	return validateSchemaProxy(contract.ResponseSchema, body, label)
}

func (s *Spec) ValidateResource(resource manifest.Resource, method string) error {
	contracts, err := s.ResourceContracts()
	if err != nil {
		return err
	}
	resourceType := strings.TrimSpace(resource.Type)
	contract, ok := contracts[resourceType]
	if !ok {
		return fmt.Errorf("unknown resource type %s", resourceType)
	}

	resolvedMethod := strings.ToUpper(strings.TrimSpace(method))
	if resolvedMethod == "" {
		for _, candidate := range []string{http.MethodPost, http.MethodPut, http.MethodPatch} {
			if _, exists := contract.Operations[candidate]; exists {
				resolvedMethod = candidate
				break
			}
		}
	}
	if resolvedMethod == "" {
		return nil
	}

	op, ok := contract.Operations[resolvedMethod]
	if !ok {
		return fmt.Errorf("resource %s does not support %s", resourceType, resolvedMethod)
	}
	params, _ := s.Resolver().PathParams(resource, op)
	input := RequestValidation{PathParams: params, Body: resource.Data}
	if err := validateOperationInput(op, input); err != nil {
		return fmt.Errorf("%s %s: %w", resolvedMethod, resourceType, err)
	}
	return nil
}

func ValidateResources(spec *Spec, resources []manifest.Resource, deleteMode bool) error {
	if spec == nil || len(resources) == 0 {
		return nil
	}
	contracts, err := spec.ResourceContracts()
	if err != nil {
		return err
	}
	validationErrors := make([]error, 0)
	for idx := range resources {
		method := ""
		if deleteMode || resources[idx].Delete {
			method = http.MethodDelete
		}
		if method == http.MethodDelete {
			contract, ok := contracts[resources[idx].Type]
			if ok && contract.Operations != nil {
				if _, hasDelete := contract.Operations[http.MethodDelete]; !hasDelete {
					continue
				}
			}
		}
		if err := spec.ValidateResource(resources[idx], method); err != nil {
			validationErrors = append(validationErrors, fmt.Errorf("resources[%d]: %w", idx, err))
		}
	}
	if len(validationErrors) == 0 {
		return nil
	}
	return errors.Join(validationErrors...)
}

// ValidateManifest validates both resources and relationship operations in a
// single call. It returns a joined error containing all validation failures.
func (s *Spec) ValidateManifest(resources []manifest.Resource, relationships []manifest.RelationshipOperation, deleteMode bool) error {
	var errs []error
	if err := ValidateResources(s, resources, deleteMode); err != nil {
		errs = append(errs, err)
	}
	if err := ValidateRelationshipOperations(s, relationships); err != nil {
		errs = append(errs, err)
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

func buildOperationContract(path, method string, operation *v3.Operation, item *v3.PathItem) OperationContract {
	contract := OperationContract{
		Path:    path,
		Method:  strings.ToUpper(method),
		Summary: operationSummary(operation),
	}
	if operation != nil && operation.Deprecated != nil {
		contract.Deprecated = *operation.Deprecated
	}

	params := make(map[string]ParameterContract)
	appendParams := func(list []*v3.Parameter) {
		for _, parameter := range list {
			if parameter == nil || parameter.Name == "" || parameter.In == "" {
				continue
			}
			key := parameter.In + ":" + parameter.Name
			params[key] = ParameterContract{
				Name:        parameter.Name,
				In:          parameter.In,
				Description: parameter.Description,
				Required:    parameter.Required != nil && *parameter.Required,
				SchemaName:  schemaNameFromProxy(parameter.Schema),
				Schema:      parameter.Schema,
			}
		}
	}
	if item != nil {
		appendParams(item.Parameters)
	}
	if operation != nil {
		appendParams(operation.Parameters)
	}
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	contract.Parameters = make([]ParameterContract, 0, len(keys))
	for _, key := range keys {
		contract.Parameters = append(contract.Parameters, params[key])
	}

	if operation != nil && operation.RequestBody != nil && operation.RequestBody.Content != nil {
		contract.RequestBodyRequired = operation.RequestBody.Required != nil && *operation.RequestBody.Required
		if media := firstJSONMedia(operation.RequestBody.Content); media != nil {
			contract.RequestBodySchema = media.Schema
			contract.RequestBodySchemaName = schemaNameFromProxy(media.Schema)
		}
	}

	if operation != nil && operation.Responses != nil {
		for responsePair := operation.Responses.Codes.First(); responsePair != nil; responsePair = responsePair.Next() {
			if !strings.HasPrefix(responsePair.Key(), "2") {
				continue
			}
			response := responsePair.Value()
			if response == nil || response.Content == nil {
				continue
			}
			if media := firstJSONMedia(response.Content); media != nil && media.Schema != nil {
				contract.ResponseSchema = media.Schema
				contract.ResponseSchemaName = schemaNameFromProxy(media.Schema)
				break
			}
		}
	}

	return contract
}

func extractResourceContract(operation *v3.Operation) (string, string, *base.SchemaProxy) {
	if operation == nil {
		return "", "", nil
	}
	if operation.RequestBody != nil && operation.RequestBody.Content != nil {
		if media := firstJSONMedia(operation.RequestBody.Content); media != nil && media.Schema != nil {
			name := schemaNameFromProxy(media.Schema)
			if name != "" {
				return normalizeSchemaResourceType(name), name, media.Schema
			}
		}
	}
	if operation.Responses != nil {
		for responsePair := operation.Responses.Codes.First(); responsePair != nil; responsePair = responsePair.Next() {
			response := responsePair.Value()
			if response == nil || response.Content == nil {
				continue
			}
			if media := firstJSONMedia(response.Content); media != nil && media.Schema != nil {
				proxy := media.Schema
				name := schemaNameFromProxy(proxy)
				if name == "" {
					schema, err := proxy.BuildSchema()
					if err == nil && schema != nil && schema.Items != nil && schema.Items.IsA() && schema.Items.A != nil {
						proxy = schema.Items.A
						name = schemaNameFromProxy(proxy)
					}
				}
				if name != "" {
					return normalizeSchemaResourceType(name), name, proxy
				}
			}
		}
	}
	return "", "", nil
}

func normalizeSchemaResourceType(schemaName string) string {
	resourceType := schemaName
	for _, suffix := range []string{"Representation", "DTO", "VO", "Resource", "Model"} {
		resourceType = strings.TrimSuffix(resourceType, suffix)
	}
	return strings.ToLower(resourceType)
}

func schemaNameFromProxy(proxy *base.SchemaProxy) string {
	if proxy == nil {
		return ""
	}
	if ref := proxy.GetReference(); ref != "" {
		return refName(ref)
	}
	return ""
}

var priorityRules = []struct {
	check func(contract OperationContract, resourceType string) bool
	score int
}{
	{func(c OperationContract, _ string) bool {
		return c.Method == http.MethodGet && isCollectionEndpoint(c.Path)
	}, 100},
	{func(c OperationContract, r string) bool {
		return c.Method == http.MethodGet && strings.Contains(c.Path, "/{"+r+"-")
	}, -5},
	{func(c OperationContract, r string) bool {
		return strings.Contains(c.Path, "/admin/realms/{realm}/"+r+"s")
	}, 50},
	{func(c OperationContract, _ string) bool {
		return strings.Contains(c.Path, "/clients/") || strings.Contains(c.Path, "/roles/") || strings.Contains(c.Path, "/scope-mappings/")
	}, -30},
	{func(c OperationContract, _ string) bool { return strings.Contains(c.Path, "/roles/{role-name}") }, 40},
	{func(c OperationContract, _ string) bool { return strings.Contains(c.Path, "/roles-by-id/") }, -20},
	{func(c OperationContract, _ string) bool { return strings.Contains(c.Path, "/composites") }, -10},
}

func operationPriority(contract OperationContract, resourceType string) int {
	pathSegments := len(strings.Split(contract.Path, "/"))
	score := (20 - pathSegments) * 5
	for _, rule := range priorityRules {
		if rule.check(contract, resourceType) {
			score += rule.score
		}
	}
	return score
}

func inferResourceTypeFromPath(path string) string {
	if path == "" {
		return ""
	}
	segments := strings.Split(path, "/")
	var candidate string
	for _, segment := range segments {
		segment = strings.TrimSpace(segment)
		if segment == "" || segment == "admin" || segment == "realms" {
			continue
		}
		if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
			continue
		}
		candidate = segment
	}
	if candidate == "" {
		return ""
	}
	candidate = strings.ToLower(candidate)
	candidate = strings.TrimSuffix(candidate, "-by-id")
	candidate = strings.TrimSuffix(candidate, "-by-name")
	candidate = strings.TrimSuffix(candidate, "-export")
	candidate = strings.TrimSuffix(candidate, "-import")
	candidate = strings.TrimSuffix(candidate, "-mappings")
	candidate = strings.TrimSuffix(candidate, "-mapping")
	candidate = strings.TrimSuffix(candidate, "-composites")
	if strings.HasSuffix(candidate, "ies") {
		return strings.TrimSuffix(candidate, "ies") + "y"
	}
	candidate = strings.TrimSuffix(candidate, "s")
	return candidate
}

func isCollectionEndpoint(path string) bool {
	pathParts := strings.Split(path, "/")
	if len(pathParts) == 0 {
		return false
	}
	lastPart := pathParts[len(pathParts)-1]
	return !strings.HasPrefix(lastPart, "{") || !strings.HasSuffix(lastPart, "}")
}

func validateOperationInput(contract OperationContract, input RequestValidation) error {
	validationErrors := make([]error, 0)
	for _, parameter := range contract.Parameters {
		value, exists := lookupParameterValue(parameter, input)
		if parameter.Required && !exists {
			validationErrors = append(validationErrors, fmt.Errorf("missing required %s parameter %s", parameter.In, parameter.Name))
			continue
		}
		if !exists || value == "" || parameter.Schema == nil {
			continue
		}
		coerced := coerceParameterValue(value, parameter.Schema)
		if err := validateSchemaProxy(parameter.Schema, coerced, parameter.In+" parameter "+parameter.Name); err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	if contract.RequestBodyRequired && input.Body == nil {
		validationErrors = append(validationErrors, fmt.Errorf("request body is required"))
	}
	if contract.RequestBodySchema != nil && input.Body != nil {
		if err := validateSchemaProxy(contract.RequestBodySchema, input.Body, contract.Path+" request body"); err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}
	return errors.Join(validationErrors...)
}

func lookupParameterValue(parameter ParameterContract, input RequestValidation) (string, bool) {
	var source map[string]string
	switch parameter.In {
	case "path":
		source = input.PathParams
	case "query":
		source = input.QueryParams
	case "header":
		source = input.HeaderParams
	case "cookie":
		source = input.CookieParams
	default:
		return "", false
	}
	if source == nil {
		return "", false
	}
	value, ok := source[parameter.Name]
	return value, ok
}

// coerceParameterValue parses a string parameter value into the schema type
// expected by the OpenAPI spec. Query, path, header, and cookie parameters are
// transported as strings, so they must be coerced before schema validation.
func coerceParameterValue(value string, proxy *base.SchemaProxy) interface{} {
	if proxy == nil {
		return value
	}
	schema, err := proxy.BuildSchema()
	if err != nil {
		return value
	}
	for _, t := range schema.Type {
		switch t {
		case "integer":
			if i, err := strconv.ParseInt(value, 10, 64); err == nil {
				return i
			}
		case "number":
			if f, err := strconv.ParseFloat(value, 64); err == nil {
				return f
			}
		case "boolean":
			if b, err := strconv.ParseBool(value); err == nil {
				return b
			}
		}
	}
	return value
}

func validateSchemaProxy(proxy *base.SchemaProxy, value interface{}, label string) error {
	if proxy == nil {
		return nil
	}
	schema, err := proxy.BuildSchema()
	if err != nil {
		return fmt.Errorf("build schema for %s: %w", label, err)
	}
	validationErrors := validateValue(schema, value, "$")
	if len(validationErrors) == 0 {
		return nil
	}
	wrapped := make([]error, 0, len(validationErrors))
	for _, validationErr := range validationErrors {
		wrapped = append(wrapped, fmt.Errorf("%s: %w", label, validationErr))
	}
	return errors.Join(wrapped...)
}

func validateValue(schema *base.Schema, value interface{}, path string) []error {
	if schema == nil {
		return nil
	}
	if value == nil {
		if schema.Nullable != nil && *schema.Nullable {
			return nil
		}
		if len(schema.Type) == 0 {
			return nil
		}
		return []error{fmt.Errorf("%s: value is null", path)}
	}

	validationErrors := make([]error, 0)
	if len(schema.Enum) > 0 && !matchesEnum(schema.Enum, value) {
		validationErrors = append(validationErrors, fmt.Errorf("%s: value is not in enum", path))
	}
	if schema.Format == "binary" && containsType(schema.Type, "string") {
		return validationErrors
	}

	for _, one := range schema.AllOf {
		validationErrors = append(validationErrors, validateProxy(one, value, path)...)
	}

	if len(schema.OneOf) > 0 {
		matches := 0
		for _, one := range schema.OneOf {
			if len(validateProxy(one, value, path)) == 0 {
				matches++
			}
		}
		if matches != 1 {
			validationErrors = append(validationErrors, fmt.Errorf("%s: oneOf matched %d schemas", path, matches))
		}
	}

	if len(schema.AnyOf) > 0 {
		matches := 0
		for _, one := range schema.AnyOf {
			if len(validateProxy(one, value, path)) == 0 {
				matches++
			}
		}
		if matches == 0 {
			validationErrors = append(validationErrors, fmt.Errorf("%s: value does not match anyOf schemas", path))
		}
	}

	typeName := inferTypeName(value)
	if len(schema.Type) > 0 && !containsType(schema.Type, typeName) {
		validationErrors = append(validationErrors, fmt.Errorf("%s: expected %s, got %s", path, strings.Join(schema.Type, "/"), typeName))
		return validationErrors
	}

	switch typeName {
	case "object":
		mapped, ok := normalizeObject(value)
		if !ok {
			validationErrors = append(validationErrors, fmt.Errorf("%s: expected object", path))
			return validationErrors
		}
		validationErrors = append(validationErrors, validateObject(schema, mapped, path)...)
	case "array":
		items, ok := normalizeSlice(value)
		if !ok {
			validationErrors = append(validationErrors, fmt.Errorf("%s: expected array", path))
			return validationErrors
		}
		validationErrors = append(validationErrors, validateArray(schema, items, path)...)
	case "string":
		validationErrors = append(validationErrors, validateString(schema, fmt.Sprint(value), path)...)
	case "integer", "number":
		validationErrors = append(validationErrors, validateNumber(schema, value, path, typeName)...)
	case "boolean":
		if _, ok := value.(bool); !ok {
			validationErrors = append(validationErrors, fmt.Errorf("%s: expected boolean", path))
		}
	}

	return validationErrors
}

func validateProxy(proxy *base.SchemaProxy, value interface{}, path string) []error {
	if proxy == nil {
		return nil
	}
	schema, err := proxy.BuildSchema()
	if err != nil {
		return []error{fmt.Errorf("%s: build nested schema: %w", path, err)}
	}
	return validateValue(schema, value, path)
}

func validateObject(schema *base.Schema, value map[string]interface{}, path string) []error {
	validationErrors := make([]error, 0)
	for _, key := range schema.Required {
		if _, ok := value[key]; !ok {
			validationErrors = append(validationErrors, fmt.Errorf("%s.%s: required property missing", path, key))
		}
	}

	known := make(map[string]*base.SchemaProxy)
	if schema.Properties != nil {
		for pair := schema.Properties.First(); pair != nil; pair = pair.Next() {
			known[pair.Key()] = pair.Value()
		}
	}

	for key, item := range value {
		if proxy, ok := known[key]; ok {
			validationErrors = append(validationErrors, validateProxy(proxy, item, path+"."+key)...)
			continue
		}

		if schema.PatternProperties != nil {
			matched := false
			for pair := schema.PatternProperties.First(); pair != nil; pair = pair.Next() {
				compiled, err := regexp.Compile(pair.Key())
				if err != nil || !compiled.MatchString(key) {
					continue
				}
				matched = true
				validationErrors = append(validationErrors, validateProxy(pair.Value(), item, path+"."+key)...)
			}
			if matched {
				continue
			}
		}

		if schema.AdditionalProperties != nil {
			if schema.AdditionalProperties.IsB() && !schema.AdditionalProperties.B {
				validationErrors = append(validationErrors, fmt.Errorf("%s.%s: additional property not allowed", path, key))
				continue
			}
			if schema.AdditionalProperties.IsA() && schema.AdditionalProperties.A != nil {
				validationErrors = append(validationErrors, validateProxy(schema.AdditionalProperties.A, item, path+"."+key)...)
			}
		}
	}

	if schema.MinProperties != nil && int64(len(value)) < *schema.MinProperties {
		validationErrors = append(validationErrors, fmt.Errorf("%s: expected at least %d properties", path, *schema.MinProperties))
	}
	if schema.MaxProperties != nil && int64(len(value)) > *schema.MaxProperties {
		validationErrors = append(validationErrors, fmt.Errorf("%s: expected at most %d properties", path, *schema.MaxProperties))
	}
	return validationErrors
}

func validateArray(schema *base.Schema, items []interface{}, path string) []error {
	validationErrors := make([]error, 0)
	if schema.MinItems != nil && int64(len(items)) < *schema.MinItems {
		validationErrors = append(validationErrors, fmt.Errorf("%s: expected at least %d items", path, *schema.MinItems))
	}
	if schema.MaxItems != nil && int64(len(items)) > *schema.MaxItems {
		validationErrors = append(validationErrors, fmt.Errorf("%s: expected at most %d items", path, *schema.MaxItems))
	}
	if schema.Items != nil && schema.Items.IsA() && schema.Items.A != nil {
		for idx, item := range items {
			validationErrors = append(validationErrors, validateProxy(schema.Items.A, item, fmt.Sprintf("%s[%d]", path, idx))...)
		}
	}
	return validationErrors
}

func validateString(schema *base.Schema, value, path string) []error {
	validationErrors := make([]error, 0)
	if schema.MinLength != nil && int64(len(value)) < *schema.MinLength {
		validationErrors = append(validationErrors, fmt.Errorf("%s: expected minimum length %d", path, *schema.MinLength))
	}
	if schema.MaxLength != nil && int64(len(value)) > *schema.MaxLength {
		validationErrors = append(validationErrors, fmt.Errorf("%s: expected maximum length %d", path, *schema.MaxLength))
	}
	if schema.Pattern != "" {
		compiled, err := regexp.Compile(schema.Pattern)
		if err == nil && !compiled.MatchString(value) {
			validationErrors = append(validationErrors, fmt.Errorf("%s: value does not match pattern %s", path, schema.Pattern))
		}
	}
	return validationErrors
}

func validateNumber(schema *base.Schema, raw interface{}, path, typeName string) []error {
	validationErrors := make([]error, 0)
	value, ok := normalizeNumber(raw)
	if !ok {
		return []error{fmt.Errorf("%s: expected %s", path, typeName)}
	}
	if typeName == "integer" && math.Mod(value, 1) != 0 {
		validationErrors = append(validationErrors, fmt.Errorf("%s: expected integer", path))
	}
	if schema.Minimum != nil && value < *schema.Minimum {
		validationErrors = append(validationErrors, fmt.Errorf("%s: expected minimum %v", path, *schema.Minimum))
	}
	if schema.Maximum != nil && value > *schema.Maximum {
		validationErrors = append(validationErrors, fmt.Errorf("%s: expected maximum %v", path, *schema.Maximum))
	}
	if schema.MultipleOf != nil && *schema.MultipleOf != 0 {
		quotient := value / *schema.MultipleOf
		if math.Mod(quotient, 1) != 0 {
			validationErrors = append(validationErrors, fmt.Errorf("%s: expected multipleOf %v", path, *schema.MultipleOf))
		}
	}
	return validationErrors
}

func matchesEnum(enum []*yaml.Node, value interface{}) bool {
	encodedValue, err := json.Marshal(value)
	if err != nil {
		return false
	}
	for _, node := range enum {
		var candidate interface{}
		if node == nil || node.Decode(&candidate) != nil {
			continue
		}
		encodedCandidate, err := json.Marshal(candidate)
		if err == nil && bytesEqual(encodedValue, encodedCandidate) {
			return true
		}
	}
	return false
}

func bytesEqual(left, right []byte) bool {
	if len(left) != len(right) {
		return false
	}
	for idx := range left {
		if left[idx] != right[idx] {
			return false
		}
	}
	return true
}

func inferTypeName(value interface{}) string {
	raw := reflect.ValueOf(value)
	if raw.IsValid() {
		switch raw.Kind() {
		case reflect.Map:
			return "object"
		case reflect.Slice, reflect.Array:
			if _, ok := value.([]byte); ok {
				return "string"
			}
			return "array"
		}
	}

	switch typed := value.(type) {
	case map[string]interface{}:
		return "object"
	case []interface{}:
		return "array"
	case string:
		return "string"
	case bool:
		return "boolean"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "integer"
	case float32:
		if math.Mod(float64(typed), 1) == 0 {
			return "integer"
		}
		return "number"
	case float64:
		if math.Mod(typed, 1) == 0 {
			return "integer"
		}
		return "number"
	default:
		return "object"
	}
}

func containsType(types []string, want string) bool {
	for _, item := range types {
		if item == want {
			return true
		}
		if item == "number" && want == "integer" {
			return true
		}
	}
	return false
}

func normalizeJSON[T any](value interface{}, direct func(interface{}) (T, bool)) (T, bool) {
	if v, ok := direct(value); ok {
		return v, true
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		var zero T
		return zero, false
	}
	var decoded T
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		var zero T
		return zero, false
	}
	return decoded, true
}

func normalizeObject(value interface{}) (map[string]interface{}, bool) {
	return normalizeJSON(value, func(v interface{}) (map[string]interface{}, bool) {
		m, ok := v.(map[string]interface{})
		return m, ok
	})
}

func normalizeSlice(value interface{}) ([]interface{}, bool) {
	return normalizeJSON(value, func(v interface{}) ([]interface{}, bool) {
		s, ok := v.([]interface{})
		return s, ok
	})
}

func normalizeNumber(raw interface{}) (float64, bool) {
	switch typed := raw.(type) {
	case int:
		return float64(typed), true
	case int8:
		return float64(typed), true
	case int16:
		return float64(typed), true
	case int32:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case uint:
		return float64(typed), true
	case uint8:
		return float64(typed), true
	case uint16:
		return float64(typed), true
	case uint32:
		return float64(typed), true
	case uint64:
		return float64(typed), true
	case float32:
		return float64(typed), true
	case float64:
		return typed, true
	case json.Number:
		parsed, err := typed.Float64()
		return parsed, err == nil
	case string:
		parsed, err := json.Number(typed).Float64()
		return parsed, err == nil
	default:
		return 0, false
	}
}
