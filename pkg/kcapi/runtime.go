// runtime.go drives Keycloak through the supplied OpenAPI spec: it resolves
// resource operations from the spec, injects auth, and validates requests and
// responses against the spec. It is the HTTP transport layer of the kcapi
// library.
package kcapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/rs/zerolog/log"
	"github.com/thedataflows/keycloak-cli/pkg/auth"
)

type RuntimeConfig struct {
	BaseURL  string
	SpecPath string
	Timeout  time.Duration
	// Spec optionally supplies a pre-loaded spec, overriding SpecPath. Side-car
	// overrides (relationship/field/built-in YAML files) are only loaded from
	// SpecPath, so they are skipped when a pre-built spec is injected.
	Spec *Spec
	// HTTP optionally overrides the HTTP client used for Keycloak requests.
	HTTP *http.Client
}

type TokenProvider interface {
	AccessToken(ctx context.Context, baseURL, accessToken, refreshToken string) (string, error)
}

type requestEditor func(ctx context.Context, req *http.Request) error

type RuntimeClient struct {
	baseURL    string
	httpClient *http.Client
	authEditor requestEditor
	spec       *Spec
}

func NewRuntimeClient(config RuntimeConfig, tokens TokenProvider) (*RuntimeClient, error) {
	baseURL := strings.TrimSpace(config.BaseURL)
	if baseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	spec := config.Spec
	if spec == nil {
		var err error
		spec, err = NewSpec(config.SpecPath)
		if err != nil {
			return nil, fmt.Errorf("load spec: %w", err)
		}

		if err := InstallDefaultRegistry(config.SpecPath); err != nil {
			return nil, fmt.Errorf("load relationship overrides: %w", err)
		}

		if err := InstallDefaultFieldOverrides(config.SpecPath); err != nil {
			return nil, fmt.Errorf("load field overrides: %w", err)
		}

		if err := InstallDefaultBuiltInResources(config.SpecPath); err != nil {
			return nil, fmt.Errorf("load built-in resources: %w", err)
		}
	}

	httpClient := config.HTTP
	if httpClient == nil {
		httpClient = &http.Client{Timeout: config.Timeout}
	}

	client := &RuntimeClient{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: httpClient,
		spec:       spec,
	}

	client.authEditor = func(ctx context.Context, req *http.Request) error {
		accessToken, err := tokens.AccessToken(ctx, baseURL, os.Getenv(auth.AccessTokenEnvVar), os.Getenv(auth.RefreshTokenEnvVar))
		if err != nil {
			return fmt.Errorf("get access token: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		return nil
	}

	return client, nil
}

func (r *RuntimeClient) Spec() *Spec {
	if r == nil {
		return nil
	}
	return r.spec
}

// declaredQueryParams drops parameters the operation does not declare as
// query parameters. The manifest fetchers emit briefRepresentation=false and
// friends for every collection; operations that don't declare a parameter
// would otherwise receive noise the server can only ignore.
func declaredQueryParams(contract OperationContract, params map[string]string) map[string]string {
	if len(params) == 0 {
		return nil
	}
	declared := make(map[string]struct{}, len(contract.Parameters))
	for _, p := range contract.Parameters {
		if p.In == "query" {
			declared[p.Name] = struct{}{}
		}
	}
	kept := make(map[string]string, len(params))
	for k, v := range params {
		if _, ok := declared[k]; ok {
			kept[k] = v
		}
	}
	if len(kept) == 0 {
		return nil
	}
	return kept
}

func (r *RuntimeClient) FetchResources(ctx context.Context, resourceType string, scope map[string]string, params ...map[string]string) ([]Resource, error) {
	contract, err := r.spec.Resolver().ResolveResourceOperation(resourceType, "", http.MethodGet, OperationCollection)
	if err != nil {
		return nil, err
	}

	op, _, err := r.spec.Operation(contract.Path, contract.Method)
	if err != nil {
		return nil, err
	}

	queryParams := declaredQueryParams(contract, mergeQueryParams(params...))
	requestPath := r.buildPathWithOperation(contract.Path, op, scope)
	if err := r.spec.ValidateOperationRequest(contract.Path, http.MethodGet, RequestValidation{
		PathParams:  scope,
		QueryParams: queryParams,
	}); err != nil {
		return nil, err
	}

	fullPath := r.baseURL + requestPath
	if len(queryParams) > 0 {
		values := url.Values{}
		for key, value := range queryParams {
			values.Add(key, value)
		}
		fullPath += "?" + values.Encode()
	}

	log.Logger.Debug().Str("pkg", "admin").Str("url", fullPath).Str("method", http.MethodGet).Msg("Fetching resources")
	req, err := r.newAuthRequest(ctx, http.MethodGet, fullPath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, r.readErrorBody(resp)
	}

	var rawResources []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawResources); err != nil {
		return nil, err
	}
	if err := r.spec.ValidateOperationResponse(contract.Path, http.MethodGet, toInterfaceSlice(rawResources)); err != nil {
		return nil, err
	}

	resources := make([]Resource, len(rawResources))
	for i, raw := range rawResources {
		resources[i] = Resource{
			Type:  resourceType,
			Realm: extractRealm(raw, scope),
			Data:  raw,
		}
	}

	return resources, nil
}

// FetchResourcesWithParent fetches a resource collection scoped by the resource's
// parent type, so nested resources resolve to the correct endpoint.
func (r *RuntimeClient) FetchResourcesWithParent(ctx context.Context, resource Resource, params ...map[string]string) ([]Resource, error) {
	contract, err := r.spec.Resolver().ResolveResourceOperation(resource.Type, resource.ParentType, http.MethodGet, OperationCollection)
	if err != nil {
		return nil, err
	}

	paramsMap, err := r.spec.Resolver().PathParams(resource, contract)
	if err != nil {
		return nil, err
	}

	op, _, err := r.spec.Operation(contract.Path, contract.Method)
	if err != nil {
		return nil, err
	}

	queryParams := declaredQueryParams(contract, mergeQueryParams(params...))
	requestPath := r.buildPathWithOperation(contract.Path, op, paramsMap)
	if err := r.spec.ValidateOperationRequest(contract.Path, http.MethodGet, RequestValidation{
		PathParams:  paramsMap,
		QueryParams: queryParams,
	}); err != nil {
		return nil, err
	}

	fullPath := r.baseURL + requestPath
	if len(queryParams) > 0 {
		values := url.Values{}
		for key, value := range queryParams {
			values.Add(key, value)
		}
		fullPath += "?" + values.Encode()
	}

	log.Logger.Debug().Str("pkg", "admin").Str("url", fullPath).Str("method", http.MethodGet).Msg("Fetching resources")
	req, err := r.newAuthRequest(ctx, http.MethodGet, fullPath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, r.readErrorBody(resp)
	}

	var rawResources []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawResources); err != nil {
		return nil, err
	}
	if err := r.spec.ValidateOperationResponse(contract.Path, http.MethodGet, toInterfaceSlice(rawResources)); err != nil {
		return nil, err
	}

	resources := make([]Resource, len(rawResources))
	for i, raw := range rawResources {
		resources[i] = Resource{
			Type:  resource.Type,
			Realm: extractRealm(raw, paramsMap),
			Data:  raw,
		}
	}

	return resources, nil
}

func (r *RuntimeClient) FetchPathCollection(ctx context.Context, path string, scope map[string]string, params ...map[string]string) ([]map[string]interface{}, error) {
	queryParams := mergeQueryParams(params...)
	// Gate only when the spec knows the path; an unknown path keeps its
	// parameters, since nothing is known about what it accepts.
	if contract, err := r.spec.OperationContract(path, http.MethodGet); err == nil {
		queryParams = declaredQueryParams(contract, queryParams)
	}
	resolvedPath := r.buildPathWithOperation(path, nil, scope)
	if err := r.spec.ValidateOperationRequest(path, http.MethodGet, RequestValidation{
		PathParams:  scope,
		QueryParams: queryParams,
	}); err != nil {
		return nil, err
	}

	fullPath := r.baseURL + resolvedPath
	if len(queryParams) > 0 {
		values := url.Values{}
		for key, value := range queryParams {
			values.Add(key, value)
		}
		fullPath += "?" + values.Encode()
	}

	// Log like FetchResources/FetchResourcesWithParent so nested/child GET volume
	// is observable at trace/debug (ISSUE 0003 — these requests were previously
	// silent, hiding the cost of depth traversal).
	log.Logger.Debug().Str("pkg", "admin").Str("url", fullPath).Str("method", http.MethodGet).Msg("Fetching path collection")
	req, err := r.newAuthRequest(ctx, http.MethodGet, fullPath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, r.readErrorBody(resp)
	}

	var payload []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if err := r.spec.ValidateOperationResponse(path, http.MethodGet, toInterfaceSlice(payload)); err != nil {
		return nil, err
	}
	return payload, nil
}

// FetchResource resolves a single GET operation for the resource, substitutes
// the resource's identifier into the path, and returns the fetched resource.
// The second return value is false when the server responds with 404.
func (r *RuntimeClient) FetchResource(ctx context.Context, resource Resource) (Resource, bool, error) {
	contract, err := r.resolveResourceContract(resource, http.MethodGet)
	if err != nil {
		return Resource{}, false, err
	}

	params, err := r.spec.Resolver().PathParams(resource, contract)
	if err != nil {
		return Resource{}, false, err
	}

	op, _, err := r.spec.Operation(contract.Path, contract.Method)
	if err != nil {
		return Resource{}, false, err
	}

	requestPath := r.buildPathWithOperation(contract.Path, op, params)
	if err := r.spec.ValidateOperationRequest(contract.Path, http.MethodGet, RequestValidation{PathParams: params}); err != nil {
		return Resource{}, false, err
	}

	fullPath := r.baseURL + requestPath
	log.Logger.Debug().Str("pkg", "admin").Str("url", fullPath).Str("method", http.MethodGet).Msg("Fetching resource")
	req, err := r.newAuthRequest(ctx, http.MethodGet, fullPath, nil)
	if err != nil {
		return Resource{}, false, err
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return Resource{}, false, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return Resource{}, false, nil
	}
	if resp.StatusCode != http.StatusOK {
		return Resource{}, false, r.readErrorBody(resp)
	}

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return Resource{}, false, err
	}
	if err := r.spec.ValidateOperationResponse(contract.Path, http.MethodGet, raw); err != nil {
		return Resource{}, false, err
	}

	return Resource{
		Type:  resource.Type,
		Realm: extractRealm(raw, params),
		Data:  raw,
	}, true, nil
}

func (r *RuntimeClient) CreateResource(ctx context.Context, resource Resource) (int, string, error) {
	contract, err := r.resolveResourceContract(resource, http.MethodPost)
	if err != nil {
		return 0, "", err
	}

	params, err := r.spec.Resolver().PathParams(resource, contract)
	if err != nil {
		return 0, "", err
	}

	op, _, err := r.spec.Operation(contract.Path, contract.Method)
	if err != nil {
		return 0, "", err
	}

	requestPath := r.buildPathWithOperation(contract.Path, op, params)
	r.sanitizeResourcePayload(&resource, http.MethodPost, contract)
	if err := r.spec.ValidateOperationRequest(contract.Path, http.MethodPost, RequestValidation{PathParams: params, Body: resource.Data}); err != nil {
		return 0, "", err
	}

	fullURL := r.baseURL + requestPath
	body, err := json.Marshal(resource.Data)
	if err != nil {
		return 0, "", err
	}
	req, err := r.newAuthRequest(ctx, http.MethodPost, fullURL, bytes.NewReader(body))
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusMultipleChoices {
		return resp.StatusCode, extractCreatedID(resp), nil
	}
	return resp.StatusCode, "", r.readErrorBody(resp)
}

func extractCreatedID(resp *http.Response) string {
	var respBody map[string]interface{}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&respBody); decodeErr == nil {
		if id, ok := respBody["id"].(string); ok && id != "" {
			return id
		}
	}
	if loc := resp.Header.Get("Location"); loc != "" {
		parts := strings.Split(strings.TrimRight(loc, "/"), "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	}
	return ""
}

func (r *RuntimeClient) UpdateResource(ctx context.Context, resource Resource) (int, error) {
	return r.resourceOperation(ctx, resource, http.MethodPut)
}

func (r *RuntimeClient) DeleteResource(ctx context.Context, resource Resource) (int, error) {
	contract, err := r.resolveResourceContract(resource, http.MethodDelete)
	if err != nil {
		return 0, err
	}

	params, err := r.spec.Resolver().PathParams(resource, contract)
	if err != nil {
		return 0, err
	}

	op, _, err := r.spec.Operation(contract.Path, contract.Method)
	if err != nil {
		return 0, err
	}

	requestPath := r.buildPathWithOperation(contract.Path, op, params)
	if err := r.spec.ValidateOperationRequest(contract.Path, http.MethodDelete, RequestValidation{PathParams: params}); err != nil {
		return 0, err
	}

	fullURL := r.baseURL + requestPath
	req, err := r.newAuthRequest(ctx, http.MethodDelete, fullURL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusMultipleChoices {
		return resp.StatusCode, nil
	}
	return resp.StatusCode, r.readErrorBody(resp)
}

func (r *RuntimeClient) ExecuteRelationship(ctx context.Context, rel RelationshipOperation) (int, error) {
	if r.spec == nil {
		return 0, fmt.Errorf("spec not initialized")
	}

	method := strings.ToUpper(strings.TrimSpace(rel.Method))
	if method == "" {
		return 0, fmt.Errorf("relationship method is required")
	}
	templatePath := strings.TrimSpace(rel.Template)
	if templatePath == "" {
		return 0, fmt.Errorf("relationship template is required")
	}
	if _, _, err := r.spec.Operation("/admin/realms/"+templatePath, method); err != nil {
		return 0, err
	}

	actualPath := strings.TrimSpace(rel.Path)
	if actualPath == "" {
		return 0, fmt.Errorf("relationship path is required")
	}

	var body io.Reader
	var requestBody interface{}
	if len(rel.Data) > 0 {
		body = bytes.NewReader(rel.Data)
		if err := json.Unmarshal(rel.Data, &requestBody); err != nil {
			return 0, err
		}
	}
	if err := r.spec.ValidateOperationRequest("/admin/realms/"+templatePath, method, RequestValidation{
		PathParams: rel.PathParams,
		Body:       requestBody,
	}); err != nil {
		return 0, err
	}

	fullURL := r.baseURL + "/admin/realms/" + actualPath
	req, err := r.newAuthRequest(ctx, method, fullURL, body)
	if err != nil {
		return 0, err
	}
	if len(rel.Data) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return resp.StatusCode, nil
	}
	return resp.StatusCode, r.readErrorBody(resp)
}

func (r *RuntimeClient) resourceOperation(ctx context.Context, resource Resource, method string) (int, error) {
	contract, err := r.resolveResourceContract(resource, method)
	if err != nil {
		return 0, err
	}

	params, err := r.spec.Resolver().PathParams(resource, contract)
	if err != nil {
		return 0, err
	}

	op, _, err := r.spec.Operation(contract.Path, contract.Method)
	if err != nil {
		return 0, err
	}

	requestPath := r.buildPathWithOperation(contract.Path, op, params)
	r.sanitizeResourcePayload(&resource, method, contract)
	if err := r.spec.ValidateOperationRequest(contract.Path, method, RequestValidation{PathParams: params, Body: resource.Data}); err != nil {
		return 0, err
	}
	fullURL := r.baseURL + requestPath
	body, err := json.Marshal(resource.Data)
	if err != nil {
		return 0, err
	}

	req, err := r.newAuthRequest(ctx, method, fullURL, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusMultipleChoices {
		return resp.StatusCode, nil
	}
	return resp.StatusCode, r.readErrorBody(resp)
}

// resolveResourceContract resolves the operation contract for a resource and method.
// It handles the realm special case and uses OperationCollection for create
// operations to avoid picking nested single-resource endpoints.
func (r *RuntimeClient) resolveResourceContract(resource Resource, method string) (OperationContract, error) {
	if resource.Type == "realm" {
		switch method {
		case http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete:
			return OperationContract{Path: "/admin/realms/{realm}", Method: method}, nil
		default:
			return OperationContract{Path: "/admin/realms", Method: method}, nil
		}
	}

	shape := OperationSingle
	if method == http.MethodPost {
		// Creation is always a collection endpoint. Use OperationCollection to avoid
		// picking nested single-resource endpoints (e.g. client scope mappings).
		shape = OperationCollection
	}

	return r.spec.Resolver().ResolveResourceOperation(resource.Type, resource.ParentType, method, shape)
}

func mergeQueryParams(params ...map[string]string) map[string]string {
	merged := make(map[string]string)
	for _, paramSet := range params {
		for key, value := range paramSet {
			merged[key] = value
		}
	}
	if len(merged) == 0 {
		return nil
	}
	return merged
}

func toInterfaceSlice(items []map[string]interface{}) []interface{} {
	converted := make([]interface{}, 0, len(items))
	for _, item := range items {
		converted = append(converted, item)
	}
	return converted
}

func (r *RuntimeClient) buildPathWithOperation(path string, operation *v3.Operation, scope map[string]string) string {
	if strings.Contains(path, "/admin/realms") && !strings.Contains(path, "{realm}") {
		return path
	}

	values := make(map[string]string, len(scope))
	maps.Copy(values, scope)
	if operation != nil && operation.Parameters != nil {
		for _, parameter := range operation.Parameters {
			if parameter == nil || parameter.In != "path" {
				continue
			}
			if value, exists := scope[parameter.Name]; exists {
				values[parameter.Name] = value
			}
		}
	}

	// Lenient by contract: manifest flows validate the request right after
	// building the path, so unfilled placeholders are surfaced there instead
	// of here. The partial result is kept for that error message.
	result, _ := instantiatePath(path, values)
	return result
}

// pathPlaceholderNames returns the distinct {name} placeholders of a path
// template in order of first appearance.
func pathPlaceholderNames(template string) []string {
	var names []string
	seen := make(map[string]struct{})
	for {
		open := strings.Index(template, "{")
		if open < 0 {
			return names
		}
		rest := template[open+1:]
		close := strings.Index(rest, "}")
		if close < 0 {
			return names
		}
		name := rest[:close]
		template = rest[close+1:]
		if name == "" {
			continue
		}
		if _, dup := seen[name]; dup {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}
}

// instantiatePath replaces every {name} placeholder of the template with the
// URL-path-escaped value from values. It is the single substitution
// implementation shared by the manifest runtime flows (via
// buildPathWithOperation) and Invoke. A placeholder with no entry in values
// is left in place and reported; the returned string still carries every
// substitution made so far.
func instantiatePath(template string, values map[string]string) (string, error) {
	var missing []string
	result := template
	for _, name := range pathPlaceholderNames(template) {
		value, exists := values[name]
		if !exists {
			missing = append(missing, name)
			continue
		}
		result = strings.ReplaceAll(result, "{"+name+"}", url.PathEscape(value))
	}
	if len(missing) > 0 {
		return result, fmt.Errorf("unfilled path placeholder(s): %s", strings.Join(missing, ", "))
	}
	return result, nil
}

// BaseURL returns the trimmed base URL every request path is joined against.
func (r *RuntimeClient) BaseURL() string {
	return r.baseURL
}

// Authorize stamps the request with the current access token from the
// configured token provider. It is the auth half of newAuthRequest, exposed
// for callers that build requests themselves.
func (r *RuntimeClient) Authorize(ctx context.Context, req *http.Request) error {
	if r == nil || r.authEditor == nil {
		return nil
	}
	return r.authEditor(ctx, req)
}

func (r *RuntimeClient) sanitizeResourcePayload(resource *Resource, method string, contract OperationContract) {
	if resource == nil {
		return
	}
	switch method {
	case http.MethodPost, http.MethodPut:
		// Clone the data map so sanitization does not mutate the caller's
		// copy. Go maps are reference types; without cloning, deleting id or
		// parent-reference fields here would leak back to apply logic that
		// still needs those fields for path resolution and conflict handling.
		if resource.Data != nil {
			resource.Data = maps.Clone(resource.Data)
		}

		// Keycloak rejects client-provided "id" in the body for some resource types
		// during creation (e.g. identity providers do not expose an id field in their
		// representation). The idMap is populated from the server response instead.
		// For updates, the id is usually needed in the body so Keycloak can identify
		// the entity (e.g. protocol mappers).
		if method == http.MethodPost {
			switch resource.Type {
			case "authenticationflow", "client", "group", "identityprovider", "realm", "protocolmapper":
				delete(resource.Data, "id")
			}
		}
		// The realm field in a non-realm representation is not part of the request schema.
		if resource.Type != "realm" {
			delete(resource.Data, "realm")
		}
		// Parent-reference fields are needed to route to the correct nested endpoint
		// but are not part of the child's request schema.
		for _, field := range r.spec.Resolver().ParentReferenceFieldNames(resource.Type, contract) {
			delete(resource.Data, field)
		}
		// A same-type parent binding (a group nested under a group) is invisible to
		// the single-resource contract above: its {group-id} placeholder is the
		// self id, so the parent field (groupId) is never listed there and would
		// leak into an update body — Keycloak rejects "Unrecognized field groupId".
		// The binding is a property of the nested-create endpoint, so strip the
		// create contract's parent-reference fields on every write when a parent
		// type is set (ISSUE 0005).
		if resource.ParentType != "" {
			if createContract, err := r.spec.Resolver().ResolveResourceOperation(resource.Type, resource.ParentType, http.MethodPost, OperationCollection); err == nil {
				for _, field := range r.spec.Resolver().ParentReferenceFieldNames(resource.Type, createContract) {
					delete(resource.Data, field)
				}
			}
		}
	}
}

func extractRealm(data map[string]interface{}, scope map[string]string) string {
	if realm, ok := data["realm"].(string); ok {
		return realm
	}
	return scope["realm"]
}

func (r *RuntimeClient) newAuthRequest(ctx context.Context, method, rawURL string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, err
	}
	if err := r.Authorize(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (r *RuntimeClient) readErrorBody(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	return newHTTPError(resp.StatusCode, body)
}
