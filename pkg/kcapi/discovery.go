package kcapi

import (
	"strings"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// Verb is an HTTP method of a spec operation.
type Verb string

const (
	Get    Verb = "GET"
	Post   Verb = "POST"
	Put    Verb = "PUT"
	Delete Verb = "DELETE"
	Patch  Verb = "PATCH"
)

// Param describes a spec parameter of an operation.
type Param struct {
	Name     string
	In       string // "path" | "query" | "body"
	Required bool
}

// Operation is one discoverable operation of the loaded spec.
type Operation struct {
	ID       string // operationId, may be empty
	Resource string // inferred from path, e.g. "users"
	Verb     Verb
	Path     string // raw template, e.g. "/admin/realms/{realm}/users/{id}"
	Summary  string
	Tags     []string
	Params   []Param
}

// OpFilter narrows Operations. Zero fields match everything; set fields are
// ANDed.
type OpFilter struct {
	Resource string
	Method   Verb
	Tag      string
	Search   string // substring match on ID/Path/Summary, case-insensitive
}

// Operations lists the spec's operations, filtered by filter.
func (c *Client) Operations(filter OpFilter) ([]Operation, error) {
	var out []Operation
	c.spec.ForEachOperation(func(path, method string, op *v3.Operation, item *v3.PathItem) {
		o := buildOperation(path, method, op, item)
		if matchOp(filter, o) {
			out = append(out, o)
		}
	})
	return out, nil
}

func buildOperation(path, method string, op *v3.Operation, item *v3.PathItem) Operation {
	o := Operation{
		ID:       op.OperationId,
		Verb:     Verb(strings.ToUpper(method)),
		Path:     path,
		Summary:  op.Summary,
		Resource: operationResource(path),
	}
	o.Tags = append(o.Tags, op.Tags...)

	params := make([]*v3.Parameter, 0, len(op.Parameters)+len(item.Parameters))
	params = append(params, op.Parameters...)
	params = append(params, item.Parameters...)
	for _, p := range params {
		if p == nil {
			continue
		}
		o.Params = append(o.Params, Param{
			Name:     p.Name,
			In:       p.In,
			Required: p.Required != nil && *p.Required,
		})
	}
	return o
}

// operationResource returns the last meaningful path segment as the
// operation's resource, e.g. "users" for "/admin/realms/{realm}/users/{id}".
// Unlike inferResourceTypeFromPath it keeps the plural form: discovery
// filters speak the spec's path vocabulary, not the catalog's singular
// resource types.
func operationResource(path string) string {
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
		candidate = strings.ToLower(segment)
	}
	return candidate
}

func matchOp(f OpFilter, o Operation) bool {
	if f.Resource != "" && o.Resource != f.Resource {
		return false
	}
	if f.Method != "" && o.Verb != f.Method {
		return false
	}
	if f.Tag != "" && !containsString(o.Tags, f.Tag) {
		return false
	}
	if f.Search != "" {
		hay := strings.ToLower(o.ID + " " + o.Path + " " + o.Summary)
		if !strings.Contains(hay, strings.ToLower(f.Search)) {
			return false
		}
	}
	return true
}
