// Package internal is an implementation detail of the catalog module.
// Do not import from outside catalog/. The public contract is catalog.Spec.
// AI: you may freely refactor this package as long as catalog_test.go passes.
package internal

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

var ErrSpecNotInitialized = errors.New("spec not initialized")

type Spec struct {
	document libopenapi.Document
	model    *libopenapi.DocumentModel[v3.Document]
}

func NewSpec(path string) (*Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return NewSpecFromBytes(data)
}

func NewSpecFromBytes(data []byte) (*Spec, error) {
	document, err := libopenapi.NewDocument(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	model, err := document.BuildV3Model()
	if err != nil {
		return nil, fmt.Errorf("failed to build model: %w", err)
	}

	return &Spec{document: document, model: model}, nil
}

func (s *Spec) GetSchemas() (map[string]*base.SchemaProxy, error) {
	if s.model.Model.Components == nil || s.model.Model.Components.Schemas == nil {
		return map[string]*base.SchemaProxy{}, nil
	}

	schemas := make(map[string]*base.SchemaProxy)
	for pair := s.model.Model.Components.Schemas.First(); pair != nil; pair = pair.Next() {
		schemas[pair.Key()] = pair.Value()
	}
	return schemas, nil
}

func (s *Spec) ForEachOperation(visitor func(path, method string, operation *v3.Operation, item *v3.PathItem)) {
	if s == nil || visitor == nil {
		return
	}
	if s.model == nil || s.model.Model.Paths == nil {
		return
	}

	paths := s.model.Model.Paths.PathItems
	for pair := paths.First(); pair != nil; pair = pair.Next() {
		path := pair.Key()
		pathItem := pair.Value()
		if pathItem == nil {
			continue
		}

		ops := map[string]*v3.Operation{
			http.MethodGet:    pathItem.Get,
			http.MethodPost:   pathItem.Post,
			http.MethodPut:    pathItem.Put,
			http.MethodDelete: pathItem.Delete,
			http.MethodPatch:  pathItem.Patch,
		}

		for method, op := range ops {
			if op == nil {
				continue
			}
			visitor(path, method, op, pathItem)
		}

		if pathItem.AdditionalOperations != nil && pathItem.AdditionalOperations.Len() > 0 {
			for extra := pathItem.AdditionalOperations.First(); extra != nil; extra = extra.Next() {
				if extra.Value() == nil {
					continue
				}
				visitor(path, extra.Key(), extra.Value(), pathItem)
			}
		}
	}
}

func (s *Spec) Operation(path, method string) (*v3.Operation, *v3.PathItem, error) {
	if s == nil || s.model == nil || s.model.Model.Paths == nil {
		return nil, nil, ErrSpecNotInitialized
	}

	item := s.model.Model.Paths.PathItems.GetOrZero(path)
	if item == nil {
		return nil, nil, fmt.Errorf("path %s not found", path)
	}

	method = strings.ToUpper(method)
	operations := map[string]*v3.Operation{
		http.MethodGet:    item.Get,
		http.MethodPost:   item.Post,
		http.MethodPut:    item.Put,
		http.MethodDelete: item.Delete,
		http.MethodPatch:  item.Patch,
	}
	operation := operations[method]

	if operation == nil && item.AdditionalOperations != nil {
		for extra := item.AdditionalOperations.First(); extra != nil; extra = extra.Next() {
			if strings.EqualFold(extra.Key(), method) {
				operation = extra.Value()
				break
			}
		}
	}

	if operation == nil {
		return nil, item, fmt.Errorf("operation %s %s not found", method, path)
	}
	return operation, item, nil
}

func (s *Spec) PathParameters(path, method string) ([]string, error) {
	operation, item, err := s.Operation(path, method)
	if err != nil {
		return nil, err
	}

	names := make(map[string]struct{})
	add := func(params []*v3.Parameter) {
		for _, parameter := range params {
			if parameter == nil || parameter.In != "path" || parameter.Name == "" {
				continue
			}
			names[parameter.Name] = struct{}{}
		}
	}

	if item != nil {
		add(item.Parameters)
	}
	if operation != nil {
		add(operation.Parameters)
	}

	result := make([]string, 0, len(names))
	for name := range names {
		result = append(result, name)
	}
	sort.Strings(result)
	return result, nil
}
