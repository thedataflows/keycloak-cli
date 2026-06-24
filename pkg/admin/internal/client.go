// Package internal is an implementation detail of the admin module.
// Do not import from outside admin/. The public contract is admin.Service.
// AI: you may freely refactor this package as long as admin_test.go passes.
package internal

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
	"github.com/thedataflows/keycloak-cli/pkg/catalog"
	"github.com/thedataflows/keycloak-cli/pkg/manifest"
)

type Config struct {
	BaseURL  string
	SpecPath string
	Timeout  time.Duration
}

type TokenProvider interface {
	AccessToken(ctx context.Context, baseURL, accessToken, refreshToken string) (string, error)
}

type requestEditor func(ctx context.Context, req *http.Request) error

type RuntimeClient struct {
	baseURL    string
	httpClient *http.Client
	authEditor requestEditor
	spec       *catalog.Spec
}

func NewRuntimeClient(config Config, tokens TokenProvider) (*RuntimeClient, error) {
	baseURL := strings.TrimSpace(config.BaseURL)
	if baseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	spec, err := catalog.NewSpec(config.SpecPath)
	if err != nil {
		return nil, fmt.Errorf("load spec: %w", err)
	}

	if err := catalog.InstallDefaultRegistry(config.SpecPath); err != nil {
		return nil, fmt.Errorf("load relationship overrides: %w", err)
	}

	if err := catalog.InstallDefaultFieldOverrides(config.SpecPath); err != nil {
		return nil, fmt.Errorf("load field overrides: %w", err)
	}

	if err := catalog.InstallDefaultBuiltInResources(config.SpecPath); err != nil {
		return nil, fmt.Errorf("load built-in resources: %w", err)
	}

	client := &RuntimeClient{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{Timeout: config.Timeout},
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

func (r *RuntimeClient) Spec() *catalog.Spec {
	if r == nil {
		return nil
	}
	return r.spec
}

func (r *RuntimeClient) FetchResources(ctx context.Context, resourceType string, scope map[string]string, params ...map[string]string) ([]manifest.Resource, error) {
	contract, err := r.spec.Resolver().ResolveResourceOperation(resourceType, "", http.MethodGet, catalog.OperationCollection)
	if err != nil {
		return nil, err
	}

	op, _, err := r.spec.Operation(contract.Path, contract.Method)
	if err != nil {
		return nil, err
	}

	queryParams := mergeQueryParams(params...)
	requestPath := r.buildPathWithOperation(contract.Path, op, scope)
	if err := r.spec.ValidateOperationRequest(contract.Path, http.MethodGet, catalog.RequestValidation{
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

	resources := make([]manifest.Resource, len(rawResources))
	for i, raw := range rawResources {
		resources[i] = manifest.Resource{
			Type:  resourceType,
			Realm: extractRealm(raw, scope),
			Data:  raw,
		}
	}

	return resources, nil
}

// FetchResourcesWithParent fetches a resource collection scoped by the resource's
// parent type, so nested resources resolve to the correct endpoint.
func (r *RuntimeClient) FetchResourcesWithParent(ctx context.Context, resource manifest.Resource, params ...map[string]string) ([]manifest.Resource, error) {
	contract, err := r.spec.Resolver().ResolveResourceOperation(resource.Type, resource.ParentType, http.MethodGet, catalog.OperationCollection)
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

	queryParams := mergeQueryParams(params...)
	requestPath := r.buildPathWithOperation(contract.Path, op, paramsMap)
	if err := r.spec.ValidateOperationRequest(contract.Path, http.MethodGet, catalog.RequestValidation{
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

	resources := make([]manifest.Resource, len(rawResources))
	for i, raw := range rawResources {
		resources[i] = manifest.Resource{
			Type:  resource.Type,
			Realm: extractRealm(raw, paramsMap),
			Data:  raw,
		}
	}

	return resources, nil
}

func (r *RuntimeClient) FetchPathCollection(ctx context.Context, path string, scope map[string]string, params ...map[string]string) ([]map[string]interface{}, error) {
	queryParams := mergeQueryParams(params...)
	resolvedPath := r.buildPathWithOperation(path, nil, scope)
	if err := r.spec.ValidateOperationRequest(path, http.MethodGet, catalog.RequestValidation{
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
func (r *RuntimeClient) FetchResource(ctx context.Context, resource manifest.Resource) (manifest.Resource, bool, error) {
	contract, err := r.resolveResourceContract(resource, http.MethodGet)
	if err != nil {
		return manifest.Resource{}, false, err
	}

	params, err := r.spec.Resolver().PathParams(resource, contract)
	if err != nil {
		return manifest.Resource{}, false, err
	}

	op, _, err := r.spec.Operation(contract.Path, contract.Method)
	if err != nil {
		return manifest.Resource{}, false, err
	}

	requestPath := r.buildPathWithOperation(contract.Path, op, params)
	if err := r.spec.ValidateOperationRequest(contract.Path, http.MethodGet, catalog.RequestValidation{PathParams: params}); err != nil {
		return manifest.Resource{}, false, err
	}

	fullPath := r.baseURL + requestPath
	log.Logger.Debug().Str("pkg", "admin").Str("url", fullPath).Str("method", http.MethodGet).Msg("Fetching resource")
	req, err := r.newAuthRequest(ctx, http.MethodGet, fullPath, nil)
	if err != nil {
		return manifest.Resource{}, false, err
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return manifest.Resource{}, false, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return manifest.Resource{}, false, nil
	}
	if resp.StatusCode != http.StatusOK {
		return manifest.Resource{}, false, r.readErrorBody(resp)
	}

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return manifest.Resource{}, false, err
	}
	if err := r.spec.ValidateOperationResponse(contract.Path, http.MethodGet, raw); err != nil {
		return manifest.Resource{}, false, err
	}

	return manifest.Resource{
		Type:  resource.Type,
		Realm: extractRealm(raw, params),
		Data:  raw,
	}, true, nil
}

func (r *RuntimeClient) CreateResource(ctx context.Context, resource manifest.Resource) (int, string, error) {
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
	if err := r.spec.ValidateOperationRequest(contract.Path, http.MethodPost, catalog.RequestValidation{PathParams: params, Body: resource.Data}); err != nil {
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

func (r *RuntimeClient) UpdateResource(ctx context.Context, resource manifest.Resource) (int, error) {
	return r.resourceOperation(ctx, resource, http.MethodPut)
}

func (r *RuntimeClient) DeleteResource(ctx context.Context, resource manifest.Resource) (int, error) {
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
	if err := r.spec.ValidateOperationRequest(contract.Path, http.MethodDelete, catalog.RequestValidation{PathParams: params}); err != nil {
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

func (r *RuntimeClient) ExecuteRelationship(ctx context.Context, rel manifest.RelationshipOperation) (int, error) {
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
	if err := r.spec.ValidateOperationRequest("/admin/realms/"+templatePath, method, catalog.RequestValidation{
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

func (r *RuntimeClient) resourceOperation(ctx context.Context, resource manifest.Resource, method string) (int, error) {
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
	if err := r.spec.ValidateOperationRequest(contract.Path, method, catalog.RequestValidation{PathParams: params, Body: resource.Data}); err != nil {
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
func (r *RuntimeClient) resolveResourceContract(resource manifest.Resource, method string) (catalog.OperationContract, error) {
	if resource.Type == "realm" {
		switch method {
		case http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete:
			return catalog.OperationContract{Path: "/admin/realms/{realm}", Method: method}, nil
		default:
			return catalog.OperationContract{Path: "/admin/realms", Method: method}, nil
		}
	}

	shape := catalog.OperationSingle
	if method == http.MethodPost {
		// Creation is always a collection endpoint. Use OperationCollection to avoid
		// picking nested single-resource endpoints (e.g. client scope mappings).
		shape = catalog.OperationCollection
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

	result := path
	for key, value := range scope {
		result = strings.ReplaceAll(result, "{"+key+"}", url.PathEscape(value))
	}

	if operation != nil && operation.Parameters != nil {
		for _, parameter := range operation.Parameters {
			if parameter == nil || parameter.In != "path" {
				continue
			}
			placeholder := "{" + parameter.Name + "}"
			if !strings.Contains(result, placeholder) {
				continue
			}
			if value, exists := scope[parameter.Name]; exists {
				result = strings.ReplaceAll(result, placeholder, url.PathEscape(value))
			}
		}
	}

	return result
}

func (r *RuntimeClient) sanitizeResourcePayload(resource *manifest.Resource, method string, contract catalog.OperationContract) {
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
	if r.authEditor != nil {
		if err := r.authEditor(ctx, req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

func (r *RuntimeClient) readErrorBody(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	return newHTTPError(resp.StatusCode, body)
}
