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
	"sort"
	"strings"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// maxInvokeResponseBytes caps how much of a response body Invoke reads.
const maxInvokeResponseBytes = 10 << 20

// P carries an Invoke call's parameters: path placeholders plus query params.
// A key that names a placeholder of the resolved operation's path template is
// substituted into the path; every other key is sent as a query parameter.
type P map[string]string

// Call is one generic Invoke request. Exactly one resolution mode must be
// set: Op alone, or Resource+Verb together. Realm is a convenience that fills
// the {realm} path placeholder when the template has one and Params does not.
type Call struct {
	Op       string      // operationId; XOR with Resource+Verb
	Resource string      // path resource, e.g. "users"; resolved with Verb
	Verb     Verb        // HTTP method for Resource+Verb mode
	Realm    string      // fills {realm} when present in the template and not in Params
	Params   P           // path + query params
	Body     interface{} // JSON-serializable; rejected on operations without a request schema
}

// resolvedOp is the outcome of resolving a Call against the spec: the
// operation's method, its raw path template, and its contract for validation.
type resolvedOp struct {
	verb     Verb
	path     string
	contract OperationContract
}

// Invoke resolves call against the loaded spec, validates it, and performs
// the request, returning the raw JSON response body. An empty (204-style)
// response body is a valid nil result. Failures surface as *Error: validation
// failures before anything is sent, taxonomy-classified errors afterwards.
func (c *Client) Invoke(ctx context.Context, call Call) (json.RawMessage, error) {
	spec, runtime := c.snapshot() // capture the pair once: in-flight calls keep it across Reload
	op, path, err := resolveCall(spec, call)
	if err != nil {
		return nil, &Error{Kind: KindValidation, Op: callLabel(call), Err: err}
	}
	if err := ValidateCall(spec, op.path, string(op.verb), call); err != nil {
		return nil, &Error{Kind: KindValidation, Op: callLabel(call), Err: err}
	}
	body, err := marshalBody(op, call)
	if err != nil {
		return nil, &Error{Kind: KindValidation, Op: callLabel(call), Err: err}
	}

	fullURL := runtime.BaseURL() + path
	if query := callQueryParams(op.path, call); len(query) > 0 {
		values := url.Values{}
		for key, value := range query {
			values.Add(key, value)
		}
		fullURL += "?" + values.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, string(op.verb), fullURL, bodyReader(body))
	if err != nil {
		return nil, &Error{Kind: KindNetwork, Op: callLabel(call), Err: err}
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if err := runtime.Authorize(ctx, req); err != nil {
		return nil, classifyError(err, 0, callLabel(call))
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, &Error{Kind: KindNetwork, Op: callLabel(call), Err: err}
	}
	defer func() { _ = resp.Body.Close() }()

	raw, readErr := io.ReadAll(io.LimitReader(resp.Body, maxInvokeResponseBytes))
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, classifyError(newHTTPError(resp.StatusCode, raw), 0, callLabel(call))
	}
	if readErr != nil {
		return nil, &Error{Kind: KindNetwork, Op: callLabel(call), Err: readErr}
	}
	if len(bytes.TrimSpace(raw)) == 0 {
		return nil, nil // 204-style empty bodies are a valid result
	}
	return json.RawMessage(raw), nil
}

// resolveCall finds the operation by ID or Resource+Verb and instantiates the
// path placeholders from Params (+Realm). Op mode requires exactly one
// operation to carry the operationId; Resource+Verb mode considers only
// operations whose every placeholder is fillable from Params+Realm and picks
// the most specific (most placeholders), shortest path on ties. It resolves
// against the spec captured by the caller's snapshot, never the live client
// field.
func resolveCall(spec *Spec, call Call) (resolvedOp, string, error) {
	if spec == nil {
		return resolvedOp{}, "", fmt.Errorf("client is not initialized")
	}

	switch {
	case call.Op != "" && call.Resource != "":
		return resolvedOp{}, "", fmt.Errorf("call must set either Op or Resource+Verb, not both")
	case call.Op != "":
		return resolveCallByOp(spec, call)
	case call.Resource != "":
		return resolveCallByResourceAndVerb(spec, call)
	default:
		return resolvedOp{}, "", fmt.Errorf("call requires either Op or Resource+Verb")
	}
}

func resolveCallByOp(spec *Spec, call Call) (resolvedOp, string, error) {
	type opMatch struct {
		path   string
		method string
	}
	var matches []opMatch
	spec.ForEachOperation(func(path, method string, operation *v3.Operation, item *v3.PathItem) {
		if operation != nil && operation.OperationId == call.Op {
			matches = append(matches, opMatch{path: path, method: strings.ToUpper(method)})
		}
	})
	if len(matches) == 0 {
		return resolvedOp{}, "", fmt.Errorf("unknown operation %q", call.Op)
	}
	distinct := make(map[string]struct{}, len(matches))
	for _, match := range matches {
		distinct[match.method+" "+match.path] = struct{}{}
	}
	if len(distinct) > 1 {
		return resolvedOp{}, "", fmt.Errorf("operationId %q is ambiguous: %d operations carry it", call.Op, len(distinct))
	}

	match := matches[0]
	contract, err := spec.OperationContract(match.path, match.method)
	if err != nil {
		return resolvedOp{}, "", err
	}
	op := resolvedOp{verb: Verb(match.method), path: match.path, contract: contract}
	path, err := instantiatePath(match.path, callParamValues(call))
	if err != nil {
		return resolvedOp{}, "", err
	}
	return op, path, nil
}

func resolveCallByResourceAndVerb(spec *Spec, call Call) (resolvedOp, string, error) {
	if call.Verb == "" {
		return resolvedOp{}, "", fmt.Errorf("resource %q requires a Verb", call.Resource)
	}
	resource := strings.ToLower(strings.TrimSpace(call.Resource))
	values := callParamValues(call)

	type candidate struct {
		path         string
		method       string
		placeholders int
	}
	var candidates []candidate
	spec.ForEachOperation(func(path, method string, operation *v3.Operation, item *v3.PathItem) {
		if !strings.EqualFold(method, string(call.Verb)) {
			return
		}
		if operationResource(path) != resource {
			return
		}
		placeholders := pathPlaceholderNames(path)
		for _, name := range placeholders {
			if _, ok := values[name]; !ok {
				return // not fillable from Params+Realm
			}
		}
		candidates = append(candidates, candidate{path: path, method: strings.ToUpper(method), placeholders: len(placeholders)})
	})
	if len(candidates) == 0 {
		return resolvedOp{}, "", fmt.Errorf("no %s operation on resource %q has all path placeholders fillable from the given params", strings.ToUpper(string(call.Verb)), call.Resource)
	}
	// Most specific (most placeholders) wins; ties break to the shortest
	// path, then lexicographically, so the choice is deterministic.
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].placeholders != candidates[j].placeholders {
			return candidates[i].placeholders > candidates[j].placeholders
		}
		if len(candidates[i].path) != len(candidates[j].path) {
			return len(candidates[i].path) < len(candidates[j].path)
		}
		return candidates[i].path < candidates[j].path
	})
	picked := candidates[0]

	contract, err := spec.OperationContract(picked.path, picked.method)
	if err != nil {
		return resolvedOp{}, "", err
	}
	op := resolvedOp{verb: Verb(picked.method), path: picked.path, contract: contract}
	path, err := instantiatePath(picked.path, values)
	if err != nil {
		return resolvedOp{}, "", err
	}
	return op, path, nil
}

// callParamValues merges the call's parameters with the Realm convenience:
// Realm fills "realm" only when the params do not already carry it.
func callParamValues(call Call) map[string]string {
	values := make(map[string]string, len(call.Params)+1)
	maps.Copy(values, call.Params)
	if call.Realm != "" {
		if _, exists := values["realm"]; !exists {
			values["realm"] = call.Realm
		}
	}
	return values
}

// callQueryParams returns the params that are not path placeholders of the
// template; Invoke sends them as the query string.
func callQueryParams(template string, call Call) map[string]string {
	placeholders := pathPlaceholderNames(template)
	query := make(map[string]string, len(call.Params))
	for key, value := range call.Params {
		if isPathPlaceholder(placeholders, key) {
			continue
		}
		query[key] = value
	}
	if len(query) == 0 {
		return nil
	}
	return query
}

func isPathPlaceholder(placeholders []string, name string) bool {
	for _, placeholder := range placeholders {
		if placeholder == name {
			return true
		}
	}
	return false
}

// callLabel names the call in errors: "invoke getUser" or "invoke GET users".
func callLabel(call Call) string {
	switch {
	case call.Op != "":
		return "invoke " + call.Op
	case call.Resource != "":
		return "invoke " + strings.ToUpper(string(call.Verb)) + " " + call.Resource
	default:
		return "invoke"
	}
}

// marshalBody serializes the call's body, refusing a body when the resolved
// operation declares no request schema (it would be silently dropped).
func marshalBody(op resolvedOp, call Call) ([]byte, error) {
	if call.Body == nil {
		return nil, nil
	}
	if op.contract.RequestBodySchema == nil {
		return nil, fmt.Errorf("%s %s does not accept a request body", op.verb, op.path)
	}
	return json.Marshal(call.Body)
}

// bodyReader wraps the marshaled body for the request: no bytes means an
// explicit empty body rather than a nil one.
func bodyReader(b []byte) io.Reader {
	if len(b) == 0 {
		return http.NoBody
	}
	return bytes.NewReader(b)
}
